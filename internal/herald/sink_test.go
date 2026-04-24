package herald_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"logdrift/internal/herald"
)

func TestWriterSink_WritesMessage(t *testing.T) {
	var buf bytes.Buffer
	s := herald.NewWriterSink(&buf)
	if err := s.Send("test alert"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "test alert") {
		t.Errorf("expected message in output, got %q", buf.String())
	}
}

func TestWriterSink_DefaultsToStderr(t *testing.T) {
	s := herald.NewWriterSink(nil)
	if s == nil {
		t.Fatal("expected non-nil sink")
	}
}

func TestWriterSink_AppendsNewline(t *testing.T) {
	var buf bytes.Buffer
	s := herald.NewWriterSink(&buf)
	_ = s.Send("line")
	if !strings.HasSuffix(buf.String(), "\n") {
		t.Errorf("expected trailing newline, got %q", buf.String())
	}
}

func TestNewWebhookSink_EmptyURLReturnsError(t *testing.T) {
	_, err := herald.NewWebhookSink("")
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestWebhookSink_Send_Success(t *testing.T) {
	var received string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r.Body)
		received = buf.String()
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	s, err := herald.NewWebhookSink(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := s.Send("anomaly detected"); err != nil {
		t.Fatalf("send failed: %v", err)
	}
	if received != "anomaly detected" {
		t.Errorf("server got %q", received)
	}
}

func TestWebhookSink_Send_Non2xxReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	s, _ := herald.NewWebhookSink(ts.URL)
	if err := s.Send("msg"); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}
