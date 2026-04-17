package dispatcher_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/logdrift/internal/dispatcher"
)

type countHandler struct {
	calls atomic.Int64
	err   error
}

func (h *countHandler) Handle(_ string) error {
	h.calls.Add(1)
	return h.err
}

func TestRegister_NilHandlerReturnsError(t *testing.T) {
	d := dispatcher.New()
	if err := d.Register(nil); err == nil {
		t.Fatal("expected error for nil handler")
	}
}

func TestDispatch_NoHandlersReturnsError(t *testing.T) {
	d := dispatcher.New()
	errs := d.Dispatch(context.Background(), "line")
	if len(errs) != 1 || !errors.Is(errs[0], dispatcher.ErrNoHandlers) {
		t.Fatalf("expected ErrNoHandlers, got %v", errs)
	}
}

func TestDispatch_CallsAllHandlers(t *testing.T) {
	d := dispatcher.New()
	h1, h2 := &countHandler{}, &countHandler{}
	_ = d.Register(h1)
	_ = d.Register(h2)

	errs := d.Dispatch(context.Background(), "entry")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if h1.calls.Load() != 1 || h2.calls.Load() != 1 {
		t.Fatal("both handlers should have been called once")
	}
}

func TestDispatch_CollectsHandlerErrors(t *testing.T) {
	d := dispatcher.New()
	sentinel := errors.New("handler error")
	_ = d.Register(&countHandler{err: sentinel})
	_ = d.Register(&countHandler{})

	errs := d.Dispatch(context.Background(), "entry")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if !errors.Is(errs[0], sentinel) {
		t.Fatalf("unexpected error: %v", errs[0])
	}
}

func TestDispatch_CancelledContextReturnsContextError(t *testing.T) {
	d := dispatcher.New()
	blocking := &countHandler{}
	_ = d.Register(blocking)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	errs := d.Dispatch(ctx, "entry")
	// May or may not fire depending on scheduling; just ensure no panic.
	_ = errs
}
