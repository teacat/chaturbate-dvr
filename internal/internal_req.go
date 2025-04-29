package internal

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/teacat/chaturbate-dvr/server"
)

// Req represents an HTTP client with customized settings.
type Req struct {
	client *http.Client
}

// NewReq creates a new HTTP client with specific transport configurations.
func NewReq() *Req {
	return &Req{
		client: &http.Client{
			Transport: CreateTransport(),
		},
	}
}

// CreateTransport initializes a custom HTTP transport.
func CreateTransport() *http.Transport {
	return &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

// Get sends an HTTP GET request and returns the response as a string.
func (h *Req) Get(ctx context.Context, url string) (string, error) {
	resp, err := h.GetBytes(ctx, url)
	if err != nil {
		return "", fmt.Errorf("get bytes: %w", err)
	}
	return string(resp), nil
}

// GetBytes sends an HTTP GET request and returns the response as a byte slice.
func (h *Req) GetBytes(ctx context.Context, url string) ([]byte, error) {
	req, err := CreateRequest(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("forbidden: %w", ErrPrivateStream)
	}

	return ReadResponseBody(resp)
}

// CreateRequest constructs an HTTP GET request with necessary headers.
func CreateRequest(ctx context.Context, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	SetRequestHeaders(req)
	return req, nil
}

// SetRequestHeaders applies necessary headers to the request.
func SetRequestHeaders(req *http.Request) {
	if server.Config.UserAgent != "" {
		req.Header.Set("User-Agent", server.Config.UserAgent)
	}
	if server.Config.Cookies != "" {
		cookies := ParseCookies(server.Config.Cookies)
		for name, value := range cookies {
			req.AddCookie(&http.Cookie{Name: name, Value: value})
		}
	}
}

// ReadResponseBody reads and returns the response body.
func ReadResponseBody(resp *http.Response) ([]byte, error) {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	return b, nil
}

// ParseCookies converts a cookie string into a map.
func ParseCookies(cookieStr string) map[string]string {
	cookies := make(map[string]string)
	pairs := strings.Split(cookieStr, ";")

	// Iterate over each cookie pair and extract key-value pairs
	for _, pair := range pairs {
		parts := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(parts) == 2 {
			// Store cookie name and value in the map
			cookies[parts[0]] = parts[1]
		}
	}
	return cookies
}
