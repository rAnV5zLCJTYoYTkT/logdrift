// Package limiter provides a concurrency limiter that caps the number of
// goroutines allowed to process log lines simultaneously.
package limiter

import (
	"errors"
	"sync"
)

// Limiter controls concurrent access using a semaphore.
type Limiter struct {
	sem chan struct{}
	mu  sync.Mutex
	active int
}

// New creates a Limiter that allows at most maxConcurrent simultaneous
// acquisitions. Returns an error if maxConcurrent is less than 1.
func New(maxConcurrent int) (*Limiter, error) {
	if maxConcurrent < 1 {
		return nil, errors.New("limiter: maxConcurrent must be at least 1")
	}
	return &Limiter{
		sem: make(chan struct{}, maxConcurrent),
	}, nil
}

// Acquire blocks until a slot is available, then claims it.
func (l *Limiter) Acquire() {
	l.sem <- struct{}{}
	l.mu.Lock()
	l.active++
	l.mu.Unlock()
}

// TryAcquire attempts to claim a slot without blocking.
// Returns true if successful, false if all slots are occupied.
func (l *Limiter) TryAcquire() bool {
	select {
	case l.sem <- struct{}{}:
		l.mu.Lock()
		l.active++
		l.mu.Unlock()
		return true
	default:
		return false
	}
}

// Release frees a previously acquired slot.
func (l *Limiter) Release() {
	<-l.sem
	l.mu.Lock()
	l.active--
	l.mu.Unlock()
}

// Active returns the number of currently held slots.
func (l *Limiter) Active() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.active
}

// Cap returns the maximum number of concurrent slots.
func (l *Limiter) Cap() int {
	return cap(l.sem)
}
