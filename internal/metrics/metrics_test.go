package metrics_test

import (
	"sync"
	"testing"

	"github.com/yourorg/logdrift/internal/metrics"
)

func TestCounter_IncAndValue(t *testing.T) {
	var c metrics.Counter
	if c.Value() != 0 {
		t.Fatalf("expected 0, got %d", c.Value())
	}
	c.Inc()
	c.Inc()
	if c.Value() != 2 {
		t.Fatalf("expected 2, got %d", c.Value())
	}
}

func TestCounter_Add(t *testing.T) {
	var c metrics.Counter
	c.Add(10)
	if c.Value() != 10 {
		t.Fatalf("expected 10, got %d", c.Value())
	}
}

func TestCounter_Reset(t *testing.T) {
	var c metrics.Counter
	c.Add(5)
	c.Reset()
	if c.Value() != 0 {
		t.Fatalf("expected 0 after reset, got %d", c.Value())
	}
}

func TestCounter_ConcurrentInc(t *testing.T) {
	var c metrics.Counter
	var wg sync.WaitGroup
	const goroutines = 100
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() { defer wg.Done(); c.Inc() }()
	}
	wg.Wait()
	if c.Value() != goroutines {
		t.Fatalf("expected %d, got %d", goroutines, c.Value())
	}
}

func TestRegistry_Counter_CreateAndReuse(t *testing.T) {
	r := metrics.NewRegistry()
	c1 := r.Counter("lines")
	c1.Inc()
	c2 := r.Counter("lines")
	if c1 != c2 {
		t.Fatal("expected same counter instance")
	}
	if c2.Value() != 1 {
		t.Fatalf("expected 1, got %d", c2.Value())
	}
}

func TestRegistry_Snapshot(t *testing.T) {
	r := metrics.NewRegistry()
	r.Counter("parsed").Add(3)
	r.Counter("anomalies").Add(1)
	snap := r.Snapshot()
	if snap["parsed"] != 3 {
		t.Fatalf("expected parsed=3, got %d", snap["parsed"])
	}
	if snap["anomalies"] != 1 {
		t.Fatalf("expected anomalies=1, got %d", snap["anomalies"])
	}
}
