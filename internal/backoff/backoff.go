// Package backoff provides an exponential back-off strategy with jitter
// for retrying operations that may transiently fail (e.g. alert delivery,
// remote log sinks).
package backoff

import (
	"errors"
	"math"
	"math/rand"
	"sync"
	"time"
)

// Backoff holds per-key retry state.
type Backoff struct {
	mu       sync.Mutex
	base     time.Duration
	max      time.Duration
	factor   float64
	attempts map[string]int
}

// New returns a Backoff with the given base delay, maximum delay and
// multiplicative factor. Returns an error when any argument is invalid.
func New(base, max time.Duration, factor float64) (*Backoff, error) {
	if base <= 0 {
		return nil, errors.New("backoff: base delay must be positive")
	}
	if max < base {
		return nil, errors.New("backoff: max delay must be >= base delay")
	}
	if factor < 1 {
		return nil, errors.New("backoff: factor must be >= 1")
	}
	return &Backoff{
		base:     base,
		max:      max,
		factor:   factor,
		attempts: make(map[string]int),
	}, nil
}

// Next returns the next delay for the given key and increments its attempt
// counter. A small random jitter (up to 10 % of the computed delay) is added
// to avoid thundering-herd scenarios.
func (b *Backoff) Next(key string) time.Duration {
	b.mu.Lock()
	defer b.mu.Unlock()

	n := b.attempts[key]
	b.attempts[key] = n + 1

	delay := float64(b.base) * math.Pow(b.factor, float64(n))
	if delay > float64(b.max) {
		delay = float64(b.max)
	}
	jitter := delay * 0.1 * rand.Float64() //nolint:gosec
	return time.Duration(delay + jitter)
}

// Reset clears the attempt counter for the given key.
func (b *Backoff) Reset(key string) {
	b.mu.Lock()
	delete(b.attempts, key)
	b.mu.Unlock()
}

// Attempts returns the current attempt count for the given key.
func (b *Backoff) Attempts(key string) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.attempts[key]
}
