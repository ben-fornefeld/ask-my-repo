package ranking

import (
	"context"
	"rankmyrepo/internal/parser"
)

type RankedChunk struct {
	ParsedChunk parser.ParsedChunk
	Score       float64
}

type RankingEngine interface {
	RankChunks(ctx context.Context, query string, chunks map[string]parser.ParsedChunk) ([]RankedChunk, error)
}

type RankingRequest struct {
	Query    string
	RepoPath string
}

type RankingResponse struct {
	Chunks     []RankedChunk
	TotalScore float64
}
