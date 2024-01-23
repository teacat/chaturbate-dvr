package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/teacat/chaturbate-dvr/chaturbate"
	"github.com/urfave/cli/v2"
)

//=======================================================
// Request & Response
//=======================================================

type PauseChannelRequest struct {
	Username string `json:"username"`
}

type PauseChannelResponse struct {
}

//=======================================================
// Factory
//=======================================================

type PauseChannelHandler struct {
	chaturbate *chaturbate.Manager
	cli        *cli.Context
}

func NewPauseChannelHandler(c *chaturbate.Manager, cli *cli.Context) *PauseChannelHandler {
	return &PauseChannelHandler{c, cli}
}

//=======================================================
// Handle
//=======================================================

func (h *PauseChannelHandler) Handle(c *gin.Context) {
	var req *PauseChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := h.chaturbate.PauseChannel(req.Username); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, &PauseChannelResponse{})
}
