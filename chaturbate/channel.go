package chaturbate

import (
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	regexpRoomDossier = regexp.MustCompile(`window\.initialRoomDossier = "(.*?)"`)
)

type roomDossier struct {
	HLSSource string `json:"hls_source"`
}

type resolution struct {
	framerate map[int]string // key: framerate, value: url
	width     int
}

type Channel struct {
	Username           string
	ChannelURL         string
	filenamePattern    string
	LastStreamedAt     string
	Interval           int
	Framerate          int
	Resolution         int
	ResolutionFallback string
	SegmentDuration    int // Seconds
	SplitDuration      int // Minutes
	SegmentFilesize    int // Bytes
	SplitFilesize      int // MB
	IsOnline           bool
	IsPaused           bool
	isStopped          bool
	Logs               []string
	logType            logType

	bufferLock   sync.Mutex
	buffer       map[int][]byte
	bufferIndex  int
	segmentIndex int
	segmentUseds []string
	rootURL      string
	sourceURL    string
	retries      int
	file         *os.File

	sessionPattern map[string]any
	splitIndex     int

	PauseChannel  chan bool
	UpdateChannel chan *Update
	ResumeChannel chan bool
}

// Run
func (w *Channel) Run() {
	if w.Username == "" {
		w.log(logTypeError, "username is empty, use `-u USERNAME` to specify")
		return
	}

	for {
		if w.IsPaused {
			w.log(logTypeInfo, "channel is paused")
			<-w.ResumeChannel // blocking
			w.log(logTypeInfo, "channel is resumed")
		}
		if w.isStopped {
			w.log(logTypeInfo, "channel is stopped")
			break
		}

		body, err := w.requestChannelBody()
		if err != nil {
			w.log(logTypeError, "body request error: %w", err)
		}
		if strings.Contains(body, "playlist.m3u8") {
			w.IsOnline = true
			w.LastStreamedAt = time.Now().Format("2006-01-02 15:04:05")
			w.log(logTypeInfo, "channel is online, start fetching...")

			if err := w.record(body); err != nil { // blocking
				w.log(logTypeError, "record error: %w", err)
			}
			continue // this excutes when recording is over/interrupted
		}
		w.IsOnline = false

		// close file when offline so user can move/delete it
		if w.file != nil {
			if err := w.releaseFile(); err != nil {
				w.log(logTypeError, "release file: %w", err)
			}
		}

		w.log(logTypeInfo, "channel is offline, check again %d min(s) later", w.Interval)
		<-time.After(time.Duration(w.Interval) * time.Minute) // minutes cooldown to check online status
	}
}

func (w *Channel) Pause() {
	w.IsPaused = true
	w.resetSession()
}

func (w *Channel) Resume() {
	w.IsPaused = false
	select {
	case w.ResumeChannel <- true:
	default:
	}
}

func (w *Channel) Stop() {
	w.isStopped = true
}

func (w *Channel) SegmentDurationStr() string {
	return DurationStr(w.SegmentDuration)
}

func (w *Channel) SplitDurationStr() string {
	return DurationStr(w.SplitDuration * 60)
}

func (w *Channel) SegmentFilesizeStr() string {
	return ByteStr(w.SegmentFilesize)
}

func (w *Channel) SplitFilesizeStr() string {
	return MBStr(w.SplitFilesize)
}

func (w *Channel) Filename() string {
	if w.file == nil {
		return ""
	}
	return w.file.Name()
}
