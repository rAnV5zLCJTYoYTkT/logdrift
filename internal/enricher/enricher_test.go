package enricher_test

import (
	"testing"

	"github.com/yourorg/logdrift/internal/enricher"
	"github.com/yourorg/logdrift/internal.LogLine{
		Message: msg,
		Logger:  logger,
		Level:   "INFO",
	}
}

func TestEnrich_FingerprintIsStable(t *testing.T) {
	e := enricher.New()
	a := e.Enrich(baseLogLine("request took 123ms", "app {
		t.Errorf("expected same fingerprint for structurally identical messages, got %q vs %q",
			a.Fingerprint, b.Fingerprint)
	}
}

func TestEnrich_FingerprintDiffersForDifferentMessages(t *testing.T) {
	e := enricher.New()
	a := e.Enrich(baseLogLine("request took 123ms", "app"))
	b := e.Enrich(baseLogLine("connection refused", "app"))

	if a.Fingerprint == b.Fingerprint {
		t.Error("expected different fingerprints for different messages")
	}
}

func TestEnrich_FingerprintLength(t *testing.T) {
	e := enricher.New()
	ent := e.Enrich(baseLogLine("some log message", "svc"))

	if len(ent.Fingerprint) != 8 {
		t.Errorf("expected fingerprint length 8, got %d", len(ent.Fingerprint))
	}
}

func TestEnrich_ServiceNameExtracted(t *testing.T) {
	e := enricher.New()
	ent := e.Enrich(baseLogLine("ok", "payments.processor"))

	if ent.ServiceName != "payments" {
		t.Errorf("expected service name %q, got %q", "payments", ent.ServiceName)
	}
}

func TestEnrich_ServiceNameNoHierarchy(t *testing.T) {
	e := enricher.New()
	ent := e.Enrich(baseLogLine("ok", "monolith"))

	if ent.ServiceName != "monolith" {
		t.Errorf("expected service name %q, got %q", "monolith", ent.ServiceName)
	}
}

func TestEnrich_PreservesOriginalLogLine(t *testing.T) {
	e := enricher.New()
	ll := baseLogLine("hello world", "app.core")
	ent := e.Enrich(ll)

	if ent.Message != ll.Message || ent.Level != ll.Level || ent.Logger != ll.Logger {
		t.Error("original LogLine fields were not preserved in Entry")
	}
}
