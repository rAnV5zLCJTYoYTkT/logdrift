// Package throttle provides rate-limiting for alert emissions,
// preventing alert storms when a metric stays anomalous for an
// extended period.
package throttle

import (
	"sync"
	"time"
)

// Throttle limits how frequently an action may be taken for a given key.
// It is safe for concurrent use.
type Throttle struct {
	mu       sync.Mutex
	last     map[string]time.Time
	interval time.Duration
	now      func() time.Time
}

// New creates a Throttle that allows at most one action per key within
// the given interval. An interval of zero disables throttling (every
// call to Allow returns true).
func New(interval time.Duration) *Throttle {
	return &Throttle{
		last:     make(map[string]time.Time),
		interval: interval,
		now:      time.Now,
	}
}

// Allow reports whether the action identified by key is permitted at
// the current time. If allowed, the internal timestamp for key is
// updated so that subsequent calls within the interval are denied.
func (t *Throttle) Allow(key string) bool {
	if t.interval <= 0 {
		return true
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	if last, ok := t.last[key]; ok {
		if now.Sub(last) < t.interval {
			return false
		}
	}
	t.last[key] = now
	return true
}

// Reset removes the recorded timestamp for key, allowing the next
// call to Allow for that key to succeed immediately.
func (t *Throttle) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.last, key)
}

// Len returns the number of keys currently tracked by the throttle.
func (t *Throttle) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.last)
}
