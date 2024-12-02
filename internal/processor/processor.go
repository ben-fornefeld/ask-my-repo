package processor

import (
	"context"
	"rankmyrepo/internal/completion"
	"rankmyrepo/internal/parser"
	"rankmyrepo/internal/ranking"
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

func (p *Processor) ProcessRankingRequest(ctx context.Context, req *ranking.RankingRequest) (*ranking.RankingResponse, error) {
	parsedChunks, err := p.parser.ParseRepository(ctx, req.RepoPath, req.IgnorePatterns)
	if err != nil {
		return nil, err
	}

	rankedChunks, err := p.ranker.RankChunks(ctx, req.Query, parsedChunks, req.ScoreThreshold)
	if err != nil {
		return nil, err
	}

	completion, err := p.completion.Run(ctx, req.Query, rankedChunks)
	if err != nil {
		return nil, err
	}

	result := ranking.RankingResponse{
		Chunks:     rankedChunks,
		Completion: completion,
	}

	return &result, nil
}
