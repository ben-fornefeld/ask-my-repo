package parser

import (
	"context"
	"rankmyrepo/internal/parser"
	"testing"
)

func TestParseRepository(t *testing.T) {
	p, err := parser.NewParser(make([]string, 0))
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}
	defer p.Cleanup()

	repoURL := "https://github.com/ben-fornefeld/neo"
	chunks, err := p.ParseRepository(context.Background(), repoURL)
	if err != nil {
		t.Fatalf("failed to parse repository: %v", err)
	}

	if len(chunks) == 0 {
		t.Error("expected chunks to be returned, got empty slice")
	}

	t.Logf("Found %d chunks in repository:", len(chunks))
	for i, chunk := range chunks {
		t.Logf("Chunk %s:", i)
		t.Logf("  File: %s", chunk.FilePath)
		t.Logf("  Content length: %d bytes", len(chunk.Content))

		if len(chunk.Content) > 100 {
			t.Logf("  Content preview: %s...", chunk.Content[:100])
		} else {
			t.Logf("  Content: %s", chunk.Content)
		}
		t.Log("---")
	}
}
