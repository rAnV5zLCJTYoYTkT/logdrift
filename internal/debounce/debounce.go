// Package debounce provides a mechanism to suppress rapid repeated events,
// emitting only the final event after a quiet period has elapsed.
package debounce

import (
	"sync"
	"time"
)

// Debouncer delays forwarding of values until no new values have arrived
// for the configured wait duration. It is safe for concurrent use.
type Debouncer struct {
	wait  time.Duration
	mu    sync.Mutex
	timers map[string]*time.Timer
}

// New creates a Debouncer that waits for the given duration of inactivity
// before invoking the callback. It returns an error if wait is zero or negative.
func New(wait time.Duration) (*Debouncer, error) {
	if wait <= 0 {
		return nil, ErrInvalidWait
	}
	return &Debouncer{
		wait:   wait,
		timers: make(map[string]*time.Timer),
	}, nil
}

// Submit schedules fn to be called after the debounce wait period for the
// given key. If Submit is called again for the same key before the timer
// fires, the previous timer is cancelled and the wait period resets.
func (d *Debouncer) Submit(key string, fn func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[key]; ok {
		t.Stop()
	}

	d.timers[key] = time.AfterFunc(d.wait, func() {
		d.mu.Lock()
		delete(d.timers, key)
		d.mu.Unlock()
		fn()
	})
}

// Flush cancels any pending timer for the given key and invokes the callback
// immediately. If no timer is pending, Flush is a no-op.
func (d *Debouncer) Flush(key string, fn func()) {
	d.mu.Lock()
	t, ok := d.timers[key]
	if ok {
		t.Stop()
		delete(d.timers, key)
	}
	d.mu.Unlock()

	if ok {
		fn()
	}
}

// Pending returns the number of keys with active pending timers.
func (d *Debouncer) Pending() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.timers)
}
