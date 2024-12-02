package parser

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/sabhiram/go-gitignore"
)

type Parser struct {
	tempDir       string
	textMimeTypes map[string]bool
}

func NewParser(textMimeTypes map[string]bool) (*Parser, error) {
	tempDir, err := os.MkdirTemp("", "repo-parser-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	return &Parser{
		tempDir:       tempDir,
		textMimeTypes: textMimeTypes,
	}, nil
}

// ParseRepository parses the repository at the given URL and returns a map of RawChunk
// with the file paths as keys and the content as values.

// NOTE: This currently only supports Public GitHub repositories.

func (p *Parser) ParseRepository(ctx context.Context, repoURL string, ignorePatterns []string) (map[string]ParsedChunk, error) {
	repoDir := filepath.Join(p.tempDir, filepath.Base(repoURL))

	_, err := git.PlainCloneContext(ctx, repoDir, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to clone repository: %w", err)
	}

	defer func() {
		if err := os.RemoveAll(repoDir); err != nil {
			fmt.Printf("warning: failed to clean up repository directory: %v\n", err)
		}
	}()

	for i, pattern := range ignorePatterns {
		if !strings.HasPrefix(pattern, "/") {
			ignorePatterns[i] = "**/" + pattern
		}
	}

	ignore := ignore.CompileIgnoreLines(ignorePatterns...)

	chunks := make(map[string]ParsedChunk, 0)

	filepath.WalkDir(repoDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(repoDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		if ignore.MatchesPath(relPath) {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		isText, err := p.IsTextFile(file)
		if err != nil {
			return fmt.Errorf("failed to check if file is text: %w", err)
		}
		if !isText {
			return nil
		}

		content, err := io.ReadAll(file)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		chunk := ParsedChunk{
			FilePath: path,
			Content:  string(content),
		}

		chunks[path] = chunk

		return nil
	})

	return chunks, nil
}

func (p *Parser) Cleanup() error {
	if p.tempDir != "" {
		if err := os.RemoveAll(p.tempDir); err != nil {
			return fmt.Errorf("failed to cleanup temporary directory: %w", err)
		}
	}
	return nil
}
