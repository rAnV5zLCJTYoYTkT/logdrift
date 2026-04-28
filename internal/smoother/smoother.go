// Package smoother applies exponential moving average smoothing to
// a stream of numeric observations, reducing noise before anomaly
// detection decisions are made downstream.
package smoother

import (
	"errors"
	"sync"
)

// ErrInvalidAlpha is returned when the smoothing factor is out of range.
var ErrInvalidAlpha = errors.New("smoother: alpha must be in the range (0, 1]")

// ErrUnknownKey is returned when Scale is called before any observation
// has been recorded for the given key.
var ErrUnknownKey = errors.New("smoother: no observations recorded for key")

// Smoother maintains a per-key exponential moving average.
type Smoother struct {
	alpha  float64
	mu     sync.Mutex
	values map[string]float64
}

// New creates a Smoother with the given smoothing factor alpha.
// alpha must be in the range (0, 1]; a value of 1.0 means no smoothing
// (the EMA equals the latest observation).
func New(alpha float64) (*Smoother, error) {
	if alpha <= 0 || alpha > 1 {
		return nil, ErrInvalidAlpha
	}
	return &Smoother{
		alpha:  alpha,
		values: make(map[string]float64),
	}, nil
}

// Observe incorporates a new observation for the given key and returns
// the updated exponential moving average.
func (s *Smoother) Observe(key string, value float64) float64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	prev, ok := s.values[key]
	var ema float64
	if !ok {
		ema = value
	} else {
		ema = s.alpha*value + (1-s.alpha)*prev
	}
	s.values[key] = ema
	return ema
}

// Value returns the current EMA for the given key without updating it.
func (s *Smoother) Value(key string) (float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, ok := s.values[key]
	if !ok {
		return 0, ErrUnknownKey
	}
	return v, nil
}

// Reset removes all stored state for the given key.
func (s *Smoother) Reset(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.values, key)
}

// Len returns the number of keys currently tracked.
func (s *Smoother) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.values)
}
