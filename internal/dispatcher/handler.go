package dispatcher

import (
	"fmt"
	"io"
	"os"
)

// WriterHandler is a Handler that writes each entry to an io.Writer.
type WriterHandler struct {
	w io.Writer
}

// NewWriterHandler creates a WriterHandler. Defaults to os.Stdout if w is nil.
func NewWriterHandler(w io.Writer) *WriterHandler {
	if w == nil {
		w = os.Stdout
	}
	return &WriterHandler{w: w}
}

// Handle writes the entry followed by a newline.
func (h *WriterHandler) Handle(entry string) error {
	_, err := fmt.Fprintln(h.w, entry)
	return err
}

// FuncHandler adapts a plain function to the Handler interface.
type FuncHandler struct {
	fn func(string) error
}

// NewFuncHandler wraps fn as a Handler. Returns an error if fn is nil.
func NewFuncHandler(fn func(string) error) (*FuncHandler, error) {
	if fn == nil {
		return nil, fmt.Errorf("dispatcher: func must not be nil")
	}
	return &FuncHandler{fn: fn}, nil
}

// Handle calls the wrapped function.
func (h *FuncHandler) Handle(entry string) error {
	return h.fn(entry)
}
