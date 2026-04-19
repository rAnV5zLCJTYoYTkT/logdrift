package topology

import (
	"testing"
)

func TestRecord_ValidEdge(t *testing.T) {
	g := New()
	if err := g.Record("api", "db"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	edges := g.Edges()
	if len(edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(edges))
	}
	if edges[0].Caller != "api" || edges[0].Callee != "db" || edges[0].Count != 1 {
		t.Errorf("unexpected edge: %+v", edges[0])
	}
}

func TestRecord_EmptyCallerReturnsError(t *testing.T) {
	g := New()
	if err := g.Record("", "db"); err == nil {
		t.Fatal("expected error for empty caller")
	}
}

func TestRecord_EmptyCalleeReturnsError(t *testing.T) {
	g := New()
	if err := g.Record("api", ""); err == nil {
		t.Fatal("expected error for empty callee")
	}
}

func TestRecord_SameCallerCalleeReturnsError(t *testing.T) {
	g := New()
	if err := g.Record("api", "api"); err == nil {
		t.Fatal("expected error when caller equals callee")
	}
}

func TestRecord_CountAccumulates(t *testing.T) {
	g := New()
	for i := 0; i < 5; i++ {
		_ = g.Record("api", "cache")
	}
	edges := g.Edges()
	if len(edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(edges))
	}
	if edges[0].Count != 5 {
		t.Errorf("expected count 5, got %d", edges[0].Count)
	}
}

func TestCallees_ReturnsDownstream(t *testing.T) {
	g := New()
	_ = g.Record("api", "db")
	_ = g.Record("api", "cache")
	callees := g.Callees("api")
	if len(callees) != 2 {
		t.Errorf("expected 2 callees, got %d", len(callees))
	}
}

func TestCallees_UnknownCallerReturnsEmpty(t *testing.T) {
	g := New()
	if c := g.Callees("unknown"); len(c) != 0 {
		t.Errorf("expected empty, got %v", c)
	}
}

func TestReset_ClearsGraph(t *testing.T) {
	g := New()
	_ = g.Record("api", "db")
	g.Reset()
	if edges := g.Edges(); len(edges) != 0 {
		t.Errorf("expected empty graph after reset, got %d edges", len(edges))
	}
}
