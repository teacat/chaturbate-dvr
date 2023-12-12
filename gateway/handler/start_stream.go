package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/teacat/chaturbate-dvr/domain/stream"
)

//=======================================================
// Request & Response
//=======================================================

type StartStreamRequest struct {
	Username           string
	Resolution         int
	ResolutionFallback stream.ResolutionFallback
	Framerate          int
	SplitByFilesize    int
	SplitByDuration    int
}

type StartStreamResponse struct {
}

//=======================================================
// Factory
//=======================================================

type StartStreamHandler struct {
	manager stream.Manager
}

func NewStartStreamHandler(manager stream.Manager) *StartStreamHandler {
	return &StartStreamHandler{
		manager: manager,
	}
}

//=======================================================
// Handle
//=======================================================

func (h *StartStreamHandler) Handle(ctx *gin.Context, req *StartStreamRequest) (*StartStreamResponse, error) {
	if err := h.manager.AddStream(
		req.Username,
		req.ResolutionFallback,
		req.Resolution,
		req.Framerate,
		req.SplitByFilesize,
		req.SplitByDuration,
		false,
	); err != nil {
		return nil, fmt.Errorf("add stream: %w", err)
	}
	return &StartStreamResponse{}, nil
}
