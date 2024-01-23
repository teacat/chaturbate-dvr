package chaturbate

import (
	"fmt"
	"time"
)

// log
func (w *Channel) log(message string, v ...interface{}) {
	updateLog := fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), fmt.Errorf(message, v...))
	consoleLog := fmt.Sprintf("[%s] [%s] %s", time.Now().Format("2006-01-02 15:04:05"), w.Username, fmt.Errorf(message, v...))

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

	w.UpdateChannel <- update

	fmt.Println(consoleLog)

	w.Logs = append(w.Logs, updateLog)

	// Only keep the last 100 logs in memory.
	if len(w.Logs) > 100 {
		w.Logs = w.Logs[len(w.Logs)-100:]
	}
}

// isDuplicateSegment returns true if the segment is already been fetched.
func (w *Channel) isSegmentFetched(url string) bool {
	for _, v := range w.segmentUseds {
		if url[len(url)-10:] == v {
			return true
		}
	}
	w.segmentUseds = append(w.segmentUseds, url[len(url)-10:])
	return false
}
