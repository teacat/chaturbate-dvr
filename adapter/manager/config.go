package manager

import (
	streamDomain "github.com/teacat/chaturbate-dvr/domain/stream"
)

// Config is a configuration for the manager, loads from config.json.
type Config struct {
	Streams []*ConfigStream `json:"streams"`
}

// ConfigStream is a configuration for a stream.
type ConfigStream struct {
	Username           string                          `json:"username"`
	Resolution         int                             `json:"resolution"`
	ResolutionFallback streamDomain.ResolutionFallback `json:"resolution_fallback"`
	Framerate          int                             `json:"framerate"`
	SplitByDuration    int                             `json:"split_by_duration"`
	SplitByFilesize    int                             `json:"split_by_filesize"`
	IsPaused           bool                            `json:"is_paused"`
}
