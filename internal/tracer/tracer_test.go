package tracer

import (
	"testing"
	"time"
)

var now = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestNew_InvalidTTL(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero TTL")
	}
}

func TestNew_Valid(t *testing.T) {
	tr, err := New(time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.Len() != 0 {
		t.Fatal("expected empty tracer")
	}
}

func TestAdd_EmptyTraceIDReturnsError(t *testing.T) {
	tr, _ := New(time.Minute)
	err := tr.Add("", "msg", now)
	if err == nil {
		t.Fatal("expected error for empty traceID")
	}
}

func TestAdd_GroupsByTraceID(t *testing.T) {
	tr, _ := New(time.Minute)
	_ = tr.Add("abc", "first", now)
	_ = tr.Add("abc", "second", now.Add(time.Second))
	_ = tr.Add("xyz", "other", now)

	if tr.Len() != 2 {
		t.Fatalf("expected 2 traces, got %d", tr.Len())
	}
	trace, ok := tr.Get("abc")
	if !ok {
		t.Fatal("expected trace abc")
	}
	if len(trace.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(trace.Entries))
	}
}

func TestGet_UnknownTraceReturnsFalse(t *testing.T) {
	tr, _ := New(time.Minute)
	_, ok := tr.Get("missing")
	if ok {
		t.Fatal("expected false for unknown trace")
	}
}

func TestEvict_RemovesStaleTraces(t *testing.T) {
	tr, _ := New(time.Minute)
	_ = tr.Add("old", "msg", now)
	_ = tr.Add("new", "msg", now.Add(2*time.Minute))

	evicted := tr.Evict(now.Add(2 * time.Minute))
	if len(evicted) != 1 || evicted[0].ID != "old" {
		t.Fatalf("expected old trace evicted, got %v", evicted)
	}
	if tr.Len() != 1 {
		t.Fatalf("expected 1 remaining trace, got %d", tr.Len())
	}
}

func TestEvict_RetainsActiveTraces(t *testing.T) {
	tr, _ := New(time.Minute)
	_ = tr.Add("active", "msg", now)
	evicted := tr.Evict(now.Add(30 * time.Second))
	if len(evicted) != 0 {
		t.Fatalf("expected no evictions, got %d", len(evicted))
	}
}

func TestTrace_TimestampsTracked(t *testing.T) {
	tr, _ := New(time.Minute)
	_ = tr.Add("t1", "a", now)
	_ = tr.Add("t1", "b", now.Add(5*time.Second))
	trace, _ := tr.Get("t1")
	if !trace.First.Equal(now) {
		t.Errorf("expected First=%v, got %v", now, trace.First)
	}
	if !trace.Last.Equal(now.Add(5 * time.Second)) {
		t.Errorf("expected Last=%v, got %v", now.Add(5*time.Second), trace.Last)
	}
}
