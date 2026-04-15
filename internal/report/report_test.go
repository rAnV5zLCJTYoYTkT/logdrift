package report_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logdrift/internal/report"
)

func TestNewRecorder_DefaultsToStdout(t *testing.T) {
	r := report.NewRecorder(nil)
	if r == nil {
		t.Fatal("expected non-nil recorder")
	}
}

func TestRecorder_Counts(t *testing.T) {
	var buf bytes.Buffer
	r := report.NewRecorder(&buf)

	r.IncProcessed()
	r.IncProcessed()
	r.IncSkipped()
	r.RecordAlert("GET /api")

	s := r.GetSummary()
	if s.LinesProcessed != 2 {
		t.Errorf("expected 2 processed, got %d", s.LinesProcessed)
	}
	if s.LinesSkipped != 1 {
		t.Errorf("expected 1 skipped, got %d", s.LinesSkipped)
	}
	if s.AlertsEmitted != 1 {
		t.Errorf("expected 1 alert, got %d", s.AlertsEmitted)
	}
	if len(s.AnomalyKeys) != 1 || s.AnomalyKeys[0] != "GET /api" {
		t.Errorf("unexpected anomaly keys: %v", s.AnomalyKeys)
	}
}

func TestRecorder_Flush_ContainsSummaryFields(t *testing.T) {
	var buf bytes.Buffer
	r := report.NewRecorder(&buf)

	r.IncProcessed()
	r.IncProcessed()
	r.IncProcessed()
	r.IncSkipped()
	r.RecordAlert("POST /login")
	r.Flush()

	out := buf.String()
	for _, want := range []string{"processed", "skipped", "alerts", "POST /login", "summary"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestRecorder_Duration_NonNegative(t *testing.T) {
	var buf bytes.Buffer
	r := report.NewRecorder(&buf)
	r.Flush()
	s := r.GetSummary()
	if s.Duration() < 0 {
		t.Errorf("duration should be non-negative, got %s", s.Duration())
	}
}

func TestRecorder_Flush_NoAnomaliesSection(t *testing.T) {
	var buf bytes.Buffer
	r := report.NewRecorder(&buf)
	r.IncProcessed()
	r.Flush()

	out := buf.String()
	if strings.Contains(out, "anomalies") {
		t.Errorf("expected no anomalies section when none recorded, got:\n%s", out)
	}
}
