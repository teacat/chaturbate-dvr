package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/teacat/chaturbate-dvr/chaturbate"
	"github.com/urfave/cli/v2"
)

//=======================================================
// Request & Response
//=======================================================

type CreateChannelRequest struct {
	Username           string `json:"username"`
	Framerate          int    `json:"framerate"`
	FilenamePattern    string `json:"filename_pattern"`
	Resolution         int    `json:"resolution"`
	ResolutionFallback string `json:"resolution_fallback"`
	SplitDuration      int    `json:"split_duration"`
	SplitFilesize      int    `json:"split_filesize"`
	Interval           int    `json:"interval"`
}

type CreateChannelResponse struct {
}

//=======================================================
// Factory
//=======================================================

type CreateChannelHandler struct {
	chaturbate *chaturbate.Manager
	cli        *cli.Context
}

func NewCreateChannelHandler(c *chaturbate.Manager, cli *cli.Context) *CreateChannelHandler {
	return &CreateChannelHandler{c, cli}
}

//=======================================================
// Handle
//=======================================================

func (h *CreateChannelHandler) Handle(c *gin.Context) {
	var req *CreateChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	usernames := strings.Split(req.Username, ",")
	for _, username := range usernames {
		if err := h.chaturbate.CreateChannel(&chaturbate.Config{
			Username:           strings.TrimSpace(username),
			Framerate:          req.Framerate,
			Resolution:         req.Resolution,
			ResolutionFallback: req.ResolutionFallback,
			FilenamePattern:    req.FilenamePattern,
			SplitDuration:      req.SplitDuration,
			SplitFilesize:      req.SplitFilesize,
			Interval:           req.Interval,
		}); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	if err := h.chaturbate.SaveChannels(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, &CreateChannelResponse{})
}
