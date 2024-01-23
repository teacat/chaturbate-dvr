package handler

import "fmt"

func DurationStr(seconds int) string {
	hours := seconds / 3600
	seconds %= 3600
	minutes := seconds / 60
	seconds %= 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func ByteStr(bytes int) string {
	return fmt.Sprintf("%.2f MiB", float64(bytes)/1024/1024)
}

func KBStr(kibs int) string {
	return fmt.Sprintf("%.2f MiB", float64(kibs)/1024)
}

func MBStr(mibs int) string {
	return fmt.Sprintf("%.2f MiB", float64(mibs))
}
