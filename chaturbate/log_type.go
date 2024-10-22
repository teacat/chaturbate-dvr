package chaturbate

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

type LogType string

type LogLevelRequest struct {
	LogLevel LogType `json:"log_level" binding:"required"`
}

// Define the log types
const (
	LogTypeDebug   LogType = "DEBUG"
	LogTypeInfo    LogType = "INFO"
	LogTypeWarning LogType = "WARN"
	LogTypeError   LogType = "ERROR"
)

// Global log level with mutex protection
var (
	globalLogLevel LogType
	logMutex       sync.RWMutex // Protects global log level access
)

// UnmarshalJSON ensures that LogType is properly parsed from JSON.
func (l *LogType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	parsed := LogType(strings.ToUpper(s))
	switch parsed {
	case LogTypeDebug, LogTypeInfo, LogTypeWarning, LogTypeError:
		*l = parsed
		return nil
	default:
		return fmt.Errorf("invalid log level: %s", s)
	}
}

// InitGlobalLogLevel initializes the global log level from settings.
func InitGlobalLogLevel(initialLevel LogType) {
	SetGlobalLogLevel(initialLevel)
}

// SetGlobalLogLevel updates the global log level
func SetGlobalLogLevel(level LogType) {
	logMutex.Lock()
	defer logMutex.Unlock()
	globalLogLevel = level
}

// GetGlobalLogLevel retrieves the current global log level
func GetGlobalLogLevel() LogType {
	logMutex.RLock()
	defer logMutex.RUnlock()
	return globalLogLevel
}
