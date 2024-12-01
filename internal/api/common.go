package api

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"details,omitempty"`
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
