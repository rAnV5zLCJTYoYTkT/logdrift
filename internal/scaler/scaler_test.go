package scaler_test

import (
	"testing"

	"github.com/yourusername/logdrift/internal/scaler"
)

func TestObserve_EmptyKeyReturnsError(t *testing.T) {
	s := scaler.New()
	if err := s.Observe("", 1.0); err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestScale_EmptyKeyReturnsError(t *testing.T) {
	s := scaler.New()
	if _, err := s.Scale("", 1.0); err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestScale_NoObservationsReturnsError(t *testing.T) {
	s := scaler.New()
	if _, err := s.Scale("latency", 42.0); err == nil {
		t.Fatal("expected error when no observations recorded, got nil")
	}
}

func TestScale_MinMaxSameReturnsZero(t *testing.T) {
	s := scaler.New()
	_ = s.Observe("latency", 5.0)
	v, err := s.Scale("latency", 5.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 0 {
		t.Fatalf("expected 0 when min==max, got %v", v)
	}
}

func TestScale_NormalisesCorrectly(t *testing.T) {
	s := scaler.New()
	for _, v := range []float64{0, 50, 100} {
		_ = s.Observe("latency", v)
	}
	cases := []struct {
		input float64
		want  float64
	}{
		{0, 0.0},
		{50, 0.5},
		{100, 1.0},
	}
	for _, tc := range cases {
		got, err := s.Scale("latency", tc.input)
		if err != nil {
			t.Fatalf("Scale(%v): unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("Scale(%v): got %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestScale_DifferentKeysAreIndependent(t *testing.T) {
	s := scaler.New()
	_ = s.Observe("a", 0)
	_ = s.Observe("a", 100)
	_ = s.Observe("b", 0)
	_ = s.Observe("b", 10)

	va, _ := s.Scale("a", 50)
	vb, _ := s.Scale("b", 5)
	if va != 0.5 {
		t.Errorf("key a: got %v, want 0.5", va)
	}
	if vb != 0.5 {
		t.Errorf("key b: got %v, want 0.5", vb)
	}
}

func TestReset_ClearsObservations(t *testing.T) {
	s := scaler.New()
	_ = s.Observe("latency", 10)
	s.Reset("latency")
	if _, err := s.Scale("latency", 10); err == nil {
		t.Fatal("expected error after reset, got nil")
	}
}
