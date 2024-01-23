package chaturbate

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/grafov/m3u8"
	"github.com/samber/lo"
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
	Framerate          int
	Resolution         int
	ResolutionFallback string
	SegmentDuration    int
	SplitDuration      int
	SegmentFilesize    int
	SplitFilesize      int
	IsOnline           bool
	IsPaused           bool
	isStopped          bool
	Logs               []string

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

	UpdateChannel chan *Update
	ResumeChannel chan bool
}

// Run
func (w *Channel) Run() {
	for {
		if w.IsPaused {
			<-w.ResumeChannel // blocking
		}
		if w.isStopped {
			break
		}

		body, err := w.requestChannelBody()
		if err != nil {
			w.log("Error occurred while requesting channel body: %w", err)
		}
		if strings.Contains(body, "playlist.m3u8") {
			w.IsOnline = true
			w.LastStreamedAt = time.Now().Format("2006-01-02 15:04:05")
			w.log("Channel is online.")

			if err := w.record(body); err != nil { // blocking
				w.log("Error occurred when start recording: %w", err)
			}
			continue // this excutes when recording is over/interrupted
		}
		w.IsOnline = false
		w.log("Channel is offline.")
		<-time.After(1 * time.Minute) // 1 minute cooldown to check online status
	}
}

// requestChannelBody
func (w *Channel) requestChannelBody() (string, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}

	resp, err := client.Get(w.ChannelURL)
	if err != nil {
		return "", fmt.Errorf("client get: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}

	return string(body), nil
}

// record
func (w *Channel) record(body string) error {
	w.resetSession()

	if err := w.newFile(); err != nil {
		return fmt.Errorf("new file: %w", err)
	}

	rootURL, sourceURL, err := w.resolveSource(body)
	if err != nil {
		return fmt.Errorf("request hls: %w", err)
	}
	w.rootURL = rootURL
	w.sourceURL = sourceURL

	go w.mergeSegments()
	w.fetchSegments() // blocking

	return nil
}

func (w *Channel) resetSession() {
	w.buffer = make(map[int][]byte)
	w.bufferLock = sync.Mutex{}
	w.bufferIndex = 0
	w.segmentIndex = 0
	w.segmentUseds = []string{}
	w.rootURL = ""
	w.sourceURL = ""
	w.retries = 0
	w.SegmentFilesize = 0
	w.SegmentDuration = 0
	w.splitIndex = 0
	w.sessionPattern = nil
}

func (w *Channel) resolveSource(body string) (string, string, error) {
	// Find the room dossier.
	matches := regexpRoomDossier.FindAllStringSubmatch(body, -1)

	// Get the HLS source from the room dossier.
	var roomData roomDossier
	data, err := strconv.Unquote(strings.Replace(strconv.Quote(string(matches[0][1])), `\\u`, `\u`, -1))
	if err != nil {
		return "", "", fmt.Errorf("unquote unicode: %w", err)
	}
	if err := json.Unmarshal([]byte(data), &roomData); err != nil {
		return "", "", fmt.Errorf("unmarshal json: %w", err)
	}

	// Get the HLS source.
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}

	resp, err := client.Get(roomData.HLSSource)
	if err != nil {
		return "", "", fmt.Errorf("client get: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusForbidden:
			return "", "", fmt.Errorf("received status code %d, the stream is private?", resp.StatusCode)
		default:
			return "", "", fmt.Errorf("received status code %d", resp.StatusCode)
		}
	}
	defer resp.Body.Close()

	m3u8Body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("read body: %w", err)
	}

	// Decode the m3u8 file.
	p, _, err := m3u8.DecodeFrom(bytes.NewReader(m3u8Body), true)
	if err != nil {
		return "", "", fmt.Errorf("decode m3u8: %w", err)
	}
	playlist, ok := p.(*m3u8.MasterPlaylist)
	if !ok {
		return "", "", fmt.Errorf("cast to master playlist")
	}

	var resolutions []*resolution
	for _, v := range playlist.Variants {
		width := strings.Split(v.Resolution, "x")[1] // 1920x1080 -> 1080
		fps := 30
		if strings.Contains(v.Name, "FPS:60.0") {
			fps = 60
		}
		variant, ok := lo.Find(resolutions, func(v *resolution) bool {
			return strconv.Itoa(v.width) == width
		})
		if ok {
			variant.framerate[fps] = v.URI
			continue
		}
		widthInt, err := strconv.Atoi(width)
		if err != nil {
			return "", "", fmt.Errorf("convert width string to int: %w", err)
		}
		resolutions = append(resolutions, &resolution{
			framerate: map[int]string{fps: v.URI},
			width:     widthInt,
		})
	}

	variant, ok := lo.Find(resolutions, func(v *resolution) bool {
		return v.width == w.Resolution
	})
	// Fallback to the nearest resolution if the preferred resolution is not found.
	if !ok {
		switch w.ResolutionFallback {
		case ResolutionFallbackDownscale:
			variant = lo.MaxBy(lo.Filter(resolutions, func(v *resolution, _ int) bool {
				log.Println(v.width, w.Resolution)
				return v.width < w.Resolution
			}), func(v, max *resolution) bool {
				return v.width > max.width
			})
		case ResolutionFallbackUpscale:
			variant = lo.MinBy(lo.Filter(resolutions, func(v *resolution, _ int) bool {
				return v.width > w.Resolution
			}), func(v, min *resolution) bool {
				return v.width < min.width
			})
		}
	}
	if variant == nil {
		return "", "", fmt.Errorf("no available variant")
	}

	url, ok := variant.framerate[w.Framerate]
	// If the framerate is not found, fallback to the first found framerate, this block pretends there're only 30 and 60 fps.
	// no complex logic here, im lazy.
	if ok {
		w.log("Framerate %d is used.", w.Framerate)
	} else {
		for k, v := range variant.framerate {
			url = v
			w.log("Framerate %d is not found, fallback to %d.", w.Framerate, k)
			break
		}
	}

	rootURL := strings.TrimSuffix(roomData.HLSSource, "playlist.m3u8")
	sourceURL := rootURL + url
	return rootURL, sourceURL, nil
}

