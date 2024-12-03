package completion

import (
	"context"
	"rankmyrepo/internal/ranking"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/packages/ssestream"
)

type Completion struct {
	anthropicClient *anthropic.Client
}

func NewCompletion(anthropicClient *anthropic.Client) *Completion {
	return &Completion{
		anthropicClient: anthropicClient,
	}
}

func (c *Completion) Run(ctx context.Context, query string, chunks []ranking.RankedChunk) *ssestream.Stream[anthropic.MessageStreamEvent] {
	prompt := buildCompletionPrompt(query, chunks)

	messageParams := anthropic.MessageNewParams{
		Model: anthropic.F(anthropic.ModelClaude3_5SonnetLatest),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		}),
		MaxTokens: anthropic.Int(8000),
	}

	stream := c.anthropicClient.Messages.NewStreaming(ctx, messageParams)

	return stream
}
