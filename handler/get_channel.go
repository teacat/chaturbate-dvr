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

type GetChannelRequest struct {
	Username string `json:"username"`
}

type GetChannelResponse struct {
	Username        string   `json:"username"`
	ChannelURL      string   `json:"channel_url"`
	Filename        string   `json:"filename"`
	LastStreamedAt  string   `json:"last_streamed_at"`
	SegmentDuration string   `json:"segment_duration"`
	SplitDuration   string   `json:"split_duration"`
	SegmentFilesize string   `json:"segment_filesize"`
	SplitFilesize   string   `json:"split_filesize"`
	IsOnline        bool     `json:"is_online"`
	IsPaused        bool     `json:"is_paused"`
	Logs            []string `json:"logs"`
}

//=======================================================
// Factory
//=======================================================

type GetChannelHandler struct {
	chaturbate *chaturbate.Manager
	cli        *cli.Context
}

func NewGetChannelHandler(c *chaturbate.Manager, cli *cli.Context) *GetChannelHandler {
	return &GetChannelHandler{c, cli}
}

//=======================================================
// Handle
//=======================================================

func (h *GetChannelHandler) Handle(c *gin.Context) {
	var req *GetChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	channel, err := h.chaturbate.GetChannel(req.Username)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, &GetChannelResponse{
		Username:        channel.Username,
		ChannelURL:      channel.ChannelURL,
		Filename:        channel.Filename(),
		LastStreamedAt:  channel.LastStreamedAt,
		SegmentDuration: channel.SegmentDurationStr(),
		SplitDuration:   channel.SplitDurationStr(),
		SegmentFilesize: channel.SegmentFilesizeStr(),
		SplitFilesize:   channel.SplitFilesizeStr(),
		IsOnline:        channel.IsOnline,
		IsPaused:        channel.IsPaused,
		Logs:            channel.Logs,
	})
}
