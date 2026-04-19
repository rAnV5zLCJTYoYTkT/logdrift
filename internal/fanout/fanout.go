// Package fanout broadcasts a single log entry to multiple downstream handlers.
package fanout

import (
	"errors"
	"fmt"
)

// Handler is any sink that can receive a string log entry.
type Handler interface {
	Handle(line string) error
}

// Fanout broadcasts each line to all registered handlers.
type Fanout struct {
	handlers []Handler
}

// New creates a Fanout with at least one handler.
func New(handlers ...Handler) (*Fanout, error) {
	if len(handlers) == 0 {
		return nil, errors.New("fanout: at least one handler is required")
	}
	for i, h := range handlers {
		if h == nil {
			return nil, fmt.Errorf("fanout: handler at index %d is nil", i)
		}
	}
	return &Fanout{handlers: handlers}, nil
}

// Send delivers line to every registered handler, collecting all errors.
func (f *Fanout) Send(line string) error {
	var errs []error
	for _, h := range f.handlers {
		if err := h.Handle(line); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("fanout: %d handler(s) failed: %w", len(errs), errors.Join(errs...))
}

// Len returns the number of registered handlers.
func (f *Fanout) Len() int { return len(f.handlers) }
