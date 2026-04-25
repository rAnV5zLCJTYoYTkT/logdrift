// Package trimmer provides a log entry field trimmer that removes leading
// and trailing whitespace from selected fields and optionally collapses
// internal runs of whitespace to a single space.
package trimmer

import (
	"errors"
	"strings"
)

// ErrNoFields is returned when New is called with an empty field list.
var ErrNoFields = errors.New("trimmer: at least one field name is required")

// ErrEmptyField is returned when an empty string is present in the field list.
var ErrEmptyField = errors.New("trimmer: field name must not be empty")

// Trimmer removes extraneous whitespace from named fields in a key-value map.
type Trimmer struct {
	fields  map[string]struct{}
	collapse bool
}

// Option is a functional option for Trimmer.
type Option func(*Trimmer)

// WithCollapse enables collapsing of internal whitespace runs to a single space.
func WithCollapse() Option {
	return func(t *Trimmer) {
		t.collapse = true
	}
}

// New creates a Trimmer that operates on the given field names.
func New(fields []string, opts ...Option) (*Trimmer, error) {
	if len(fields) == 0 {
		return nil, ErrNoFields
	}
	m := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		if f == "" {
			return nil, ErrEmptyField
		}
		m[f] = struct{}{}
	}
	t := &Trimmer{fields: m}
	for _, o := range opts {
		o(t)
	}
	return t, nil
}

// Apply returns a new map with whitespace trimmed from the configured fields.
// Fields not listed during construction are copied unchanged.
func (t *Trimmer) Apply(entry map[string]string) map[string]string {
	out := make(map[string]string, len(entry))
	for k, v := range entry {
		if _, ok := t.fields[k]; ok {
			v = strings.TrimSpace(v)
			if t.collapse {
				v = strings.Join(strings.Fields(v), " ")
			}
		}
		out[k] = v
	}
	return out
}
