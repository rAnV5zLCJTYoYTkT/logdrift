// Package aggregator groups log entries into time-bucketed windows
// and emits per-bucket summaries for downstream anomaly detection.
package aggregator

import (
	"sync"
	"time"

	"github.com/user/logdrift/internal/parser"
)

// Bucket holds aggregated statistics for a single time window.
type Bucket struct {
	Start    time.Time
	End      time.Time
	Count    int
	Errors   int
	AvgLatency float64
	latencySum float64
}

// Aggregator accumulates log lines into fixed-duration buckets.
type Aggregator struct {
	mu       sync.Mutex
	window   time.Duration
	current  *Bucket
	finished []*Bucket
}

// New creates an Aggregator with the given bucket window duration.
// It returns an error if the window is non-positive.
func New(window time.Duration) (*Aggregator, error) {
	if window <= 0 {
		return nil, ErrInvalidWindow
	}
	return &Aggregator{window: window}, nil
}

// Add incorporates a parsed log line into the current bucket,
// rotating to a new bucket when the window has elapsed.
func (a *Aggregator) Add(line parser.LogLine, now time.Time) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.current == nil {
		a.current = a.newBucket(now)
	}

	if now.After(a.current.End) {
		a.finished = append(a.finished, a.current)
		a.current = a.newBucket(now)
	}

	a.current.Count++
	if line.Level == "ERROR" || line.Level == "FATAL" {
		a.current.Errors++
	}
	if line.Latency > 0 {
		a.current.latencySum += line.Latency
		if a.current.Count > 0 {
			a.current.AvgLatency = a.current.latencySum / float64(a.current.Count)
		}
	}
}

// Drain returns and clears all completed buckets.
func (a *Aggregator) Drain() []*Bucket {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := a.finished
	a.finished = nil
	return out
}

func (a *Aggregator) newBucket(now time.Time) *Bucket {
	start := now.Truncate(a.window)
	return &Bucket{
		Start: start,
		End:   start.Add(a.window),
	}
}
