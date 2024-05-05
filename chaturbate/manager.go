package chaturbate

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
)

const (
	ResolutionFallbackUpscale   = "up"
	ResolutionFallbackDownscale = "down"
)

var (
	ErrChannelNotFound  = errors.New("channel not found")
	ErrChannelExists    = errors.New("channel already exists")
	ErrChannelNotPaused = errors.New("channel not paused")
	ErrChannelIsPaused  = errors.New("channel is paused")
	ErrListenNotFound   = errors.New("listen not found")
)

// Config
type Config struct {
	Username           string
	FilenamePattern    string
	Framerate          int
	Resolution         int
	ResolutionFallback string
	SplitDuration      int
	SplitFilesize      int
	Interval           int
}

// Manager
type Manager struct {
	cli      *cli.Context
	Channels map[string]*Channel
	Updates  map[string]chan *Update
}

// NewManager
func NewManager(c *cli.Context) *Manager {
	return &Manager{
		cli:      c,
		Channels: map[string]*Channel{},
		Updates:  map[string]chan *Update{},
	}
}

// PauseChannel
func (m *Manager) PauseChannel(username string) error {
	v, ok := m.Channels[username]
	if !ok {
		return ErrChannelNotFound
	}
	if v.IsPaused { // no-op
		return nil
	}
	v.Pause()
	return nil
}

// ResumeChannel
func (m *Manager) ResumeChannel(username string) error {
	v, ok := m.Channels[username]
	if !ok {
		return ErrChannelNotFound
	}
	if !v.IsPaused { // no-op
		return nil
	}
	v.Resume()
	return nil
}

// DeleteChannel
func (m *Manager) DeleteChannel(username string) error {
	v, ok := m.Channels[username]
	if !ok {
		return ErrChannelNotFound
	}
	v.Stop()
	delete(m.Channels, username)
	return nil
}

// CreateChannel
func (m *Manager) CreateChannel(conf *Config) error {
	_, ok := m.Channels[conf.Username]
	if ok {
		return ErrChannelExists
	}
	c := &Channel{
		Username:           conf.Username,
		ChannelURL:         "https://chaturbate.com/" + conf.Username,
		filenamePattern:    conf.FilenamePattern,
		Framerate:          conf.Framerate,
		Resolution:         conf.Resolution,
		ResolutionFallback: conf.ResolutionFallback,
		Interval:           conf.Interval,
		LastStreamedAt:     "-",
		SegmentDuration:    0,
		SplitDuration:      conf.SplitDuration,
		SegmentFilesize:    0,
		SplitFilesize:      conf.SplitFilesize,
		IsOnline:           false,
		IsPaused:           false,
		isStopped:          false,
		Logs:               []string{},
		UpdateChannel:      make(chan *Update),
		ResumeChannel:      make(chan bool),
		logType:            logType(m.cli.String("log-level")),
	}
	go func() {
		for update := range c.UpdateChannel {
			for _, v := range m.Updates {
				if v != nil {
					v <- update
				}
			}
		}
	}()
	m.Channels[conf.Username] = c
	c.log(logTypeInfo, "channel created")
	go c.Run()
	return nil
}

// ListChannels
func (m *Manager) ListChannels() ([]*Channel, error) {
	var channels []*Channel
	for _, v := range m.Channels {
		channels = append(channels, v)
	}
	return channels, nil
}

// GetChannel
func (m *Manager) GetChannel(username string) (*Channel, error) {
	v, ok := m.Channels[username]
	if !ok {
		return nil, ErrChannelNotFound
	}
	return v, nil
}

// ListenUpdate
func (m *Manager) ListenUpdate() (<-chan *Update, string) {
	c := make(chan *Update)
	id := uuid.New().String()
	m.Updates[id] = c
	return c, id
}

// StopListenUpdate
func (m *Manager) StopListenUpdate(id string) error {
	v, ok := m.Updates[id]
	if !ok {
		return ErrListenNotFound
	}
	delete(m.Updates, id)
	close(v)
	return nil
}

// SaveChannels
func (m *Manager) SaveChannels() error {
	configs := make([]*Config, 0)
	for _, v := range m.Channels {
		configs = append(configs, &Config{
			Username:           v.Username,
			Framerate:          v.Framerate,
			Resolution:         v.Resolution,
			ResolutionFallback: v.ResolutionFallback,
			FilenamePattern:    v.filenamePattern,
			SplitDuration:      v.SplitDuration,
			SplitFilesize:      v.SplitFilesize,
			Interval:           v.Interval,
		})
	}
	b, err := json.MarshalIndent(configs, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile("chaturbate_channels.json", b, 0777)
}

// LoadChannels
func (m *Manager) LoadChannels() error {
	b, err := os.ReadFile("chaturbate_channels.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var configs []*Config
	if err := json.Unmarshal(b, &configs); err != nil {
		return err
	}
	for _, v := range configs {
		if err := m.CreateChannel(v); err != nil {
			return err
		}
	}
	return nil
}
