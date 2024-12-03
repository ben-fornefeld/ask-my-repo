package common

import (
	"rankmyrepo/internal/parser"
	"rankmyrepo/internal/ranking"
)

// Query

type QueryEventType string

const (
	EventTypeRankingParsed    QueryEventType = "ranking.parsed"
	EventTypeRankingRanked    QueryEventType = "ranking.ranked"
	EventTypeCompletionDelta  QueryEventType = "completion.delta"
	EventTypeError            QueryEventType = "error"
)

type QueryResponseChunk struct {
	Type        QueryEventType `json:"type"`
	ParsedChunk *parser.ParsedChunk `json:"parsed_chunk,omitempty"`
	RankedChunk *ranking.RankedChunk `json:"ranked_chunk,omitempty"`
	Completion  string        `json:"completion,omitempty"`
	Error       string        `json:"error,omitempty"`
}