// Package headroom tracks how close a rolling metric is to a configured
// ceiling and reports the remaining capacity as a normalised [0, 1] fraction.
package headroom

import (
	"errors"
	"sync"
)

// ErrInvalidCeiling is returned when the ceiling is not a positive number.
var ErrInvalidCeiling = errors.New("headroom: ceiling must be greater than zero")

// ErrUnknownKey is returned when a key has no recorded observations.
var ErrUnknownKey = errors.New("headroom: no observations for key")

// Tracker keeps per-key peak observations and computes remaining headroom
// relative to a fixed ceiling.
type Tracker struct {
	mu      sync.Mutex
	ceiling float64
	peak    map[string]float64
}

// New returns a Tracker that measures headroom against ceiling.
// ceiling must be greater than zero.
func New(ceiling float64) (*Tracker, error) {
	if ceiling <= 0 {
		return nil, ErrInvalidCeiling
	}
	return &Tracker{
		ceiling: ceiling,
		peak:    make(map[string]float64),
	}, nil
}

// Observe records a new value for key, updating the peak if value exceeds the
// current maximum.
func (t *Tracker) Observe(key string, value float64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if value > t.peak[key] {
		t.peak[key] = value
	}
}

// Headroom returns the remaining capacity for key as a fraction of the
// ceiling in the range [0, 1].  A value of 1.0 means the metric has never
// been observed; 0.0 means the peak has reached or exceeded the ceiling.
// ErrUnknownKey is returned when no observation has been recorded.
func (t *Tracker) Headroom(key string) (float64, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	peak, ok := t.peak[key]
	if !ok {
		return 0, ErrUnknownKey
	}
	remaining := (t.ceiling - peak) / t.ceiling
	if remaining < 0 {
		remaining = 0
	}
	return remaining, nil
}

// Reset clears all recorded observations.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.peak = make(map[string]float64)
}
