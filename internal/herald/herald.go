// Package herald provides a multi-channel alert dispatcher that fans out
// anomaly notifications to one or more named sinks (e.g. stderr, webhook, file).
package herald

import (
	"errors"
	"fmt"
	"sync"
)

// Sink is a destination that can receive a formatted alert message.
type Sink interface {
	Send(msg string) error
}

// Herald dispatches messages to a registered set of named sinks.
type Herald struct {
	mu    sync.RWMutex
	sinks map[string]Sink
}

// New returns a Herald with no sinks registered.
func New() *Herald {
	return &Herald{sinks: make(map[string]Sink)}
}

// Register adds a named sink. Returns an error if name is empty or sink is nil.
func (h *Herald) Register(name string, s Sink) error {
	if name == "" {
		return errors.New("herald: sink name must not be empty")
	}
	if s == nil {
		return errors.New("herald: sink must not be nil")
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.sinks[name] = s
	return nil
}

// Unregister removes a named sink. It is a no-op if the name is not found.
func (h *Herald) Unregister(name string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.sinks, name)
}

// Dispatch sends msg to every registered sink, collecting any errors.
func (h *Herald) Dispatch(msg string) error {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if len(h.sinks) == 0 {
		return errors.New("herald: no sinks registered")
	}
	var errs []error
	for name, s := range h.sinks {
		if err := s.Send(msg); err != nil {
			errs = append(errs, fmt.Errorf("herald: sink %q: %w", name, err))
		}
	}
	if len(errs) == 1 {
		return errs[0]
	}
	if len(errs) > 1 {
		return fmt.Errorf("herald: %d sink(s) failed: %v", len(errs), errs)
	}
	return nil
}

// Len returns the number of currently registered sinks.
func (h *Herald) Len() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.sinks)
}
