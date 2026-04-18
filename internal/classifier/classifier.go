// Package classifier assigns a category label to a log entry based on
// its severity, message fingerprint, and optional keyword rules.
package classifier

import (
	"errors"
	"strings"
)

// Category represents a broad log entry classification.
type Category string

const (
	CategoryNoise    Category = "noise"
	CategoryInfo     Category = "info"
	CategoryWarning  Category = "warning"
	CategoryError    Category = "error"
	CategoryCritical Category = "critical"
)

// Rule maps a set of keywords to a Category. The first matching rule wins.
type Rule struct {
	Keywords []string
	Category Category
}

// Classifier categorises log messages.
type Classifier struct {
	rules []Rule
}

// New returns a Classifier with the provided rules.
// Rules are evaluated in order; the first keyword match wins.
func New(rules []Rule) (*Classifier, error) {
	if len(rules) == 0 {
		return nil, errors.New("classifier: at least one rule is required")
	}
	for i, r := range rules {
		if len(r.Keywords) == 0 {
			return nil, errors.New("classifier: rule has no keywords")
		}
		if r.Category == "" {
			return nil, errors.New("classifier: rule missing category")
		}
		_ = i
	}
	return &Classifier{rules: rules}, nil
}

// Classify returns the Category for the given message.
// If no rule matches, CategoryNoise is returned.
func (c *Classifier) Classify(message string) Category {
	lower := strings.ToLower(message)
	for _, r := range c.rules {
		for _, kw := range r.Keywords {
			if strings.Contains(lower, strings.ToLower(kw)) {
				return r.Category
			}
		}
	}
	return CategoryNoise
}
