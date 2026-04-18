// Package normalizer collapses dynamic tokens in log messages into
// canonical placeholders, making messages comparable across requests.
package normalizer

import (
	"regexp"
	"strings"
)

var (
	reHex    = regexp.MustCompile(`\b[0-9a-fA-F]{8,}\b`)
	reUUID   = regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
	reIP     = regexp.MustCompile(`\b\d{1,3}(?:\.\d{1,3}){3}\b`)
	reNum    = regexp.MustCompile(`\b\d+\b`)
	reQuoted = regexp.MustCompile(`"[^"]{1,120}"`)
)

// Option configures a Normalizer.
type Option func(*Normalizer)

// WithQuotedStrings enables replacement of quoted string literals.
func WithQuotedStrings() Option {
	return func(n *Normalizer) { n.quotedStrings = true }
}

// Normalizer replaces variable tokens with stable placeholders.
type Normalizer struct {
	quotedStrings bool
}

// New returns a Normalizer configured with opts.
func New(opts ...Option) *Normalizer {
	n := &Normalizer{}
	for _, o := range opts {
		o(n)
	}
	return n
}

// Normalize returns a canonical form of msg.
func (n *Normalizer) Normalize(msg string) string {
	s := reUUID.ReplaceAllString(msg, "<uuid>")
	s = reIP.ReplaceAllString(s, "<ip>")
	s = reHex.ReplaceAllString(s, "<hex>")
	if n.quotedStrings {
		s = reQuoted.ReplaceAllString(s, "\"<str>\"")
	}
	s = reNum.ReplaceAllString(s, "<N>")
	s = strings.TrimSpace(s)
	return s
}
