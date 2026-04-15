package aggregator

import (
	"testing"
	"time"

	"github.com/user/logdrift/internal/parser"
)

func TestNew_InvalidWindow(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero window, got nil")
	}
	_, err = New(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative window, got nil")
	}
}

func TestNew_ValidWindow(t *testing.T) {
	a, err := New(time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil aggregator")
	}
}

func TestAdd_CountsLines(t *testing.T) {
	a, _ := New(time.Minute)
	now := time.Now()
	line := parser.LogLine{Level: "INFO", Message: "ok"}
	a.Add(line, now)
	a.Add(line, now)

	if a.current == nil {
		t.Fatal("expected current bucket to be set")
	}
	if a.current.Count != 2 {
		t.Fatalf("expected count 2, got %d", a.current.Count)
	}
}

func TestAdd_CountsErrors(t *testing.T) {
	a, _ := New(time.Minute)
	now := time.Now()
	a.Add(parser.LogLine{Level: "INFO"}, now)
	a.Add(parser.LogLine{Level: "ERROR"}, now)
	a.Add(parser.LogLine{Level: "FATAL"}, now)

	if a.current.Errors != 2 {
		t.Fatalf("expected 2 errors, got %d", a.current.Errors)
	}
}

func TestAdd_BucketRotation(t *testing.T) {
	a, _ := New(time.Minute)
	t0 := time.Now().Truncate(time.Minute)

	a.Add(parser.LogLine{Level: "INFO"}, t0)
	// advance past the window boundary
	a.Add(parser.LogLine{Level: "INFO"}, t0.Add(2*time.Minute))

	buckets := a.Drain()
	if len(buckets) != 1 {
		t.Fatalf("expected 1 finished bucket, got %d", len(buckets))
	}
	if buckets[0].Count != 1 {
		t.Fatalf("expected count 1 in finished bucket, got %d", buckets[0].Count)
	}
}

func TestDrain_ClearsFinished(t *testing.T) {
	a, _ := New(time.Minute)
	t0 := time.Now().Truncate(time.Minute)
	a.Add(parser.LogLine{Level: "INFO"}, t0)
	a.Add(parser.LogLine{Level: "INFO"}, t0.Add(2*time.Minute))

	a.Drain() // first drain captures the finished bucket
	second := a.Drain()
	if len(second) != 0 {
		t.Fatalf("expected empty drain, got %d buckets", len(second))
	}
}

func TestAdd_AvgLatency(t *testing.T) {
	a, _ := New(time.Minute)
	now := time.Now()
	a.Add(parser.LogLine{Level: "INFO", Latency: 100}, now)
	a.Add(parser.LogLine{Level: "INFO", Latency: 200}, now)

	if a.current.AvgLatency == 0 {
		t.Fatal("expected non-zero avg latency")
	}
}
