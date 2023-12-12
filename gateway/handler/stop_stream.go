package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/teacat/chaturbate-dvr/domain/stream"
)

//=======================================================
// Request & Response
//=======================================================

type StopStreamRequest struct {
	Username string
}

type StopStreamResponse struct {
}

//=======================================================
// Factory
//=======================================================

type StopStreamHandler struct {
	manager stream.Manager
}

func NewStopStreamHandler(manager stream.Manager) *StopStreamHandler {
	return &StopStreamHandler{
		manager: manager,
	}
}

//=======================================================
// Handle
//=======================================================

func (h *StopStreamHandler) Handle(ctx *gin.Context, req *StopStreamRequest) (*StopStreamResponse, error) {
	if err := h.manager.StopStream(req.Username); err != nil {
		return nil, err
	}
	return &StopStreamResponse{}, nil
}
