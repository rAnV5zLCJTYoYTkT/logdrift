// Package retention implements a time-based log entry retention filter for
// logdrift. It provides a Policy type that encodes a TTL duration and a Filter
// that accepts or rejects log timestamps based on whether they fall within the
// configured retention window.
//
// Usage:
//
//	p, err := retention.ParseTTL("24h")
//	if err != nil { ... }
//	f, err := retention.New(p)
//	if err != nil { ... }
//	if f.Allow(entry.Timestamp) {
//	    // process entry
//	}
package retention
