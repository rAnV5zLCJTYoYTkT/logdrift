// Package redactor provides PII and sensitive-data scrubbing for log lines
// before they are processed or emitted by logdrift.
package redactor

import (
	"regexp"
	"strings"
)

// rule pairs a compiled pattern with its replacement string.
type rule struct {
	pattern     *regexp.Regexp
	replacement string
}

// Redactor scrubs sensitive tokens from log messages.
type Redactor struct {
	rules []rule
}

// defaultRules contains built-in patterns for common PII / secrets.
var defaultRules = []struct {
	pattern     string
	replacement string
}{
	{`\b[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}\b`, "[EMAIL]"},
	{`\b(?:\d{1,3}\.){3}\d{1,3}\b`, "[IP]"},
	{`\b(?:4[0-9]{12}(?:[0-9]{3})?|5[1-5][0-9]{14}|3[47][0-9]{13})\b`, "[CARD]"},
	{`(?i)(?:password|passwd|secret|token|api[_-]?key)\s*[=:]\s*\S+`, "[REDACTED]"},
	{`\b[0-9a-fA-F]{32,}\b`, "[HASH]"},
}

// New creates a Redactor pre-loaded with built-in rules. Additional custom
// patterns (raw regex strings) and their replacements may be supplied.
func New(extra ...PatternReplacement) (*Redactor, error) {
	r := &Redactor{}
	for _, d := range defaultRules {
		compiled, err := regexp.Compile(d.pattern)
		if err != nil {
			return nil, err
		}
		r.rules = append(r.rules, rule{pattern: compiled, replacement: d.replacement})
	}
	for _, e := range extra {
		compiled, err := regexp.Compile(e.Pattern)
		if err != nil {
			return nil, err
		}
		replacement := e.Replacement
		if replacement == "" {
			replacement = "[REDACTED]"
		}
		r.rules = append(r.rules, rule{pattern: compiled, replacement: replacement})
	}
	return r, nil
}

// PatternReplacement bundles a raw regex pattern with its replacement text.
type PatternReplacement struct {
	Pattern     string
	Replacement string
}

// Scrub applies all rules to msg and returns the sanitised string.
func (r *Redactor) Scrub(msg string) string {
	for _, rl := range r.rules {
		msg = rl.pattern.ReplaceAllString(msg, rl.replacement)
	}
	return strings.TrimSpace(msg)
}
