package manager

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/grafov/m3u8"
	"github.com/samber/lo"
	"github.com/teacat/chaturbate-dvr/domain/stream"
)

var (
	regexpRoomDossier = regexp.MustCompile(`window\.initialRoomDossier = "(.*?)"`)
)

type roomDossier struct {
	HLSSource string `json:"hls_source"`
}

type source struct {
	framerate map[int]string // key: framerate, value: url
	size      int
}

func parseHLS(resolution, framerate int, resFallback stream.ResolutionFallback, body string) (string, error) {
	// Find the room dossier.
	matches := regexpRoomDossier.FindAllStringSubmatch(body, -1)

	// Get the HLS source from the room dossier.
	var roomData roomDossier
	data, err := strconv.Unquote(strings.Replace(strconv.Quote(string(matches[0][1])), `\\u`, `\u`, -1))
	if err != nil {
		return "", fmt.Errorf("unquote unicode: %w", err)
	}
	if err := json.Unmarshal([]byte(data), &roomData); err != nil {
		return "", fmt.Errorf("unmarshal json: %w", err)
	}

	// Get the HLS source.
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(roomData.HLSSource)
	if err != nil {
		return "", fmt.Errorf("client get: %w", err)
	}
	defer resp.Body.Close()

	m3u8Body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}
	if resp.StatusCode == http.StatusForbidden {
		return "", fmt.Errorf("received status code %d", resp.StatusCode)
	}

	// Decode the m3u8 file.
	p, _, err := m3u8.DecodeFrom(bytes.NewReader(m3u8Body), true)
	if err != nil {
		return "", fmt.Errorf("decode m3u8: %w", err)
	}
	playlist, ok := p.(*m3u8.MasterPlaylist)
	if !ok {
		return "", fmt.Errorf("cast to master playlist")
	}

	var sources []*source

	//
	for _, v := range playlist.Variants {
	}

	variant, ok := lo.Find(sources, func(v *source) bool {
		return v.size == resolution
	})
	// If the variant is not found, we fallback to the nearest resolution.
	if !ok {
		switch resFallback {
		case stream.ResolutionFallbackDownscale:
			variant = lo.MinBy(sources, func(v, min *source) bool {
				// return v.size < resolution && v.size < min.size
				return math.Abs(float64(v.size-resolution)) < math.Abs(float64(min.size-resolution))
			})
		case stream.ResolutionFallbackUpscale:
			variant = lo.MaxBy(sources, func(v, max *source) bool {
				return math.Abs(float64(v.size-resolution)) > math.Abs(float64(max.size-resolution))
			})
		}
	}
	if variant == nil {
		return "", fmt.Errorf("variant not found")
	}

	uri, ok := variant.framerate[framerate]
	// If the framerate is not found, we fallback to the nearest framerate.
	if !ok {
		for _, v := range variant.framerate {
			uri = v
			// TODO: log fps
			break
		}
	}

	baseURL := strings.TrimSuffix(roomData.HLSSource, "playlist.m3u8")
	return baseURL + uri, nil
}
