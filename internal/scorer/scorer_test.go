package scorer_test

import (
	"testing"

	"github.com/user/logdrift/internal/scorer"
	"github.com/user/logdrift/internal/severity"
)

func entry(level string, latency, mean, stddev, errRate float64) scorer.Entry {
	return scorer.Entry{
		Level:     severity.Parse(level),
		Latency:   latency,
		Mean:      mean,
		StdDev:    stddev,
		ErrorRate: errRate,
	}
}

func TestScore_ZeroForBenignEntry(t *testing.T) {
	s := scorer.New()
	v := s.Score(entry("info", 0, 0, 0, 0))
	if v != 0 {
		t.Fatalf("expected 0, got %f", v)
	}
}

func TestScore_InRange(t *testing.T) {
	s := scorer.New()
	v := s.Score(entry("error", 500, 100, 50, 0.8))
	if v < 0 || v > 1 {
		t.Fatalf("score out of range: %f", v)
	}
}

func TestScore_HighLatencyIncreasesScore(t *testing.T) {
	s := scorer.New()
	low := s.Score(entry("info", 110, 100, 10, 0))
	high := s.Score(entry("info", 200, 100, 10, 0))
	if high <= low {
		t.Fatalf("expected higher score for larger deviation, got low=%f high=%f", low, high)
	}
}

func TestScore_HighSeverityIncreasesScore(t *testing.T) {
	s := scorer.New()
	info := s.Score(entry("info", 0, 0, 0, 0))
	err := s.Score(entry("error", 0, 0, 0, 0))
	if err <= info {
		t.Fatalf("expected error > info, got info=%f error=%f", info, err)
	}
}

func TestScore_HighErrorRateIncreasesScore(t *testing.T) {
	s := scorer.New()
	low := s.Score(entry("info", 0, 0, 0, 0.1))
	high := s.Score(entry("info", 0, 0, 0, 0.9))
	if high <= low {
		t.Fatalf("expected higher error rate to raise score, got low=%f high=%f", low, high)
	}
}

func TestScore_CustomWeights(t *testing.T) {
	s := scorer.New(
		scorer.WithLatencyWeight(1),
		scorer.WithSeverityWeight(0),
		scorer.WithErrorRateWeight(0),
	)
	v := s.Score(entry("fatal", 0, 0, 0, 0))
	// latency absent → latency component 0; severity/error ignored
	if v != 0 {
		t.Fatalf("expected 0 with only latency weight and no latency, got %f", v)
	}
}

func TestScore_NeverExceedsOne(t *testing.T) {
	s := scorer.New()
	v := s.Score(entry("fatal", 1e9, 0, 1, 1.0))
	if v > 1 {
		t.Fatalf("score exceeded 1: %f", v)
	}
}
