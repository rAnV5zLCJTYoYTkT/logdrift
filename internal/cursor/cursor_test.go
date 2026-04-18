package cursor

import (
	"testing"
)

func TestNew_EmptyFileReturnsError(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for empty file")
	}
}

func TestNew_ValidFile(t *testing.T) {
	c, err := New("app.log")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.File() != "app.log" {
		t.Errorf("expected app.log, got %s", c.File())
	}
	if c.Offset() != 0 {
		t.Errorf("expected offset 0, got %d", c.Offset())
	}
}

func TestSet_UpdatesOffset(t *testing.T) {
	c, _ := New("app.log")
	if err := c.Set(128); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Offset() != 128 {
		t.Errorf("expected 128, got %d", c.Offset())
	}
}

func TestSet_NegativeOffsetReturnsError(t *testing.T) {
	c, _ := New("app.log")
	if err := c.Set(-1); err != ErrNegativeOffset {
		t.Errorf("expected ErrNegativeOffset, got %v", err)
	}
}

func TestAdvance_MovesForward(t *testing.T) {
	c, _ := New("app.log")
	_ = c.Set(100)
	if err := c.Advance(50); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Offset() != 150 {
		t.Errorf("expected 150, got %d", c.Offset())
	}
}

func TestAdvance_NegativeDeltaBelowZeroReturnsError(t *testing.T) {
	c, _ := New("app.log")
	_ = c.Set(10)
	if err := c.Advance(-20); err != ErrNegativeOffset {
		t.Errorf("expected ErrNegativeOffset, got %v", err)
	}
	if c.Offset() != 10 {
		t.Errorf("offset should be unchanged: got %d", c.Offset())
	}
}

func TestReset_SetsOffsetToZero(t *testing.T) {
	c, _ := New("app.log")
	_ = c.Set(999)
	c.Reset()
	if c.Offset() != 0 {
		t.Errorf("expected 0 after reset, got %d", c.Offset())
	}
}
