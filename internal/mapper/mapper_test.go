package mapper_test

import (
	"testing"
	"time"

	"github.com/yourorg/logdrift/internal/mapper"
	"github.com/yourorg/logdrift/internal/parser"
	"github.com/yourorg/logdrift/internal/severity"
)

func baseLine() parser.LogLine {
	return parser.LogLine{
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Level:     "error",
		Message:   "something failed",
		Service:   "api",
		Latency:   12.5,
		Raw:       "2024-01-01T00:00:00Z error api something failed latency=12.5",
	}
}

func TestMap_ValidLine(t *testing.T) {
	m := mapper.New()
	e, err := m.Map(baseLine())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Level != severity.Error {
		t.Errorf("level: got %v, want error", e.Level)
	}
	if e.Service != "api" {
		t.Errorf("service: got %q, want api", e.Service)
	}
	if e.Latency != 12.5 {
		t.Errorf("latency: got %v, want 12.5", e.Latency)
	}
}

func TestMap_UnknownLevelReturnsError(t *testing.T) {
	m := mapper.New()
	l := baseLine()
	l.Level = "verbose"
	_, err := m.Map(l)
	if err == nil {
		t.Fatal("expected error for unknown level")
	}
}

func TestMap_EmptyServiceUsesDefault(t *testing.T) {
	m := mapper.New(mapper.WithDefaultService("fallback"))
	l := baseLine()
	l.Service = ""
	e, err := m.Map(l)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Service != "fallback" {
		t.Errorf("service: got %q, want fallback", e.Service)
	}
}

func TestMap_DefaultServiceIsUnknown(t *testing.T) {
	m := mapper.New()
	l := baseLine()
	l.Service = ""
	e, err := m.Map(l)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Service != "unknown" {
		t.Errorf("service: got %q, want unknown", e.Service)
	}
}

func TestMap_RawPreserved(t *testing.T) {
	m := mapper.New()
	l := baseLine()
	e, err := m.Map(l)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Raw != l.Raw {
		t.Errorf("raw: got %q, want %q", e.Raw, l.Raw)
	}
}
