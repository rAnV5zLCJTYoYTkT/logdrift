// Package watchdog monitors pipeline health by tracking error rates
// and emitting alerts when thresholds are exceeded within a rolling window.
package watchdog

import (
	"errors"
	"sync"
	"time"
)

// ErrInvalidThreshold is returned when the error rate threshold is out of range.
var ErrInvalidThreshold = errors.New("watchdog: threshold must be between 0 and 1")

// ErrInvalidWindow is returned when the window duration is non-positive.
var ErrInvalidWindow = errors.New("watchdog: window must be positive")

// AlertFunc is called when the error rate exceeds the configured threshold.
type AlertFunc func(rate float64)

// Watchdog tracks error and total event counts within a rolling time window.
type Watchdog struct {
	mu        sync.Mutex
	window    time.Duration
	threshold float64
	buckets   []bucket
	now       func() time.Time
	onAlert   AlertFunc
}

type bucket struct {
	at     time.Time
	total  int
	errors int
}

// New creates a Watchdog that fires onAlert when the rolling error rate
// exceeds threshold (0–1) over the given window duration.
func New(window time.Duration, threshold float64, onAlert AlertFunc) (*Watchdog, error) {
	if window <= 0 {
		return nil, ErrInvalidWindow
	}
	if threshold < 0 || threshold > 1 {
		return nil, ErrInvalidThreshold
	}
	return &Watchdog{
		window:    window,
		threshold: threshold,
		onAlert:   onAlert,
		now:       time.Now,
	}, nil
}

// Record records one event. isError should be true if the event represents
// an error condition. It returns true if an alert was fired.
func (w *Watchdog) Record(isError bool) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := w.now()
	w.evict(now)

	var b *bucket
	if len(w.buckets) > 0 {
		last := &w.buckets[len(w.buckets)-1]
		if now.Sub(last.at) < time.Second {
			b = last
		}
	}
	if b == nil {
		w.buckets = append(w.buckets, bucket{at: now})
		b = &w.buckets[len(w.buckets)-1]
	}

	b.total++
	if isError {
		b.errors++
	}

	rate := w.rate()
	if rate > w.threshold && w.onAlert != nil {
		w.onAlert(rate)
		return true
	}
	return false
}

// Rate returns the current rolling error rate.
func (w *Watchdog) Rate() float64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict(w.now())
	return w.rate()
}

func (w *Watchdog) rate() float64 {
	var total, errs int
	for _, b := range w.buckets {
		total += b.total
		errs += b.errors
	}
	if total == 0 {
		return 0
	}
	return float64(errs) / float64(total)
}

func (w *Watchdog) evict(now time.Time) {
	cutoff := now.Add(-w.window)
	i := 0
	for i < len(w.buckets) && w.buckets[i].at.Before(cutoff) {
		i++
	}
	w.buckets = w.buckets[i:]
}
