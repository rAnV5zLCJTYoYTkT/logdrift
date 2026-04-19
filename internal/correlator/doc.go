// Package correlator groups log entries that share a common correlation key
// (such as a trace ID or request ID) into named groups. Groups are held in
// memory until they have been idle for longer than the configured TTL, at
// which point they are returned by Flush for downstream processing.
//
// Typical usage:
//
//	c, err := correlator.New(5 * time.Second)
//	c.Add(correlator.Entry{Key: traceID, Message: line, Timestamp: ts})
//	for _, g := range c.Flush() {
//		// process g.Entries
//	}
package correlator
