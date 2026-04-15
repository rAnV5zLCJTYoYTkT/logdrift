// Package dedup provides a deduplication filter that suppresses repeated
// identical log messages within a configurable time window.
package dedup

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

// Filter tracks recently seen message fingerprints and suppresses duplicates
// within the configured TTL window.
type Filter struct {
	mu      sync.Mutex
	seen    map[string]time.Time
	ttl     time.Duration
	nowFunc func() time.Time
}

// New creates a Filter that deduplicates messages within the given TTL.
// A zero TTL disables deduplication (all messages pass through).
func New(ttl time.Duration) *Filter {
	return &Filter{
		seen:    make(map[string]time.Time),
		ttl:     ttl,
		nowFunc: time.Now,
	}
}

// IsDuplicate returns true if the given message was seen within the TTL window.
// It records the message as seen on first encounter or after expiry.
func (f *Filter) IsDuplicate(message string) bool {
	if f.ttl == 0 {
		return false
	}

	key := fingerprint(message)
	now := f.nowFunc()

	f.mu.Lock()
	defer f.mu.Unlock()

	if last, ok := f.seen[key]; ok && now.Sub(last) < f.ttl {
		return true
	}

	f.seen[key] = now
	return false
}

// Evict removes all expired entries from the internal map. It is safe to call
// periodically to bound memory usage.
func (f *Filter) Evict() {
	now := f.nowFunc()

	f.mu.Lock()
	defer f.mu.Unlock()

	for k, t := range f.seen {
		if now.Sub(t) >= f.ttl {
			delete(f.seen, k)
		}
	}
}

// Len returns the current number of tracked fingerprints.
func (f *Filter) Len() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.seen)
}

func fingerprint(s string) string {
	h := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", h)
}