func (w *Channel) mergeSegments() {
	var segmentRetries int

	for {
		if w.IsPaused || w.isStopped {
			break
		}
		if segmentRetries > 5 {
			w.log("Segment #%d error, the segment has been skipped.", w.bufferIndex)
			w.bufferIndex++
			segmentRetries = 0
			continue
		}
		if len(w.buffer) == 0 {
			<-time.After(1 * time.Second)
			continue
		}
		buf, ok := w.buffer[w.bufferIndex]
		if !ok {
			segmentRetries++
			<-time.After(time.Duration(segmentRetries) * time.Second)
			continue
		}
		lens, err := w.file.Write(buf)
		if err != nil {
			w.log("Error occurred while writing segment #%d to file: %v", w.bufferIndex, err)
			w.retries++
			continue
		}
		w.log("Segment #%d written to file.", w.bufferIndex)

		w.SegmentFilesize += lens
		segmentRetries = 0

		if w.SplitFilesize > 0 && w.SegmentFilesize >= w.SplitFilesize*1024*1024 {
			w.log("File size has exceeded, creating new file.")

			if err := w.nextFile(); err != nil {
				w.log("Error occurred while creating file for next part: %v", err)
				break
			}
		}

		if w.SplitDuration > 0 && w.SegmentDuration >= w.SplitDuration*60 {
			w.log("Duration has exceeded, creating new file.")

			if err := w.nextFile(); err != nil {
				w.log("Error occurred while creating file for next part: %v", err)
				break
			}
		}

		w.bufferLock.Lock()
		delete(w.buffer, w.bufferIndex)
		w.bufferLock.Unlock()

		w.bufferIndex++
	}
}

func (w *Channel) requestChunks() ([]*m3u8.MediaSegment, float64, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}

	resp, err := client.Get(w.sourceURL)
	if err != nil {
		return nil, 3, fmt.Errorf("client get: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusForbidden:
			return nil, 3, fmt.Errorf("received status code %d, the stream is private?", resp.StatusCode)
		default:
			return nil, 3, fmt.Errorf("received status code %d", resp.StatusCode)
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 3, fmt.Errorf("read body: %w", err)
	}

	p, _, err := m3u8.DecodeFrom(bytes.NewReader(body), true)
	if err != nil {
		return nil, 3, fmt.Errorf("decode m3u8: %w", err)
	}
	playlist, ok := p.(*m3u8.MediaPlaylist)
	if !ok {
		return nil, 3, fmt.Errorf("cast to media playlist")
	}
	chunks := lo.Filter(playlist.Segments, func(v *m3u8.MediaSegment, _ int) bool {
		return v != nil
	})

	//log.Println(playlist.TargetDuration)

	return chunks, 1, nil
}

// fetchSegments
func (w *Channel) fetchSegments() {
	var disconnectRetries int

	for {
		if w.IsPaused || w.isStopped {
			break
		}

		chunks, wait, err := w.requestChunks()
		if err != nil {
			if disconnectRetries > 10 {
				w.IsOnline = false
				break
			}

			w.log("Error occurred while parsing m3u8: %v", err)
			disconnectRetries++

			<-time.After(time.Duration(wait) * time.Second)
			continue
		}

		if disconnectRetries > 0 {
			w.log("Stream is online")
			w.IsOnline = true
			disconnectRetries = 0
		}

		for _, v := range chunks {
			if w.isSegmentFetched(v.URI) {
				continue
			}

			go w.requestSegment(v.URI, w.segmentIndex)
			w.SegmentDuration += int(v.Duration)
			w.segmentIndex++
		}
		<-time.After(time.Duration(wait) * time.Second)
	}
}

func (w *Channel) requestSegment(url string, index int) error {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}

	resp, err := client.Get(w.rootURL + url)
	if err != nil {
		return fmt.Errorf("client get: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}

	w.log("Segment #%d fetched.", index)

	w.bufferLock.Lock()
	w.buffer[index] = body
	w.bufferLock.Unlock()

	return nil
}

func (w *Channel) Pause() {
	w.IsPaused = true
	w.resetSession()
	w.log("Channel was paused.")
}

func (w *Channel) Resume() {
	w.IsPaused = false
	w.ResumeChannel <- true //BUG:
	w.log("Channel was resumed.")
}

func (w *Channel) Stop() {
	w.isStopped = true
	w.log("Channel was stopped.")
}
