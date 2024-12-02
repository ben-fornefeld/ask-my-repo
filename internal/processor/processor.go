package processor

import (
	"context"
	"rankmyrepo/internal/parser"
	"rankmyrepo/internal/ranking"
)

type Processor struct {
	parser *parser.Parser
	ranker *ranking.Engine
}

func NewProcessor(parser *parser.Parser, ranker *ranking.Engine) *Processor {
	return &Processor{
		parser: parser,
		ranker: ranker,
	}
}

func (p *Processor) ProcessRankingRequest(ctx context.Context, req *ranking.RankingRequest) (*ranking.RankingResponse, error) {
	parsedChunks, err := p.parser.ParseRepository(ctx, req.RepoPath)
	if err != nil {
		return nil, err
	}

	rankedChunks, err := p.ranker.RankChunks(ctx, req.Query, parsedChunks)
	if err != nil {
		return nil, err
	}

	result := ranking.RankingResponse{
		Chunks:     rankedChunks,
		TotalScore: 0,
	}

	return &result, nil
}
