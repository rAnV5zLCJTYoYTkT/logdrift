package reaper

import (
	"testing"
	"time"
)

func TestNew_InvalidIdle(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero idle timeout")
	}
	_, err = New(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative idle timeout")
	}
}

func TestNew_Valid(t *testing.T) {
	r, err := New(time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil reaper")
	}
}

func TestTouch_And_Len(t *testing.T) {
	r, _ := New(time.Minute)
	r.Touch("a")
	r.Touch("b")
	if got := r.Len(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestSweep_NoExpiredEntries(t *testing.T) {
	r, _ := New(time.Minute)
	r.Touch("x")
	expired := r.Sweep()
	if len(expired) != 0 {
		t.Fatalf("expected no expired entries, got %v", expired)
	}
	if r.Len() != 1 {
		t.Fatal("entry should still be present")
	}
}

func TestSweep_RemovesExpiredEntries(t *testing.T) {
	r, _ := New(time.Minute)

	past := time.Now().Add(-2 * time.Minute)
	r.mu.Lock()
	r.entries["old"] = entry{lastSeen: past}
	r.mu.Unlock()

	r.Touch("fresh")

	expired := r.Sweep()
	if len(expired) != 1 || expired[0] != "old" {
		t.Fatalf("expected [old], got %v", expired)
	}
	if r.Len() != 1 {
		t.Fatal("fresh entry should remain")
	}
}

func TestSweep_EmptyStoreReturnsNil(t *testing.T) {
	r, _ := New(time.Second)
	if got := r.Sweep(); len(got) != 0 {
		t.Fatalf("expected empty result, got %v", got)
	}
}

func TestTouch_RefreshesTimestamp(t *testing.T) {
	r, _ := New(time.Minute)

	past := time.Now().Add(-2 * time.Minute)
	r.mu.Lock()
	r.entries["key"] = entry{lastSeen: past}
	r.mu.Unlock()

	// Refresh the key so it is no longer stale.
	r.Touch("key")

	expired := r.Sweep()
	if len(expired) != 0 {
		t.Fatalf("expected no expired entries after touch, got %v", expired)
	}
}
