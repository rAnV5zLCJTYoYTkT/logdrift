package pipeline_test

import (
	"context"
	"io"
	"sync"
	"testing"

	"github.com/yourorg/logdrift/internal/alert"
	"github.com/yourorg/logdrift/internal/pipeline"
)

type testNotifier struct {
	mu     sync.Mutex
	alerts []alert.Alert
}

func newTestNotifier() *testNotifier { return &testNotifier{} }
func (n *testNotifier) Notify(a alert.Alert) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.alerts = append(n.alerts, a)
}
func (n *testNotifier) Len() int {
	n.mu.Lock()
	defer n.mu.Unlock()
	return len(n.alerts)
}

func TestPipeline_New_InvalidWindow(t *testing.T) {
	_, err := pipeline.New(1, 2.0, newTestNotifier())
	if err == nil {
		t.Fatal("expected error for window < 2")
	}
}

func TestPipeline_Run_NoLines(t *testing.T) {
	n := newTestNotifier()
	p, _ := pipeline.New(5, 2.0, n)
	ch := make(chan string)
	close(ch)
	if err := p.Run(context.Background(), ch, io.Discard); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.Len() != 0 {
		t.Fatalf("expected 0 alerts, got %d", n.Len())
	}
}

func TestPipeline_Run_SkipsInvalidLines(t *testing.T) {
	n := newTestNotifier()
	p, _ := pipeline.New(5, 2.0, n)
	ch := make(chan string, 2)
	ch <- "not a log line"
	ch <- "also bad"
	close(ch)
	p.Run(context.Background(), ch, io.Discard) //nolint:errcheck
	snap := p.Registry().Snapshot()
	if snap["lines.skipped"] != 2 {
		t.Fatalf("expected 2 skipped, got %d", snap["lines.skipped"])
	}
}

func TestPipeline_Run_EmitsAlertOnAnomaly(t *testing.T) {
	n := newTestNotifier()
	p, _ := pipeline.New(4, 1.0, n)
	ch := make(chan string, 10)
	// Seed baseline with consistent latency then spike.
	for _, l := range []string{
		`2024-01-01T00:00:00Z INFO  GET /a latency=10ms`,
		`2024-01-01T00:00:01Z INFO  GET /b latency=11ms`,
		`2024-01-01T00:00:02Z INFO  GET /c latency=10ms`,
		`2024-01-01T00:00:03Z INFO  GET /d latency=10ms`,
		`2024-01-01T00:00:04Z ERROR GET /e latency=9999ms`,
	} {
		ch <- l
	}
	close(ch)
	p.Run(context.Background(), ch, io.Discard) //nolint:errcheck
	if n.Len() == 0 {
		t.Fatal("expected at least one anomaly alert")
	}
}
