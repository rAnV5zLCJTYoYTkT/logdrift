// Package ratelimit provides a token-bucket style rate limiter for
// controlling how frequently anomaly alerts are emitted per log source key.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter enforces a maximum number of events per interval per key.
type Limiter struct {
	mu       sync.Mutex
	buckets  map[string]*bucket
	max      int
	interval time.Duration
}

type bucket struct {
	tokens    int
	lastReset time.Time
}

// New creates a Limiter that allows at most maxTokens events per interval
// for each distinct key. maxTokens must be >= 1 and interval must be > 0.
func New(maxTokens int, interval time.Duration) (*Limiter, error) {
	if maxTokens < 1 {
		return nil, ErrInvalidMaxTokens
	}
	if interval <= 0 {
		return nil, ErrInvalidInterval
	}
	return &Limiter{
		buckets:  make(map[string]*bucket),
		max:      maxTokens,
		interval: interval,
	}, nil
}

// Allow reports whether the event identified by key is permitted under the
// current rate limit. It consumes one token from the bucket for key.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	b, ok := l.buckets[key]
	if !ok {
		l.buckets[key] = &bucket{tokens: l.max - 1, lastReset: now}
		return true
	}

	if now.Sub(b.lastReset) >= l.interval {
		b.tokens = l.max
		b.lastReset = now
	}

	if b.tokens <= 0 {
		return false
	}
	b.tokens--
	return true
}

// Reset clears the bucket for the given key, restoring full capacity.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.buckets, key)
}

// Errors returned by New.
var (
	ErrInvalidMaxTokens = limiterError("maxTokens must be >= 1")
	ErrInvalidInterval  = limiterError("interval must be > 0")
)

type limiterError string

func (e limiterError) Error() string { return string(e) }
