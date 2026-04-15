// Package severity provides utilities for ranking and comparing log severity levels.
package severity

import (
	"fmt"
	"strings"
)

// Level represents a numeric severity rank.
type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
	Fatal
)

var levelNames = map[Level]string{
	Debug: "DEBUG",
	Info:  "INFO",
	Warn:  "WARN",
	Error: "ERROR",
	Fatal: "FATAL",
}

var nameToLevel = map[string]Level{
	"DEBUG": Debug,
	"INFO":  Info,
	"WARN":  Warn,
	"WARNING": Warn,
	"ERROR": Error,
	"FATAL": Fatal,
}

// String returns the canonical name for the level.
func (l Level) String() string {
	if name, ok := levelNames[l]; ok {
		return name
	}
	return fmt.Sprintf("LEVEL(%d)", int(l))
}

// Parse converts a string to a Level. Returns an error if the string is
// not a recognised severity name.
func Parse(s string) (Level, error) {
	norm := strings.ToUpper(strings.TrimSpace(s))
	if l, ok := nameToLevel[norm]; ok {
		return l, nil
	}
	return Debug, fmt.Errorf("severity: unknown level %q", s)
}

// AtLeast reports whether l is at least as severe as min.
func AtLeast(l, min Level) bool {
	return l >= min
}

// Rank returns the numeric rank of the level, useful for external comparisons.
func Rank(l Level) int {
	return int(l)
}
