package internal

import (
	"fmt"
	"regexp"
	"strconv"
)

// FormatDuration converts a float64 duration (in seconds) to h:m:s format.
func FormatDuration(duration float64) string {
	if duration == 0 {
		return ""
	}
	var (
		hours   = int(duration) / 3600
		minutes = (int(duration) % 3600) / 60
		seconds = int(duration) % 60
	)
	return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
}

// FormatFilesize converts an int filesize in bytes to a human-readable string (KB, MB, GB).
func FormatFilesize(filesize int) string {
	if filesize == 0 {
		return ""
	}
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case filesize >= GB:
		return fmt.Sprintf("%.2f GB", float64(filesize)/float64(GB))
	case filesize >= MB:
		return fmt.Sprintf("%.2f MB", float64(filesize)/float64(MB))
	case filesize >= KB:
		return fmt.Sprintf("%.2f KB", float64(filesize)/float64(KB))
	default:
		return fmt.Sprintf("%d bytes", filesize)
	}
}

// SegmentSeq extracts the segment sequence number from a filename.
func SegmentSeq(filename string) int {
	re := regexp.MustCompile(`_(\d+)\.ts$`)
	match := re.FindStringSubmatch(filename)

	if len(match) > 1 {
		number, err := strconv.Atoi(match[1])
		if err == nil {
			return number
		}
	}
	return -1
}
