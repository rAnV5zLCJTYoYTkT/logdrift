// Package eviction provides a time-based cache eviction policy that removes
// entries exceeding a configurable TTL, suitable for bounding memory use in
// long-running log processing pipelines.
package eviction

import (
	"errors"
	"sync"
	"time"
)

// ErrInvalidTTL is returned when a zero or negative TTL is provided.
var ErrInvalidTTL = errors.New("eviction: TTL must be greater than zero")

type entry struct {
	addedAt time.Time
}

// Cache tracks when keys were first seen and evicts them after the TTL.
type Cache struct {
	mu      sync.Mutex
	entries map[string]entry
	ttl     time.Duration
	now     func() time.Time
}

// New creates a Cache with the given TTL.
func New(ttl time.Duration) (*Cache, error) {
	if ttl <= 0 {
		return nil, ErrInvalidTTL
	}
	return &Cache{
		entries: make(map[string]entry),
		ttl:     ttl,
		now:     time.Now,
	}, nil
}

// Track records a key and returns true if the key is new (not yet evicted).
// If the key exists and has not expired it returns false. If expired, the key
// is reset and true is returned.
func (c *Cache) Track(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.now()
	if e, ok := c.entries[key]; ok {
		if now.Sub(e.addedAt) < c.ttl {
			return false
		}
	}
	c.entries[key] = entry{addedAt: now}
	return true
}

// Evict removes all entries whose TTL has elapsed.
func (c *Cache) Evict() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.now()
	removed := 0
	for k, e := range c.entries {
		if now.Sub(e.addedAt) >= c.ttl {
			delete(c.entries, k)
			removed++
		}
	}
	return removed
}

// Len returns the current number of tracked keys.
func (c *Cache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.entries)
}
