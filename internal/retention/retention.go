// Package retention provides a policy-based log entry retention filter.
// Entries older than the configured TTL are considered expired and dropped.
package retention

import (
	"errors"
	"sync"
	"time"
)

// Policy holds retention configuration.
type Policy struct {
	TTL time.Duration
}

// Filter discards log entries that fall outside the retention window.
type Filter struct {
	mu     sync.Mutex
	policy Policy
	now    func() time.Time
}

// New creates a Filter with the given Policy.
// Returns an error if TTL is zero or negative.
func New(p Policy) (*Filter, error) {
	if p.TTL <= 0 {
		return nil, errors.New("retention: TTL must be positive")
	}
	return &Filter{policy: p, now: time.Now}, nil
}

// Allow returns true if ts is within the retention window.
func (f *Filter) Allow(ts time.Time) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	cutoff := f.now().Add(-f.policy.TTL)
	return ts.After(cutoff)
}

// TTL returns the configured retention duration.
func (f *Filter) TTL() time.Duration {
	return f.policy.TTL
}

// SetClock replaces the internal clock; intended for testing.
func (f *Filter) SetClock(fn func() time.Time) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.now = fn
}
