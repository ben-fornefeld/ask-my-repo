package ranking

import (
	"context"
	"fmt"
	"log"
	"rankmyrepo/internal/parser"
	"sort"
	"sync"

	"github.com/replicate/replicate-go"
)

type Engine struct {
	r8         *replicate.Client
	maxWorkers int
}

func NewEngine(r8 *replicate.Client, maxWorkers int) *Engine {
	return &Engine{
		r8:         r8,
		maxWorkers: maxWorkers,
	}
}

func (e *Engine) RankChunks(ctx context.Context, query string, chunks map[string]parser.ParsedChunk, scoreThreshold float64) ([]RankedChunk, error) {
	log.Printf("Starting to rank %d chunks for query: %s", len(chunks), query)

	// worker pool for parallel ranking
	results := make(chan RankedChunk, len(chunks))
	errors := make(chan error, len(chunks))

	// semaphore to limit concurrent LLM calls
	sem := make(chan struct{}, e.maxWorkers)

	for _, chunk := range chunks {
		go func(c parser.ParsedChunk) {
			sem <- struct{}{}
			defer func() { <-sem }()

			log.Printf("Ranking chunk %s", c.FilePath)

			score, err := e.RankSingleChunk(ctx, query, c)
			if err != nil {
				errors <- err
				return
			}

			ranked := RankedChunk{
				ParsedChunk: c,
				Score:       score,
			}
			log.Printf("Ranked chunk %s with score %.2f", c.FilePath, score)
			results <- ranked
		}(chunk)
	}

	// collect results
	rankedChunks := make([]RankedChunk, 0, len(chunks))
	for i := 0; i < len(chunks); i++ {
		select {
		case chunk := <-results:
			rankedChunks = append(rankedChunks, chunk)
		case err := <-errors:
			return nil, err
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	filteredChunks := make([]RankedChunk, 0, len(rankedChunks))
	for _, chunk := range rankedChunks {
		if chunk.Score >= scoreThreshold {
			filteredChunks = append(filteredChunks, chunk)
		}
	}

	sort.Slice(filteredChunks, func(i, j int) bool {
		return filteredChunks[i].Score < filteredChunks[j].Score
	})

	log.Printf("Finished ranking %d chunks. After Threshold: %d", len(rankedChunks), len(filteredChunks))
	return filteredChunks, nil
}

func (e *Engine) RankSingleChunk(ctx context.Context, query string, chunk parser.ParsedChunk) (float64, error) {
	prompt := buildRankingPrompt(query, chunk)

	model := "meta/meta-llama-3-70b-instruct"

	input := replicate.PredictionInput{
		"prompt":        prompt,
		"system_prompt": systemPrompt,
		"temperature":   0.1,
	}

	output, err := e.r8.Run(ctx, model, input, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to run model: %w", err)
	}

	tokens, ok := output.([]interface{})
	if !ok {
		return 0, fmt.Errorf("unexpected output type from model: got %T, want []interface{}", output)
	}

	var result string
	for _, token := range tokens {
		if str, ok := token.(string); ok {
			result += str
		}
	}

	return parseScore(result)
}

func (e *Engine) RankChunksStream(ctx context.Context, query string, chunks map[string]parser.ParsedChunk, scoreThreshold float64, parsedChan chan<- parser.ParsedChunk, rankedChan chan<- RankedChunk) error {
	log.Printf("Starting to rank %d chunks for query: %s", len(chunks), query)

	errors := make(chan error, len(chunks))
	sem := make(chan struct{}, e.maxWorkers)

	var wg sync.WaitGroup
	for _, chunk := range chunks {
		wg.Add(1)
		go func(c parser.ParsedChunk) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			parsedChan <- c

			score, err := e.RankSingleChunk(ctx, query, c)
			if err != nil {
				errors <- err
				return
			}

			if score >= scoreThreshold {
				ranked := RankedChunk{
					ParsedChunk: c,
					Score:      score,
				}
				select {
				case rankedChan <- ranked:
				case <-ctx.Done():
					errors <- ctx.Err()
				}
			}
		}(chunk)
	}

	go func() {
		wg.Wait()
		close(errors)
	}()

	for err := range errors {
		if err != nil {
			return err
		}
	}

	return nil
}
