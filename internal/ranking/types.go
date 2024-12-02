package ranking

import "context"

type RankedChunk struct {
	Content   string
	FilePath  string
	StartLine int
	EndLine   int
	Language  string
	Symbols   []string
	Score     float64
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
