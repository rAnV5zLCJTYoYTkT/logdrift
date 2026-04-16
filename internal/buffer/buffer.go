// Package buffer provides a fixed-size ring buffer for storing recent log lines.
package buffer

import (
	"errors"
	"sync"
)

// ErrInvalidCapacity is returned when capacity is less than 1.
var ErrInvalidCapacity = errors.New("buffer: capacity must be at least 1")

// RingBuffer holds the last N log lines in a thread-safe ring buffer.
type RingBuffer struct {
	mu   sync.Mutex
	data []string
	head int
	size int
	cap  int
}

// New creates a RingBuffer with the given capacity.
func New(capacity int) (*RingBuffer, error) {
	if capacity < 1 {
		return nil, ErrInvalidCapacity
	}
	return &RingBuffer{
		data: make([]string, capacity),
		cap:  capacity,
	}, nil
}

// Add inserts a line into the buffer, overwriting the oldest entry when full.
func (r *RingBuffer) Add(line string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[r.head%r.cap] = line
	r.head++
	if r.size < r.cap {
		r.size++
	}
}

// Snapshot returns a copy of buffered lines in insertion order (oldest first).
func (r *RingBuffer) Snapshot() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]string, r.size)
	start := 0
	if r.size == r.cap {
		start = r.head % r.cap
	}
	for i := 0; i < r.size; i++ {
		out[i] = r.data[(start+i)%r.cap]
	}
	return out
}

// Len returns the number of lines currently stored.
func (r *RingBuffer) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.size
}

// Cap returns the maximum capacity of the buffer.
func (r *RingBuffer) Cap() int { return r.cap }
