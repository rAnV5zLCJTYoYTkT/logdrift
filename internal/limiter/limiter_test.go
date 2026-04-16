package limiter

import (
	"sync"
	"testing"
)

func TestNew_InvalidConcurrency(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for maxConcurrent=0")
	}
}

func TestNew_Valid(t *testing.T) {
	l, err := New(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.Cap() != 3 {
		t.Fatalf("expected cap 3, got %d", l.Cap())
	}
}

func TestAcquireRelease_Active(t *testing.T) {
	l, _ := New(2)
	l.Acquire()
	if l.Active() != 1 {
		t.Fatalf("expected active=1, got %d", l.Active())
	}
	l.Acquire()
	if l.Active() != 2 {
		t.Fatalf("expected active=2, got %d", l.Active())
	}
	l.Release()
	if l.Active() != 1 {
		t.Fatalf("expected active=1 after release, got %d", l.Active())
	}
	l.Release()
	if l.Active() != 0 {
		t.Fatalf("expected active=0, got %d", l.Active())
	}
}

func TestTryAcquire_SucceedsWhenSlotAvailable(t *testing.T) {
	l, _ := New(1)
	if !l.TryAcquire() {
		t.Fatal("expected TryAcquire to succeed")
	}
	if l.TryAcquire() {
		t.Fatal("expected TryAcquire to fail when full")
	}
	l.Release()
	if !l.TryAcquire() {
		t.Fatal("expected TryAcquire to succeed after release")
	}
	l.Release()
}

func TestAcquire_ConcurrentSafety(t *testing.T) {
	const workers = 20
	l, _ := New(5)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.Acquire()
			defer l.Release()
			if l.Active() > l.Cap() {
				t.Errorf("active %d exceeded cap %d", l.Active(), l.Cap())
			}
		}()
	}
	wg.Wait()
	if l.Active() != 0 {
		t.Fatalf("expected active=0 after all workers done, got %d", l.Active())
	}
}
