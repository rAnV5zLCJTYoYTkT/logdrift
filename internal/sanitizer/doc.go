// Package sanitizer provides a Sanitizer type that normalises raw log lines
// before they enter the logdrift processing pipeline.
//
// Sanitization steps applied in order:
//
//  1. Trim leading and trailing whitespace (including newlines).
//  2. Remove non-printable Unicode control characters (tab is preserved).
//  3. Truncate the result to a configurable maximum length (default 4096 bytes).
//
// Usage:
//
//	s := sanitizer.New(sanitizer.WithMaxLen(2048))
//	clean := s.Clean(rawLine)
//	if clean == "" {
//	    continue // skip blank lines
//	}
package sanitizer
