package completion

import (
	"fmt"
	"rankmyrepo/internal/ranking"
)

func buildCompletionPrompt(query string, chunks []ranking.RankedChunk) string {
	return `Answer the following user query with the additional context provided in the chunks.

	<context>` + buildContext(chunks) + `</context>

	<query>` + query + `</query>`
}

func buildContext(chunks []ranking.RankedChunk) (context string) {
	for _, chunk := range chunks {
		context += fmt.Sprintf("Chunk: %s\nContent: %s\n\n", chunk.ParsedChunk.FilePath, chunk.ParsedChunk.Content)
	}

	return context
}
