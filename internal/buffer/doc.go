// Package buffer implements a fixed-capacity, thread-safe ring buffer for
// retaining the most recent N log lines in memory.
//
// It is useful for providing context around anomalous events — when an alert
// fires, the caller can call Snapshot to retrieve the preceding lines for
// inclusion in the alert payload or report output.
//
// Usage:
//
//	buf, err := buffer.New(100)
//	if err != nil { ... }
//	buf.Add(line)
//	lines := buf.Snapshot()
package buffer
