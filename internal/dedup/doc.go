// Package dedup implements a time-windowed deduplication filter for log
// messages.
//
// A [Filter] computes a SHA-256 fingerprint of each message and suppresses
// subsequent identical messages that arrive within a configurable TTL window.
// Once the TTL elapses the fingerprint expires and the message is treated as
// new again.
//
// Usage:
//
//	f := dedup.New(30 * time.Second)
//	if !f.IsDuplicate(line.Message) {
//	    // forward the line downstream
//	}
//
// Call [Filter.Evict] periodically to release memory consumed by expired
// fingerprints.
package dedup
