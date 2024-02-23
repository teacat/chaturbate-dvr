package chaturbate

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

// filename
func (w *Channel) filename() (string, error) {
	data := w.sessionPattern
	if data == nil {
		data = map[string]any{
			"Username": w.Username,
			"Year":     time.Now().Format("2006"),
			"Month":    time.Now().Format("01"),
			"Day":      time.Now().Format("02"),
			"Hour":     time.Now().Format("15"),
			"Minute":   time.Now().Format("04"),
			"Second":   time.Now().Format("05"),
			"Sequence": 0,
		}
		w.sessionPattern = data
	} else {
		data["Sequence"] = w.splitIndex
	}
	t, err := template.New("filename").Parse(w.filenamePattern)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// newFile
func (w *Channel) newFile() error {
	filename, err := w.filename()
	if err != nil {
		return fmt.Errorf("filename pattern error: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(filename), 0777); err != nil {
		return fmt.Errorf("create folder: %w", err)
	}
	file, err := os.OpenFile(filename+".ts", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		return fmt.Errorf("cannot open file: %s: %w", filename, err)
	}
	w.log(logTypeInfo, "the stream will be saved as %s.ts", filename)
	w.file = file
	return nil
}

// releaseFile
func (w *Channel) releaseFile() error {
	if w.file == nil {
		return nil
	}
	// close the file to remove it
	if err := w.file.Close(); err != nil {
		return fmt.Errorf("close file: %s: %w", w.file.Name(), err)
	}
	// remove it if it was empty
	if w.SegmentFilesize == 0 {
		w.log(logTypeInfo, "%s was removed because it was empty", w.file.Name())

		if err := os.Remove(w.file.Name()); err != nil {
			return fmt.Errorf("remove zero file: %s: %w", w.file.Name(), err)
		}
	}
	w.file = nil
	return nil
}

// nextFile
func (w *Channel) nextFile() error {
	if err := w.releaseFile(); err != nil {
		w.log(logTypeError, "release file: %w", err)
	}

	w.splitIndex++
	w.SegmentFilesize = 0
	w.SegmentDuration = 0

	return w.newFile()
}
