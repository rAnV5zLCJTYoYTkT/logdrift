package buffer

import (
	"testing"
)

func TestNew_InvalidCapacity(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for capacity 0")
	}
	_, err = New(-1)
	if err == nil {
		t.Fatal("expected error for negative capacity")
	}
}

func TestAdd_And_Len(t *testing.T) {
	b, _ := New(3)
	if b.Len() != 0 {
		t.Fatalf("expected 0, got %d", b.Len())
	}
	b.Add("a")
	b.Add("b")
	if b.Len() != 2 {
		t.Fatalf("expected 2, got %d", b.Len())
	}
}

func TestSnapshot_Order(t *testing.T) {
	b, _ := New(3)
	b.Add("x")
	b.Add("y")
	b.Add("z")
	snap := b.Snapshot()
	if len(snap) != 3 || snap[0] != "x" || snap[2] != "z" {
		t.Fatalf("unexpected snapshot: %v", snap)
	}
}

func TestSnapshot_Eviction(t *testing.T) {
	b, _ := New(3)
	b.Add("a")
	b.Add("b")
	b.Add("c")
	b.Add("d") // evicts "a"
	snap := b.Snapshot()
	if len(snap) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(snap))
	}
	if snap[0] != "b" || snap[1] != "c" || snap[2] != "d" {
		t.Fatalf("unexpected snapshot after eviction: %v", snap)
	}
}

func TestCap(t *testing.T) {
	b, _ := New(5)
	if b.Cap() != 5 {
		t.Fatalf("expected cap 5, got %d", b.Cap())
	}
}

func TestSnapshot_Empty(t *testing.T) {
	b, _ := New(4)
	snap := b.Snapshot()
	if len(snap) != 0 {
		t.Fatalf("expected empty snapshot, got %v", snap)
	}
}
