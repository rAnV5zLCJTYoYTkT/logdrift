package headroom_test

import (
	"testing"

	"github.com/yourorg/logdrift/internal/headroom"
)

func TestNew_InvalidCeiling(t *testing.T) {
	for _, c := range []float64{0, -1, -100} {
		_, err := headroom.New(c)
		if err == nil {
			t.Fatalf("expected error for ceiling %v", c)
		}
	}
}

func TestNew_Valid(t *testing.T) {
	tr, err := headroom.New(100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil tracker")
	}
}

func TestHeadroom_UnknownKeyReturnsError(t *testing.T) {
	tr, _ := headroom.New(100)
	_, err := tr.Headroom("missing")
	if err == nil {
		t.Fatal("expected ErrUnknownKey")
	}
}

func TestHeadroom_FullCapacityAtZeroObservation(t *testing.T) {
	tr, _ := headroom.New(100)
	tr.Observe("svc", 0)
	h, err := tr.Headroom("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h != 1.0 {
		t.Fatalf("expected 1.0, got %v", h)
	}
}

func TestHeadroom_HalfCeilingReturnsHalf(t *testing.T) {
	tr, _ := headroom.New(100)
	tr.Observe("svc", 50)
	h, err := tr.Headroom("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h != 0.5 {
		t.Fatalf("expected 0.5, got %v", h)
	}
}

func TestHeadroom_ExceedsCeilingClampsToZero(t *testing.T) {
	tr, _ := headroom.New(100)
	tr.Observe("svc", 150)
	h, err := tr.Headroom("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h != 0 {
		t.Fatalf("expected 0, got %v", h)
	}
}

func TestObserve_PeakIsMonotonicallyIncreasing(t *testing.T) {
	tr, _ := headroom.New(100)
	tr.Observe("svc", 80)
	tr.Observe("svc", 40) // lower — should not replace peak
	h, _ := tr.Headroom("svc")
	if h != 0.2 {
		t.Fatalf("expected 0.2, got %v", h)
	}
}

func TestReset_ClearsAllKeys(t *testing.T) {
	tr, _ := headroom.New(100)
	tr.Observe("a", 50)
	tr.Observe("b", 75)
	tr.Reset()
	if _, err := tr.Headroom("a"); err == nil {
		t.Fatal("expected error after reset for key a")
	}
	if _, err := tr.Headroom("b"); err == nil {
		t.Fatal("expected error after reset for key b")
	}
}
