package chaturbate

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/grafov/m3u8"
	"github.com/samber/lo"
)

// requestChannelBody requests the channel page and returns the body.
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

// record starts the recording process,
// this function get called when the channel is online and back online from offline status.
//
// this is a blocking function until fetching segments gone wrong (or nothing to fetch, aka offline).
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

// resetSession resets the session data,
// usually called when the channel is online or paused to resumed.
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

// resolveSource resolves the HLS source from the channel page.
// the HLS Source is a list that contains all the available resolutions and framerates.
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
			return "", "", fmt.Errorf("ticket/private stream?")
		default:
			return "", "", fmt.Errorf("status code %d", resp.StatusCode)
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
		return "", "", fmt.Errorf("no available resolution")
	}
	w.log(LogTypeInfo, "resolution %dp is used", variant.width)

	url, ok := variant.framerate[w.Framerate]
	// If the framerate is not found, fallback to the first found framerate, this block pretends there're only 30 and 60 fps.
	// no complex logic here, im lazy.
	if ok {
		w.log(LogTypeInfo, "framerate %dfps is used", w.Framerate)
	} else {
		for k, v := range variant.framerate {
			url = v
			w.log(LogTypeWarning, "framerate %dfps not found, fallback to %dfps", w.Framerate, k)
			w.Framerate = k
			break
		}
	}

	rootURL := strings.TrimSuffix(roomData.HLSSource, "playlist.m3u8")
	sourceURL := rootURL + url
	return rootURL, sourceURL, nil
}

// mergeSegments runs in the background and merges segments from the buffer to the file.
func (w *Channel) mergeSegments() {
	var segmentRetries int
	startTime := time.Now() // Track the start time of the current segment.

	for {
		if w.IsPaused || w.isStopped {
			break
		}

		// Handle segment retries if not found.
		if segmentRetries > 5 {
			w.log(LogTypeWarning, "segment #%d not found in buffer, skipped", w.bufferIndex)
			w.bufferIndex++
			segmentRetries = 0
			continue
		}

		// If buffer is empty, wait and retry.
		if len(w.buffer) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		// Retrieve segment from buffer.
		w.bufferLock.Lock()
		buf, ok := w.buffer[w.bufferIndex]
		w.bufferLock.Unlock()

		if !ok {
			segmentRetries++
			time.Sleep(time.Duration(segmentRetries) * time.Second)
			continue
		}

		// Write the segment to the file.
		lens, err := w.file.Write(buf)
		if err != nil {
			w.log(LogTypeError, "segment #%d written error: %v", w.bufferIndex, err)
			w.retries++
			continue
		}

		// Update segment size and log progress.
		w.SegmentFilesize += lens
		w.log(LogTypeInfo, "segment #%d written", w.bufferIndex)
		w.log(LogTypeDebug, "duration: %s, size: %s", DurationStr(w.SegmentDuration), ByteStr(w.SegmentFilesize))

		// Check if the file size limit has been reached.
		if w.SplitFilesize > 0 && w.SegmentFilesize >= w.SplitFilesize*1024*1024 {
			w.log(LogTypeInfo, "filesize exceeded, creating new file")

			if err := w.nextFile(startTime); err != nil {
				w.log(LogTypeError, "next file error: %v", err)
				break
			}

			startTime = time.Now() // Reset start time for the new segment.
		}

		// Check if the duration limit has been reached.
		elapsed := int(time.Since(startTime).Minutes())
		if w.SplitDuration > 0 && elapsed >= w.SplitDuration {
			w.log(LogTypeInfo, "duration exceeded, creating new file")

			if err := w.nextFile(startTime); err != nil {
				w.log(LogTypeError, "next file error: %v", err)
				break
			}

			startTime = time.Now() // Reset start time for the new segment.
		}

		// Remove the processed segment from the buffer.
		w.bufferLock.Lock()
		delete(w.buffer, w.bufferIndex)
		w.bufferLock.Unlock()

		w.bufferIndex++    // Move to the next segment.
		segmentRetries = 0 // Reset retries for the next segment.
	}
}

// fetchSegments is a blocking function,
// it will keep asking the segment list for the latest segments.
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

			w.log(LogTypeError, "segment list error, will try again [%d/10]: %v", disconnectRetries, err)
			disconnectRetries++

			<-time.After(time.Duration(wait) * time.Second)
			continue
		}

		if disconnectRetries > 0 {
			w.log(LogTypeInfo, "channel is back online!")
			w.IsOnline = true
			disconnectRetries = 0
		}

		for _, v := range chunks {
			if w.isSegmentFetched(v.URI) {
				continue
			}

			go func(index int, uri string) {
				if err := w.requestSegment(uri, index); err != nil {
					w.log(LogTypeError, "segment #%d request error, ignored: %v", index, err)
					return
				}
			}(w.segmentIndex, v.URI)
			w.SegmentDuration += int(v.Duration)
			w.segmentIndex++
		}
		<-time.After(time.Duration(wait) * time.Second)
	}
}

// requestChunks requests the segment list from the HLS source,
// the same segment list will be updated every few seconds from chaturbate.
func (w *Channel) requestChunks() ([]*m3u8.MediaSegment, float64, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}

	if w.sourceURL == "" {
		return nil, 0, fmt.Errorf("channel seems to be paused?")
	}

	resp, err := client.Get(w.sourceURL)
	if err != nil {
		return nil, 3, fmt.Errorf("client get: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusForbidden:
			return nil, 3, fmt.Errorf("ticket/private stream?")
		default:
			return nil, 3, fmt.Errorf("status code %d", resp.StatusCode)
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
	return chunks, playlist.TargetDuration, nil
}

// requestSegment requests the specific single segment and put it into the buffer.
// the mergeSegments function will merge the segment from buffer to the file in the backgrond.
func (w *Channel) requestSegment(url string, index int) error {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}

	if w.rootURL == "" {
		return fmt.Errorf("channel seems to be paused?")
	}

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

	w.log(LogTypeDebug, "segment #%d fetched", index)

	w.bufferLock.Lock()
	w.buffer[index] = body
	w.bufferLock.Unlock()

	return nil
}
