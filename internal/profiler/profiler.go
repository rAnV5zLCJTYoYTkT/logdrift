// Package profiler tracks per-key event rate profiles and flags keys
// whose recent rate deviates significantly from their historical baseline.
package profiler

import (
	"errors"
	"sync"
	"time"
)

// Profile holds a rolling mean and sample count for a single key.
type Profile struct {
	mean    float64
	count   float64
	lastSeen time.Time
}

// Profiler maintains per-key rate profiles.
type Profiler struct {
	mu        sync.Mutex
	profiles  map[string]*Profile
	threshold float64 // z-score-like multiplier
	ttl       time.Duration
}

// New creates a Profiler. threshold is the factor above mean that triggers
// an anomaly (e.g. 3.0 means 3× the running mean). ttl evicts stale keys.
func New(threshold float64, ttl time.Duration) (*Profiler, error) {
	if threshold <= 0 {
		return nil, errors.New("profiler: threshold must be positive")
	}
	if ttl <= 0 {
		return nil, errors.New("profiler: ttl must be positive")
	}
	return &Profiler{
		profiles:  make(map[string]*Profile),
		threshold: threshold,
		ttl:       ttl,
	}, nil
}

// Observe records a new rate sample for key and returns true if the sample
// is anomalous relative to the historical mean.
func (p *Profiler) Observe(key string, rate float64) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	p.evictLocked(now)

	pr, ok := p.profiles[key]
	if !ok {
		p.profiles[key] = &Profile{mean: rate, count: 1, lastSeen: now}
		return false
	}

	pr.lastSeen = now
	pr.count++
	// exponential moving average (α = 2/(count+1) capped at 0.1)
	alpha := 2.0 / (pr.count + 1)
	if alpha < 0.1 {
		alpha = 0.1
	}
	prev := pr.mean
	pr.mean = prev + alpha*(rate-prev)

	if prev == 0 {
		return false
	}
	return rate > prev*p.threshold
}

// evictLocked removes profiles that have not been seen within ttl.
func (p *Profiler) evictLocked(now time.Time) {
	for k, pr := range p.profiles {
		if now.Sub(pr.lastSeen) > p.ttl {
			delete(p.profiles, k)
		}
	}
}

// Len returns the number of active profiles.
func (p *Profiler) Len() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.profiles)
}
