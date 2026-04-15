package sampler_test

import (
	"testing"
	"time"

	"github.com/yourorg/logdrift/internal/sampler"
)

func TestAllow_FirstCallAlwaysAllowed(t *testing.T) {
	s := sampler.New(5 * time.Second)
	if !s.Allow("key1") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallWithinCooldownBlocked(t *testing.T) {
	s := sampler.New(5 * time.Second)
	s.Allow("key1")
	if s.Allow("key1") {
		t.Fatal("expected second call within cooldown to be blocked")
	}
}

func TestAllow_DifferentKeysIndependent(t *testing.T) {
	s := sampler.New(5 * time.Second)
	s.Allow("key1")
	if !s.Allow("key2") {
		t.Fatal("expected different key to be allowed independently")
	}
}

func TestAllow_ZeroCooldownAlwaysAllows(t *testing.T) {
	s := sampler.New(0)
	for i := 0; i < 5; i++ {
		if !s.Allow("key1") {
			t.Fatalf("expected call %d to be allowed with zero cooldown", i)
		}
	}
}

func TestAllow_AllowedAfterCooldownExpires(t *testing.T) {
	s := sampler.New(20 * time.Millisecond)
	s.Allow("key1")
	time.Sleep(30 * time.Millisecond)
	if !s.Allow("key1") {
		t.Fatal("expected key to be allowed after cooldown expires")
	}
}

func TestEvict_RemovesExpiredKeys(t *testing.T) {
	s := sampler.New(20 * time.Millisecond)
	s.Allow("key1")
	s.Allow("key2")

	if s.Len() != 2 {
		t.Fatalf("expected 2 keys, got %d", s.Len())
	}

	time.Sleep(30 * time.Millisecond)
	s.Evict()

	if s.Len() != 0 {
		t.Fatalf("expected 0 keys after eviction, got %d", s.Len())
	}
}

func TestEvict_KeepsActiveKeys(t *testing.T) {
	s := sampler.New(5 * time.Second)
	s.Allow("key1")
	s.Evict()

	if s.Len() != 1 {
		t.Fatalf("expected 1 active key to remain, got %d", s.Len())
	}
}
