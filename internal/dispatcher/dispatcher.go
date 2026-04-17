// Package dispatcher fans out log entries to multiple handlers concurrently.
package dispatcher

import (
	"context"
	"fmt"
	"sync"
)

// Handler processes a single log entry string.
type Handler interface {
	Handle(entry string) error
}

// Dispatcher sends each entry to all registered handlers.
type Dispatcher struct {
	handlers []Handler
	mu       sync.RWMutex
}

// ErrNoHandlers is returned when Dispatch is called with no handlers registered.
var ErrNoHandlers = fmt.Errorf("dispatcher: no handlers registered")

// New creates an empty Dispatcher.
func New() *Dispatcher {
	return &Dispatcher{}
}

// Register adds a handler. Returns an error if handler is nil.
func (d *Dispatcher) Register(h Handler) error {
	if h == nil {
		return fmt.Errorf("dispatcher: handler must not be nil")
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers = append(d.handlers, h)
	return nil
}

// Dispatch sends entry to all handlers concurrently and collects errors.
// Cancelling ctx stops waiting for remaining handlers.
func (d *Dispatcher) Dispatch(ctx context.Context, entry string) []error {
	d.mu.RLock()
	handlers := make([]Handler, len(d.handlers))
	copy(handlers, d.handlers)
	d.mu.RUnlock()

	if len(handlers) == 0 {
		return []error{ErrNoHandlers}
	}

	type result struct {
		err error
	}
	results := make(chan result, len(handlers))
	var wg sync.WaitGroup

	for _, h := range handlers {
		wg.Add(1)
		go func(h Handler) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				results <- result{ctx.Err()}
			default:
				results <- result{h.Handle(entry)}
			}
		}(h)
	}

	wg.Wait()
	close(results)

	var errs []error
	for r := range results {
		if r.err != nil {
			errs = append(errs, r.err)
		}
	}
	return errs
}
