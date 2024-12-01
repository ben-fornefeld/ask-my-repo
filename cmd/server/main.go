package main

import (
	"log"
	"net/http"
	"rankmyrepo/internal/api"
)

func main() {
	// TODO : init components

	handler := api.NewHandler()

	http.HandleFunc("/query", handler.RankCode)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
