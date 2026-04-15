package ratelimit_test

import (
	"testing"
	"time"

	"github.com/yourorg/logdrift/internal/ratelimit"
)

func TestNew_InvalidMaxTokens(t *testing.T) {
	_, err := ratelimit.New(0, time.Second)
	if err == nil {
		t.Fatal("expected error for maxTokens=0, got nil")
	}
}

func TestNew_InvalidInterval(t *testing.T) {
	_, err := ratelimit.New(1, 0)
	if err == nil {
		t.Fatal("expected error for interval=0, got nil")
	}
}

func TestAllow_FirstCallsConsumeTokens(t *testing.T) {
	l, err := ratelimit.New(3, time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i := 0; i < 3; i++ {
		if !l.Allow("src") {
			t.Fatalf("call %d: expected Allow=true", i+1)
		}
	}
	if l.Allow("src") {
		t.Fatal("expected Allow=false after exhausting tokens")
	}
}

func TestAllow_DifferentKeysAreIndependent(t *testing.T) {
	l, _ := ratelimit.New(1, time.Hour)

	if !l.Allow("a") {
		t.Fatal("key a: first call should be allowed")
	}
	if !l.Allow("b") {
		t.Fatal("key b: first call should be allowed")
	}
	if l.Allow("a") {
		t.Fatal("key a: second call should be denied")
	}
}

func TestAllow_RefillsAfterInterval(t *testing.T) {
	l, _ := ratelimit.New(1, 20*time.Millisecond)

	if !l.Allow("x") {
		t.Fatal("first call should be allowed")
	}
	if l.Allow("x") {
		t.Fatal("second call within interval should be denied")
	}

	time.Sleep(30 * time.Millisecond)

	if !l.Allow("x") {
		t.Fatal("call after interval should be allowed")
	}
}

func TestReset_RestoresCapacity(t *testing.T) {
	l, _ := ratelimit.New(1, time.Hour)

	l.Allow("k")
	if l.Allow("k") {
		t.Fatal("expected deny before reset")
	}

	l.Reset("k")
	if !l.Allow("k") {
		t.Fatal("expected allow after reset")
	}
}
