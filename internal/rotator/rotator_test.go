package rotator

import (
	"os"
	"testing"
	"time"
)

func TestNew_InvalidInterval(t *testing.T) {
	f := tempFile(t)
	_, err := New(f, 0)
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestNew_MissingFile(t *testing.T) {
	_, err := New("/nonexistent/path/file.log", time.Second)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestNew_Valid(t *testing.T) {
	f := tempFile(t)
	r, err := New(f, time.Millisecond*100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Interval() != time.Millisecond*100 {
		t.Fatalf("expected interval 100ms, got %v", r.Interval())
	}
}

func TestCheck_NoRotation(t *testing.T) {
	f := tempFile(t)
	r, err := New(f, time.Millisecond*50)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Append data — size grows, no rotation expected.
	if err := os.WriteFile(f, []byte("hello\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := r.Check(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestCheck_TruncationDetected(t *testing.T) {
	f := tempFile(t)
	if err := os.WriteFile(f, []byte("some initial content\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	r, err := New(f, time.Millisecond*50)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Truncate the file.
	if err := os.WriteFile(f, []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := r.Check(); err != ErrRotated {
		t.Fatalf("expected ErrRotated, got %v", err)
	}
}

func TestCheck_RotationResetsBaseline(t *testing.T) {
	f := tempFile(t)
	if err := os.WriteFile(f, []byte("data\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	r, err := New(f, time.Millisecond*50)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Truncate to trigger rotation.
	_ = os.WriteFile(f, []byte(""), 0o644)
	_ = r.Check() // consumes ErrRotated and resets baseline
	// Next check with no further change should be clean.
	if err := r.Check(); err != nil {
		t.Fatalf("expected nil after baseline reset, got %v", err)
	}
}

// tempFile creates a temporary file and returns its path.
func tempFile(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "rotator-*.log")
	if err != nil {
		t.Fatal(err)
	}
	_ = f.Close()
	return f.Name()
}
