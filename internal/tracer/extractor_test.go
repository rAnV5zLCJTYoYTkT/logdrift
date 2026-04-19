package tracer

import "testing"

func TestNewExtractor_EmptyPatternReturnsError(t *testing.T) {
	_, err := NewExtractor("", "trace_id")
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNewExtractor_EmptyGroupReturnsError(t *testing.T) {
	_, err := NewExtractor(`trace=(?P<trace_id>\w+)`, "")
	if err == nil {
		t.Fatal("expected error for empty group")
	}
}

func TestNewExtractor_InvalidPatternReturnsError(t *testing.T) {
	_, err := NewExtractor(`[invalid`, "trace_id")
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestExtract_MatchesNamedGroup(t *testing.T) {
	ex, err := NewExtractor(`trace=(?P<trace_id>[a-f0-9]+)`, "trace_id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := ex.Extract("2024-01-01 INFO trace=deadbeef request completed")
	if got != "deadbeef" {
		t.Errorf("expected deadbeef, got %q", got)
	}
}

func TestExtract_NoMatchReturnsEmpty(t *testing.T) {
	ex, _ := NewExtractor(`trace=(?P<trace_id>[a-f0-9]+)`, "trace_id")
	got := ex.Extract("INFO no trace here")
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestExtract_WrongGroupNameReturnsEmpty(t *testing.T) {
	ex, _ := NewExtractor(`trace=(?P<tid>[a-f0-9]+)`, "trace_id")
	got := ex.Extract("trace=abc123")
	if got != "" {
		t.Errorf("expected empty for wrong group name, got %q", got)
	}
}

func TestExtract_MultipleGroupsPicksCorrect(t *testing.T) {
	ex, _ := NewExtractor(`req=(?P<req_id>\d+) trace=(?P<trace_id>[a-f0-9]+)`, "trace_id")
	got := ex.Extract("req=42 trace=cafebabe done")
	if got != "cafebabe" {
		t.Errorf("expected cafebabe, got %q", got)
	}
}
