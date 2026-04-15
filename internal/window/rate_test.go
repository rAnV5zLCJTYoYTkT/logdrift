package window

import (
	"testing"
	"time"
)

func TestRate_EmptyWindow(t *testing.T) {
	s, _ := New(5, time.Minute)
	if r := s.Rate(); r != 0 {
		t.Fatalf("expected 0 rate on empty window, got %f", r)
	}
}

func TestRate_NonZeroAfterAdd(t *testing.T) {
	s, _ := New(60, time.Minute)
	s.Add(60)
	if r := s.Rate(); r <= 0 {
		t.Fatalf("expected positive rate, got %f", r)
	}
}

func TestPeak_EmptyWindow(t *testing.T) {
	s, _ := New(5, time.Minute)
	if p := s.Peak(); p != 0 {
		t.Fatalf("expected 0 peak on empty window, got %d", p)
	}
}

func TestPeak_SingleBucket(t *testing.T) {
	s, _ := New(5, time.Minute)
	s.Add(42)
	if p := s.Peak(); p != 42 {
		t.Fatalf("expected peak 42, got %d", p)
	}
}

func TestPeak_MultipleBuckets(t *testing.T) {
	s, _ := New(10, time.Second)
	s.Add(5)
	time.Sleep(15 * time.Millisecond)
	s.Add(20)
	time.Sleep(15 * time.Millisecond)
	s.Add(3)
	if p := s.Peak(); p < 20 {
		t.Fatalf("expected peak >= 20, got %d", p)
	}
}

func TestPeak_EvictedBucketsNotCounted(t *testing.T) {
	s, _ := New(5, 50*time.Millisecond)
	s.Add(100)
	time.Sleep(60 * time.Millisecond)
	s.Add(1)
	if p := s.Peak(); p > 1 {
		t.Fatalf("expected peak 1 after eviction, got %d", p)
	}
}
