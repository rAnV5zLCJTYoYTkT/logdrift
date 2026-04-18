package summarizer_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/logdrift/internal/summarizer"
)

func TestNew_DefaultsToStdout(t *testing.T) {
	s := summarizer.New(nil)
	if s == nil {
		t.Fatal("expected non-nil summarizer")
	}
}

func TestAdd_IncrementsTotal(t *testing.T) {
	s := summarizer.New(new(bytes.Buffer))
	s.Add(summarizer.Entry{Level: "error", Service: "api", Message: "boom"})
	s.Add(summarizer.Entry{Level: "info", Service: "api", Message: "ok"})
	snap := s.Snapshot()
	if snap.Total != 2 {
		t.Fatalf("expected total 2, got %d", snap.Total)
	}
}

func TestSnapshot_ByLevel(t *testing.T) {
	s := summarizer.New(new(bytes.Buffer))
	s.Add(summarizer.Entry{Level: "error", Service: "svc"})
	s.Add(summarizer.Entry{Level: "error", Service: "svc"})
	s.Add(summarizer.Entry{Level: "warn", Service: "svc"})
	snap := s.Snapshot()
	if snap.ByLevel["error"] != 2 {
		t.Errorf("expected 2 errors, got %d", snap.ByLevel["error"])
	}
	if snap.ByLevel["warn"] != 1 {
		t.Errorf("expected 1 warn, got %d", snap.ByLevel["warn"])
	}
}

func TestSnapshot_ByService(t *testing.T) {
	s := summarizer.New(new(bytes.Buffer))
	s.Add(summarizer.Entry{Level: "info", Service: "auth"})
	s.Add(summarizer.Entry{Level: "info", Service: "billing"})
	s.Add(summarizer.Entry{Level: "info", Service: "auth"})
	snap := s.Snapshot()
	if snap.ByService["auth"] != 2 {
		t.Errorf("expected auth=2, got %d", snap.ByService["auth"])
	}
	if snap.ByService["billing"] != 1 {
		t.Errorf("expected billing=1, got %d", snap.ByService["billing"])
	}
}

func TestFlush_ResetsCounters(t *testing.T) {
	buf := new(bytes.Buffer)
	s := summarizer.New(buf)
	s.Add(summarizer.Entry{Level: "info", Service: "api"})
	s.Flush()
	snap := s.Snapshot()
	if snap.Total != 0 {
		t.Errorf("expected total 0 after flush, got %d", snap.Total)
	}
}

func TestFlush_OutputContainsTotalAndLevel(t *testing.T) {
	buf := new(bytes.Buffer)
	s := summarizer.New(buf)
	s.Add(summarizer.Entry{Level: "error", Service: "svc"})
	s.Flush()
	out := buf.String()
	if !strings.Contains(out, "total=1") {
		t.Errorf("expected 'total=1' in output, got: %s", out)
	}
	if !strings.Contains(out, "level.error=1") {
		t.Errorf("expected 'level.error=1' in output, got: %s", out)
	}
	if !strings.Contains(out, "service.svc=1") {
		t.Errorf("expected 'service.svc=1' in output, got: %s", out)
	}
}
