// Package scorer assigns a numeric anomaly score to a log entry
// based on its severity, latency deviation, and error rate.
package scorer

import (
	"math"

	"github.com/user/logdrift/internal/severity"
)

// Entry holds the fields required to compute a score.
type Entry struct {
	Level      severity.Level
	Latency    float64 // milliseconds; 0 means absent
	Mean       float64 // rolling mean latency
	StdDev     float64 // rolling stddev latency
	ErrorRate  float64 // 0.0–1.0
}

// Scorer computes anomaly scores.
type Scorer struct {
	latencyWeight   float64
	severityWeight  float64
	errorRateWeight float64
}

// Option configures a Scorer.
type Option func(*Scorer)

// WithLatencyWeight overrides the latency component weight (default 0.5).
func WithLatencyWeight(w float64) Option {
	return func(s *Scorer) { s.latencyWeight = w }
}

// WithSeverityWeight overrides the severity component weight (default 0.3).
func WithSeverityWeight(w float64) Option {
	return func(s *Scorer) { s.severityWeight = w }
}

// WithErrorRateWeight overrides the error-rate component weight (default 0.2).
func WithErrorRateWeight(w float64) Option {
	return func(s *Scorer) { s.errorRateWeight = w }
}

// New returns a Scorer with default weights.
func New(opts ...Option) *Scorer {
	s := &Scorer{
		latencyWeight:   0.5,
		severityWeight:  0.3,
		errorRateWeight: 0.2,
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Score returns a value in [0, 1] representing how anomalous the entry is.
func (s *Scorer) Score(e Entry) float64 {
	latencyScore := s.latencyComponent(e)
	severityScore := s.severityComponent(e.Level)
	errorScore := clamp(e.ErrorRate)

	total := latencyScore*s.latencyWeight +
		severityScore*s.severityWeight +
		errorScore*s.errorRateWeight

	return clamp(total)
}

func (s *Scorer) latencyComponent(e Entry) float64 {
	if e.Latency <= 0 || e.StdDev <= 0 {
		return 0
	}
	zScore := math.Abs(e.Latency-e.Mean) / e.StdDev
	// sigmoid-like mapping: 3-sigma → ~0.95
	return clamp(1 - math.Exp(-zScore/3))
}

func (s *Scorer) severityComponent(l severity.Level) float64 {
	rank := severity.Rank(l)
	max := severity.Rank(severity.Level("fatal"))
	if max == 0 {
		return 0
	}
	return clamp(float64(rank) / float64(max))
}

func clamp(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
