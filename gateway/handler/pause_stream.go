package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/teacat/chaturbate-dvr/domain/stream"
)

//=======================================================
// Request & Response
//=======================================================

type PauseStreamRequest struct {
	Username string `json:"username"`
}

type PauseStreamResponse struct {
}

//=======================================================
// Factory
//=======================================================

type PauseStreamHandler struct {
	manager stream.Manager
}

func NewPauseStreamHandler(manager stream.Manager) *PauseStreamHandler {
	return &PauseStreamHandler{
		manager: manager,
	}
}

//=======================================================
// Handle
//=======================================================

func (h *PauseStreamHandler) Handle(ctx *gin.Context, req *PauseStreamRequest) (*PauseStreamResponse, error) {
	if err := h.manager.PauseStream(req.Username); err != nil {
		return nil, fmt.Errorf("pause stream: %w", err)
	}
	return &PauseStreamResponse{}, nil
}
