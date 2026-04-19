package correlator

import (
	"testing"
	"time"
)

func TestNew_InvalidTTL(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero TTL")
	}
}

func TestNew_Valid(t *testing.T) {
	c, err := New(time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Len() != 0 {
		t.Fatal("expected empty correlator")
	}
}

func TestAdd_GroupsByKey(t *testing.T) {
	c, _ := New(time.Second)
	now := time.Now()
	c.Add(Entry{Key: "req-1", Message: "start", Timestamp: now})
	c.Add(Entry{Key: "req-1", Message: "end", Timestamp: now})
	c.Add(Entry{Key: "req-2", Message: "start", Timestamp: now})
	if c.Len() != 2 {
		t.Fatalf("expected 2 groups, got %d", c.Len())
	}
}

func TestFlush_ReturnsExpiredGroups(t *testing.T) {
	c, _ := New(50 * time.Millisecond)
	fixed := time.Now()
	c.now = func() time.Time { return fixed }
	c.Add(Entry{Key: "req-1", Message: "a", Timestamp: fixed})
	// advance clock past TTL
	c.now = func() time.Time { return fixed.Add(100 * time.Millisecond) }
	groups := c.Flush()
	if len(groups) != 1 {
		t.Fatalf("expected 1 flushed group, got %d", len(groups))
	}
	if groups[0].Key != "req-1" {
		t.Fatalf("unexpected key: %s", groups[0].Key)
	}
	if c.Len() != 0 {
		t.Fatal("expected correlator to be empty after flush")
	}
}

func TestFlush_RetainsActiveGroups(t *testing.T) {
	c, _ := New(time.Second)
	fixed := time.Now()
	c.now = func() time.Time { return fixed }
	c.Add(Entry{Key: "req-1", Message: "a", Timestamp: fixed})
	groups := c.Flush()
	if len(groups) != 0 {
		t.Fatalf("expected 0 flushed groups, got %d", len(groups))
	}
	if c.Len() != 1 {
		t.Fatal("expected active group to remain")
	}
}

func TestFlush_MultipleEntries(t *testing.T) {
	c, _ := New(50 * time.Millisecond)
	fixed := time.Now()
	c.now = func() time.Time { return fixed }
	c.Add(Entry{Key: "r", Message: "x", Timestamp: fixed})
	c.Add(Entry{Key: "r", Message: "y", Timestamp: fixed})
	c.now = func() time.Time { return fixed.Add(200 * time.Millisecond) }
	groups := c.Flush()
	if len(groups[0].Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(groups[0].Entries))
	}
}
