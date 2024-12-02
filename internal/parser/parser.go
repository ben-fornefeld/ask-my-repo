package parser

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

type Parser struct {
	tempDir           string
	filteredFileTypes []string
}

func NewParser(filteredFileTypes []string) (*Parser, error) {
	tempDir, err := os.MkdirTemp("", "repo-parser-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	return &Parser{
		tempDir:           tempDir,
		filteredFileTypes: filteredFileTypes,
	}, nil
}

func (p *Parser) ParseRepository(ctx context.Context, repoURL string) (map[string]RawChunk, error) {
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

	chunks := make(map[string]RawChunk, 0)

	filepath.WalkDir(repoDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// TODO: filter file types

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		chunk := RawChunk{
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
