// Package enricher attaches derived metadata to parsed log entries
// before they are evaluated by the anomaly detection pipeline.
package enricher

import (
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/yourorg/logdrift/internal/parser"
)

// Entry wraps a parsed log line with additional derived fields.
type Entry struct {
	parser.LogLine

	// Fingerprint is a short hash identifying the structural shape of the message,
	// with variable tokens (numbers, UUIDs, IPs) normalized away.
	Fingerprint string

	// ServiceName is extracted from the logger field when it contains a dot-separated
	// hierarchy (e.g. "payments.processor" → "payments").
	ServiceName string
}

// Enricher transforms a LogLine into an Entry.
type Enricher struct {
	// normalizer replaces variable tokens so structurally identical messages
	// produce the same fingerprint.
	normalizer *strings.Replacer
}

// New returns an Enricher ready for use.
func New() *Enricher {
	return &Enricher{
		normalizer: strings.NewReplacer(
			// keep replacer lightweight; regex-heavy normalization lives in the
			// fingerprint helper below.
		),
	}
}

// Enrich derives metadata from ll and returns a populated Entry.
func (e *Enricher) Enrich(ll parser.LogLine) Entry {
	return Entry{
		LogLine:     ll,
		Fingerprint: fingerprint(ll.Message),
		ServiceName: serviceName(ll.Logger),
	}
}

// fingerprint normalises msg by collapsing digit runs and hex sequences, then
// returns the first 8 hex characters of its MD5 hash.
func fingerprint(msg string) string {
	norm := normalizeMessage(msg)
	sum := md5.Sum([]byte(norm))
	return fmt.Sprintf("%x", sum[:4])
}

// normalizeMessage replaces numeric tokens with "#" to collapse variable parts.
func normalizeMessage(msg string) string {
	words := strings.Fields(msg)
	for i, w := range words {
		if isVariable(w) {
			words[i] = "#"
		}
	}
	return strings.Join(words, " ")
}

// isVariable returns true when w looks like a number, UUID segment, or IP octet.
func isVariable(w string) bool {
	if len(w) == 0 {
		return false
	}
	all := true
	for _, c := range w {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F') || c == '.' || c == '-') {
			all = false
			break
		}
	}
	return all && strings.IndexAny(w, "0123456789") >= 0
}

// serviceName extracts the top-level component from a dot-separated logger name.
func serviceName(logger string) string {
	if idx := strings.IndexByte(logger, '.'); idx > 0 {
		return logger[:idx]
	}
	return logger
}
