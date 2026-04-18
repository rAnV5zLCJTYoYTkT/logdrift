package retention_test

import (
	"testing"
	"time"

	"github.com/user/logdrift/internal/retention"
)

func TestNew_InvalidTTL(t *testing.T) {
	_, err := retention.New(retention.Policy{TTL: 0})
	if err == nil {
		t.Fatal("expected error for zero TTL")
	}
	_, err = retention.New(retention.Policy{TTL: -time.Second})
	if err == nil {
		t.Fatal("expected error for negative TTL")
	}
}

func TestNew_ValidTTL(t *testing.T) {
	f, err := retention.New(retention.Policy{TTL: time.Minute})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.TTL() != time.Minute {
		t.Errorf("expected TTL=1m, got %v", f.TTL())
	}
}

func TestAllow_WithinTTL(t *testing.T) {
	f, _ := retention.New(retention.Policy{TTL: time.Hour})
	now := time.Now()
	f.SetClock(func() time.Time { return now })
	if !f.Allow(now.Add(-30 * time.Minute)) {
		t.Error("expected entry within TTL to be allowed")
	}
}

func TestAllow_ExpiredEntry(t *testing.T) {
	f, _ := retention.New(retention.Policy{TTL: time.Hour})
	now := time.Now()
	f.SetClock(func() time.Time { return now })
	if f.Allow(now.Add(-2 * time.Hour)) {
		t.Error("expected expired entry to be denied")
	}
}

func TestAllow_ExactlyAtCutoff(t *testing.T) {
	f, _ := retention.New(retention.Policy{TTL: time.Hour})
	now := time.Now()
	f.SetClock(func() time.Time { return now })
	// exactly at cutoff is not after cutoff
	if f.Allow(now.Add(-time.Hour)) {
		t.Error("expected entry at exact cutoff to be denied")
	}
}

func TestAllow_ClockAdvances(t *testing.T) {
	f, _ := retention.New(retention.Policy{TTL: time.Minute})
	base := time.Now()
	ts := base.Add(-30 * time.Second)

	f.SetClock(func() time.Time { return base })
	if !f.Allow(ts) {
		t.Error("expected allowed at base time")
	}

	f.SetClock(func() time.Time { return base.Add(31 * time.Second) })
	if f.Allow(ts) {
		t.Error("expected denied after clock advances past TTL")
	}
}
