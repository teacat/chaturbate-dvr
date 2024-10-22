package chaturbate

import (
	"fmt"
	"time"
)

// log
func (w *Channel) log(typ LogType, message string, v ...interface{}) {
	// Check the global log level
	currentLogLevel := GetGlobalLogLevel()

	switch currentLogLevel {
	case LogTypeInfo:
		if typ == LogTypeDebug {
			return
		}
	case LogTypeWarning:
		if typ == LogTypeDebug || typ == LogTypeInfo {
			return
		}
	case LogTypeError:
		if typ == LogTypeDebug || typ == LogTypeInfo || typ == LogTypeWarning {
			return
		}
	}

	updateLog := fmt.Sprintf("[%s] [%s] %s", time.Now().Format("2006-01-02 15:04:05"), typ, fmt.Sprintf(message, v...))
	consoleLog := fmt.Sprintf("[%s] [%s] [%s] %s", time.Now().Format("2006-01-02 15:04:05"), typ, w.Username, fmt.Sprintf(message, v...))

	update := &Update{
		Username:        w.Username,
		Log:             updateLog,
		IsPaused:        w.IsPaused,
		IsOnline:        w.IsOnline,
		LastStreamedAt:  w.LastStreamedAt,
		SegmentDuration: w.SegmentDuration,
		SegmentFilesize: w.SegmentFilesize,
	}

	if w.file != nil {
		update.Filename = w.file.Name()
	}

	select {
	case w.UpdateChannel <- update:
	default:
	}

	fmt.Println(consoleLog)

	w.Logs = append(w.Logs, updateLog)

	// Only keep the last 100 logs in memory.
	if len(w.Logs) > 100 {
		w.Logs = w.Logs[len(w.Logs)-100:]
	}
}

// isSegmentFetched returns true if the segment has been fetched.
func (w *Channel) isSegmentFetched(url string) bool {
	for _, v := range w.segmentUseds {
		if url == v {
			return true
		}
	}
	if len(w.segmentUseds) > 100 {
		w.segmentUseds = w.segmentUseds[len(w.segmentUseds)-30:]
	}
	w.segmentUseds = append(w.segmentUseds, url)
	return false
}

func DurationStr(seconds int) string {
	hours := seconds / 3600
	seconds %= 3600
	minutes := seconds / 60
	seconds %= 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func MBStr(mibs int) string {
	return fmt.Sprintf("%.2f MiB", float64(mibs))
}

func ByteStr(bytes int) string {
	return fmt.Sprintf("%.2f MiB", float64(bytes)/1024/1024)
}
