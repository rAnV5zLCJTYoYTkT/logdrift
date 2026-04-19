// Package router dispatches log entries to named output sinks based on
// configurable routing rules. Each rule matches on severity level and an
// optional substring pattern; the first matching rule wins.
package router

import (
	"errors"
	"fmt"
	"strings"

	"github.com/user/logdrift/internal/severity"
)

// Entry is the minimal log record the router operates on.
type Entry struct {
	Level   severity.Level
	Message string
	Service string
}

// Rule describes a single routing condition.
type Rule struct {
	// Sink is the destination name (e.g. "alerts", "audit").
	Sink string
	// MinLevel is the minimum severity that triggers this rule.
	MinLevel severity.Level
	// Contains, when non-empty, requires the message to contain the substring.
	Contains string
}

// Router routes entries to sink names.
type Router struct {
	rules []Rule
}

// New returns a Router configured with the provided rules.
// At least one rule must be supplied.
func New(rules []Rule) (*Router, error) {
	if len(rules) == 0 {
		return nil, errors.New("router: at least one rule is required")
	}
	for i, r := range rules {
		if strings.TrimSpace(r.Sink) == "" {
			return nil, fmt.Errorf("router: rule %d has an empty sink name", i)
		}
	}
	return &Router{rules: rules}, nil
}

// Route returns the sink name for the given entry.
// It returns the first rule whose MinLevel and Contains constraints are
// satisfied. If no rule matches, an empty string and false are returned.
func (r *Router) Route(e Entry) (string, bool) {
	for _, rule := range r.rules {
		if !severity.AtLeast(e.Level, rule.MinLevel) {
			continue
		}
		if rule.Contains != "" && !strings.Contains(e.Message, rule.Contains) {
			continue
		}
		return rule.Sink, true
	}
	return "", false
}
