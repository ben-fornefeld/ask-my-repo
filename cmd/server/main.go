package main

import (
	"log"
	"net/http"
	"os"
	"rankmyrepo/internal/api"
	"rankmyrepo/internal/parser"
	"rankmyrepo/internal/processor"
	"rankmyrepo/internal/ranking"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	textMimeTypes := map[string]bool{
		"text/":                     true,
		"application/json":          true,
		"application/xml":           true,
		"application/x-yaml":        true,
		"application/toml":          true,
		"application/x-javascript":  true,
		"application/x-shellscript": true,
	}

	parserIgnorePatterns := []string{"*.min.js"}

	parser, err := parser.NewParser(parserIgnorePatterns, textMimeTypes)
	if err != nil {
		log.Fatal(err)
	}
	defer parser.Cleanup()

	llmClient := anthropic.NewClient(option.WithAPIKey(os.Getenv("ANTHROPIC_API_KEY")))

	ranker := ranking.NewEngine(llmClient, 5)

	processor := processor.NewProcessor(parser, ranker)

	handler, err := api.NewHandler(processor)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/query", handler.RankCode)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
