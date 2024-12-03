package processor

import (
	"context"
	"rankmyrepo/internal/common"
	"rankmyrepo/internal/completion"
	"rankmyrepo/internal/parser"
	"rankmyrepo/internal/ranking"

	"github.com/anthropics/anthropic-sdk-go"
)

type Processor struct {
	parser     *parser.Parser
	ranker     *ranking.Engine
	completion *completion.Completion
}

func NewProcessor(parser *parser.Parser, ranker *ranking.Engine, compcompletion *completion.Completion) *Processor {
	return &Processor{
		parser:     parser,
		ranker:     ranker,
		completion: compcompletion,
	}
}

func (p *Processor) ProcessRankingRequestStream(ctx context.Context, req *ranking.RankingRequest, resultChan chan<- common.QueryResponseChunk) error {
	parsedChunks, err := p.parser.ParseRepository(ctx, req.RepoPath, req.IgnorePatterns)
	if err != nil {
		return err
	}

	bufferSize := len(parsedChunks)
	rankingParsedChan := make(chan parser.ParsedChunk, bufferSize)
	rankingRankedChan := make(chan ranking.RankedChunk, bufferSize)
	rankingErrChan := make(chan error, 1)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		defer close(rankingParsedChan)
		defer close(rankingRankedChan)

		if err := p.ranker.RankChunksStream(ctx, req.Query, parsedChunks, req.ScoreThreshold, rankingParsedChan, rankingRankedChan); err != nil {
			rankingErrChan <- err
			cancel()
			return
		}
		close(rankingErrChan)
	}()

	var rankedChunks []ranking.RankedChunk

	for chunk := range rankingParsedChan {
		resultChan <- common.QueryResponseChunk{
			Type:        common.EventTypeRankingParsed,
			ParsedChunk: &chunk,
		}
	}
	for chunk := range rankingRankedChan {
		rankedChunks = append(rankedChunks, chunk)
		resultChan <- common.QueryResponseChunk{
			Type:        common.EventTypeRankingRanked,
			RankedChunk: &chunk,
		}
	}

	if err := <-rankingErrChan; err != nil {
		return err
	}

	stream := p.completion.Run(ctx, req.Query, rankedChunks)

	for stream.Next() {
		event := stream.Current()

		switch delta := event.Delta.(type) {
		case anthropic.ContentBlockDeltaEventDelta:
			if delta.Text != "" {
				resultChan <- common.QueryResponseChunk{
					Type:       common.EventTypeCompletionDelta,
					Completion: delta.Text,
				}
			}
		}
	}

	if stream.Err() != nil {
		return stream.Err()
	}

	return nil
}
