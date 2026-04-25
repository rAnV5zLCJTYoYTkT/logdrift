// Package scaler normalises numeric metrics into a [0, 1] range using
// min-max scaling so that heterogeneous signals (latency, error-rate, …)
// can be compared on a common scale.
package scaler

import (
	"errors"
	"math"
	"sync"
)

// ErrEmptyKey is returned when an empty metric key is supplied.
var ErrEmptyKey = errors.New("scaler: key must not be empty")

// ErrNoObservations is returned when Scale is called before any values have
// been recorded for the given key.
var ErrNoObservations = errors.New("scaler: no observations recorded for key")

type stats struct {
	min float64
	max float64
	set bool
}

// Scaler tracks per-key min/max bounds and maps raw values into [0, 1].
type Scaler struct {
	mu   sync.RWMutex
	keys map[string]*stats
}

// New returns a ready-to-use Scaler.
func New() *Scaler {
	return &Scaler{keys: make(map[string]*stats)}
}

// Observe records a raw value for the given key, expanding the tracked
// min/max bounds as necessary.
func (s *Scaler) Observe(key string, value float64) error {
	if key == "" {
		return ErrEmptyKey
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	st, ok := s.keys[key]
	if !ok {
		st = &stats{}
		s.keys[key] = st
	}
	if !st.set {
		st.min = value
		st.max = value
		st.set = true
		return nil
	}
	st.min = math.Min(st.min, value)
	st.max = math.Max(st.max, value)
	return nil
}

// Scale returns the min-max normalised value in [0, 1] for the given key.
// When min == max the function returns 0 to avoid a division-by-zero.
func (s *Scaler) Scale(key string, value float64) (float64, error) {
	if key == "" {
		return 0, ErrEmptyKey
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	st, ok := s.keys[key]
	if !ok || !st.set {
		return 0, ErrNoObservations
	}
	span := st.max - st.min
	if span == 0 {
		return 0, nil
	}
	return (value - st.min) / span, nil
}

// Reset clears all recorded observations for the given key.
func (s *Scaler) Reset(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.keys, key)
}
