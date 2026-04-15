// Package metrics provides counters and gauges for tracking
// runtime statistics about log processing in logdrift.
package metrics

import (
	"sync"
	"sync/atomic"
)

// Counter is a monotonically increasing integer counter.
type Counter struct {
	val uint64
}

// Inc increments the counter by 1.
func (c *Counter) Inc() { atomic.AddUint64(&c.val, 1) }

// Add increments the counter by n.
func (c *Counter) Add(n uint64) { atomic.AddUint64(&c.val, n) }

// Value returns the current counter value.
func (c *Counter) Value() uint64 { return atomic.LoadUint64(&c.val) }

// Reset sets the counter back to zero.
func (c *Counter) Reset() { atomic.StoreUint64(&c.val, 0) }

// Registry holds a named set of counters.
type Registry struct {
	mu       sync.RWMutex
	counters map[string]*Counter
}

// NewRegistry creates an empty Registry.
func NewRegistry() *Registry {
	return &Registry{counters: make(map[string]*Counter)}
}

// Counter returns the named counter, creating it if necessary.
func (r *Registry) Counter(name string) *Counter {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c, ok := r.counters[name]; ok {
		return c
	}
	c := &Counter{}
	r.counters[name] = c
	return c
}

// Snapshot returns a copy of all counter values keyed by name.
func (r *Registry) Snapshot() map[string]uint64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]uint64, len(r.counters))
	for name, c := range r.counters {
		out[name] = c.Value()
	}
	return out
}
