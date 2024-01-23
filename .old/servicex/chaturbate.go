package service

import "errors"

var (
	ErrChannelNotFound  = errors.New("channel not found")
	ErrChannelExists    = errors.New("channel already exists")
	ErrChannelNotPaused = errors.New("channel not paused")
	ErrChannelIsPaused  = errors.New("channel is paused")
	ErrListenNotFound   = errors.New("listen not found")
)

const (
	ResolutionFallbackUpscale   = "up"
	ResolutionFallbackDownscale = "down"
)

type Chaturbate interface {
	GetChannel(username string) (ChaturbateChannel, error)
	CreateChannel(config *ChaturbateConfig) error
	DeleteChannel(username string) error
	PauseChannel(username string) error
	ResumeChannel(username string) error
	ListChannels() ([]ChaturbateChannel, error)
	ListenUpdate() (<-chan *Update, string)
	StopListenUpdate(id string) error
}

type ChaturbateChannel interface {
	Username() string
	ChannelURL() string
	FilenamePattern() string
	Framerate() int
	Resolution() int
	ResolutionFallback() string
	LastStreamedAt() string
	Filename() string
	SegmentDuration() int
	SplitDuration() int
	SegmentFilesize() int
	SplitFilesize() int
	IsOnline() bool
	IsPaused() bool
	Logs() []string
}

type ChaturbateConfig struct {
	Username           string
	FilenamePattern    string
	Framerate          int
	Resolution         int
	ResolutionFallback string
	SplitDuration      int
	SplitFilesize      int
}

type Update struct {
	Username        string `json:"username"`
	Log             string `json:"log"`
	IsPaused        bool   `json:"is_paused"`
	IsOnline        bool   `json:"is_online"`
	IsStopped       bool   `json:"is_stopped"`
	Filename        string `json:"filename"`
	LastStreamedAt  string `json:"last_streamed_at"`
	SegmentDuration string `json:"segment_duration"`
	SegmentFilesize string `json:"segment_filesize"`
}
