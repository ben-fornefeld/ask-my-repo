package ranking

import (
	"context"
	"sort"
)

type Engine struct {
	maxWorkers int
}

func NewEngine(maxWorkers int) *Engine {
	return &Engine{
		maxWorkers: maxWorkers,
	}
}

func (e *Engine) RankChunks(ctx context.Context, query string, chunks []CodeChunk) ([]CodeChunk, error) {
	// worker pool for parallel ranking
	results := make(chan CodeChunk, len(chunks))
	errors := make(chan error, len(chunks))

	// semaphore to limit concurrent LLM calls
	sem := make(chan struct{}, e.maxWorkers)

	for _, chunk := range chunks {
		go func(c CodeChunk) {
			sem <- struct{}{}
			defer func() { <-sem }()

			score, err := e.rankSingleChunk(ctx, query, c)
			if err != nil {
				errors <- err
				return
			}

			c.Score = score
			results <- c
		}(chunk)
	}

	// collect results
	rankedChunks := make([]CodeChunk, 0, len(chunks))
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

	return rankedChunks, nil
}

func (e *Engine) rankSingleChunk(ctx context.Context, query string, chunk CodeChunk) (float64, error) {
	prompt := buildRankingPrompt(query, chunk)

	// TODO: call LLM here
	response := (prompt)

	return parseScore(response)
}
