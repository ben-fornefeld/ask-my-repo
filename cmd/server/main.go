package main

import (
	"log"
	"net/http"
	"rankmyrepo/internal/api"
	"rankmyrepo/internal/parser"
)

func main() {
	parserIgnorePatterns := []string{}
	parser, err := parser.NewParser(parserIgnorePatterns)
	if err != nil {
		log.Fatal(err)
	}
	defer parser.Cleanup()

	handler := api.NewHandler()

	http.HandleFunc("/query", handler.RankCode)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
