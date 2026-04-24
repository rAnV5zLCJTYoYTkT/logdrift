package herald_test

import (
	"errors"
	"strings"
	"testing"

	"logdrift/internal/herald"
)

// recordSink captures every message sent to it.
type recordSink struct {
	msgs []string
	err  error
}

func (r *recordSink) Send(msg string) error {
	if r.err != nil {
		return r.err
	}
	r.msgs = append(r.msgs, msg)
	return nil
}

func TestRegister_EmptyNameReturnsError(t *testing.T) {
	h := herald.New()
	if err := h.Register("", &recordSink{}); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestRegister_NilSinkReturnsError(t *testing.T) {
	h := herald.New()
	if err := h.Register("s1", nil); err == nil {
		t.Fatal("expected error for nil sink")
	}
}

func TestDispatch_NoSinksReturnsError(t *testing.T) {
	h := herald.New()
	if err := h.Dispatch("hello"); err == nil {
		t.Fatal("expected error when no sinks registered")
	}
}

func TestDispatch_DeliveredToAllSinks(t *testing.T) {
	h := herald.New()
	a, b := &recordSink{}, &recordSink{}
	_ = h.Register("a", a)
	_ = h.Register("b", b)

	if err := h.Dispatch("alert"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(a.msgs) != 1 || a.msgs[0] != "alert" {
		t.Errorf("sink a: got %v", a.msgs)
	}
	if len(b.msgs) != 1 || b.msgs[0] != "alert" {
		t.Errorf("sink b: got %v", b.msgs)
	}
}

func TestDispatch_SinkErrorPropagated(t *testing.T) {
	h := herald.New()
	_ = h.Register("bad", &recordSink{err: errors.New("boom")})
	err := h.Dispatch("msg")
	if err == nil {
		t.Fatal("expected error from failing sink")
	}
	if !strings.Contains(err.Error(), "bad") {
		t.Errorf("error should mention sink name, got: %v", err)
	}
}

func TestUnregister_RemovesSink(t *testing.T) {
	h := herald.New()
	s := &recordSink{}
	_ = h.Register("s", s)
	h.Unregister("s")
	if h.Len() != 0 {
		t.Fatalf("expected 0 sinks, got %d", h.Len())
	}
}

func TestLen_ReflectsRegistrations(t *testing.T) {
	h := herald.New()
	if h.Len() != 0 {
		t.Fatal("expected 0 initially")
	}
	_ = h.Register("x", &recordSink{})
	if h.Len() != 1 {
		t.Fatal("expected 1 after register")
	}
}
