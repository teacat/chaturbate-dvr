package internal

import "errors"

var (
	ErrChannelExists     = errors.New("channel exists")
	ErrChannelNotFound   = errors.New("channel not found")
	ErrCloudflareBlocked = errors.New("blocked by Cloudflare; try with `-cookies` and `-user-agent`")
	ErrChannelOffline    = errors.New("channel offline")
	ErrPrivateStream     = errors.New("channel went offline")
	ErrPaused            = errors.New("channel paused")
	ErrStopped           = errors.New("channel stopped")
)
