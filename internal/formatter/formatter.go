// Package formatter provides utilities for rendering log anomaly alerts
// and report summaries into human-readable or structured output formats.
package formatter

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Format represents the output format for alerts and reports.
type Format int

const (
	// Text renders output as plain human-readable text.
	Text Format = iota
	// JSON renders output as a JSON object.
	JSON
)

// ParseFormat converts a string (e.g. from config) to a Format constant.
// Returns Text and an error if the value is unrecognised.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "text", "":
		return Text, nil
	case "json":
		return JSON, nil
	default:
		return Text, fmt.Errorf("unknown format %q: must be \"text\" or \"json\"", s)
	}
}

// Entry is a structured representation of a single anomaly event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	ZScore    float64   `json:"z_score"`
	Mean      float64   `json:"mean"`
	StdDev    float64   `json:"std_dev"`
}

// Formatter renders Entry values into a chosen format.
type Formatter struct {
	fmt Format
}

// New returns a Formatter configured for the given Format.
func New(f Format) *Formatter {
	return &Formatter{fmt: f}
}

// Render converts an Entry to a string in the configured format.
// For JSON format any marshalling error falls back to text.
func (f *Formatter) Render(e Entry) string {
	switch f.fmt {
	case JSON:
		b, err := json.Marshal(e)
		if err == nil {
			return string(b)
		}
		fallthrough
	default:
		return fmt.Sprintf(
			"[%s] %s | z=%.2f mean=%.2f stddev=%.2f | %s",
			e.Timestamp.UTC().Format(time.RFC3339),
			e.Level,
			e.ZScore,
			e.Mean,
			e.StdDev,
			e.Message,
		)
	}
}
