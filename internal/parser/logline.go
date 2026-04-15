package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// LogLevel represents the severity level of a log entry.
type LogLevel string

const (
	LevelDebug LogLevel = "DEBUG"
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
	LevelFatal LogLevel = "FATAL"
)

// LogLine represents a parsed log entry.
type LogLine struct {
	Timestamp time.Time
	Level     LogLevel
	Message   string
	LatencyMs float64
}

// defaultPattern matches lines like:
// 2024-01-15T10:23:45Z INFO  some message latency=123ms
var defaultPattern = regexp.MustCompile(
	`^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z)\s+(DEBUG|INFO|WARN|ERROR|FATAL)\s+(.+?)(?:\s+latency=(\d+)ms)?$`,
)

// Parse attempts to parse a raw log line string into a LogLine.
// Returns an error if the line does not match the expected format.
func Parse(raw string) (*LogLine, error) {
	matches := defaultPattern.FindStringSubmatch(raw)
	if matches == nil {
		return nil, fmt.Errorf("parser: line does not match expected format: %q", raw)
	}

	ts, err := time.Parse(time.RFC3339, matches[1])
	if err != nil {
		return nil, fmt.Errorf("parser: invalid timestamp %q: %w", matches[1], err)
	}

	var latency float64
	if matches[4] != "" {
		latency, err = strconv.ParseFloat(matches[4], 64)
		if err != nil {
			return nil, fmt.Errorf("parser: invalid latency value %q: %w", matches[4], err)
		}
	}

	return &LogLine{
		Timestamp: ts,
		Level:     LogLevel(matches[2]),
		Message:   matches[3],
		LatencyMs: latency,
	}, nil
}
