package api

import (
	"net/http"
	"time"
)

type RankingResponse struct {
	// TODO: add results
	ProcessTime time.Duration `json:"process_time"`
	RequestID   string        `json:"request_id"`
}

type Handler struct {
	// TODO: add processor
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RankCode(w http.ResponseWriter, r *http.Request) {
	//	TODO: add api route handler
}
