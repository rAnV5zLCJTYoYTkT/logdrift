// Package filter provides log level and pattern-based filtering
// for the logdrift pipeline, allowing users to restrict analysis
// to specific log levels or message patterns.
package filter

import (
	"regexp"

	"github.com/user/logdrift/internal/parser"
)

// Level represents a minimum log level threshold for filtering.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

var levelNames = map[string]Level{
	"DEBUG": LevelDebug,
	"INFO":  LevelInfo,
	"WARN":  LevelWarn,
	"ERROR": LevelError,
}

// ParseLevel converts a string level name to a Level constant.
// Returns LevelDebug and false if the name is unrecognised.
func ParseLevel(name string) (Level, bool) {
	l, ok := levelNames[name]
	return l, ok
}

// Filter decides whether a parsed log line should be processed
// by the pipeline.
type Filter struct {
	minLevel Level
	pattern  *regexp.Regexp
}

// Options configures a Filter.
type Options struct {
	// MinLevel is the minimum log level to allow through (inclusive).
	MinLevel Level
	// Pattern, when non-nil, restricts lines to those whose message
	// matches the compiled regular expression.
	Pattern *regexp.Regexp
}

// New returns a Filter configured with the supplied Options.
func New(opts Options) *Filter {
	return &Filter{
		minLevel: opts.MinLevel,
		pattern:  opts.Pattern,
	}
}

// Allow returns true when line passes all configured filter criteria.
func (f *Filter) Allow(line parser.LogLine) bool {
	lineLevel, ok := levelNames[line.Level]
	if !ok {
		return false
	}
	if lineLevel < f.minLevel {
		return false
	}
	if f.pattern != nil && !f.pattern.MatchString(line.Message) {
		return false
	}
	return true
}
