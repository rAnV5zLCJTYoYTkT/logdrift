package tracer

import (
	"regexp"
	"strings"
)

// Extractor pulls a trace ID from a raw log line using a regex pattern.
type Extractor struct {
	re    *regexp.Regexp
	group string
}

// NewExtractor compiles the given pattern. The pattern must contain a named
// capture group whose name is passed as group (e.g. "trace_id").
func NewExtractor(pattern, group string) (*Extractor, error) {
	if strings.TrimSpace(pattern) == "" {
		return nil, errEmptyPattern
	}
	if strings.TrimSpace(group) == "" {
		return nil, errEmptyGroup
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &Extractor{re: re, group: group}, nil
}

// Extract returns the trace ID found in line, or empty string if none.
func (e *Extractor) Extract(line string) string {
	match := e.re.FindStringSubmatch(line)
	if match == nil {
		return ""
	}
	for i, name := range e.re.SubexpNames() {
		if name == e.group && i < len(match) {
			return match[i]
		}
	}
	return ""
}
