package retention_test

import (
	"testing"
	"time"

	"github.com/user/logdrift/internal/retention"
)

func TestParseTTL_Hours(t *testing.T) {
	p, err := retention.ParseTTL("2h")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.TTL != 2*time.Hour {
		t.Errorf("expected 2h, got %v", p.TTL)
	}
}

func TestParseTTL_Minutes(t *testing.T) {
	p, err := retention.ParseTTL("45m")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.TTL != 45*time.Minute {
		t.Errorf("expected 45m, got %v", p.TTL)
	}
}

func TestParseTTL_Days(t *testing.T) {
	p, err := retention.ParseTTL("3d")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.TTL != 72*time.Hour {
		t.Errorf("expected 72h, got %v", p.TTL)
	}
}

func TestParseTTL_InvalidString(t *testing.T) {
	_, err := retention.ParseTTL("banana")
	if err == nil {
		t.Fatal("expected error for invalid duration")
	}
}

func TestParseTTL_ZeroDuration(t *testing.T) {
	_, err := retention.ParseTTL("0s")
	if err == nil {
		t.Fatal("expected error for zero duration")
	}
}
