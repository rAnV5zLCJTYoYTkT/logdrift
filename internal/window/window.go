// Package window provides a sliding time-window counter for grouping
// log events into discrete, evictable buckets.
package window

import (
	"errors"
	"sync"
	"time"
)

// Bucket holds the count of events that arrived within a single time slot.
type Bucket struct {
	Timestamp time.Time
	Count     int64
}

// Slider maintains a fixed-size ring of time buckets and evicts stale ones
// on every call to Add or Buckets.
type Slider struct {
	mu       sync.Mutex
	size     int
	duration time.Duration
	buckets  []Bucket
}

// New creates a Slider with the given number of buckets spread over duration.
// Returns an error when size < 1 or duration <= 0.
func New(size int, duration time.Duration) (*Slider, error) {
	if size < 1 {
		return nil, errors.New("window: size must be at least 1")
	}
	if duration <= 0 {
		return nil, errors.New("window: duration must be positive")
	}
	return &Slider{
		size:     size,
		duration: duration,
		buckets:  make([]Bucket, 0, size),
	}, nil
}

// Add increments the current bucket's counter by delta and evicts expired buckets.
func (s *Slider) Add(delta int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.evict()
	now := time.Now()
	if len(s.buckets) == 0 {
		s.buckets = append(s.buckets, Bucket{Timestamp: now, Count: delta})
		return
	}
	last := &s.buckets[len(s.buckets)-1]
	slotWidth := s.duration / time.Duration(s.size)
	if now.Sub(last.Timestamp) < slotWidth {
		last.Count += delta
	} else {
		if len(s.buckets) >= s.size {
			s.buckets = s.buckets[1:]
		}
		s.buckets = append(s.buckets, Bucket{Timestamp: now, Count: delta})
	}
}

// Total returns the sum of all counts across active buckets.
func (s *Slider) Total() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.evict()
	var total int64
	for _, b := range s.buckets {
		total += b.Count
	}
	return total
}

// Buckets returns a snapshot of the current active buckets.
func (s *Slider) Buckets() []Bucket {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.evict()
	snap := make([]Bucket, len(s.buckets))
	copy(snap, s.buckets)
	return snap
}

// evict removes buckets older than the configured duration. Must be called with mu held.
func (s *Slider) evict() {
	cutoff := time.Now().Add(-s.duration)
	i := 0
	for i < len(s.buckets) && s.buckets[i].Timestamp.Before(cutoff) {
		i++
	}
	s.buckets = s.buckets[i:]
}
