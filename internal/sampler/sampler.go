// Package sampler provides rate-limiting and sampling utilities for
// controlling how frequently alerts are emitted for repeated anomalies.
package sampler

import (
	"sync"
	"time"
)

// Sampler tracks the last emission time per key and suppresses duplicates
// within a configurable cooldown window.
type Sampler struct {
	mu       sync.Mutex
	cooldown time.Duration
	lastSeen map[string]time.Time
}

// New creates a Sampler with the given cooldown duration.
// A cooldown of zero means every event is allowed through.
func New(cooldown time.Duration) *Sampler {
	return &Sampler{
		cooldown: cooldown,
		lastSeen: make(map[string]time.Time),
	}
}

// Allow returns true if the key has not been seen within the cooldown window.
// It updates the last-seen timestamp when it returns true.
func (s *Sampler) Allow(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	if last, ok := s.lastSeen[key]; ok {
		if now.Sub(last) < s.cooldown {
			return false
		}
	}
	s.lastSeen[key] = now
	return true
}

// Evict removes keys whose last-seen time is older than the cooldown, freeing
// memory for long-running processes. Call periodically if key cardinality is
// high.
func (s *Sampler) Evict() {
	s.mu.Lock()
	defer s.mu.Unlock()

	cutoff := time.Now().Add(-s.cooldown)
	for k, t := range s.lastSeen {
		if t.Before(cutoff) {
			delete(s.lastSeen, k)
		}
	}
}

// Len returns the number of keys currently tracked.
func (s *Sampler) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.lastSeen)
}
