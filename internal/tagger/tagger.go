// Package tagger assigns categorical tags to log lines based on
// configurable keyword rules, enabling downstream grouping and filtering.
package tagger

import (
	"strings"

	"github.com/user/logdrift/internal/parser"
)

// Rule maps a tag name to a set of keywords that trigger it.
type Rule struct {
	Tag      string
	Keywords []string
}

// Tagger holds compiled tagging rules.
type Tagger struct {
	rules []Rule
}

// New creates a Tagger from the provided rules.
// Returns an error if any rule has an empty tag or no keywords.
func New(rules []Rule) (*Tagger, error) {
	for _, r := range rules {
		if strings.TrimSpace(r.Tag) == "" {
			return nil, ErrEmptyTag
		}
		if len(r.Keywords) == 0 {
			return nil, ErrNoKeywords
		}
	}
	return &Tagger{rules: rules}, nil
}

// Tag returns all tags whose keywords appear in the log line's message.
// Matching is case-insensitive.
func (t *Tagger) Tag(line parser.LogLine) []string {
	lower := strings.ToLower(line.Message)
	seen := make(map[string]struct{})
	var tags []string
	for _, r := range t.rules {
		for _, kw := range r.Keywords {
			if strings.Contains(lower, strings.ToLower(kw)) {
				if _, ok := seen[r.Tag]; !ok {
					seen[r.Tag] = struct{}{}
					tags = append(tags, r.Tag)
				}
				break
			}
		}
	}
	return tags
}
