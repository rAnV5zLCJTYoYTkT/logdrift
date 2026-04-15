// Package sanitizer provides utilities for normalising raw log lines
// before they are processed by the pipeline. It trims whitespace,
// enforces a maximum line length, and strips non-printable characters.
package sanitizer

import (
	"strings"
	"unicode"
)

const defaultMaxLen = 4096

// Sanitizer cleans raw log lines.
type Sanitizer struct {
	maxLen int
}

// Option is a functional option for Sanitizer.
type Option func(*Sanitizer)

// WithMaxLen sets the maximum byte length a line may have after trimming.
// Lines longer than this are truncated. A value of 0 restores the default.
func WithMaxLen(n int) Option {
	return func(s *Sanitizer) {
		if n > 0 {
			s.maxLen = n
		}
	}
}

// New returns a Sanitizer configured with the supplied options.
func New(opts ...Option) *Sanitizer {
	s := &Sanitizer{maxLen: defaultMaxLen}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Clean returns a sanitized copy of line.
// It trims leading/trailing whitespace, removes non-printable runes, and
// truncates the result to the configured maximum length.
// An empty string is returned when the input is blank after cleaning.
func (s *Sanitizer) Clean(line string) string {
	line = strings.TrimSpace(line)
	if line == "" {
		return ""
	}

	// Strip non-printable, non-space control characters.
	var b strings.Builder
	b.Grow(len(line))
	for _, r := range line {
		if unicode.IsPrint(r) || r == '\t' {
			b.WriteRune(r)
		}
	}

	result := b.String()
	if len(result) > s.maxLen {
		// Truncate at a rune boundary.
		result = string([]rune(result)[:s.maxLen])
	}
	return result
}
