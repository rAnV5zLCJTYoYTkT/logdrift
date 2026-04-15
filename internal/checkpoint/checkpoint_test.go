package checkpoint_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/logdrift/internal/checkpoint"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestLoad_NoCheckpointReturnsError(t *testing.T) {
	s := checkpoint.New(tempPath(t))
	_, err := s.Load()
	if !errors.Is(err, checkpoint.ErrNoCheckpoint) {
		t.Fatalf("expected ErrNoCheckpoint, got %v", err)
	}
}

func TestSave_ThenLoad_RoundTrip(t *testing.T) {
	s := checkpoint.New(tempPath(t))
	want := checkpoint.State{File: "/var/log/app.log", Offset: 4096}

	if err := s.Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestSave_OverwritesPreviousState(t *testing.T) {
	s := checkpoint.New(tempPath(t))
	_ = s.Save(checkpoint.State{File: "a.log", Offset: 100})
	_ = s.Save(checkpoint.State{File: "a.log", Offset: 200})

	got, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Offset != 200 {
		t.Errorf("expected offset 200, got %d", got.Offset)
	}
}

func TestDelete_RemovesFile(t *testing.T) {
	p := tempPath(t)
	s := checkpoint.New(p)
	_ = s.Save(checkpoint.State{File: "x.log", Offset: 1})

	if err := s.Delete(); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := os.Stat(p); !errors.Is(err, os.ErrNotExist) {
		t.Error("expected file to be removed")
	}
}

func TestDelete_MissingFileReturnsError(t *testing.T) {
	s := checkpoint.New(tempPath(t))
	if err := s.Delete(); err == nil {
		t.Error("expected error when deleting non-existent checkpoint")
	}
}

func TestLoad_CorruptedFileReturnsError(t *testing.T) {
	p := tempPath(t)
	_ = os.WriteFile(p, []byte("not-json"), 0o644)
	s := checkpoint.New(p)
	_, err := s.Load()
	if err == nil {
		t.Error("expected error for corrupted checkpoint file")
	}
}
