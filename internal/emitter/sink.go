package emitter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// WriterSink writes each event as a JSON line to an io.Writer.
type WriterSink struct {
	w io.Writer
}

// NewWriterSink returns a WriterSink that writes to w.
// If w is nil it defaults to os.Stderr.
func NewWriterSink(w io.Writer) *WriterSink {
	if w == nil {
		w = os.Stderr
	}
	return &WriterSink{w: w}
}

// Accept serialises ev as JSON and writes it followed by a newline.
func (s *WriterSink) Accept(ev Event) error {
	b, err := json.Marshal(ev)
	if err != nil {
		return fmt.Errorf("emitter: marshal event: %w", err)
	}
	_, err = fmt.Fprintf(s.w, "%s\n", b)
	return err
}

// FuncSink wraps a plain function so it satisfies the Sink interface.
type FuncSink struct {
	fn func(Event) error
}

// NewFuncSink returns a FuncSink backed by fn.
// Returns an error if fn is nil.
func NewFuncSink(fn func(Event) error) (*FuncSink, error) {
	if fn == nil {
		return nil, fmt.Errorf("emitter: func must not be nil")
	}
	return &FuncSink{fn: fn}, nil
}

// Accept calls the wrapped function with ev.
func (s *FuncSink) Accept(ev Event) error {
	return s.fn(ev)
}
