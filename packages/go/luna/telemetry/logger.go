// Package telemetry provides logging and metrics for the Luna SDK.
package telemetry

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"time"
)

// LogLevel represents logging verbosity
type LogLevel int

const (
	LogLevelError LogLevel = 50
	LogLevelWarn  LogLevel = 40
	LogLevelInfo  LogLevel = 30
	LogLevelDebug LogLevel = 20
	LogLevelTrace LogLevel = 10
)

// Logger is the interface for SDK logging
type Logger interface {
	Error(message string, context map[string]interface{})
	Warn(message string, context map[string]interface{})
	Info(message string, context map[string]interface{})
	Debug(message string, context map[string]interface{})
	Trace(message string, context map[string]interface{})
}

// ConsoleLogger logs to stdout/stderr with JSON formatting
type ConsoleLogger struct {
	level    LogLevel
	redactRe []*regexp.Regexp
}

// NewConsoleLogger creates a new console logger
func NewConsoleLogger(level LogLevel) *ConsoleLogger {
	return &ConsoleLogger{
		level: level,
		redactRe: []*regexp.Regexp{
			regexp.MustCompile(`(?i)api[_-]?key`),
			regexp.MustCompile(`(?i)authorization`),
			regexp.MustCompile(`(?i)x-luna-api-key`),
			regexp.MustCompile(`(?i)password`),
			regexp.MustCompile(`(?i)secret`),
			regexp.MustCompile(`(?i)token`),
			regexp.MustCompile(`(?i)bearer`),
			// POPIA / SA Specific
			regexp.MustCompile(`(?i)id[_-]?number`),
			regexp.MustCompile(`(?i)tax[_-]?ref`),
			regexp.MustCompile(`(?i)registration[_-]?number`),
			regexp.MustCompile(`(?i)account[_-]?number`),
			regexp.MustCompile(`(?i)cvv`),
			regexp.MustCompile(`(?i)pan`),
		},
	}
}

func (l *ConsoleLogger) Error(message string, context map[string]interface{}) {
	l.log(LogLevelError, "ERROR", message, context)
}

func (l *ConsoleLogger) Warn(message string, context map[string]interface{}) {
	l.log(LogLevelWarn, "WARN", message, context)
}

func (l *ConsoleLogger) Info(message string, context map[string]interface{}) {
	l.log(LogLevelInfo, "INFO", message, context)
}

func (l *ConsoleLogger) Debug(message string, context map[string]interface{}) {
	l.log(LogLevelDebug, "DEBUG", message, context)
}

func (l *ConsoleLogger) Trace(message string, context map[string]interface{}) {
	l.log(LogLevelTrace, "TRACE", message, context)
}

func (l *ConsoleLogger) log(level LogLevel, levelStr, message string, context map[string]interface{}) {
	if level < l.level {
		return
	}

	entry := map[string]interface{}{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"level":     levelStr,
		"message":   message,
		"sdk":       "luna-sdk",
		"version":   "1.0.0",
		"language":  "go",
	}

	if context != nil {
		entry["context"] = l.sanitize(context)
	}

	output, _ := json.Marshal(entry)

	if level >= LogLevelError {
		fmt.Fprintln(os.Stderr, string(output))
	} else {
		fmt.Println(string(output))
	}
}

func (l *ConsoleLogger) sanitize(obj map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range obj {
		if l.isSensitiveKey(key) {
			result[key] = "[REDACTED]"
		} else if nested, ok := value.(map[string]interface{}); ok {
			result[key] = l.sanitize(nested)
		} else {
			result[key] = value
		}
	}

	return result
}

func (l *ConsoleLogger) isSensitiveKey(key string) bool {
	for _, re := range l.redactRe {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}

var _ Logger = (*ConsoleLogger)(nil)
