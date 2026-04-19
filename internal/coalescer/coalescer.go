// Package coalescer merges repeated log entries within a time window into a
// single representative entry with an occurrence count attached.
package coalescer

import (
	"errors"
	"sync"
	"time"
)

// Entry holds a deduplicated log message and the number of times it was seen.
type Entry struct {
	Message string
	Count   int
	First   time.Time
	Last    time.Time
}

// Coalescer groups identical messages that arrive within a sliding window.
type Coalescer struct {
	mu      sync.Mutex
	window  time.Duration
	buckets map[string]*Entry
}

// New returns a Coalescer that merges messages seen within window.
// window must be positive.
func New(window time.Duration) (*Coalescer, error) {
	if window <= 0 {
		return nil, errors.New("coalescer: window must be positive")
	}
	return &Coalescer{
		window:  window,
		buckets: make(map[string]*Entry),
	}, nil
}

// Add records a message occurrence. It returns the current Entry for that
// message so callers can decide whether to forward or suppress it.
func (c *Coalescer) Add(msg string, at time.Time) Entry {
	c.mu.Lock()
	defer c.mu.Unlock()

	if e, ok := c.buckets[msg]; ok && at.Sub(e.Last) <= c.window {
		e.Count++
		e.Last = at
		return *e
	}

	e := &Entry{Message: msg, Count: 1, First: at, Last: at}
	c.buckets[msg] = e
	return *e
}

// Flush returns all coalesced entries and resets internal state.
func (c *Coalescer) Flush() []Entry {
	c.mu.Lock()
	defer c.mu.Unlock()

	out := make([]Entry, 0, len(c.buckets))
	for _, e := range c.buckets {
		out = append(out, *e)
	}
	c.buckets = make(map[string]*Entry)
	return out
}

// Len returns the number of distinct messages currently tracked.
func (c *Coalescer) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.buckets)
}
