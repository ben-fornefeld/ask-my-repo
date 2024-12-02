package ranking

import (
	"context"
	"rankmyrepo/internal/parser"
)

type RankedChunk struct {
	ParsedChunk parser.ParsedChunk
	Score       float64
}

// TODO: use interface in engine.go
type RankingEngine interface {
	RankChunks(ctx context.Context, query string, chunks []RankedChunk) ([]RankedChunk, error)
}

type RankingRequest struct {
	Query    string
	RepoPath string
}

type RankingResponse struct {
	Chunks     []RankedChunk
	TotalScore float64
}
