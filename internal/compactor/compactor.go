// Package compactor merges repeated log entries within a sliding window,
// collapsing identical fingerprints into a single entry with an occurrence count.
package compactor

import (
	"errors"
	"sync"
	"time"
)

// Entry holds a compacted log line and its occurrence metadata.
type Entry struct {
	Message   string
	Fingerprint string
	First     time.Time
	Last      time.Time
	Count     int
}

// Compactor collapses repeated messages within a TTL window.
type Compactor struct {
	mu      sync.Mutex
	ttl     time.Duration
	buckets map[string]*Entry
	now     func() time.Time
}

// New creates a Compactor with the given TTL. Returns an error if ttl <= 0.
func New(ttl time.Duration) (*Compactor, error) {
	if ttl <= 0 {
		return nil, errors.New("compactor: ttl must be positive")
	}
	return &Compactor{
		ttl:     ttl,
		buckets: make(map[string]*Entry),
		now:     time.Now,
	}, nil
}

// Add records a message with its fingerprint. Returns the current Entry and
// whether this is a new (non-duplicate) occurrence that should be forwarded.
func (c *Compactor) Add(fingerprint, message string) (Entry, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.now()
	c.evict(now)

	if e, ok := c.buckets[fingerprint]; ok {
		e.Count++
		e.Last = now
		return *e, false
	}

	e := &Entry{
		Message:     message,
		Fingerprint: fingerprint,
		First:       now,
		Last:        now,
		Count:       1,
	}
	c.buckets[fingerprint] = e
	return *e, true
}

// Flush returns all current entries and clears the compactor.
func (c *Compactor) Flush() []Entry {
	c.mu.Lock()
	defer c.mu.Unlock()

	out := make([]Entry, 0, len(c.buckets))
	for _, e := range c.buckets {
		out = append(out, *e)
	}
	c.buckets = make(map[string]*Entry)
	return out
}

// evict removes entries whose last-seen time is beyond the TTL.
func (c *Compactor) evict(now time.Time) {
	for k, e := range c.buckets {
		if now.Sub(e.Last) > c.ttl {
			delete(c.buckets, k)
		}
	}
}
