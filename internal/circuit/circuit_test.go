package circuit

import (
	"testing"
	"time"
)

func TestNew_InvalidThreshold(t *testing.T) {
	_, err := New(0, time.Second)
	if err == nil {
		t.Fatal("expected error for zero threshold")
	}
}

func TestNew_InvalidCooldown(t *testing.T) {
	_, err := New(3, 0)
	if err == nil {
		t.Fatal("expected error for zero cooldown")
	}
}

func TestNew_Valid(t *testing.T) {
	b, err := New(3, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.State() != StateClosed {
		t.Fatalf("expected StateClosed, got %v", b.State())
	}
}

func TestAllow_ClosedInitially(t *testing.T) {
	b, _ := New(3, time.Second)
	if !b.Allow() {
		t.Fatal("expected Allow to return true when closed")
	}
}

func TestRecordFailure_OpensAfterThreshold(t *testing.T) {
	b, _ := New(3, time.Second)
	for i := 0; i < 3; i++ {
		b.RecordFailure()
	}
	if b.State() != StateOpen {
		t.Fatalf("expected StateOpen after %d failures", 3)
	}
	if b.Allow() {
		t.Fatal("expected Allow to return false when open")
	}
}

func TestRecordSuccess_ResetsBreakerToClosed(t *testing.T) {
	b, _ := New(2, time.Second)
	b.RecordFailure()
	b.RecordFailure()
	if b.State() != StateOpen {
		t.Fatal("expected StateOpen")
	}
	b.RecordSuccess()
	if b.State() != StateClosed {
		t.Fatalf("expected StateClosed after success, got %v", b.State())
	}
	if !b.Allow() {
		t.Fatal("expected Allow to return true after reset")
	}
}

func TestHalfOpen_TransitionAfterCooldown(t *testing.T) {
	b, _ := New(1, 20*time.Millisecond)
	b.RecordFailure()
	if b.State() != StateOpen {
		t.Fatal("expected StateOpen")
	}
	time.Sleep(30 * time.Millisecond)
	if !b.Allow() {
		t.Fatal("expected Allow after cooldown")
	}
	if b.State() != StateHalfOpen {
		t.Fatalf("expected StateHalfOpen, got %v", b.State())
	}
}

func TestHalfOpen_FailureReopens(t *testing.T) {
	b, _ := New(1, 20*time.Millisecond)
	b.RecordFailure()
	time.Sleep(30 * time.Millisecond)
	b.Allow() // transitions to HalfOpen
	b.RecordFailure()
	if b.State() != StateOpen {
		t.Fatalf("expected StateOpen after half-open failure, got %v", b.State())
	}
}
