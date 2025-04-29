package entity

type Event = string

const (
	EventUpdate Event = "update"
	EventLog    Event = "log"
)

type ChannelConfig struct {
	IsPaused    bool   `json:"is_paused"`
	Username    string `json:"username"`
	Framerate   int    `json:"framerate"`
	Resolution  int    `json:"resolution"`
	Pattern     string `json:"pattern"`
	MaxDuration int    `json:"max_duration"`
	MaxFilesize int    `json:"max_filesize"`
}

type ChannelInfo struct {
	IsOnline     bool
	IsPaused     bool
	Username     string
	Duration     string
	Filesize     string
	Filename     string
	StreamedAt   string
	MaxDuration  string
	MaxFilesize  string
	Logs         []string
	GlobalConfig *Config // for nested template to access $.Config
}

// Config holds the configuration for the application.
type Config struct {
	Version       string
	Username      string
	AdminUsername string
	AdminPassword string
	Framerate     int
	Resolution    int
	Pattern       string
	MaxDuration   int
	MaxFilesize   int
	Port          string
	Interval      int
	Cookies       string
	UserAgent     string
	Domain        string
}
