package coalescer

import (
	"testing"
	"time"
)

func TestNew_InvalidWindow(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
	_, err = New(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative window")
	}
}

func TestNew_Valid(t *testing.T) {
	c, err := New(time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Len() != 0 {
		t.Fatalf("expected empty coalescer")
	}
}

func TestAdd_FirstCallCountIsOne(t *testing.T) {
	c, _ := New(time.Second)
	now := time.Now()
	e := c.Add("hello", now)
	if e.Count != 1 {
		t.Fatalf("expected count 1, got %d", e.Count)
	}
}

func TestAdd_RepeatedWithinWindowIncrementsCount(t *testing.T) {
	c, _ := New(time.Second)
	now := time.Now()
	c.Add("hello", now)
	e := c.Add("hello", now.Add(100*time.Millisecond))
	if e.Count != 2 {
		t.Fatalf("expected count 2, got %d", e.Count)
	}
}

func TestAdd_AfterWindowResetsCount(t *testing.T) {
	c, _ := New(500 * time.Millisecond)
	now := time.Now()
	c.Add("hello", now)
	e := c.Add("hello", now.Add(time.Second))
	if e.Count != 1 {
		t.Fatalf("expected count reset to 1, got %d", e.Count)
	}
}

func TestAdd_DifferentMessagesAreIndependent(t *testing.T) {
	c, _ := New(time.Second)
	now := time.Now()
	c.Add("foo", now)
	c.Add("bar", now)
	if c.Len() != 2 {
		t.Fatalf("expected 2 distinct entries, got %d", c.Len())
	}
}

func TestFlush_ReturnsAllEntries(t *testing.T) {
	c, _ := New(time.Second)
	now := time.Now()
	c.Add("a", now)
	c.Add("b", now)
	entries := c.Flush()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestFlush_ResetsState(t *testing.T) {
	c, _ := New(time.Second)
	c.Add("x", time.Now())
	c.Flush()
	if c.Len() != 0 {
		t.Fatal("expected empty coalescer after flush")
	}
}
