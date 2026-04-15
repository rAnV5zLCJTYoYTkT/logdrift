// Package watcher provides tail-like file watching for logdrift,
// emitting new lines as they are appended to a log file.
package watcher

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"
)

// DefaultPollInterval is how often the watcher checks for new data.
const DefaultPollInterval = 200 * time.Millisecond

// Watcher tails a file, sending new lines to Lines.
type Watcher struct {
	path         string
	pollInterval time.Duration
	Lines        chan string
}

// New creates a Watcher for the given file path.
// pollInterval controls how frequently the file is polled;
// pass 0 to use DefaultPollInterval.
func New(path string, pollInterval time.Duration) *Watcher {
	if pollInterval <= 0 {
		pollInterval = DefaultPollInterval
	}
	return &Watcher{
		path:         path,
		pollInterval: pollInterval,
		Lines:        make(chan string, 64),
	}
}

// Run opens the file, seeks to the end, and begins tailing.
// It sends each new line to Lines until ctx is cancelled.
// Run closes Lines when it returns.
func (w *Watcher) Run(ctx context.Context) error {
	f, err := os.Open(w.path)
	if err != nil {
		return err
	}
	defer func() {
		f.Close()
		close(w.Lines)
	}()

	// Seek to end so we only tail new content.
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		return err
	}

	reader := bufio.NewReader(f)
	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			for {
				line, err := reader.ReadString('\n')
				if len(line) > 0 {
					// Strip trailing newline before sending.
					if len(line) > 0 && line[len(line)-1] == '\n' {
						line = line[:len(line)-1]
					}
					select {
					case w.Lines <- line:
					case <-ctx.Done():
						return ctx.Err()
					}
				}
				if err == io.EOF {
					break
				}
				if err != nil {
					return err
				}
			}
		}
	}
}
