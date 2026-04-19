package eviction

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
	if c.Len() != 0 {
		t.Fatal("expected empty cache")
	}
}

func TestTrack_NewKeyReturnsTrue(t *testing.T) {
	c, _ := New(time.Minute)
	if !c.Track("key1") {
		t.Fatal("expected true for new key")
	}
}

func TestTrack_SameKeyWithinTTLReturnsFalse(t *testing.T) {
	c, _ := New(time.Minute)
	c.Track("key1")
	if c.Track("key1") {
		t.Fatal("expected false for key seen within TTL")
	}
}

func TestTrack_KeyAfterTTLResetsAndReturnsTrue(t *testing.T) {
	now := time.Now()
	c, _ := New(time.Second)
	c.now = func() time.Time { return now }
	c.Track("key1")
	c.now = func() time.Time { return now.Add(2 * time.Second) }
	if !c.Track("key1") {
		t.Fatal("expected true after TTL expiry")
	}
}

func TestEvict_RemovesExpiredEntries(t *testing.T) {
	now := time.Now()
	c, _ := New(time.Second)
	c.now = func() time.Time { return now }
	c.Track("a")
	c.Track("b")
	c.now = func() time.Time { return now.Add(2 * time.Second) }
	c.Track("c") // fresh
	removed := c.Evict()
	if removed != 2 {
		t.Fatalf("expected 2 removed, got %d", removed)
	}
	if c.Len() != 1 {
		t.Fatalf("expected 1 remaining, got %d", c.Len())
	}
}

func TestEvict_EmptyCacheReturnsZero(t *testing.T) {
	c, _ := New(time.Second)
	if n := c.Evict(); n != 0 {
		t.Fatalf("expected 0, got %d", n)
	}
}

func TestTrack_DifferentKeysAreIndependent(t *testing.T) {
	c, _ := New(time.Minute)
	c.Track("x")
	if !c.Track("y") {
		t.Fatal("expected true for distinct key")
	}
}
