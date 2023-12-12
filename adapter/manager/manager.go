package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	streamDomain "github.com/teacat/chaturbate-dvr/domain/stream"
)

var (
	ErrStreamExists    = errors.New("stream exists")
	ErrStreamNotExists = errors.New("stream not exists")
)

// Manager manages the streams.
type Manager struct {
	config  *Config
	streams map[string]*stream
}

// New creates a new Manager.
func New() (*Manager, error) {
	if err := os.MkdirAll("./videos", 0777); err != nil {
		return nil, fmt.Errorf("create videos directory: %w", err)
	}
	config := &Config{}

	b, err := os.ReadFile("./config.json")
	if os.IsNotExist(err) {
		b, err := json.Marshal(config)
		if err != nil {
			return nil, fmt.Errorf("marshal config: %w", err)
		}
		if err := os.WriteFile("./config.json", b, 0777); err != nil {
			return nil, fmt.Errorf("write config: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	} else {
		if err := json.Unmarshal(b, config); err != nil {
			return nil, fmt.Errorf("unmarshal config: %w", err)
		}
	}
	manager := &Manager{
		config:  config,
		streams: make(map[string]*stream),
	}
	manager.Ready()
	return manager, nil
}

// Ready prepares the manager.
func (m *Manager) Ready() error {
	for _, s := range m.config.Streams {
		if err := m.AddStream(s.Username, s.ResolutionFallback, s.Resolution, s.Framerate, s.SplitByFilesize, s.SplitByDuration, s.IsPaused); err != nil {
			return fmt.Errorf("add stream: %w", err)
		}
	}
	return nil
}

func (m *Manager) ListStreams() ([]*streamDomain.StreamDTO, error) {
	return nil, nil
}

// PauseStream pauses the stream.
func (m *Manager) PauseStream(username string) error {
	if _, ok := m.streams[username]; !ok {
		return ErrStreamNotExists
	}
	m.streams[username].pause()
	return nil
}

// AddStream adds a stream to watching list and starts watching.
func (m *Manager) AddStream(username string, resFallback streamDomain.ResolutionFallback, resolution, framerate, splitByFilesize, splitByDuration int, isPaused bool) error {
	if _, ok := m.streams[username]; ok {
		return ErrStreamExists
	}
	// TODO: Sanitize, Trim username.
	s, _, _ := newStream(username)
	if !isPaused {
		s.start()
	}
	m.streams[username] = s
	return nil
}

func (m *Manager) StopStream(username string) error {
	if _, ok := m.streams[username]; !ok {
		return ErrStreamNotExists
	}
	m.streams[username].stop()
	return nil
}

func (m *Manager) ResumeStream(username string) error {
	if _, ok := m.streams[username]; !ok {
		return ErrStreamNotExists
	}
	m.streams[username].resume()
	return nil
}

func (m *Manager) SubscribeStreams(chUpd chan<- *streamDomain.StreamUpdateDTO, chOut chan<- *streamDomain.StreamOutputDTO) error {
	return nil
}
