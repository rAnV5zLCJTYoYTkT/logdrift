package baseline

import (
	"math"
	"testing"
)

func TestNewRollingStats_InvalidCapacity(t *testing.T) {
	_, err := NewRollingStats(0)
	if err == nil {
		t.Fatal("expected error for zero capacity, got nil")
	}
}

func TestRollingStats_MeanAndStdDev(t *testing.T) {
	rs, err := NewRollingStats(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = rs.Mean()
	if err == nil {
		t.Fatal("expected error on empty window")
	}

	values := []float64{2, 4, 4, 4, 5}
	for _, v := range values {
		rs.Add(v)
	}

	mean, err := rs.Mean()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(mean-3.8) > 1e-9 {
		t.Errorf("expected mean 3.8, got %f", mean)
	}

	stddev, err := rs.StdDev()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// population stddev of [2,4,4,4,5] = ~0.9798
	expected := 0.9797958971132712
	if math.Abs(stddev-expected) > 1e-9 {
		t.Errorf("expected stddev %f, got %f", expected, stddev)
	}
}

func TestRollingStats_WindowEviction(t *testing.T) {
	rs, _ := NewRollingStats(3)

	rs.Add(10)
	rs.Add(10)
	rs.Add(10)
	rs.Add(1) // evicts first 10

	if rs.Len() != 3 {
		t.Errorf("expected window length 3, got %d",t}

	mean, _ := rs.Mean()
	// window is [10, 10, 1] -> mean = 7.0
	if math.Abs(mean-7.0) > 1e-9 {
		t.Errorf("expected mean 7.0 after eviction, got %f", mean)
	}
}

func TestRollingStats_IsAnomaly(t *testing.T) {
	rs, _ := NewRollingStats(10)

	for i := 0; i < 9; i++ {
		rs.Add(10)
	}

	// value close to mean — not anomalous
	anom, err := rs.IsAnomaly(10, 3.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if anom {
		t.Error("expected 10 to not be anomalous")
	}

	// value far from mean — anomalous
	anom, err = rs.IsAnomaly(100, 3.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !anom {
		t.Error("expected 100 to be anomalous")
	}
}
