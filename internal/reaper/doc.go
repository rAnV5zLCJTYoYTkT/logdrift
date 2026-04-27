// Package reaper provides idle-key expiry for logdrift's in-memory stores.
//
// A Reaper tracks the last time each key was "touched" and, when Sweep is
// called, returns and removes every key that has been idle for longer than
// the configured timeout.  The caller decides when to call Sweep — typically
// via a time.Ticker in a background goroutine.
//
// Example:
//
//	r, err := reaper.New(5 * time.Minute)
//	if err != nil { ... }
//
//	// Record activity.
//	r.Touch("service-a")
//
//	// Periodically evict stale keys.
//	for _, key := range r.Sweep() {
//		log.Printf("evicted idle key: %s", key)
//	}
package reaper
