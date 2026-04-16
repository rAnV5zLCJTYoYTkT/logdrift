package batcher

import (
	"sync"
	"testing"
	"time"
)

func TestNew_InvalidSize(t *testing.T) {
	_, err := New(0, time.Second, func([]string) {})
	if err != ErrInvalidSize {
		t.Fatalf("expected ErrInvalidSize, got %v", err)
	}
}

func TestNew_InvalidTimeout(t *testing.T) {
	_, err := New(10, 0, func([]string) {})
	if err != ErrInvalidTimeout {
		t.Fatalf("expected ErrInvalidTimeout, got %v", err)
	}
}

func TestAdd_FlushesOnFullBatch(t *testing.T) {
	var mu sync.Mutex
	var got [][]string
	b, err := New(3, time.Minute, func(batch []string) {
		mu.Lock()
		got = append(got, batch)
		mu.Unlock()
	})
	if err != nil {
		t.Fatal(err)
	}
	b.Add("a")
	b.Add("b")
	b.Add("c")
	time.Sleep(20 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	if len(got) != 1 {
		t.Fatalf("expected 1 flush, got %d", len(got))
	}
	if len(got[0]) != 3 {
		t.Fatalf("expected batch of 3, got %d", len(got[0]))
	}
}

func TestFlush_EmitsPartialBatch(t *testing.T) {
	var mu sync.Mutex
	var got [][]string
	b, err := New(10, time.Minute, func(batch []string) {
		mu.Lock()
		got = append(got, batch)
		mu.Unlock()
	})
	if err != nil {
		t.Fatal(err)
	}
	b.Add("x")
	b.Add("y")
	b.Flush()
	time.Sleep(20 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	if len(got) != 1 || len(got[0]) != 2 {
		t.Fatalf("expected 1 flush with 2 items, got %v", got)
	}
}

func TestFlush_EmptyBatchIsNoop(t *testing.T) {
	called := false
	b, _ := New(5, time.Minute, func([]string) { called = true })
	b.Flush()
	time.Sleep(20 * time.Millisecond)
	if called {
		t.Fatal("expected no flush on empty batch")
	}
}

func TestAdd_TimeoutFlush(t *testing.T) {
	var mu sync.Mutex
	var got [][]string
	b, err := New(100, 50*time.Millisecond, func(batch []string) {
		mu.Lock()
		got = append(got, batch)
		mu.Unlock()
	})
	if err != nil {
		t.Fatal(err)
	}
	b.Add("timer-item")
	time.Sleep(120 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	if len(got) == 0 {
		t.Fatal("expected timeout flush, got none")
	}
}
