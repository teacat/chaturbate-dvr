package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/teacat/chaturbate-dvr/domain/stream"
)

//=======================================================
// Request & Response
//=======================================================

type ResumeStreamRequest struct {
	Username string
}

type ResumeStreamResponse struct {
}

//=======================================================
// Factory
//=======================================================

type ResumeStreamHandler struct {
	manager stream.Manager
}

func NewResumeStreamHandler(manager stream.Manager) *ResumeStreamHandler {
	return &ResumeStreamHandler{
		manager: manager,
	}
}

//=======================================================
// Handle
//=======================================================

func (h *ResumeStreamHandler) Handle(ctx *gin.Context, req *ResumeStreamRequest) (*ResumeStreamResponse, error) {
	if err := h.manager.ResumeStream(req.Username); err != nil {
		return nil, fmt.Errorf("resume stream: %w", err)
	}
	return &ResumeStreamResponse{}, nil
}
