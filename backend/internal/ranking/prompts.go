package ranking

import (
	"fmt"
	"rankmyrepo/internal/parser"
	"strconv"
	"strings"
)

// TODO: use <thinking> tags for chain of thought if necessary

var systemPrompt = `You are a code ranking assistant. Your task is to analyze code chunks and assign them relevance scores based on how well they help answer the user's query. Be direct and precise in your scoring. Only output a score tag with a number between 0.0 and 1.0. Higher scores mean the code is more relevant for answering the query.`

func buildRankingPrompt(query string, chunk parser.ParsedChunk) string {
	return `Rate how relevant this code is to answering the query.
Score from 0.0 to 1.0 with max 1 decimal point.
Return ONLY <score>X</score> where X is the score.
0.0 = not relevant at all
1.0 = highly relevant

Query: ` + query + `

File: ` + chunk.FilePath + `
Code:
` + chunk.Content + `

Remember: Just return <score>X</score> with X between 0.0-1.0, one decimal max.`
}

func parseScore(response string) (float64, error) {
	startTag := "<score>"
	endTag := "</score>"

	startIndex := strings.Index(response, startTag)
	if startIndex == -1 {
		return 0, fmt.Errorf("no start tag found in response")
	}

	endIndex := strings.Index(response, endTag)
	if endIndex == -1 {
		return 0, fmt.Errorf("no end tag found in response")
	}

	scoreStr := strings.TrimSpace(response[startIndex+len(startTag) : endIndex])

	score, err := strconv.ParseFloat(scoreStr, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse score: %w", err)
	}

	if score < 0 || score > 1 {
		return 0, fmt.Errorf("score %f is outside valid range [0,1]", score)
	}

	// Round to 1 decimal place
	score = float64(int(score*10)) / 10

	return score, nil
}
