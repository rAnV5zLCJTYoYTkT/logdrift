// Package summarizer aggregates log entries into a human-readable summary
// grouped by severity level and service name.
package summarizer

import (
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
)

// Entry holds the fields summarizer cares about.
type Entry struct {
	Level   string
	Service string
	Message string
}

// Summary holds aggregated counts.
type Summary struct {
	ByLevel   map[string]int
	ByService map[string]int
	Total     int
}

// Summarizer collects entries and can flush a summary.
type Summarizer struct {
	mu        sync.Mutex
	byLevel   map[string]int
	byService map[string]int
	total     int
	out       io.Writer
}

// New returns a new Summarizer writing to out. If out is nil, os.Stdout is used.
func New(out io.Writer) *Summarizer {
	if out == nil {
		out = os.Stdout
	}
	return &Summarizer{
		byLevel:   make(map[string]int),
		byService: make(map[string]int),
		out:       out,
	}
}

// Add records a single log entry.
func (s *Summarizer) Add(e Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.byLevel[e.Level]++
	s.byService[e.Service]++
	s.total++
}

// Snapshot returns a copy of the current summary without resetting state.
func (s *Summarizer) Snapshot() Summary {
	s.mu.Lock()
	defer s.mu.Unlock()
	snap := Summary{
		ByLevel:   make(map[string]int, len(s.byLevel)),
		ByService: make(map[string]int, len(s.byService)),
		Total:     s.total,
	}
	for k, v := range s.byLevel {
		snap.ByLevel[k] = v
	}
	for k, v := range s.byService {
		snap.ByService[k] = v
	}
	return snap
}

// Flush writes the summary to the configured writer and resets all counters.
func (s *Summarizer) Flush() {
	snap := s.Snapshot()
	s.mu.Lock()
	s.byLevel = make(map[string]int)
	s.byService = make(map[string]int)
	s.total = 0
	s.mu.Unlock()

	fmt.Fprintf(s.out, "total=%d\n", snap.Total)

	levels := sortedKeys(snap.ByLevel)
	for _, l := range levels {
		fmt.Fprintf(s.out, "level.%s=%d\n", l, snap.ByLevel[l])
	}

	services := sortedKeys(snap.ByService)
	for _, svc := range services {
		fmt.Fprintf(s.out, "service.%s=%d\n", svc, snap.ByService[svc])
	}
}

func sortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
