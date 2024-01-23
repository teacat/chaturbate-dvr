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

type ResumeChannelRequest struct {
	Username string `json:"username"`
}

type ResumeChannelResponse struct {
}

//=======================================================
// Factory
//=======================================================

type ResumeChannelHandler struct {
	chaturbate *chaturbate.Manager
	cli        *cli.Context
}

func NewResumeChannelHandler(c *chaturbate.Manager, cli *cli.Context) *ResumeChannelHandler {
	return &ResumeChannelHandler{c, cli}
}

//=======================================================
// Handle
//=======================================================

func (h *ResumeChannelHandler) Handle(c *gin.Context) {
	var req *ResumeChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := h.chaturbate.ResumeChannel(req.Username); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, &ResumeChannelResponse{})
}
