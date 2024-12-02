package parser

import (
	"context"
	"os"
	"path/filepath"
	"rankmyrepo/internal/parser"
	"strings"
	"testing"
)

func TestParseRepository(t *testing.T) {
	textMimeTypes := map[string]bool{
		"text/":                     true,
		"application/json":          true,
		"application/xml":           true,
		"application/x-yaml":        true,
		"application/toml":          true,
		"application/x-javascript":  true,
		"application/x-shellscript": true,
	}

	patterns := []string{"*.md"}

	p, err := parser.NewParser(textMimeTypes)
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}
	defer p.Cleanup()

	repoURL := "https://github.com/ben-fornefeld/neo"
	chunks, err := p.ParseRepository(context.Background(), repoURL, patterns)
	if err != nil {
		t.Fatalf("failed to parse repository: %v", err)
	}

	if len(chunks) == 0 {
		t.Error("expected chunks to be returned, got empty slice")
	}

	for _, chunk := range chunks {
		if strings.HasSuffix(chunk.FilePath, ".md") {
			t.Errorf("found blacklisted .md file in chunks: %s", chunk.FilePath)
		}
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

func TestTextFileFilter(t *testing.T) {
	textMimeTypes := map[string]bool{
		"text/":                     true,
		"application/json":          true,
		"application/xml":           true,
		"application/x-yaml":        true,
		"application/toml":          true,
		"application/x-javascript":  true,
		"application/x-shellscript": true,
	}

	parser, err := parser.NewParser(textMimeTypes)
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}

	// Create temporary directory for test files
	tmpDir := t.TempDir()

	// Test case 1: Valid UTF-8 text file
	validTextPath := filepath.Join(tmpDir, "valid.txt")
	err = os.WriteFile(validTextPath, []byte("Hello, 世界!\nThis is a valid UTF-8 text file."), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Test case 2: Valid JSON file
	jsonPath := filepath.Join(tmpDir, "config.json")
	err = os.WriteFile(jsonPath, []byte(`{"name": "test", "value": 123}`), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Test case 3: Binary file (PNG image)
	binaryPath := filepath.Join(tmpDir, "binary.png")
	// PNG file header followed by minimal IHDR chunk
	pngData := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
		0x00, 0x00, 0x00, 0x0D, // Length of IHDR chunk
		0x49, 0x48, 0x44, 0x52, // "IHDR"
		0x00, 0x00, 0x00, 0x01, // Width: 1
		0x00, 0x00, 0x00, 0x01, // Height: 1
		0x08,                   // Bit depth
		0x06,                   // Color type
		0x00,                   // Compression method
		0x00,                   // Filter method
		0x00,                   // Interlace method
		0x1f, 0x15, 0xc4, 0x89, // CRC
	}
	err = os.WriteFile(binaryPath, pngData, 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Test case 4: Custom programming language file (not in mime types)
	customLangPath := filepath.Join(tmpDir, "code.xyz")
	err = os.WriteFile(customLangPath, []byte(`
	function main() {
		print("This is a custom programming language");
		// It's valid UTF-8 text but with an unknown mime type
	}
`), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Run tests
	testCases := []struct {
		name     string
		path     string
		expected bool
	}{
		{"Valid UTF-8 text file", validTextPath, true},
		{"Valid JSON file", jsonPath, true},
		{"Binary PNG file", binaryPath, false},
		{"Custom language file", customLangPath, true}, // Should be true because it's valid UTF-8
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file, err := os.Open(tc.path)
			if err != nil {
				t.Fatalf("failed to open test file: %v", err)
			}
			defer file.Close()

			isText, err := parser.IsTextFile(file)
			if err != nil {
				t.Fatalf("IsTextFile failed: %v", err)
			}

			if isText != tc.expected {
				t.Errorf("IsTextFile(%s) = %v, want %v", tc.name, isText, tc.expected)
			}
		})
	}
}
