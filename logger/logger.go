package logger

import (
	"fmt"
	"os"
	"strings"
)

type Logger struct {
	file *os.File
}

func New() (*Logger, error) {
	if err := os.MkdirAll("./conf", 0777); err != nil {
		return nil, fmt.Errorf("mkdir all conf: %w", err)
	}
	f, err := os.OpenFile("./conf/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	return &Logger{
		file: f,
	}, nil
}

// Write writes a log entry to the file and truncates the file if it exceeds 1 MB.
func (l *Logger) Write(v string) error {
	// Write the log entry
	if _, err := l.file.WriteString(v + "\n"); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	// Check the file size
	info, err := l.file.Stat()
	if err != nil {
		return fmt.Errorf("stat file: %w", err)
	}

	// If the file size exceeds 1 MB, truncate to the last 1000 lines
	const maxSize = 1 * 1024 * 1024 // 1 MB
	if info.Size() > maxSize {
		if err := l.TruncateToLastLines(1000); err != nil {
			return fmt.Errorf("truncate file: %w", err)
		}
	}

	return nil
}

// TruncateToLastLines truncates the log file to keep only the last N lines.
func (l *Logger) TruncateToLastLines(lines int) error {
	// Read the entire file content
	content, err := os.ReadFile(l.file.Name())
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	// Split the content into lines
	allLines := SplitLines(string(content))
	if len(allLines) > lines {
		allLines = allLines[len(allLines)-lines:] // Keep only the last `lines` lines
	}

	// Truncate the file and write back the last `lines` lines
	if err := l.file.Truncate(0); err != nil {
		return fmt.Errorf("truncate file: %w", err)
	}
	if _, err := l.file.Seek(0, 0); err != nil {
		return fmt.Errorf("seek file: %w", err)
	}
	if _, err := l.file.WriteString(JoinLines(allLines)); err != nil {
		return fmt.Errorf("write truncated file: %w", err)
	}

	return nil
}

// SplitLines splits a string into lines using "\n" as the delimiter.
func SplitLines(content string) []string {
	return strings.Split(content, "\n")
}

// JoinLines joins a slice of strings into a single string using "\n" as the delimiter.
func JoinLines(lines []string) string {
	return strings.Join(lines, "\n")
}
