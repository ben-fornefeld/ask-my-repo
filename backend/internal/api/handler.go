package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"rankmyrepo/internal/common"
	"rankmyrepo/internal/processor"
	"rankmyrepo/internal/ranking"
	"time"

	"github.com/gin-gonic/gin"
)

type RankingResponse struct {
	Results     *ranking.RankingResponse `json:"results"`
	ProcessTime time.Duration           `json:"process_time"`
	RequestID   string                  `json:"request_id"`
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

func (h *Handler) Query(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	var req ranking.RankingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeSSEEvent(c, common.QueryResponseChunk{
			Type:      common.EventTypeError,
			Error:     "Invalid request body",
		})
		return
	}

	resultChan := make(chan common.QueryResponseChunk)
	errChan := make(chan error, 1)

	go func() {
		err := h.processor.ProcessRankingRequestStream(c.Request.Context(), &req, resultChan)
		if err != nil {
			errChan <- err
		}
		close(resultChan)
	}()

	for {
		select {
		case chunk, ok := <-resultChan:
			if !ok {
				return
			}
			writeSSEEvent(c, chunk)
		case err := <-errChan:
			writeSSEEvent(c, common.QueryResponseChunk{
				Type: common.EventTypeError,
				Error: err.Error(),
			})
			return
		case <-c.Request.Context().Done():
			return
		}
	}
}

func writeSSEEvent(c *gin.Context, event common.QueryResponseChunk) {
	data, _ := json.Marshal(event)
	c.Writer.Write([]byte(fmt.Sprintf("data: %s\n\n", data)))
	c.Writer.Flush()
}
