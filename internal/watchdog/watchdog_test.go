package watchdog

import (
	"testing"
	"time"
)

func TestNew_InvalidThreshold(t *testing.T) {
	_, err := New(time.Minute, 1.5, nil)
	if err != ErrInvalidThreshold {
		t.Fatalf("expected ErrInvalidThreshold, got %v", err)
	}
}

func TestNew_InvalidWindow(t *testing.T) {
	_, err := New(0, 0.5, nil)
	if err != ErrInvalidWindow {
		t.Fatalf("expected ErrInvalidWindow, got %v", err)
	}
}

func TestNew_Valid(t *testing.T) {
	w, err := New(time.Minute, 0.5, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil watchdog")
	}
}

func TestRate_EmptyIsZero(t *testing.T) {
	w, _ := New(time.Minute, 0.5, nil)
	if r := w.Rate(); r != 0 {
		t.Fatalf("expected 0, got %f", r)
	}
}

func TestRecord_NoErrorsRateZero(t *testing.T) {
	w, _ := New(time.Minute, 0.5, nil)
	for i := 0; i < 5; i++ {
		w.Record(false)
	}
	if r := w.Rate(); r != 0 {
		t.Fatalf("expected 0, got %f", r)
	}
}

func TestRecord_AllErrorsRateOne(t *testing.T) {
	w, _ := New(time.Minute, 0.5, nil)
	for i := 0; i < 4; i++ {
		w.Record(true)
	}
	if r := w.Rate(); r != 1.0 {
		t.Fatalf("expected 1.0, got %f", r)
	}
}

func TestRecord_AlertFiredWhenThresholdExceeded(t *testing.T) {
	fired := false
	w, _ := New(time.Minute, 0.4, func(r float64) { fired = true })
	w.Record(true)
	w.Record(true)
	w.Record(false)
	if !fired {
		t.Fatal("expected alert to fire")
	}
}

func TestRecord_AlertNotFiredBelowThreshold(t *testing.T) {
	fired := false
	w, _ := New(time.Minute, 0.9, func(r float64) { fired = true })
	w.Record(true)
	w.Record(false)
	w.Record(false)
	w.Record(false)
	if fired {
		t.Fatal("expected alert not to fire")
	}
}

func TestRecord_WindowEviction(t *testing.T) {
	w, _ := New(50*time.Millisecond, 0.5, nil)
	// inject a fake clock
	base := time.Now()
	w.now = func() time.Time { return base }
	w.Record(true)
	w.Record(true)

	// advance past window
	w.now = func() time.Time { return base.Add(100 * time.Millisecond) }
	w.Record(false)

	if r := w.Rate(); r != 0 {
		t.Fatalf("expected 0 after eviction, got %f", r)
	}
}
