// Package counter provides a thread-safe frequency counter that tracks
// how often distinct keys appear within a sliding time window.
package counter

import (
	"sync"
	"time"
)

// entry holds the count and the timestamp of the last increment.
type entry struct {
	count     int64
	updatedAt time.Time
}

// Counter tracks hit counts per key, expiring entries after a configurable TTL.
type Counter struct {
	mu      sync.Mutex
	entries map[string]*entry
	ttl     time.Duration
}

// New creates a Counter whose entries expire after ttl.
// A zero ttl disables expiry.
func New(ttl time.Duration) *Counter {
	return &Counter{
		entries: make(map[string]*entry),
		ttl:     ttl,
	}
}

// Inc increments the counter for key and returns the updated count.
func (c *Counter) Inc(key string) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	e, ok := c.entries[key]
	if !ok || (c.ttl > 0 && now.Sub(e.updatedAt) > c.ttl) {
		c.entries[key] = &entry{count: 1, updatedAt: now}
		return 1
	}
	e.count++
	e.updatedAt = now
	return e.count
}

// Value returns the current count for key, or 0 if the entry has expired or
// does not exist.
func (c *Counter) Value(key string) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	e, ok := c.entries[key]
	if !ok {
		return 0
	}
	if c.ttl > 0 && time.Since(e.updatedAt) > c.ttl {
		delete(c.entries, key)
		return 0
	}
	return e.count
}

// Reset clears the count for key.
func (c *Counter) Reset(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Len returns the number of active (non-expired) keys.
func (c *Counter) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for k, e := range c.entries {
		if c.ttl > 0 && now.Sub(e.updatedAt) > c.ttl {
			delete(c.entries, k)
		}
	}
	return len(c.entries)
}
