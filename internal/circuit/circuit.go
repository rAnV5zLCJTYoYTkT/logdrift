// Package circuit implements a simple circuit-breaker that opens after a
// configurable number of consecutive failures and resets after a cooldown.
package circuit

import (
	"errors"
	"sync"
	"time"
)

// State represents the current state of the circuit breaker.
type State int

const (
	StateClosed   State = iota // normal operation
	StateOpen                  // failures exceeded threshold; requests rejected
	StateHalfOpen              // cooldown elapsed; one probe allowed
)

// ErrOpen is returned when the circuit is open and the call is rejected.
var ErrOpen = errors.New("circuit: breaker is open")

// Breaker is a thread-safe circuit breaker.
type Breaker struct {
	mu           sync.Mutex
	threshold    int
	cooldown     time.Duration
	consecutive  int
	state        State
	openedAt     time.Time
}

// New creates a Breaker that opens after threshold consecutive failures and
// attempts a half-open probe after cooldown.
func New(threshold int, cooldown time.Duration) (*Breaker, error) {
	if threshold <= 0 {
		return nil, errors.New("circuit: threshold must be greater than zero")
	}
	if cooldown <= 0 {
		return nil, errors.New("circuit: cooldown must be greater than zero")
	}
	return &Breaker{
		threshold: threshold,
		cooldown:  cooldown,
		state:     StateClosed,
	}, nil
}

// Allow reports whether the caller may proceed. It transitions the breaker to
// HalfOpen once the cooldown has elapsed after opening.
func (b *Breaker) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(b.openedAt) >= b.cooldown {
			b.state = StateHalfOpen
			return true
		}
		return false
	case StateHalfOpen:
		return true
	}
	return false
}

// RecordSuccess resets the breaker to closed.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.consecutive = 0
	b.state = StateClosed
}

// RecordFailure increments the failure counter and may open the breaker.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.consecutive++
	if b.state == StateHalfOpen || b.consecutive >= b.threshold {
		b.state = StateOpen
		b.openedAt = time.Now()
	}
}

// State returns the current breaker state.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
