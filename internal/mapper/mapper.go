// Package mapper translates raw parsed log lines into enriched pipeline entries.
package mapper

import (
	"fmt"
	"time"

	"github.com/yourorg/logdrift/internal/parser"
	"github.com/yourorg/logdrift/internal/severity"
)

// Entry is a normalised, enriched representation of a single log line.
type Entry struct {
	Timestamp time.Time
	Level     severity.Level
	Message   string
	Service   string
	Latency   float64
	Raw       string
}

// Mapper converts parser.LogLine values into Entry values.
type Mapper struct {
	defaultService string
}

// Option configures a Mapper.
type Option func(*Mapper)

// WithDefaultService sets the service name used when none can be inferred.
func WithDefaultService(s string) Option {
	return func(m *Mapper) { m.defaultService = s }
}

// New returns a new Mapper.
func New(opts ...Option) *Mapper {
	m := &Mapper{defaultService: "unknown"}
	for _, o := range opts {
		o(m)
	}
	return m
}

// Map converts a parser.LogLine into an Entry.
func (m *Mapper) Map(line parser.LogLine) (Entry, error) {
	lvl, err := severity.Parse(line.Level)
	if err != nil {
		return Entry{}, fmt.Errorf("mapper: unknown level %q: %w", line.Level, err)
	}
	svc := line.Service
	if svc == "" {
		svc = m.defaultService
	}
	return Entry{
		Timestamp: line.Timestamp,
		Level:     lvl,
		Message:   line.Message,
		Service:   svc,
		Latency:   line.Latency,
		Raw:       line.Raw,
	}, nil
}
