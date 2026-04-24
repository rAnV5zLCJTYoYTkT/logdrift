// Package jitter adds randomised delay to retry and backoff strategies,
// preventing thundering-herd problems when many goroutines retry simultaneously.
package jitter

import (
	"errors"
	"math/rand"
	"time"
)

// Jitter applies a random fraction of the base duration on top of a fixed
// minimum, keeping delays within [min, min+base].
type Jitter struct {
	min  time.Duration
	base time.Duration
	rng  *rand.Rand
}

// New creates a Jitter instance.
//
//   - min  is the guaranteed floor delay (must be >= 0).
//   - base is the maximum additional random component (must be > 0).
func New(min, base time.Duration) (*Jitter, error) {
	if min < 0 {
		return nil, errors.New("jitter: min must be >= 0")
	}
	if base <= 0 {
		return nil, errors.New("jitter: base must be > 0")
	}
	return &Jitter{
		min:  min,
		base: base,
		//nolint:gosec // non-cryptographic use
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

// Next returns a duration in [min, min+base).
func (j *Jitter) Next() time.Duration {
	return j.min + time.Duration(j.rng.Int63n(int64(j.base)))
}

// NextFrom returns a duration in [base*0, base*1) added to the provided
// anchor, which is useful when combining with exponential backoff values.
func (j *Jitter) NextFrom(anchor time.Duration) time.Duration {
	if anchor <= 0 {
		return j.Next()
	}
	offset := time.Duration(j.rng.Int63n(int64(j.base)))
	return anchor + offset
}

// Sleep blocks for Next() duration.
func (j *Jitter) Sleep() {
	time.Sleep(j.Next())
}
