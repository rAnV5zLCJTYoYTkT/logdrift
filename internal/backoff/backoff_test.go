package backoff

import (
	"testing"
	"time"
)

func TestNew_InvalidBase(t *testing.T) {
	_, err := New(0, time.Second, 2)
	if err == nil {
		t.Fatal("expected error for zero base delay")
	}
}

func TestNew_MaxLessThanBase(t *testing.T) {
	_, err := New(time.Second, time.Millisecond, 2)
	if err == nil {
		t.Fatal("expected error when max < base")
	}
}

func TestNew_FactorBelowOne(t *testing.T) {
	_, err := New(time.Millisecond, time.Second, 0.5)
	if err == nil {
		t.Fatal("expected error for factor < 1")
	}
}

func TestNew_Valid(t *testing.T) {
	b, err := New(10*time.Millisecond, time.Second, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil Backoff")
	}
}

func TestNext_GrowsWithAttempts(t *testing.T) {
	b, _ := New(10*time.Millisecond, time.Hour, 2)

	prev := b.Next("k")
	for i := 0; i < 4; i++ {
		next := b.Next("k")
		if next < prev {
			t.Fatalf("delay should grow: attempt %d got %v < %v", i+1, next, prev)
		}
		prev = next
	}
}

func TestNext_CappedAtMax(t *testing.T) {
	maxDelay := 50 * time.Millisecond
	b, _ := New(10*time.Millisecond, maxDelay, 10)

	for i := 0; i < 10; i++ {
		d := b.Next("k")
		// Allow 10 % jitter on top of max.
		if d > maxDelay+maxDelay/10+time.Millisecond {
			t.Fatalf("delay %v exceeds max %v (attempt %d)", d, maxDelay, i)
		}
	}
}

func TestReset_ClearsAttempts(t *testing.T) {
	b, _ := New(10*time.Millisecond, time.Second, 2)
	b.Next("k")
	b.Next("k")
	if b.Attempts("k") != 2 {
		t.Fatalf("expected 2 attempts, got %d", b.Attempts("k"))
	}
	b.Reset("k")
	if b.Attempts("k") != 0 {
		t.Fatalf("expected 0 attempts after reset, got %d", b.Attempts("k"))
	}
}

func TestNext_DifferentKeysAreIndependent(t *testing.T) {
	b, _ := New(10*time.Millisecond, time.Second, 2)
	b.Next("a")
	b.Next("a")
	b.Next("a")

	if b.Attempts("b") != 0 {
		t.Fatalf("key 'b' should have 0 attempts, got %d", b.Attempts("b"))
	}
}
