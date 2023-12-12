package handler

import "fmt"

func formatPlaytime(seconds int) string {
	hours := seconds / 3600
	seconds %= 3600
	minutes := seconds / 60
	seconds %= 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
