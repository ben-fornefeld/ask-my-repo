package ranking

import (
	"context"
	"rankmyrepo/internal/parser"
	"sort"

	"github.com/anthropics/anthropic-sdk-go"
)

type Engine struct {
	llmClient  *anthropic.Client
	maxWorkers int
}

func NewEngine(llmClient *anthropic.Client, maxWorkers int) *Engine {
	return &Engine{
		llmClient:  llmClient,
		maxWorkers: maxWorkers,
	}
}

func (e *Engine) RankChunks(ctx context.Context, query string, chunks map[string]parser.ParsedChunk) ([]RankedChunk, error) {
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

			results <- RankedChunk{
				ParsedChunk: c,
				Score:       score,
			}
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

	return rankedChunks, nil
}

func (e *Engine) rankSingleChunk(ctx context.Context, query string, chunk parser.ParsedChunk) (float64, error) {
	prompt := buildRankingPrompt(query, chunk)

	message, err := e.llmClient.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.F(anthropic.ModelClaude3_5HaikuLatest),
		MaxTokens: anthropic.F(int64(1024)),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		}),
	})
	if err != nil {
		return 0, err
	}

	return parseScore(message.ToParam().Content.String())
}
