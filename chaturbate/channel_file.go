package chaturbate

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

// filename generates the filename based on the session pattern and current split index.
func (w *Channel) filename() (string, error) {
	if w.sessionPattern == nil {
		w.sessionPattern = map[string]any{
			"Username": w.Username,
			"Year":     time.Now().Format("2006"),
			"Month":    time.Now().Format("01"),
			"Day":      time.Now().Format("02"),
			"Hour":     time.Now().Format("15"),
			"Minute":   time.Now().Format("04"),
			"Second":   time.Now().Format("05"),
			"Sequence": 0,
		}
	}

	w.sessionPattern["Sequence"] = w.splitIndex

	var buf bytes.Buffer
	tmpl, err := template.New("filename").Parse(w.filenamePattern)
	if err != nil {
		return "", fmt.Errorf("filename pattern error: %w", err)
	}
	if err := tmpl.Execute(&buf, w.sessionPattern); err != nil {
		return "", fmt.Errorf("template execution error: %w", err)
	}

	return buf.String(), nil
}

// newFile creates a new file and prepares it for writing stream data.
func (w *Channel) newFile() error {
	filename, err := w.filename()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(filename), 0777); err != nil {
		return fmt.Errorf("create folder: %w", err)
	}

	file, err := os.OpenFile(filename+".ts", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		return fmt.Errorf("cannot open file: %s: %w", filename, err)
	}

	w.log(LogTypeInfo, "the stream will be saved as %s.ts", filename)
	w.file = file
	return nil
}

// releaseFile closes the current file and removes it if empty.
func (w *Channel) releaseFile() error {
	if w.file == nil {
		return nil
	}

	if err := w.file.Close(); err != nil {
		return fmt.Errorf("close file: %s: %w", w.file.Name(), err)
	}

	if w.SegmentFilesize == 0 {
		w.log(LogTypeInfo, "%s was removed because it was empty", w.file.Name())
		if err := os.Remove(w.file.Name()); err != nil {
			return fmt.Errorf("remove zero file: %s: %w", w.file.Name(), err)
		}
	}

	w.file = nil
	return nil
}

// nextFile handles the transition to a new file segment, ensuring correct timing.
func (w *Channel) nextFile(startTime time.Time) error {
	// Release the current file before creating a new one.
	if err := w.releaseFile(); err != nil {
		w.log(LogTypeError, "release file: %v", err)
		return err
	}

	// Increment the split index for the next file.
	w.splitIndex++

	// Reset segment data.
	w.SegmentFilesize = 0

	// Calculate the actual segment duration using the elapsed time.
	elapsed := int(time.Since(startTime).Minutes())
	w.SegmentDuration = elapsed

	// Create the new file.
	return w.newFile()
}
