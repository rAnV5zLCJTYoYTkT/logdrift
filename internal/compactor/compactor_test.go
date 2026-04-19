package compactor

import (
	"testing"
	"time"
)

func TestNew_InvalidTTL(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero TTL")
	}
	_, err = New(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative TTL")
	}
}

func TestNew_Valid(t *testing.T) {
	c, err := New(time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil compactor")
	}
}

func TestAdd_FirstCallIsNew(t *testing.T) {
	c, _ := New(time.Minute)
	e, isNew := c.Add("fp1", "hello world")
	if !isNew {
		t.Fatal("expected first call to be new")
	}
	if e.Count != 1 {
		t.Fatalf("expected count 1, got %d", e.Count)
	}
}

func TestAdd_RepeatedWithinTTLNotNew(t *testing.T) {
	c, _ := New(time.Minute)
	c.Add("fp1", "hello")
	e, isNew := c.Add("fp1", "hello")
	if isNew {
		t.Fatal("expected duplicate to not be new")
	}
	if e.Count != 2 {
		t.Fatalf("expected count 2, got %d", e.Count)
	}
}

func TestAdd_DifferentFingerprintsAreIndependent(t *testing.T) {
	c, _ := New(time.Minute)
	_, n1 := c.Add("fp1", "msg1")
	_, n2 := c.Add("fp2", "msg2")
	if !n1 || !n2 {
		t.Fatal("expected both fingerprints to be new")
	}
}

func TestAdd_EvictsExpiredEntries(t *testing.T) {
	c, _ := New(50 * time.Millisecond)
	base := time.Now()
	c.now = func() time.Time { return base }
	c.Add("fp1", "msg")

	// advance past TTL
	c.now = func() time.Time { return base.Add(100 * time.Millisecond) }
	_, isNew := c.Add("fp1", "msg")
	if !isNew {
		t.Fatal("expected entry to be treated as new after TTL expiry")
	}
}

func TestFlush_ReturnsAndClearsEntries(t *testing.T) {
	c, _ := New(time.Minute)
	c.Add("fp1", "a")
	c.Add("fp2", "b")
	c.Add("fp1", "a")

	entries := c.Flush()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	after := c.Flush()
	if len(after) != 0 {
		t.Fatalf("expected empty after flush, got %d", len(after))
	}
}
