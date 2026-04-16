// Package truncator provides utilities for truncating log messages
// that exceed a configured byte length, appending a configurable suffix.
package truncator

import "errors"

// ErrInvalidMaxBytes is returned when maxBytes is not positive.
var ErrInvalidMaxBytes = errors.New("truncator: maxBytes must be greater than zero")

// Truncator trims strings to a maximum byte length.
type Truncator struct {
	maxBytes int
	suffix   string
}

// Option configures a Truncator.
type Option func(*Truncator)

// WithSuffix sets the suffix appended to truncated strings.
func WithSuffix(s string) Option {
	return func(t *Truncator) {
		t.suffix = s
	}
}

// New creates a Truncator that trims input to maxBytes.
// The suffix (default "...") is appended when truncation occurs.
// Returns ErrInvalidMaxBytes if maxBytes < 1.
func New(maxBytes int, opts ...Option) (*Truncator, error) {
	if maxBytes < 1 {
		return nil, ErrInvalidMaxBytes
	}
	t := &Truncator{
		maxBytes: maxBytes,
		suffix:   "...",
	}
	for _, o := range opts {
		o(t)
	}
	return t, nil
}

// Cut returns s unchanged if len(s) <= maxBytes, otherwise it trims s so
// that the result including the suffix fits within maxBytes bytes.
// If the suffix alone exceeds maxBytes the suffix is itself trimmed.
func (t *Truncator) Cut(s string) string {
	if len(s) <= t.maxBytes {
		return s
	}
	suf := t.suffix
	if len(suf) >= t.maxBytes {
		return suf[:t.maxBytes]
	}
	keep := t.maxBytes - len(suf)
	return s[:keep] + suf
}

// Truncated reports whether s would be truncated by Cut.
func (t *Truncator) Truncated(s string) bool {
	return len(s) > t.maxBytes
}
