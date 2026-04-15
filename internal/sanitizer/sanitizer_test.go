package sanitizer_test

import (
	"strings"
	"testing"

	"github.com/user/logdrift/internal/sanitizer"
)

func TestClean_TrimsSurroundingWhitespace(t *testing.T) {
	s := sanitizer.New()
	got := s.Clean("  hello world  ")
	if got != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", got)
	}
}

func TestClean_EmptyStringReturnsEmpty(t *testing.T) {
	s := sanitizer.New()
	if got := s.Clean(""); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestClean_WhitespaceOnlyReturnsEmpty(t *testing.T) {
	s := sanitizer.New()
	if got := s.Clean("   \t  "); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestClean_RemovesNonPrintableRunes(t *testing.T) {
	s := sanitizer.New()
	// embed a null byte and a BEL character
	input := "hello\x00world\x07"
	got := s.Clean(input)
	if got != "helloworld" {
		t.Fatalf("expected %q, got %q", "helloworld", got)
	}
}

func TestClean_PreservesTabCharacter(t *testing.T) {
	s := sanitizer.New()
	got := s.Clean("col1\tcol2")
	if got != "col1\tcol2" {
		t.Fatalf("expected tab to be preserved, got %q", got)
	}
}

func TestClean_TruncatesLongLine(t *testing.T) {
	s := sanitizer.New(sanitizer.WithMaxLen(10))
	input := strings.Repeat("a", 50)
	got := s.Clean(input)
	if len(got) != 10 {
		t.Fatalf("expected length 10, got %d", len(got))
	}
}

func TestClean_DefaultMaxLenNotExceeded(t *testing.T) {
	s := sanitizer.New()
	// 4096 chars should pass through unchanged
	input := strings.Repeat("x", 4096)
	got := s.Clean(input)
	if len(got) != 4096 {
		t.Fatalf("expected 4096, got %d", len(got))
	}
}

func TestWithMaxLen_ZeroValueKeepsDefault(t *testing.T) {
	s := sanitizer.New(sanitizer.WithMaxLen(0))
	input := strings.Repeat("b", 4096)
	got := s.Clean(input)
	if len(got) != 4096 {
		t.Fatalf("zero maxLen should keep default 4096, got %d", len(got))
	}
}
