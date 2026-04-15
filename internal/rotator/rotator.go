// Package rotator provides log file rotation detection, signalling when
// a watched file has been replaced or truncated beneath the reader.
package rotator

import (
	"errors"
	"os"
	"sync"
	"time"
)

// ErrRotated is returned when the underlying file has been rotated or truncated.
var ErrRotated = errors.New("rotator: file rotated or truncated")

// Rotator watches a file path and detects rotation (inode change) or
// truncation (size decrease). Call Check periodically; it returns
// ErrRotated when the file should be re-opened.
type Rotator struct {
	mu       sync.Mutex
	path     string
	inode    uint64
	size     int64
	interval time.Duration
}

// New creates a Rotator for the given path and poll interval.
// interval must be positive.
func New(path string, interval time.Duration) (*Rotator, error) {
	if interval <= 0 {
		return nil, errors.New("rotator: interval must be positive")
	}
	r := &Rotator{path: path, interval: interval}
	if err := r.snapshot(); err != nil {
		return nil, err
	}
	return r, nil
}

// Interval returns the configured poll interval.
func (r *Rotator) Interval() time.Duration {
	return r.interval
}

// Check inspects the file and returns ErrRotated if rotation or truncation
// is detected. On rotation the internal baseline is updated automatically
// so subsequent calls reflect the new file.
func (r *Rotator) Check() error {
	info, err := os.Stat(r.path)
	if err != nil {
		return err
	}
	inode := inoOf(info)
	size := info.Size()

	r.mu.Lock()
	defer r.mu.Unlock()

	if inode != r.inode || size < r.size {
		r.inode = inode
		r.size = size
		return ErrRotated
	}
	r.size = size
	return nil
}

func (r *Rotator) snapshot() error {
	info, err := os.Stat(r.path)
	if err != nil {
		return err
	}
	r.inode = inoOf(info)
	r.size = info.Size()
	return nil
}
