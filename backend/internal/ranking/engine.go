package ranking

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"rankmyrepo/internal/parser"
	"sync"

	"github.com/replicate/replicate-go"
)

type Engine struct {
	r8         *replicate.Client
	maxWorkers int
}

func NewEngine(r8 *replicate.Client, maxWorkers int) *Engine {
	return &Engine{
		r8:         r8,
		maxWorkers: maxWorkers,
	}
}

func (e *Engine) RankSingleChunkReplicate(ctx context.Context, query string, chunk parser.ParsedChunk) (float64, error) {
	prompt := buildRankingPrompt(query, chunk)

	model := "meta/meta-llama-3-8b-instruct"

	input := replicate.PredictionInput{
		"prompt":       prompt,
		"system_prompt": systemPrompt,
		"temperature": 0.1,
	}

	output, err := e.r8.Run(ctx, model, input, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to run model: %w", err)
	}

	tokens, ok := output.([]interface{})
	if !ok {
		return 0, fmt.Errorf("unexpected output type from model: got %T, want []interface{}", output)
	}

	var result string
	for _, token := range tokens {
		if str, ok := token.(string); ok {
			result += str
		}
	}

	return parseScore(result)
}

func (e *Engine) RankSingleChunkFireworks(ctx context.Context, query string, chunk parser.ParsedChunk) (float64, error) {
	prompt := buildRankingPrompt(query, chunk)

	requestBody := struct {
		Model            string `json:"model"`
		TopP            float64 `json:"top_p"`
		TopK            int     `json:"top_k"`
		PresencePenalty float64 `json:"presence_penalty"`
		FrequencyPenalty float64 `json:"frequency_penalty"`
		Temperature     float64 `json:"temperature"`
		Messages        []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
	}{
		Model:            "accounts/fireworks/models/llama-v3p2-3b-instruct",
		PresencePenalty:  0,
		FrequencyPenalty: 0,
		Temperature:     0.1,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return 0, fmt.Errorf("marshaling request body: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.fireworks.ai/inference/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+os.Getenv("FIREWORKS_API_KEY"))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, _ := io.ReadAll(resp.Body)

	return parseScore(string(body))
}


func (e *Engine) RankChunksStream(ctx context.Context, query string, chunks map[string]parser.ParsedChunk, scoreThreshold float64, parsedChan chan<- parser.ParsedChunk, rankedChan chan<- RankedChunk) error {
	log.Printf("Starting to rank %d chunks for query: %s", len(chunks), query)

	errors := make(chan error, len(chunks))
	sem := make(chan struct{}, e.maxWorkers)

	var wg sync.WaitGroup
	for _, chunk := range chunks {
		wg.Add(1)
		go func(c parser.ParsedChunk) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			select {
			case parsedChan <- c:
			case <-ctx.Done():
				errors <- ctx.Err()
				return
			}

			score, err := e.RankSingleChunkFireworks(ctx, query, c)

			log.Printf("Score: %f", score)

			if err != nil {
				errors <- err
				return
			}

			if score >= scoreThreshold {
				ranked := RankedChunk{
					ParsedChunk: c,
					Score:       score,
				}
				select {
				case rankedChan <- ranked:
				case <-ctx.Done():
					errors <- ctx.Err()
					return
				}
			}
		}(chunk)
	}

	go func() {
		wg.Wait()
		close(errors)
	}()

	for err := range errors {
		if err != nil {
			return err
		}
	}

	return nil
}
