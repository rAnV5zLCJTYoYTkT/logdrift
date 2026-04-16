// Package batcher groups log entries into fixed-size or time-bounded batches
// before forwarding them downstream, reducing per-entry processing overhead.
package batcher

import (
	"errors"
	"sync"
	"time"
)

// ErrInvalidSize is returned when the batch size is not positive.
var ErrInvalidSize = errors.New("batcher: size must be greater than zero")

// ErrInvalidTimeout is returned when the flush timeout is not positive.
var ErrInvalidTimeout = errors.New("batcher: timeout must be greater than zero")

// Batcher accumulates items and flushes them when the batch is full or the
// timeout elapses.
type Batcher struct {
	mu      sync.Mutex
	items   []string
	size    int
	flushFn func([]string)
	timer   *time.Timer
	timeout time.Duration
}

// New creates a Batcher that calls flushFn whenever a batch is ready.
func New(size int, timeout time.Duration, flushFn func([]string)) (*Batcher, error) {
	if size <= 0 {
		return nil, ErrInvalidSize
	}
	if timeout <= 0 {
		return nil, ErrInvalidTimeout
	}
	b := &Batcher{
		items:   make([]string, 0, size),
		size:    size,
		flushFn: flushFn,
		timeout: timeout,
	}
	b.resetTimer()
	return b, nil
}

// Add appends an item to the current batch and flushes if the batch is full.
func (b *Batcher) Add(item string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.items = append(b.items, item)
	if len(b.items) >= b.size {
		b.flush()
	}
}

// Flush forces the current batch to be emitted regardless of size.
func (b *Batcher) Flush() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.flush()
}

// flush emits the batch and resets state. Caller must hold mu.
func (b *Batcher) flush() {
	if len(b.items) == 0 {
		return
	}
	batch := make([]string, len(b.items))
	copy(batch, b.items)
	b.items = b.items[:0]
	b.resetTimer()
	go b.flushFn(batch)
}

func (b *Batcher) resetTimer() {
	if b.timer != nil {
		b.timer.Stop()
	}
	b.timer = time.AfterFunc(b.timeout, func() {
		b.Flush()
	})
}
