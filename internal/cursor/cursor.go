// Package cursor tracks the read position within a log file, enabling
// resumable tailing across restarts.
package cursor

import (
	"errors"
	"sync"
)

// ErrNegativeOffset is returned when a negative offset is provided.
var ErrNegativeOffset = errors.New("cursor: offset must be non-negative")

// Cursor holds the current byte offset for a named file.
type Cursor struct {
	mu     sync.RWMutex
	file   string
	offset int64
}

// New creates a Cursor for the given filename starting at offset 0.
// Returns an error if file is empty.
func New(file string) (*Cursor, error) {
	if file == "" {
		return nil, errors.New("cursor: file must not be empty")
	}
	return &Cursor{file: file}, nil
}

// File returns the filename associated with this cursor.
func (c *Cursor) File() string {
	return c.file
}

// Offset returns the current byte offset.
func (c *Cursor) Offset() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.offset
}

// Set updates the cursor to the given offset.
// Returns ErrNegativeOffset if offset < 0.
func (c *Cursor) Set(offset int64) error {
	if offset < 0 {
		return ErrNegativeOffset
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.offset = offset
	return nil
}

// Advance moves the cursor forward by delta bytes.
// Returns ErrNegativeOffset if the resulting offset would be negative.
func (c *Cursor) Advance(delta int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	next := c.offset + delta
	if next < 0 {
		return ErrNegativeOffset
	}
	c.offset = next
	return nil
}

// Reset sets the cursor back to zero.
func (c *Cursor) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.offset = 0
}
