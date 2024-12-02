package ranking

import (
	"context"
	"fmt"
	"log"
	"rankmyrepo/internal/parser"
	"sort"

	"github.com/replicate/replicate-go"
)

type Engine struct {
	r8             *replicate.Client
	maxWorkers     int
	scoreThreshold float64
}

func NewEngine(r8 *replicate.Client, maxWorkers int, scoreThreshold float64) *Engine {
	return &Engine{
		r8:             r8,
		maxWorkers:     maxWorkers,
		scoreThreshold: scoreThreshold,
	}
}

func (e *Engine) RankChunks(ctx context.Context, query string, chunks map[string]parser.ParsedChunk) ([]RankedChunk, error) {
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

			score, err := e.rankSingleChunk(ctx, query, c)
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

	sort.Slice(rankedChunks, func(i, j int) bool {
		return rankedChunks[i].Score > rankedChunks[j].Score
	})

	log.Printf("Finished ranking %d chunks. Top score: %.2f", len(rankedChunks), rankedChunks[0].Score)
	return rankedChunks, nil
}

// TODO: switch to more cheaper / faster model (e.g., groq llama 70b?)
func (e *Engine) rankSingleChunk(ctx context.Context, query string, chunk parser.ParsedChunk) (float64, error) {
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
		// Convert each token to string
		if str, ok := token.(string); ok {
			result += str
		}
	}

	return parseScore(result)
}
