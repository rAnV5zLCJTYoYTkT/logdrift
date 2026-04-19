// Package tracer tracks request trace IDs across log entries,
// grouping related lines by a shared trace or request identifier.
package tracer

import (
	"errors"
	"sync"
	"time"
)

// Entry holds a single log line associated with a trace.
type Entry struct {
	TraceID string
	Message string
	Timestamp time.Time
}

// Trace is a collection of entries sharing the same trace ID.
type Trace struct {
	ID      string
	Entries []Entry
	First   time.Time
	Last    time.Time
}

// Tracer groups log entries by trace ID and evicts stale traces.
type Tracer struct {
	mu     sync.Mutex
	traces map[string]*Trace
	ttl    time.Duration
}

// New creates a Tracer with the given TTL for trace eviction.
func New(ttl time.Duration) (*Tracer, error) {
	if ttl <= 0 {
		return nil, errors.New("tracer: ttl must be positive")
	}
	return &Tracer{
		traces: make(map[string]*Trace),
		ttl:    ttl,
	}, nil
}

// Add appends a log entry to the trace identified by traceID.
func (t *Tracer) Add(traceID, message string, ts time.Time) error {
	if traceID == "" {
		return errors.New("tracer: traceID must not be empty")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	tr, ok := t.traces[traceID]
	if !ok {
		tr = &Trace{ID: traceID, First: ts}
		t.traces[traceID] = tr
	}
	tr.Entries = append(tr.Entries, Entry{TraceID: traceID, Message: message, Timestamp: ts})
	tr.Last = ts
	return nil
}

// Get returns the trace for the given ID, or false if not found.
func (t *Tracer) Get(traceID string) (Trace, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	tr, ok := t.traces[traceID]
	if !ok {
		return Trace{}, false
	}
	return *tr, true
}

// Evict removes traces whose last entry is older than the TTL.
func (t *Tracer) Evict(now time.Time) []Trace {
	t.mu.Lock()
	defer t.mu.Unlock()
	var evicted []Trace
	for id, tr := range t.traces {
		if now.Sub(tr.Last) >= t.ttl {
			evicted = append(evicted, *tr)
			delete(t.traces, id)
		}
	}
	return evicted
}

// Len returns the number of active traces.
func (t *Tracer) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.traces)
}
