// Package masker provides field-level masking for log entries,
// replacing sensitive field values with a fixed placeholder.
package masker

import (
	"errors"
	"strings"
)

const defaultPlaceholder = "[MASKED]"

// Masker replaces the values of named fields in a key=value log message.
type Masker struct {
	fields      map[string]struct{}
	placeholder string
}

// Option configures a Masker.
type Option func(*Masker)

// WithPlaceholder overrides the default replacement string.
func WithPlaceholder(p string) Option {
	return func(m *Masker) {
		if p != "" {
			m.placeholder = p
		}
	}
}

// New creates a Masker that will redact the given field names.
// At least one field name must be provided.
func New(fields []string, opts ...Option) (*Masker, error) {
	if len(fields) == 0 {
		return nil, errors.New("masker: at least one field name required")
	}
	m := &Masker{
		fields:      make(map[string]struct{}, len(fields)),
		placeholder: defaultPlaceholder,
	}
	for _, f := range fields {
		if f == "" {
			return nil, errors.New("masker: field name must not be empty")
		}
		m.fields[strings.ToLower(f)] = struct{}{}
	}
	for _, o := range opts {
		o(m)
	}
	return m, nil
}

// Mask scans a log line for key=value or key="value" pairs and replaces
// the value of any matching field with the placeholder.
func (m *Masker) Mask(line string) string {
	var sb strings.Builder
	sb.Grow(len(line))
	tokens := strings.Fields(line)
	for i, tok := range tokens {
		if i > 0 {
			sb.WriteByte(' ')
		}
		eq := strings.IndexByte(tok, '=')
		if eq <= 0 {
			sb.WriteString(tok)
			continue
		}
		key := tok[:eq]
		if _, sensitive := m.fields[strings.ToLower(key)]; sensitive {
			sb.WriteString(key)
			sb.WriteByte('=')
			sb.WriteString(m.placeholder)
		} else {
			sb.WriteString(tok)
		}
	}
	return sb.String()
}
