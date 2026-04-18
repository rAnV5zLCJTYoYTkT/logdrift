// Package labeler attaches key-value labels to log entries based on
// field matching rules defined at construction time.
package labeler

import (
	"errors"
	"strings"
)

// Rule maps a field value substring to a set of labels.
type Rule struct {
	Field    string
	Contains string
	Labels   map[string]string
}

// Labeler applies label rules to an entry's fields.
type Labeler struct {
	rules []Rule
}

// New creates a Labeler from the provided rules.
// Returns an error if no rules are supplied or any rule is misconfigured.
func New(rules []Rule) (*Labeler, error) {
	if len(rules) == 0 {
		return nil, errors.New("labeler: at least one rule is required")
	}
	for i, r := range rules {
		if strings.TrimSpace(r.Field) == "" {
			return nil, fmt.Errorf("labeler: rule %d has empty field", i)
		}
		if len(r.Labels) == 0 {
			return nil, fmt.Errorf("labeler: rule %d has no labels", i)
		}
	}
	return &Labeler{rules: rules}, nil
}

// Apply evaluates each rule against fields and returns the merged label set.
// fields is a map of field name to value (e.g. {"level": "error", "service": "api"}).
func (l *Labeler) Apply(fields map[string]string) map[string]string {
	out := make(map[string]string)
	for _, r := range l.rules {
		v, ok := fields[r.Field]
		if !ok {
			continue
		}
		if r.Contains == "" || strings.Contains(strings.ToLower(v), strings.ToLower(r.Contains)) {
			for k, lv := range r.Labels {
				out[k] = lv
			}
		}
	}
	return out
}
