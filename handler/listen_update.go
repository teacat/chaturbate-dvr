package handler

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/teacat/chaturbate-dvr/chaturbate"
	"github.com/urfave/cli/v2"
)

//=======================================================
// Request & Response
//=======================================================

type ListenUpdateRequest struct {
}

//=======================================================
// Factory
//=======================================================

type ListenUpdateHandler struct {
	chaturbate *chaturbate.Manager
	cli        *cli.Context
}

func NewListenUpdateHandler(c *chaturbate.Manager, cli *cli.Context) *ListenUpdateHandler {
	return &ListenUpdateHandler{c, cli}
}

//=======================================================
// Handle
//=======================================================

func (h *ListenUpdateHandler) Handle(c *gin.Context) {
	ch, id := h.chaturbate.ListenUpdate()

	c.Stream(func(w io.Writer) bool {
		update := <-ch
		if update == nil {
			return false
		}
		c.SSEvent("message", map[string]any{
			"username":         update.Username,
			"log":              update.Log,
			"is_paused":        update.IsPaused,
			"is_online":        update.IsOnline,
			"last_streamed_at": update.LastStreamedAt,
			"segment_duration": update.SegmentDurationStr(),
			"segment_filesize": update.SegmentFilesizeStr(),
			"filename":         update.Filename,
		})
		return true
	})
	h.chaturbate.StopListenUpdate(id)
}
