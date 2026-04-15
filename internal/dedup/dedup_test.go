package dedup

import (
	"testing"
	"time"
)

func TestIsDuplicate_FirstCallNotDuplicate(t *testing.T) {
	f := New(5 * time.Second)
	if f.IsDuplicate("hello world") {
		t.Fatal("first call should not be a duplicate")
	}
}

func TestIsDuplicate_SecondCallWithinTTLIsDuplicate(t *testing.T) {
	f := New(5 * time.Second)
	f.IsDuplicate("hello world")
	if !f.IsDuplicate("hello world") {
		t.Fatal("second call within TTL should be a duplicate")
	}
}

func TestIsDuplicate_DifferentMessagesAreIndependent(t *testing.T) {
	f := New(5 * time.Second)
	f.IsDuplicate("message A")
	if f.IsDuplicate("message B") {
		t.Fatal("different messages should not deduplicate each other")
	}
}

func TestIsDuplicate_ZeroTTLAlwaysPasses(t *testing.T) {
	f := New(0)
	f.IsDuplicate("msg")
	if f.IsDuplicate("msg") {
		t.Fatal("zero TTL should never deduplicate")
	}
}

func TestIsDuplicate_AllowedAfterTTLExpires(t *testing.T) {
	base := time.Unix(1_000_000, 0)
	f := New(2 * time.Second)
	f.nowFunc = func() time.Time { return base }

	f.IsDuplicate("msg")

	f.nowFunc = func() time.Time { return base.Add(3 * time.Second) }
	if f.IsDuplicate("msg") {
		t.Fatal("message should be allowed after TTL expires")
	}
}

func TestEvict_RemovesExpiredEntries(t *testing.T) {
	base := time.Unix(1_000_000, 0)
	f := New(2 * time.Second)
	f.nowFunc = func() time.Time { return base }

	f.IsDuplicate("msg1")
	f.IsDuplicate("msg2")

	f.nowFunc = func() time.Time { return base.Add(3 * time.Second) }
	f.Evict()

	if f.Len() != 0 {
		t.Fatalf("expected 0 entries after eviction, got %d", f.Len())
	}
}

func TestLen_TracksEntries(t *testing.T) {
	f := New(10 * time.Second)
	if f.Len() != 0 {
		t.Fatal("expected empty filter")
	}
	f.IsDuplicate("a")
	f.IsDuplicate("b")
	if f.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", f.Len())
	}
}
