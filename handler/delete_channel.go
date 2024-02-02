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

type DeleteChannelRequest struct {
	Username string `json:"username"`
}

type DeleteChannelResponse struct {
}

//=======================================================
// Factory
//=======================================================

type DeleteChannelHandler struct {
	chaturbate *chaturbate.Manager
	cli        *cli.Context
}

func NewDeleteChannelHandler(c *chaturbate.Manager, cli *cli.Context) *DeleteChannelHandler {
	return &DeleteChannelHandler{c, cli}
}

//=======================================================
// Handle
//=======================================================

func (h *DeleteChannelHandler) Handle(c *gin.Context) {
	var req *DeleteChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := h.chaturbate.DeleteChannel(req.Username); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := h.chaturbate.SaveChannels(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, &DeleteChannelResponse{})
}
