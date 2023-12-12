package stream

import "fmt"

//=======================================================
// Enum
//=======================================================

type ResolutionFallback string

const (
	ResolutionFallbackUnknown   ResolutionFallback = ""
	ResolutionFallbackUpscale   ResolutionFallback = "upscale"
	ResolutionFallbackDownscale ResolutionFallback = "downscale"
)

//=======================================================
// Entity
//=======================================================

type Stream struct {
	channelURL         string
	channelUsername    string
	splitFilesize      int
	splitDuration      int
	resolution         int
	resolutionFallback ResolutionFallback
	framerate          int
}

type StreamDTO struct {
	Username             string
	LastStreamedAt       int64
	SegmentDuration      int
	SegmentDurationSplit int
	SegmentFilesize      int
	SegmentFilesizeSplit int
	IsOnline             bool
	IsPaused             bool
}

type StreamUpdateDTO struct {
	Username        string
	IsOnline        bool
	IsPaused        bool
	LastStreamedAt  string
	SegmentDuration string
	SegmentFilesize string
}

type StreamOutputDTO struct {
	Username string
	Output   string
}

//=======================================================
// Domain
//=======================================================

// Start starts the stream and recording.
func (s *Stream) Start() error {
	return nil
}

// Pause pauses the stream and keep the stream in the list.
func (s *Stream) Pause() error {
	return nil
}

// Stop stops the stream and removes the stream from the list.
func (s *Stream) Stop() error {
	return nil
}

// Resume resumes the paused stream.
func (s *Stream) Resume() error {
	return nil
}

//=======================================================
// Factory
//=======================================================

type StreamFactory struct {
	sanitizer *streamSanitizer
}

func NewStreamFactory() *StreamFactory {
	return &StreamFactory{
		sanitizer: newStreamSanitizer(),
	}
}

func (f *StreamFactory) New(username string, resFallback ResolutionFallback, resolution, framerate, splitFilesize, splitDuration int) (*Stream, error) {
	username, err := f.sanitizer.sanitizeUsername(username)
	if err != nil {
		return nil, fmt.Errorf("sanitize username: %w", err)
	}
	resolution, err = f.sanitizer.sanitizeResolution(resolution)
	if err != nil {
		return nil, fmt.Errorf("sanitize resolution: %w", err)
	}
	resFallback, err = f.sanitizer.sanitizeResolutionFallback(resFallback)
	if err != nil {
		return nil, fmt.Errorf("sanitize resolution fallback: %w", err)
	}
	framerate, err = f.sanitizer.sanitizeFramerate(framerate)
	if err != nil {
		return nil, fmt.Errorf("sanitize framerate: %w", err)
	}
	splitFilesize, err = f.sanitizer.sanitizeSplitByFilesize(splitFilesize)
	if err != nil {
		return nil, fmt.Errorf("sanitize split by filesize: %w", err)
	}
	splitDuration, err = f.sanitizer.sanitizeSplitByDuration(splitDuration)
	if err != nil {
		return nil, fmt.Errorf("sanitize split by duration: %w", err)
	}
	return &Stream{
		channelUsername:    username,
		resolution:         resolution,
		resolutionFallback: resFallback,
		framerate:          framerate,
		splitFilesize:      splitFilesize,
		splitDuration:      splitDuration,
	}, nil
}

//=======================================================
// Sanitizer
//=======================================================

type streamSanitizer struct {
}

func newStreamSanitizer() *streamSanitizer {
	return &streamSanitizer{}
}

func (s *streamSanitizer) sanitizeUsername(v string) (string, error) {
	return v, nil
}

func (s *streamSanitizer) sanitizeResolution(v int) (int, error) {
	return v, nil
}

func (s *streamSanitizer) sanitizeResolutionFallback(v ResolutionFallback) (ResolutionFallback, error) {
	return v, nil
}

func (s *streamSanitizer) sanitizeFramerate(v int) (int, error) {
	return v, nil
}

func (s *streamSanitizer) sanitizeSplitByFilesize(v int) (int, error) {
	return v, nil
}

func (s *streamSanitizer) sanitizeSplitByDuration(v int) (int, error) {
	return v, nil
}
