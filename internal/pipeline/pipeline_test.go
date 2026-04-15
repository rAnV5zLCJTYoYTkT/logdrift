package pipeline_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/logdrift/internal/alert"
	"github.com/user/logdrift/internal/pipeline"
)

func newTestNotifier(buf *bytes.Buffer) *alert.Notifier {
	return alert.NewNotifier(buf)
}

func TestPipeline_New_InvalidWindow(t *testing.T) {
	_, err := pipeline.New(pipeline.Config{
		WindowSize: 0,
		Threshold:  2.0,
		Notifier:   newTestNotifier(&bytes.Buffer{}),
	})
	if err == nil {
		t.Fatal("expected error for window size 0, got nil")
	}
}

func TestPipeline_Run_NoLines(t *testing.T) {
	buf := &bytes.Buffer{}
	p, err := pipeline.New(pipeline.Config{
		WindowSize: 10,
		Threshold:  2.0,
		Notifier:   newTestNotifier(buf),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := p.Run(strings.NewReader("")); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no alerts, got: %s", buf.String())
	}
}

func TestPipeline_Run_SkipsInvalidLines(t *testing.T) {
	buf := &bytes.Buffer{}
	p, _ := pipeline.New(pipeline.Config{
		WindowSize: 10,
		Threshold:  2.0,
		Notifier:   newTestNotifier(buf),
	})
	input := "not a valid log line\nalso invalid\n"
	if err := p.Run(strings.NewReader(input)); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no alerts for invalid lines, got: %s", buf.String())
	}
}

func TestPipeline_Run_EmitsAlertOnAnomaly(t *testing.T) {
	buf := &bytes.Buffer{}
	p, _ := pipeline.New(pipeline.Config{
		WindowSize: 10,
		Threshold:  1.0,
		Notifier:   newTestNotifier(buf),
	})
	// Feed normal lines to build baseline, then one spike
	normal := strings.Repeat("2024-01-01T00:00:00Z INFO request latency=10ms\n", 8)
	spike := "2024-01-01T00:00:01Z ERROR request latency=9999ms\n"
	if err := p.Run(strings.NewReader(normal + spike)); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected at least one alert for spike latency, got none")
	}
}
