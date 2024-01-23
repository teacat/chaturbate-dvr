package chaturbate

import (
	"fmt"
	"time"
)

type logType string

const (
	logTypeDebug   logType = "DEBUG"
	logTypeInfo    logType = "INFO"
	logTypeWarning logType = "WARN"
	logTypeError   logType = "ERROR"
)

// log
func (w *Channel) log(typ logType, message string, v ...interface{}) {
	switch w.logType {
	case logTypeInfo:
		if typ == logTypeDebug {
			return
		}
	case logTypeWarning:
		if typ == logTypeDebug || typ == logTypeInfo {
			return
		}
	case logTypeError:
		if typ == logTypeDebug || typ == logTypeInfo || typ == logTypeWarning {
			return
		}
	}

	updateLog := fmt.Sprintf("[%s] [%s] %s", time.Now().Format("2006-01-02 15:04:05"), typ, fmt.Errorf(message, v...))
	consoleLog := fmt.Sprintf("[%s] [%s] [%s] %s", time.Now().Format("2006-01-02 15:04:05"), typ, w.Username, fmt.Errorf(message, v...))

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
