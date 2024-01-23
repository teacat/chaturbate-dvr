package chaturbate

type Update struct {
	Username        string `json:"username"`
	Log             string `json:"log"`
	IsPaused        bool   `json:"is_paused"`
	IsOnline        bool   `json:"is_online"`
	IsStopped       bool   `json:"is_stopped"`
	Filename        string `json:"filename"`
	LastStreamedAt  string `json:"last_streamed_at"`
	SegmentDuration int    `json:"segment_duration"`
	SegmentFilesize int    `json:"segment_filesize"`
}

func (u *Update) SegmentDurationStr() string {
	return DurationStr(u.SegmentDuration)
}

func (u *Update) SegmentFilesizeStr() string {
	return ByteStr(u.SegmentFilesize)
}
