package handler

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/teacat/chaturbate-dvr/chaturbate"
	"github.com/urfave/cli/v2"
)

//=======================================================
// Request & Response
//=======================================================

type ListChannelsRequest struct {
}

type ListChannelsResponse struct {
	Channels []*ListChannelsResponseChannel `json:"channels"`
}

type ListChannelsResponseChannel struct {
	Username        string   `json:"username"`
	ChannelURL      string   `json:"channel_url"`
	Filename        string   `json:"filename"`
	LastStreamedAt  string   `json:"last_streamed_at"`
	SegmentDuration string   `json:"segment_duration"`
	SplitDuration   string   `json:"split_duration"`
	SegmentFilesize string   `json:"segment_filesize"`
	SplitFilesize   string   `json:"split_filesize"`
	Interval        int      `json:"interval"`
	IsOnline        bool     `json:"is_online"`
	IsPaused        bool     `json:"is_paused"`
	Logs            []string `json:"logs"`
}

//=======================================================
// Factory
//=======================================================

type ListChannelsHandler struct {
	chaturbate *chaturbate.Manager
	cli        *cli.Context
}

func NewListChannelsHandler(c *chaturbate.Manager, cli *cli.Context) *ListChannelsHandler {
	return &ListChannelsHandler{c, cli}
}

//=======================================================
// Handle
//=======================================================

// Handle processes the request to list channels, sorting by IsOnline.
func (h *ListChannelsHandler) Handle(c *gin.Context) {
	var req *ListChannelsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Fetch channels
	channels, err := h.chaturbate.ListChannels()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Sort by IsOnline: online channels first, then offline
	sort.SliceStable(channels, func(i, j int) bool {
		return channels[i].IsOnline && !channels[j].IsOnline
	})

	// Populate response
	resp := &ListChannelsResponse{
		Channels: make([]*ListChannelsResponseChannel, len(channels)),
	}
	for i, channel := range channels {
		resp.Channels[i] = &ListChannelsResponseChannel{
			Username:        channel.Username,
			ChannelURL:      channel.ChannelURL,
			Filename:        channel.Filename(),
			LastStreamedAt:  channel.LastStreamedAt,
			SegmentDuration: channel.SegmentDurationStr(),
			SplitDuration:   channel.SplitDurationStr(),
			SegmentFilesize: channel.SegmentFilesizeStr(),
			SplitFilesize:   channel.SplitFilesizeStr(),
			Interval:        channel.Interval,
			IsOnline:        channel.IsOnline,
			IsPaused:        channel.IsPaused,
			Logs:            channel.Logs,
		}
	}

	// Send the response
	c.JSON(http.StatusOK, resp)
}
