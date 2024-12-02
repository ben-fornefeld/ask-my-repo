package processor

import (
	"context"
	"rankmyrepo/internal/ranking"
)

type Processor struct {
	// TODO: add parser
	ranker ranking.RankingEngine
}

func NewProcessor(ranker ranking.RankingEngine) *Processor {
	return &Processor{
		ranker: ranker,
	}
}

func (p *Processor) ProcessRankingRequest(ctx context.Context, req *ranking.RankingRequest) (*ranking.RankingResponse, error) {
	// TODO: chunk repository

	// prefilter chunks

	// rank chunks

	// assemble context

	return nil, nil
}
