package ranking

import (
	"fmt"
	"strconv"
	"strings"
)

// TODO: use <thinking> tags for chain of thought

func buildRankingPrompt(query string, chunk RankedChunk) string {
	return `You are a code ranking assistant. Analyze the code snippet's relevance to the query.
Return ONLY a score between 0.0 and 1.0 wrapped in XML tags <score></score>.
1.0 means highly relevant, 0.0 means not relevant at all.
DO NOT include any other text or explanations in your response.

<query>` + query + `</query>

<code>` + chunk.Content + `</code>

Remember: Respond ONLY with <score>X</score> where X is a number between 0.0 and 1.0`
}

func parseScore(response string) (float64, error) {
	// Find content between <score> and </score> tags
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

	// Extract the score string
	scoreStr := strings.TrimSpace(response[startIndex+len(startTag) : endIndex])

	// Parse the score as float64
	score, err := strconv.ParseFloat(scoreStr, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse score: %w", err)
	}

	// Validate score range
	if score < 0 || score > 1 {
		return 0, fmt.Errorf("score %f is outside valid range [0,1]", score)
	}

	return score, nil
}
