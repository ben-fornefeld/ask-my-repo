package ranking

import "context"

type CodeChunk struct {
	Content   string
	FilePath  string
	StartLine int
	EndLine   int
	Language  string
	Symbols   []string
	Score     float64
}

// TODO: use interface in engine.go
type RankingInterface interface {
	RankChunks(ctx context.Context, query string, chunks []CodeChunk) ([]CodeChunk, error)
}
