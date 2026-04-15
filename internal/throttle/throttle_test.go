package throttle

import (
	"testing"
	"time"
)

func TestAllow_FirstCallAlwaysAllowed(t *testing.T) {
	th := New(5 * time.Second)
	if !th.Allow("key1") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallWithinIntervalDenied(t *testing.T) {
	th := New(5 * time.Second)
	th.Allow("key1")
	if th.Allow("key1") {
		t.Fatal("expected second call within interval to be denied")
	}
}

func TestAllow_AllowedAfterIntervalExpires(t *testing.T) {
	now := time.Now()
	th := New(5 * time.Second)
	th.now = func() time.Time { return now }

	th.Allow("key1")

	// advance clock beyond the interval
	th.now = func() time.Time { return now.Add(6 * time.Second) }
	if !th.Allow("key1") {
		t.Fatal("expected call after interval to be allowed")
	}
}

func TestAllow_ZeroIntervalAlwaysAllows(t *testing.T) {
	th := New(0)
	for i := 0; i < 10; i++ {
		if !th.Allow("key") {
			t.Fatalf("expected zero-interval throttle to always allow (iteration %d)", i)
		}
	}
}

func TestAllow_DifferentKeysAreIndependent(t *testing.T) {
	th := New(5 * time.Second)
	th.Allow("a")
	if !th.Allow("b") {
		t.Fatal("expected different key to be allowed")
	}
}

func TestReset_AllowsImmediateRetry(t *testing.T) {
	th := New(5 * time.Second)
	th.Allow("key1")
	th.Reset("key1")
	if !th.Allow("key1") {
		t.Fatal("expected Allow after Reset to succeed")
	}
}

func TestLen_TracksActiveKeys(t *testing.T) {
	th := New(5 * time.Second)
	if th.Len() != 0 {
		t.Fatalf("expected 0 keys, got %d", th.Len())
	}
	th.Allow("a")
	th.Allow("b")
	if th.Len() != 2 {
		t.Fatalf("expected 2 keys, got %d", th.Len())
	}
	th.Reset("a")
	if th.Len() != 1 {
		t.Fatalf("expected 1 key after reset, got %d", th.Len())
	}
}
