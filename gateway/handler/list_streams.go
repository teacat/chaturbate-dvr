package handler

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/teacat/chaturbate-dvr/domain/stream"
)

//=======================================================
// Request & Response
//=======================================================

type ListStreamsRequest struct {
}

type ListStreamsResponse struct {
	Streams []*ListStreamsResponseStream `json:"streams"`
}

type ListStreamsResponseStream struct {
	Username             string `json:"username"`
	ChannelURL           string `json:"channel_url"`
	SavedTo              string `json:"saved_to"`
	LastStreamedAt       string `json:"last_streamed_at"`
	SegmentDuration      string `json:"segment_duration"`
	SegmentDurationSplit string `json:"segment_duration_split"`
	SegmentFilesize      string `json:"segment_filesize"`
	SegmentFilesizeSplit string `json:"segment_filesize_split"`
	IsOnline             bool   `json:"is_online"`
	IsPaused             bool   `json:"is_paused"`
}

//=======================================================
// Factory
//=======================================================

type ListStreamsHandler struct {
	manager stream.Manager
}

func NewListStreamsHandler(manager stream.Manager) *ListStreamsHandler {
	return &ListStreamsHandler{
		manager: manager,
	}
}

//=======================================================
// Handle
//=======================================================

func (h *ListStreamsHandler) Handle(ctx *gin.Context, req *ListStreamsRequest) (*ListStreamsResponse, error) {
	streams, err := h.manager.ListStreams()
	if err != nil {
		return nil, fmt.Errorf("list streams: %w", err)
	}
	resp := &ListStreamsResponse{
		Streams: make([]*ListStreamsResponseStream, len(streams)),
	}
	for i, s := range streams {
		resp.Streams[i] = &ListStreamsResponseStream{
			Username:             s.Username,
			ChannelURL:           fmt.Sprintf("https://chaturbate.com/%s/", s.Username),
			SavedTo:              fmt.Sprintf("./videos/%s/", s.Username),
			LastStreamedAt:       time.Unix(s.LastStreamedAt, 0).Format("2006-01-02 15:04:05"),
			SegmentDuration:      formatPlaytime(s.SegmentDuration),
			SegmentDurationSplit: formatPlaytime(s.SegmentDurationSplit),
			SegmentFilesize:      fmt.Sprintf("%d MB", s.SegmentFilesize),
			SegmentFilesizeSplit: fmt.Sprintf("%d MB", s.SegmentFilesizeSplit),
			IsOnline:             s.IsOnline,
			IsPaused:             s.IsPaused,
		}
	}
	return resp, nil
}
