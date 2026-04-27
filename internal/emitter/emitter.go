// Package emitter provides a structured event emitter that broadcasts
// anomaly detection results to one or more registered sinks.
package emitter

import (
	"errors"
	"sync"
	"time"
)

// Event represents a single anomaly or noteworthy log event emitted by the pipeline.
type Event struct {
	Timestamp time.Time
	Service   string
	Level     string
	Message   string
	Score     float64
	Tags      []string
}

// Sink is a destination that receives emitted events.
type Sink interface {
	Accept(Event) error
}

// Emitter broadcasts events to all registered sinks.
type Emitter struct {
	mu    sync.RWMutex
	sinks map[string]Sink
}

// New returns an initialised Emitter with no sinks registered.
func New() *Emitter {
	return &Emitter{sinks: make(map[string]Sink)}
}

// Register adds a named sink. Returns an error if name is empty or sink is nil.
func (e *Emitter) Register(name string, s Sink) error {
	if name == "" {
		return errors.New("emitter: sink name must not be empty")
	}
	if s == nil {
		return errors.New("emitter: sink must not be nil")
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.sinks[name] = s
	return nil
}

// Deregister removes a previously registered sink by name.
func (e *Emitter) Deregister(name string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.sinks, name)
}

// Emit delivers ev to every registered sink. All sinks are called even if one
// returns an error; a combined error is returned when any sink fails.
func (e *Emitter) Emit(ev Event) error {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if len(e.sinks) == 0 {
		return errors.New("emitter: no sinks registered")
	}
	var errs []error
	for _, s := range e.sinks {
		if err := s.Accept(ev); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// Len returns the number of registered sinks.
func (e *Emitter) Len() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.sinks)
}
