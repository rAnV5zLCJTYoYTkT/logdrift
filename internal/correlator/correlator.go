// Package correlator groups related log entries by a shared trace or request ID.
package correlator

import (
	"errors"
	"sync"
	"time"
)

// Entry holds a single log line associated with a correlation key.
type Entry struct {
	Key       string
	Message   string
	Timestamp time.Time
}

// Group is a collection of entries sharing the same correlation key.
type Group struct {
	Key     string
	Entries []Entry
}

// Correlator buffers entries by key and flushes groups older than TTL.
type Correlator struct {
	mu      sync.Mutex
	groups  map[string]*Group
	times   map[string]time.Time
	ttl     time.Duration
	now     func() time.Time
}

// New creates a Correlator with the given TTL.
func New(ttl time.Duration) (*Correlator, error) {
	if ttl <= 0 {
		return nil, errors.New("correlator: ttl must be positive")
	}
	return &Correlator{
		groups: make(map[string]*Group),
		times:  make(map[string]time.Time),
		ttl:    ttl,
		now:    time.Now,
	}, nil
}

// Add appends an entry to the group identified by key.
func (c *Correlator) Add(e Entry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.groups[e.Key]; !ok {
		c.groups[e.Key] = &Group{Key: e.Key}
	}
	c.groups[e.Key].Entries = append(c.groups[e.Key].Entries, e)
	c.times[e.Key] = c.now()
}

// Flush returns and removes all groups whose last update exceeds the TTL.
func (c *Correlator) Flush() []Group {
	c.mu.Lock()
	defer c.mu.Unlock()
	cutoff := c.now().Add(-c.ttl)
	var out []Group
	for k, t := range c.times {
		if t.Before(cutoff) {
			out = append(out, *c.groups[k])
			delete(c.groups, k)
			delete(c.times, k)
		}
	}
	return out
}

// Len returns the number of active correlation groups.
func (c *Correlator) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.groups)
}
