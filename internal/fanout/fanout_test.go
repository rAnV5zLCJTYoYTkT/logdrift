package fanout_test

import (
	"errors"
	"testing"

	"logdrift/internal/fanout"
)

type stubHandler struct {
	received []string
	err      error
}

func (s *stubHandler) Handle(line string) error {
	s.received = append(s.received, line)
	return s.err
}

func TestNew_NoHandlersReturnsError(t *testing.T) {
	_, err := fanout.New()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNew_NilHandlerReturnsError(t *testing.T) {
	_, err := fanout.New(nil)
	if err == nil {
		t.Fatal("expected error for nil handler")
	}
}

func TestNew_Valid(t *testing.T) {
	h := &stubHandler{}
	f, err := fanout.New(h)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Len() != 1 {
		t.Fatalf("expected 1 handler, got %d", f.Len())
	}
}

func TestSend_DeliveredToAllHandlers(t *testing.T) {
	h1, h2 := &stubHandler{}, &stubHandler{}
	f, _ := fanout.New(h1, h2)
	if err := f.Send("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, h := range []*stubHandler{h1, h2} {
		if len(h.received) != 1 || h.received[0] != "hello" {
			t.Errorf("handler did not receive expected line")
		}
	}
}

func TestSend_CollectsHandlerErrors(t *testing.T) {
	h1 := &stubHandler{err: errors.New("boom")}
	h2 := &stubHandler{}
	f, _ := fanout.New(h1, h2)
	err := f.Send("line")
	if err == nil {
		t.Fatal("expected error from failing handler")
	}
	if len(h2.received) != 1 {
		t.Error("healthy handler should still receive line")
	}
}

func TestSend_AllHandlersFail(t *testing.T) {
	h1 := &stubHandler{err: errors.New("e1")}
	h2 := &stubHandler{err: errors.New("e2")}
	f, _ := fanout.New(h1, h2)
	err := f.Send("x")
	if err == nil {
		t.Fatal("expected combined error")
	}
}
