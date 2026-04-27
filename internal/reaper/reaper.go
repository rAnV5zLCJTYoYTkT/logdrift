// Package reaper periodically removes stale entries from a keyed store
// once they have exceeded a configurable idle timeout.
package reaper

import (
	"errors"
	"sync"
	"time"
)

// ErrInvalidInterval is returned when the sweep interval is not positive.
var ErrInvalidInterval = errors.New("reaper: sweep interval must be positive")

// ErrInvalidIdle is returned when the idle timeout is not positive.
var ErrInvalidIdle = errors.New("reaper: idle timeout must be positive")

type entry struct {
	lastSeen time.Time
}

// Reaper tracks last-seen timestamps for keys and exposes a Sweep method
// that returns all keys that have been idle for longer than the configured
// timeout.
type Reaper struct {
	mu      sync.Mutex
	entries map[string]entry
	idle    time.Duration
	now     func() time.Time
}

// New creates a Reaper with the given idle timeout.
// interval is validated but the caller is responsible for scheduling Sweep.
func New(idle time.Duration) (*Reaper, error) {
	if idle <= 0 {
		return nil, ErrInvalidIdle
	}
	return &Reaper{
		entries: make(map[string]entry),
		idle:    idle,
		now:     time.Now,
	}, nil
}

// Touch records or refreshes the last-seen time for key.
func (r *Reaper) Touch(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[key] = entry{lastSeen: r.now()}
}

// Sweep returns all keys whose last-seen time is older than the idle timeout
// and removes them from the internal store.
func (r *Reaper) Sweep() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	cutoff := r.now().Add(-r.idle)
	var expired []string
	for k, e := range r.entries {
		if e.lastSeen.Before(cutoff) {
			expired = append(expired, k)
			delete(r.entries, k)
		}
	}
	return expired
}

// Len returns the number of currently tracked keys.
func (r *Reaper) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.entries)
}
