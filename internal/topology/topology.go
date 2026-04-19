// Package topology tracks service-to-service call relationships observed
// in log lines, building a lightweight dependency graph at runtime.
package topology

import (
	"errors"
	"sync"
)

// Edge represents a directed call from one service to another.
type Edge struct {
	Caller string
	Callee string
	Count  int64
}

// Graph holds observed service dependencies.
type Graph struct {
	mu    sync.RWMutex
	edges map[string]map[string]int64
}

// New returns an empty Graph.
func New() *Graph {
	return &Graph{
		edges: make(map[string]map[string]int64),
	}
}

// Record registers a call from caller to callee.
func (g *Graph) Record(caller, callee string) error {
	if caller == "" {
		return errors.New("topology: caller must not be empty")
	}
	if callee == "" {
		return errors.New("topology: callee must not be empty")
	}
	if caller == callee {
		return errors.New("topology: caller and callee must differ")
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.edges[caller] == nil {
		g.edges[caller] = make(map[string]int64)
	}
	g.edges[caller][callee]++
	return nil
}

// Edges returns a snapshot of all recorded edges sorted by caller name.
func (g *Graph) Edges() []Edge {
	g.mu.RLock()
	defer g.mu.RUnlock()
	var out []Edge
	for caller, callees := range g.edges {
		for callee, count := range callees {
			out = append(out, Edge{Caller: caller, Callee: callee, Count: count})
		}
	}
	return out
}

// Callees returns all services that the given caller has been observed calling.
func (g *Graph) Callees(caller string) []string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	var out []string
	for callee := range g.edges[caller]{
		out = append(out, callee)
	}
	return out
}

// Reset clears all recorded edges.
func (g *Graph) Reset() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.edges = make(map[string]map[string]int64)
}
