package profiler

import (
	"testing"
	"time"
)

func TestNew_InvalidThreshold(t *testing.T) {
	_, err := New(0, time.Minute)
	if err == nil {
		t.Fatal("expected error for zero threshold")
	}
}

func TestNew_InvalidTTL(t *testing.T) {
	_, err := New(3.0, 0)
	if err == nil {
		t.Fatal("expected error for zero ttl")
	}
}

func TestNew_Valid(t *testing.T) {
	p, err := New(3.0, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Len() != 0 {
		t.Fatal("expected empty profiler")
	}
}

func TestObserve_FirstCallNotAnomalous(t *testing.T) {
	p, _ := New(2.0, time.Minute)
	if p.Observe("svc", 100.0) {
		t.Fatal("first observation should never be anomalous")
	}
}

func TestObserve_NormalRateNotAnomalous(t *testing.T) {
	p, _ := New(3.0, time.Minute)
	for i := 0; i < 10; i++ {
		p.Observe("svc", 50.0)
	}
	if p.Observe("svc", 55.0) {
		t.Fatal("slightly elevated rate should not be anomalous")
	}
}

func TestObserve_SpikeIsAnomalous(t *testing.T) {
	p, _ := New(3.0, time.Minute)
	for i := 0; i < 20; i++ {
		p.Observe("svc", 10.0)
	}
	if !p.Observe("svc", 500.0) {
		t.Fatal("50× spike should be anomalous")
	}
}

func TestObserve_DifferentKeysAreIndependent(t *testing.T) {
	p, _ := New(3.0, time.Minute)
	for i := 0; i < 10; i++ {
		p.Observe("a", 10.0)
		p.Observe("b", 200.0)
	}
	if p.Len() != 2 {
		t.Fatalf("expected 2 profiles, got %d", p.Len())
	}
}

func TestObserve_TTLEvictsStaleKey(t *testing.T) {
	p, _ := New(3.0, 10*time.Millisecond)
	p.Observe("stale", 1.0)
	time.Sleep(20 * time.Millisecond)
	// trigger eviction via a new observe on a different key
	p.Observe("fresh", 1.0)
	if p.Len() != 1 {
		t.Fatalf("expected stale key evicted, got Len=%d", p.Len())
	}
}
