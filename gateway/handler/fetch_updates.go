package handler

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/teacat/chaturbate-dvr/domain/stream"
)

//=======================================================
// Request & Response
//=======================================================

type FetchUpdatesRequest struct {
}

type FetchUpdatesResponse struct {
	Updates []*FetchUpdatesResponseUpdate `json:"updates"`
	Outputs []*FetchUpdatesResponseOutput `json:"outputs"`
}

type FetchUpdatesResponseUpdate struct {
	Username        string `json:"username"`
	IsOnline        bool   `json:"is_online"`
	IsPaused        bool   `json:"is_paused"`
	LastStreamedAt  string `json:"last_streamed_at"`
	SegmentDuration string `json:"segment_duration"`
	SegmentFilesize string `json:"segment_filesize"`
}

type FetchUpdatesResponseOutput struct {
	Username string `json:"username"`
	Output   string `json:"output"`
}

//=======================================================
// Factory
//=======================================================

type FetchUpdatesHandler struct {
	manager stream.Manager
}

func NewFetchUpdatesHandler() *FetchUpdatesHandler {
	return &FetchUpdatesHandler{}
}

//=======================================================
// Handle
//=======================================================

func (h *FetchUpdatesHandler) Handle(ctx *gin.Context, req *FetchUpdatesRequest) (*FetchUpdatesResponse, error) {
	chUpd := make(chan *stream.StreamUpdateDTO)
	chOut := make(chan *stream.StreamOutputDTO)

	if err := h.manager.SubscribeStreams(chUpd, chOut); err != nil {
		return nil, fmt.Errorf("subscribe pool: %w", err)
	}

	for {
		select {
		case upd := <-chUpd:
			b, err := json.Marshal(&FetchUpdatesResponseUpdate{
				Username:        upd.Username,
				IsOnline:        upd.IsOnline,
				IsPaused:        upd.IsPaused,
				LastStreamedAt:  upd.LastStreamedAt,
				SegmentDuration: upd.SegmentDuration,
				SegmentFilesize: upd.SegmentFilesize,
			})
			if err != nil {
				return nil, fmt.Errorf("marshal update: %w", err)
			}
			ctx.SSEvent("update", b)

		case out := <-chOut:
			b, err := json.Marshal(&FetchUpdatesResponseOutput{
				Username: out.Username,
				Output:   out.Output,
			})
			if err != nil {
				return nil, fmt.Errorf("marshal output: %w", err)
			}
			ctx.SSEvent("output", b)
		}
	}

	return nil, nil
}
