package counter_test

import (
	"testing"
	"time"

	"github.com/logdrift/internal/counter"
)

func TestInc_StartsAtOne(t *testing.T) {
	c := counter.New(0)
	if got := c.Inc("key"); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestInc_Accumulates(t *testing.T) {
	c := counter.New(0)
	c.Inc("key")
	c.Inc("key")
	if got := c.Inc("key"); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestValue_UnknownKeyReturnsZero(t *testing.T) {
	c := counter.New(0)
	if got := c.Value("missing"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestValue_ReturnsCurrentCount(t *testing.T) {
	c := counter.New(0)
	c.Inc("k")
	c.Inc("k")
	if got := c.Value("k"); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestReset_ClearsKey(t *testing.T) {
	c := counter.New(0)
	c.Inc("k")
	c.Reset("k")
	if got := c.Value("k"); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestInc_ExpiresAfterTTL(t *testing.T) {
	c := counter.New(20 * time.Millisecond)
	c.Inc("k")
	c.Inc("k")
	time.Sleep(30 * time.Millisecond)
	if got := c.Inc("k"); got != 1 {
		t.Fatalf("expected 1 after TTL expiry, got %d", got)
	}
}

func TestValue_ExpiredKeyReturnsZero(t *testing.T) {
	c := counter.New(20 * time.Millisecond)
	c.Inc("k")
	time.Sleep(30 * time.Millisecond)
	if got := c.Value("k"); got != 0 {
		t.Fatalf("expected 0 for expired key, got %d", got)
	}
}

func TestLen_CountsActiveKeys(t *testing.T) {
	c := counter.New(0)
	c.Inc("a")
	c.Inc("b")
	c.Inc("c")
	if got := c.Len(); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestLen_ExcludesExpiredKeys(t *testing.T) {
	c := counter.New(20 * time.Millisecond)
	c.Inc("a")
	c.Inc("b")
	time.Sleep(30 * time.Millisecond)
	if got := c.Len(); got != 0 {
		t.Fatalf("expected 0 after expiry, got %d", got)
	}
}
