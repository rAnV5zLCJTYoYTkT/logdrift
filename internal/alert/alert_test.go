package alert

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestAlert_String(t *testing.T) {
	a := Alert{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Level:     LevelWarn,
		Message:   "latency spike detected",
		Value:     320.5,
		Mean:      100.0,
		StdDev:    15.0,
	}

	got := a.String()

	if !strings.Contains(got, "WARN") {
		t.Errorf("expected WARN in output, got: %s", got)
	}
	if !strings.Contains(got, "latency spike detected") {
		t.Errorf("expected message in output, got: %s", got)
	}
	if !strings.Contains(got, "320.50") {
		t.Errorf("expected value 320.50 in output, got: %s", got)
	}
}

func TestNewAlert_Level(t *testing.T) {
	tests := []struct {
		name  string
		level Level
	}{
		{"warn level", LevelWarn},
		{"error level", LevelError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAlert(tt.level, 1.0, 2.0, 3.0, "test")
			if a.Level != tt.level {
				t.Errorf("expected level %s, got %s", tt.level, a.Level)
			}
			if a.Timestamp.IsZero() {
				t.Error("expected non-zero timestamp")
			}
		})
	}
}

func TestNotifier_Notify(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(&buf)

	a := NewAlert(LevelError, 500.0, 100.0, 20.0, "error rate exceeded threshold")
	if err := n.Notify(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "ERROR") {
		t.Errorf("expected ERROR in notifier output, got: %s", output)
	}
	if !strings.Contains(output, "error rate exceeded threshold") {
		t.Errorf("expected message in notifier output, got: %s", output)
	}
}

func TestNewNotifier_DefaultsToStderr(t *testing.T) {
	n := NewNotifier(nil)
	if n.out == nil {
		t.Error("expected non-nil writer when nil is passed")
	}
}
