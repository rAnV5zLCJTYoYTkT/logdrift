// Package splitter splits a raw log line into labelled fields
// using a configurable delimiter and field mapping.
package splitter

import (
	"errors"
	"strings"
)

// ErrEmptyDelimiter is returned when an empty delimiter is provided.
var ErrEmptyDelimiter = errors.New("splitter: delimiter must not be empty")

// ErrNoFields is returned when no field names are provided.
var ErrNoFields = errors.New("splitter: at least one field name is required")

// Splitter splits raw lines into named fields.
type Splitter struct {
	delimiter string
	fields    []string
}

// New creates a Splitter that splits on delimiter and maps columns to fields.
func New(delimiter string, fields []string) (*Splitter, error) {
	if delimiter == "" {
		return nil, ErrEmptyDelimiter
	}
	if len(fields) == 0 {
		return nil, ErrNoFields
	}
	return &Splitter{delimiter: delimiter, fields: fields}, nil
}

// Split parses a raw line and returns a map of field name to value.
// Extra columns beyond the defined fields are ignored.
// Missing columns result in empty string values.
func (s *Splitter) Split(line string) map[string]string {
	parts := strings.SplitN(line, s.delimiter, len(s.fields))
	out := make(map[string]string, len(s.fields))
	for i, name := range s.fields {
		if i < len(parts) {
			out[name] = strings.TrimSpace(parts[i])
		} else {
			out[name] = ""
		}
	}
	return out
}
