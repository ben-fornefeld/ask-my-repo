package api

type APIError struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"details,omitempty"`
}

