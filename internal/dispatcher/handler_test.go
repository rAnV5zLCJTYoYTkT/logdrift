package dispatcher_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/logdrift/internal/dispatcher"
)

func TestWriterHandler_WritesEntry(t *testing.T) {
	var buf bytes.Buffer
	h := dispatcher.NewWriterHandler(&buf)
	if err := h.Handle("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "hello") {
		t.Fatalf("expected 'hello' in output, got %q", buf.String())
	}
}

func TestWriterHandler_DefaultsToStdout(t *testing.T) {
	h := dispatcher.NewWriterHandler(nil)
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestFuncHandler_NilFuncReturnsError(t *testing.T) {
	_, err := dispatcher.NewFuncHandler(nil)
	if err == nil {
		t.Fatal("expected error for nil func")
	}
}

func TestFuncHandler_CallsFunc(t *testing.T) {
	var got string
	h, err := dispatcher.NewFuncHandler(func(s string) error {
		got = s
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := h.Handle("test-entry"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "test-entry" {
		t.Fatalf("expected 'test-entry', got %q", got)
	}
}

func TestFuncHandler_PropagatesError(t *testing.T) {
	sentinel := errors.New("fn error")
	h, _ := dispatcher.NewFuncHandler(func(_ string) error { return sentinel })
	if err := h.Handle("x"); !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}
