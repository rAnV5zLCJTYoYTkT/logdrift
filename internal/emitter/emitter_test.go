package emitter_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/logdrift/logdrift/internal/emitter"
)

func sampleEvent() emitter.Event {
	return emitter.Event{
		Timestamp: time.Now(),
		Service:   "api",
		Level:     "error",
		Message:   "latency spike detected",
		Score:     0.87,
		Tags:      []string{"anomaly", "latency"},
	}
}

func TestRegister_NilSinkReturnsError(t *testing.T) {
	e := emitter.New()
	if err := e.Register("x", nil); err == nil {
		t.Fatal("expected error for nil sink")
	}
}

func TestRegister_EmptyNameReturnsError(t *testing.T) {
	e := emitter.New()
	sink := emitter.NewWriterSink(nil)
	if err := e.Register("", sink); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestEmit_NoSinksReturnsError(t *testing.T) {
	e := emitter.New()
	if err := e.Emit(sampleEvent()); err == nil {
		t.Fatal("expected error when no sinks registered")
	}
}

func TestEmit_DeliveredToAllSinks(t *testing.T) {
	e := emitter.New()
	var buf1, buf2 bytes.Buffer
	_ = e.Register("a", emitter.NewWriterSink(&buf1))
	_ = e.Register("b", emitter.NewWriterSink(&buf2))

	if err := e.Emit(sampleEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, buf := range []*bytes.Buffer{&buf1, &buf2} {
		if !strings.Contains(buf.String(), "latency spike detected") {
			t.Error("sink did not receive event message")
		}
	}
}

func TestDeregister_RemovesSink(t *testing.T) {
	e := emitter.New()
	var buf bytes.Buffer
	_ = e.Register("only", emitter.NewWriterSink(&buf))
	e.Deregister("only")
	if e.Len() != 0 {
		t.Fatalf("expected 0 sinks, got %d", e.Len())
	}
}

func TestEmit_SinkErrorCollected(t *testing.T) {
	e := emitter.New()
	failSink, _ := emitter.NewFuncSink(func(emitter.Event) error {
		return errors.New("sink failure")
	})
	_ = e.Register("fail", failSink)
	err := e.Emit(sampleEvent())
	if err == nil {
		t.Fatal("expected error from failing sink")
	}
	if !strings.Contains(err.Error(), "sink failure") {
		t.Errorf("unexpected error text: %v", err)
	}
}

func TestNewFuncSink_NilFuncReturnsError(t *testing.T) {
	_, err := emitter.NewFuncSink(nil)
	if err == nil {
		t.Fatal("expected error for nil func")
	}
}

func TestWriterSink_DefaultsToStderr(t *testing.T) {
	// Constructing with nil must not panic.
	sink := emitter.NewWriterSink(nil)
	if sink == nil {
		t.Fatal("expected non-nil sink")
	}
}
