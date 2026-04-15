package parser

import (
	"testing"
	"time"
)

func TestParse_ValidLineWithLatency(t *testing.T) {
	raw := "2024-01-15T10:23:45Z INFO  user login successful latency=42ms"
	line, err := Parse(raw)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	expectedTS, _ := time.Parse(time.RFC3339, "2024-01-15T10:23:45Z")
	if !line.Timestamp.Equal(expectedTS) {
		t.Errorf("expected timestamp %v, got %v", expectedTS, line.Timestamp)
	}
	if line.Level != LevelInfo {
		t.Errorf("expected level INFO, got %q", line.Level)
	}
	if line.Message != "user login successful" {
		t.Errorf("unexpected message: %q", line.Message)
	}
	if line.LatencyMs != 42.0 {
		t.Errorf("expected latency 42.0, got %f", line.LatencyMs)
	}
}

func TestParse_ValidLineWithoutLatency(t *testing.T) {
	raw := "2024-01-15T10:23:45Z ERROR  database connection failed"
	line, err := Parse(raw)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if line.Level != LevelError {
		t.Errorf("expected level ERROR, got %q", line.Level)
	}
	if line.LatencyMs != 0 {
		t.Errorf("expected latency 0, got %f", line.LatencyMs)
	}
}

func TestParse_AllLevels(t *testing.T) {
	levels := []LogLevel{LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal}
	for _, lvl := range levels {
		raw := "2024-01-15T10:23:45Z " + string(lvl) + "  test message"
		line, err := Parse(raw)
		if err != nil {
			t.Err %q: unexpected error: %v", lvl, err)
			continue
		}
		if line.Level != lvl {
			t.Errorf("expected level %q, got %q", lvl, line.Level)
		}
	}
}

func TestParse_InvalidLine(t *testing.T) {
	cases := []string{
		"",
		"not a log line at all",
		"2024-01-15T10:23:45Z UNKNOWN message",
		"INFO  missing timestamp",
	}
	for _, raw := range cases {
		_, err := Parse(raw)
		if err == nil {
			t.Errorf("expected error for input %q, got nil", raw)
		}
	}
}
