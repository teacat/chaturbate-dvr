package chaturbate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/grafov/m3u8"
	"github.com/samber/lo"
	"github.com/teacat/chaturbate-dvr/internal"
	"github.com/teacat/chaturbate-dvr/server"
)

// roomDossierRegexp is used to extract the room dossier information from the HTML response.
var roomDossierRegexp = regexp.MustCompile(`window\.initialRoomDossier = "(.*?)"`)

// Client represents an API client for interacting with Chaturbate.
type Client struct {
	Req *internal.Req
}

// NewClient initializes and returns a new Client instance.
func NewClient() *Client {
	return &Client{
		Req: internal.NewReq(),
	}
}

// GetStream fetches the stream information for a given username.
func (c *Client) GetStream(ctx context.Context, username string) (*Stream, error) {
	return FetchStream(ctx, c.Req, username)
}

// FetchStream retrieves the streaming data from the given username's page.
func FetchStream(ctx context.Context, client *internal.Req, username string) (*Stream, error) {
	body, err := client.Get(ctx, fmt.Sprintf("%s%s", server.Config.Domain, username))
	if err != nil {
		return nil, fmt.Errorf("failed to get page body: %w", err)
	}

	// Ensure that the playlist.m3u8 file is present in the response
	if !strings.Contains(body, "playlist.m3u8") {
		return nil, internal.ErrChannelOffline
	}

	return ParseStream(body)
}

// ParseStream extracts the HLS source URL from the given page body.
func ParseStream(body string) (*Stream, error) {
	matches := roomDossierRegexp.FindStringSubmatch(body)
	if len(matches) == 0 {
		return nil, errors.New("room dossier not found")
	}

	// Decode Unicode escape sequences in the extracted JSON string
	sourceData, err := strconv.Unquote(strings.Replace(strconv.Quote(matches[1]), `\\u`, `\u`, -1))
	if err != nil {
		return nil, fmt.Errorf("failed to decode unicode: %w", err)
	}

	// Unmarshal JSON to extract HLS source URL
	var room struct {
		HLSSource string `json:"hls_source"`
	}
	if err := json.Unmarshal([]byte(sourceData), &room); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &Stream{HLSSource: room.HLSSource}, nil
}

// Stream represents an HLS stream source.
type Stream struct {
	HLSSource string
}

// GetPlaylist retrieves the playlist corresponding to the given resolution and framerate.
func (s *Stream) GetPlaylist(ctx context.Context, resolution, framerate int) (*Playlist, error) {
	return FetchPlaylist(ctx, s.HLSSource, resolution, framerate)
}

// FetchPlaylist fetches and decodes the HLS playlist file.
func FetchPlaylist(ctx context.Context, hlsSource string, resolution, framerate int) (*Playlist, error) {
	if hlsSource == "" {
		return nil, errors.New("HLS source is empty")
	}

	resp, err := internal.NewReq().Get(ctx, hlsSource)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch HLS source: %w", err)
	}

	return ParsePlaylist(resp, hlsSource, resolution, framerate)
}

// ParsePlaylist decodes the M3U8 playlist and extracts the variant streams.
func ParsePlaylist(resp, hlsSource string, resolution, framerate int) (*Playlist, error) {
	p, _, err := m3u8.DecodeFrom(strings.NewReader(resp), true)
	if err != nil {
		return nil, fmt.Errorf("failed to decode m3u8 playlist: %w", err)
	}

	masterPlaylist, ok := p.(*m3u8.MasterPlaylist)
	if !ok {
		return nil, errors.New("invalid master playlist format")
	}

	return PickPlaylist(masterPlaylist, hlsSource, resolution, framerate)
}

// Playlist represents an HLS playlist containing variant streams.
type Playlist struct {
	PlaylistURL string
	RootURL     string
	Resolution  int
	Framerate   int
}

