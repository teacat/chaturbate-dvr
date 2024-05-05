package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
)

//=======================================================
// Request & Response
//=======================================================

type GetSettingsHandlerRequest struct {
}

type GetSettingsHandlerResponse struct {
	Version            string `json:"version"`
	Framerate          int    `json:"framerate"`
	Resolution         int    `json:"resolution"`
	ResolutionFallback string `json:"resolution_fallback"`
	FilenamePattern    string `json:"filename_pattern"`
	SplitDuration      int    `json:"split_duration"`
	SplitFilesize      int    `json:"split_filesize"`
	Interval           int    `json:"interval"`
	LogLevel           string `json:"log_level"`
	Port               string `json:"port"`
	GUI                string `json:"gui"`
}

//=======================================================
// Factory
//=======================================================

type GetSettingsHandler struct {
	cli *cli.Context
}

func NewGetSettingsHandler(cli *cli.Context) *GetSettingsHandler {
	return &GetSettingsHandler{cli}
}

//=======================================================
// Handle
//=======================================================

func (h *GetSettingsHandler) Handle(c *gin.Context) {
	var req *GetSettingsHandlerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, &GetSettingsHandlerResponse{
		Version:            h.cli.App.Version,
		Framerate:          h.cli.Int("framerate"),
		Resolution:         h.cli.Int("resolution"),
		ResolutionFallback: h.cli.String("resolution-fallback"),
		FilenamePattern:    h.cli.String("filename-pattern"),
		SplitDuration:      h.cli.Int("split-duration"),
		SplitFilesize:      h.cli.Int("split-filesize"),
		Interval:           h.cli.Int("interval"),
		LogLevel:           h.cli.String("log-level"),
		Port:               h.cli.String("port"),
		GUI:                h.cli.String("gui"),
	})
}
