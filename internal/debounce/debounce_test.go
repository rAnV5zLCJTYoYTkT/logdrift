package debounce_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/logdrift/internal/debounce"
)

func TestNew_InvalidWait(t *testing.T) {
	_, err := debounce.New(0)
	if err == nil {
		t.Fatal("expected error for zero wait duration")
	}
	_, err = debounce.New(-time.Millisecond)
	if err == nil {
		t.Fatal("expected error for negative wait duration")
	}
}

func TestNew_ValidWait(t *testing.T) {
	d, err := debounce.New(10 * time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil Debouncer")
	}
}

func TestSubmit_CallbackFiredAfterWait(t *testing.T) {
	d, _ := debounce.New(20 * time.Millisecond)

	var called int32
	d.Submit("key", func() { atomic.StoreInt32(&called, 1) })

	time.Sleep(50 * time.Millisecond)
	if atomic.LoadInt32(&called) != 1 {
		t.Fatal("expected callback to have been called")
	}
}

func TestSubmit_RapidCallsOnlyFireOnce(t *testing.T) {
	d, _ := debounce.New(30 * time.Millisecond)

	var count int32
	for i := 0; i < 5; i++ {
		d.Submit("key", func() { atomic.AddInt32(&count, 1) })
		time.Sleep(5 * time.Millisecond)
	}

	time.Sleep(60 * time.Millisecond)
	if n := atomic.LoadInt32(&count); n != 1 {
		t.Fatalf("expected callback called once, got %d", n)
	}
}

func TestSubmit_DifferentKeysAreIndependent(t *testing.T) {
	d, _ := debounce.New(20 * time.Millisecond)

	var a, b int32
	d.Submit("a", func() { atomic.StoreInt32(&a, 1) })
	d.Submit("b", func() { atomic.StoreInt32(&b, 1) })

	time.Sleep(50 * time.Millisecond)
	if atomic.LoadInt32(&a) != 1 || atomic.LoadInt32(&b) != 1 {
		t.Fatal("expected both callbacks to fire independently")
	}
}

func TestFlush_InvokesPendingCallbackImmediately(t *testing.T) {
	d, _ := debounce.New(200 * time.Millisecond)

	var called int32
	d.Submit("key", func() { atomic.StoreInt32(&called, 1) })

	d.Flush("key", func() { atomic.StoreInt32(&called, 1) })

	if atomic.LoadInt32(&called) != 1 {
		t.Fatal("expected Flush to invoke callback immediately")
	}
	if d.Pending() != 0 {
		t.Fatal("expected no pending timers after Flush")
	}
}

func TestPending_ReturnsActiveTimerCount(t *testing.T) {
	d, _ := debounce.New(200 * time.Millisecond)

	d.Submit("x", func() {})
	d.Submit("y", func() {})

	if n := d.Pending(); n != 2 {
		t.Fatalf("expected 2 pending, got %d", n)
	}
}