// Resolution represents a video resolution and its corresponding framerate.
type Resolution struct {
	Framerate map[int]string // [framerate]url
	Width     int
}

// PickPlaylist selects the best matching variant stream based on resolution and framerate.
func PickPlaylist(masterPlaylist *m3u8.MasterPlaylist, baseURL string, resolution, framerate int) (*Playlist, error) {
	resolutions := map[int]*Resolution{}

	// Extract available resolutions and framerates from the master playlist
	for _, v := range masterPlaylist.Variants {
		parts := strings.Split(v.Resolution, "x")
		if len(parts) != 2 {
			continue
		}
		width, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("parse resolution: %w", err)
		}
		framerateVal := 30
		if strings.Contains(v.Name, "FPS:60.0") {
			framerateVal = 60
		}
		if _, exists := resolutions[width]; !exists {
			resolutions[width] = &Resolution{Framerate: map[int]string{}, Width: width}
		}
		resolutions[width].Framerate[framerateVal] = v.URI
	}

	// Find exact match for requested resolution
	variant, exists := resolutions[resolution]
	if !exists {
		// Filter resolutions below the requested resolution
		candidates := lo.Filter(lo.Values(resolutions), func(r *Resolution, _ int) bool {
			return r.Width < resolution
		})
		// Pick the highest resolution among the candidates
		variant = lo.MaxBy(candidates, func(a, b *Resolution) bool {
			return a.Width > b.Width
		})
	}
	if variant == nil {
		return nil, fmt.Errorf("resolution not found")
	}

	var (
		finalResolution = variant.Width
		finalFramerate  = framerate
	)
	// Select the desired framerate, or fallback to the first available framerate
	playlistURL, exists := variant.Framerate[framerate]
	if !exists {
		for fr, url := range variant.Framerate {
			playlistURL = url
			finalFramerate = fr
			break
		}
	}

	return &Playlist{
		PlaylistURL: strings.TrimSuffix(baseURL, "playlist.m3u8") + playlistURL,
		RootURL:     strings.TrimSuffix(baseURL, "playlist.m3u8"),
		Resolution:  finalResolution,
		Framerate:   finalFramerate,
	}, nil
}

// WatchHandler is a function type that processes video segments.
type WatchHandler func(b []byte, duration float64) error

// WatchSegments continuously fetches and processes video segments.
func (p *Playlist) WatchSegments(ctx context.Context, handler WatchHandler) error {
	var (
		client  = internal.NewReq()
		lastSeq = -1
	)

	for {
		// Fetch the latest playlist
		resp, err := client.Get(ctx, p.PlaylistURL)
		if err != nil {
			return fmt.Errorf("get playlist: %w", err)
		}
		pl, _, err := m3u8.DecodeFrom(strings.NewReader(resp), true)
		if err != nil {
			return fmt.Errorf("decode from: %w", err)
		}
		playlist, ok := pl.(*m3u8.MediaPlaylist)
		if !ok {
			return fmt.Errorf("cast to media playlist")
		}

		// Process new segments
		for _, v := range playlist.Segments {
			if v == nil {
				continue
			}
			seq := internal.SegmentSeq(v.URI)
			if seq == -1 || seq <= lastSeq {
				continue
			}
			lastSeq = seq

			// Fetch segment data with retry mechanism
			pipeline := func() ([]byte, error) {
				return client.GetBytes(ctx, fmt.Sprintf("%s%s", p.RootURL, v.URI))
			}

			resp, err := retry.DoWithData(
				pipeline,
				retry.Context(ctx),
				retry.Attempts(3),
				retry.Delay(600*time.Millisecond),
				retry.DelayType(retry.FixedDelay),
			)
			if err != nil {
				break
			}

			// Process the segment using the provided handler
			if err := handler(resp, v.Duration); err != nil {
				return fmt.Errorf("handler: %w", err)
			}
		}

		<-time.After(1 * time.Second) // time.Duration(playlist.TargetDuration)
	}
}
