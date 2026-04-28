package smoother

import (
	"math"
	"testing"
)

func TestNew_InvalidAlpha(t *testing.T) {
	cases := []float64{-1, 0, 1.1, 2}
	for _, a := range cases {
		_, err := New(a)
		if err == nil {
			t.Errorf("expected error for alpha=%v", a)
		}
	}
}

func TestNew_ValidAlpha(t *testing.T) {
	for _, a := range []float64{0.1, 0.5, 1.0} {
		s, err := New(a)
		if err != nil {
			t.Fatalf("unexpected error for alpha=%v: %v", a, err)
		}
		if s == nil {
			t.Fatal("expected non-nil Smoother")
		}
	}
}

func TestObserve_FirstCallReturnsValue(t *testing.T) {
	s, _ := New(0.5)
	got := s.Observe("svc", 42.0)
	if got != 42.0 {
		t.Fatalf("expected 42.0, got %v", got)
	}
}

func TestObserve_EMAConverges(t *testing.T) {
	s, _ := New(0.5)
	s.Observe("k", 100)
	ema := s.Observe("k", 0)
	// EMA = 0.5*0 + 0.5*100 = 50
	if math.Abs(ema-50.0) > 1e-9 {
		t.Fatalf("expected 50.0, got %v", ema)
	}
}

func TestObserve_AlphaOneNoSmoothing(t *testing.T) {
	s, _ := New(1.0)
	s.Observe("k", 999)
	got := s.Observe("k", 7)
	if got != 7 {
		t.Fatalf("expected 7, got %v", got)
	}
}

func TestValue_UnknownKeyReturnsError(t *testing.T) {
	s, _ := New(0.3)
	_, err := s.Value("missing")
	if err == nil {
		t.Fatal("expected error for unknown key")
	}
}

func TestValue_KnownKeyReturnsEMA(t *testing.T) {
	s, _ := New(0.5)
	s.Observe("x", 10)
	v, err := s.Value("x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 10 {
		t.Fatalf("expected 10, got %v", v)
	}
}

func TestReset_RemovesKey(t *testing.T) {
	s, _ := New(0.5)
	s.Observe("y", 5)
	s.Reset("y")
	_, err := s.Value("y")
	if err == nil {
		t.Fatal("expected error after reset")
	}
}

func TestLen_TracksKeys(t *testing.T) {
	s, _ := New(0.5)
	if s.Len() != 0 {
		t.Fatal("expected 0 initial keys")
	}
	s.Observe("a", 1)
	s.Observe("b", 2)
	if s.Len() != 2 {
		t.Fatalf("expected 2, got %d", s.Len())
	}
	s.Reset("a")
	if s.Len() != 1 {
		t.Fatalf("expected 1 after reset, got %d", s.Len())
	}
}
