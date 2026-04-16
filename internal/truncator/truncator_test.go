package truncator

import (
	"strings"
	"testing"
)

func TestNew_InvalidMaxBytes(t *testing.T) {
	_, err := New(0)
	if err != ErrInvalidMaxBytes {
		t.Fatalf("expected ErrInvalidMaxBytes, got %v", err)
	}
}

func TestNew_Valid(t *testing.T) {
	tr, err := New(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil Truncator")
	}
}

func TestCut_ShortStringUnchanged(t *testing.T) {
	tr, _ := New(20)
	out := tr.Cut("hello")
	if out != "hello" {
		t.Fatalf("expected 'hello', got %q", out)
	}
}

func TestCut_ExactLengthUnchanged(t *testing.T) {
	tr, _ := New(5)
	out := tr.Cut("hello")
	if out != "hello" {
		t.Fatalf("expected 'hello', got %q", out)
	}
}

func TestCut_TruncatesWithDefaultSuffix(t *testing.T) {
	tr, _ := New(10)
	out := tr.Cut("hello world!")
	if len(out) > 10 {
		t.Fatalf("result too long: %d", len(out))
	}
	if !strings.HasSuffix(out, "...") {
		t.Fatalf("expected suffix '...', got %q", out)
	}
}

func TestCut_CustomSuffix(t *testing.T) {
	tr, _ := New(10, WithSuffix("[cut]"))
	out := tr.Cut("hello world!")
	if len(out) > 10 {
		t.Fatalf("result too long: %d", len(out))
	}
	if !strings.HasSuffix(out, "[cut]") {
		t.Fatalf("expected suffix '[cut]', got %q", out)
	}
}

func TestCut_SuffixLargerThanMax(t *testing.T) {
	tr, _ := New(2, WithSuffix("..."))
	out := tr.Cut("hello")
	if len(out) > 2 {
		t.Fatalf("result too long: %d", len(out))
	}
}

func TestTruncated_True(t *testing.T) {
	tr, _ := New(4)
	if !tr.Truncated("hello") {
		t.Fatal("expected Truncated to return true")
	}
}

func TestTruncated_False(t *testing.T) {
	tr, _ := New(10)
	if tr.Truncated("hello") {
		t.Fatal("expected Truncated to return false")
	}
}
