package jitter

import (
	"testing"
	"time"
)

func TestNew_NegativeMinReturnsError(t *testing.T) {
	_, err := New(-time.Millisecond, time.Millisecond)
	if err == nil {
		t.Fatal("expected error for negative min")
	}
}

func TestNew_ZeroBaseReturnsError(t *testing.T) {
	_, err := New(0, 0)
	if err == nil {
		t.Fatal("expected error for zero base")
	}
}

func TestNew_NegativeBaseReturnsError(t *testing.T) {
	_, err := New(0, -time.Millisecond)
	if err == nil {
		t.Fatal("expected error for negative base")
	}
}

func TestNew_Valid(t *testing.T) {
	j, err := New(10*time.Millisecond, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if j == nil {
		t.Fatal("expected non-nil Jitter")
	}
}

func TestNext_WithinBounds(t *testing.T) {
	min := 10 * time.Millisecond
	base := 100 * time.Millisecond
	j, _ := New(min, base)

	for i := 0; i < 200; i++ {
		d := j.Next()
		if d < min {
			t.Fatalf("iteration %d: got %v, want >= %v", i, d, min)
		}
		if d >= min+base {
			t.Fatalf("iteration %d: got %v, want < %v", i, d, min+base)
		}
	}
}

func TestNext_ZeroMinWithinBounds(t *testing.T) {
	base := 50 * time.Millisecond
	j, _ := New(0, base)

	for i := 0; i < 100; i++ {
		d := j.Next()
		if d < 0 {
			t.Fatalf("got negative duration: %v", d)
		}
		if d >= base {
			t.Fatalf("got %v, want < %v", d, base)
		}
	}
}

func TestNextFrom_AnchorIsRespected(t *testing.T) {
	anchor := 200 * time.Millisecond
	base := 50 * time.Millisecond
	j, _ := New(0, base)

	for i := 0; i < 100; i++ {
		d := j.NextFrom(anchor)
		if d < anchor {
			t.Fatalf("iteration %d: got %v, want >= %v", i, d, anchor)
		}
		if d >= anchor+base {
			t.Fatalf("iteration %d: got %v, want < %v", i, d, anchor+base)
		}
	}
}

func TestNextFrom_NonPositiveAnchorFallsBackToNext(t *testing.T) {
	min := 5 * time.Millisecond
	base := 20 * time.Millisecond
	j, _ := New(min, base)

	d := j.NextFrom(0)
	if d < min {
		t.Fatalf("got %v, want >= %v", d, min)
	}
}
