package main

import (
	"log"
	"os"
	"rankmyrepo/internal/api"
	"rankmyrepo/internal/completion"
	"rankmyrepo/internal/parser"
	"rankmyrepo/internal/processor"
	"rankmyrepo/internal/ranking"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/gin-gonic/gin"
	"github.com/replicate/replicate-go"
)

func main() {
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

	ranker := ranking.NewEngine(r8, 50)

	anthropicClient := anthropic.NewClient(option.WithAPIKey(os.Getenv("ANTHROPIC_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}

	completion := completion.NewCompletion(anthropicClient)

	processor := processor.NewProcessor(parser, ranker, completion)

	handler, err := api.NewHandler(processor)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	r.POST("/query", handler.Query)

	log.Printf("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
