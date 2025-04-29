package internal

import "errors"

var (
	ErrChannelExists     = errors.New("channel already exists")
	ErrChannelNotFound   = errors.New("channel not found")
	ErrCloudflareBlocked = errors.New("channel was blocked by Cloudflare, try again with `--cookies` and `--user-agent` options")
	ErrChannelOffline    = errors.New("channel is offline")
	ErrPrivateStream     = errors.New("possibly private stream, try `--cookies` option to login")
	ErrPaused            = errors.New("channel is paused")
	ErrStopped           = errors.New("channel is stopped")
)
