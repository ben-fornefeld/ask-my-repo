package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"rankmyrepo/internal/processor"
	"rankmyrepo/internal/ranking"
	"time"

	"github.com/google/uuid"
)

type RankingResponse struct {
	Results     *ranking.RankingResponse `json:"results"`
	ProcessTime time.Duration            `json:"process_time"`
	RequestID   string                   `json:"request_id"`
}

type Handler struct {
	processor *processor.Processor
}

func NewHandler(processor *processor.Processor) (*Handler, error) {
	if processor == nil {
		return nil, errors.New("processor cannot be nil")
	}

	return &Handler{
		processor: processor,
	}, nil
}

func (h *Handler) RankCode(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := uuid.New().String()

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-ID", requestID)

	// 1. validation
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed,
			errors.New("only POST method is allowed"))
		return
	}

	// 2. request parsing
	var req ranking.RankingRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest,
			errors.New("invalid request body"))
		return
	}
	defer r.Body.Close()

	// 3. request validation
	if err := h.validateRequest(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, err)
		return
	}

	// 4. process request
	result, err := h.processor.ProcessRankingRequest(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			h.writeError(w, http.StatusGatewayTimeout, err)
		default:
			h.writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	// 5. prepare and send response
	response := RankingResponse{
		Results:     result,
		ProcessTime: time.Since(start),
		RequestID:   requestID,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.writeError(w, http.StatusInternalServerError,
			errors.New("failed to encode response"))
		return
	}
}

func (h *Handler) validateRequest(req *ranking.RankingRequest) error {
	if req.Query == "" {
		return errors.New("query cannot be empty")
	}
	if req.RepoPath == "" {
		return errors.New("repository path cannot be empty")
	}
	if req.IgnorePatterns == nil {
		return errors.New("ignore patterns cannot be empty")
	}
	if req.ScoreThreshold <= 0 || req.ScoreThreshold > 1 {
		return errors.New("score threshold must be between 0 and 1 (inclusive)")
	}
	return nil
}

// TODO: add requestId and logging

func (h *Handler) writeError(w http.ResponseWriter, status int, err error) {
	apiError := APIError{
		Error:   http.StatusText(status),
		Code:    status,
		Message: err.Error(),
	}

	w.WriteHeader(status)

	json.NewEncoder(w).Encode(apiError)
}
