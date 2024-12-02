package main

import (
	"log"
	"net/http"
	"rankmyrepo/internal/api"
	"rankmyrepo/internal/parser"
	"rankmyrepo/internal/processor"
	"rankmyrepo/internal/ranking"

	"github.com/joho/godotenv"
	"github.com/replicate/replicate-go"
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

	parser, err := parser.NewParser(textMimeTypes)
	if err != nil {
		log.Fatal(err)
	}
	defer parser.Cleanup()

	r8, err := replicate.NewClient(replicate.WithTokenFromEnv())
	if err != nil {
		log.Fatal(err)
	}

	ranker := ranking.NewEngine(r8, 100, 0.4)

	processor := processor.NewProcessor(parser, ranker)

	handler, err := api.NewHandler(processor)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/query", handler.RankCode)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
