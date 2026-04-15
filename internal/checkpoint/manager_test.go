package checkpoint_test

import (
	"testing"

	"github.com/user/logdrift/internal/checkpoint"
)

func TestResume_NoCheckpointReturnsZero(t *testing.T) {
	store := checkpoint.New(tempPath(t))
	m := checkpoint.NewManager(store, "/var/log/app.log")

	offset, err := m.Resume()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if offset != 0 {
		t.Errorf("expected 0, got %d", offset)
	}
}

func TestResume_ReturnsSavedOffset(t *testing.T) {
	store := checkpoint.New(tempPath(t))
	m := checkpoint.NewManager(store, "/var/log/app.log")

	if err := m.Commit(8192); err != nil {
		t.Fatalf("Commit: %v", err)
	}

	offset, err := m.Resume()
	if err != nil {
		t.Fatalf("Resume: %v", err)
	}
	if offset != 8192 {
		t.Errorf("expected 8192, got %d", offset)
	}
}

func TestResume_DifferentFilenameReturnsZero(t *testing.T) {
	store := checkpoint.New(tempPath(t))
	m1 := checkpoint.NewManager(store, "/var/log/app.log")
	m2 := checkpoint.NewManager(store, "/var/log/other.log")

	_ = m1.Commit(512)

	offset, err := m2.Resume()
	if err != nil {
		t.Fatalf("Resume: %v", err)
	}
	if offset != 0 {
		t.Errorf("expected 0 for different file, got %d", offset)
	}
}

func TestReset_ClearsCheckpoint(t *testing.T) {
	store := checkpoint.New(tempPath(t))
	m := checkpoint.NewManager(store, "/var/log/app.log")

	_ = m.Commit(1024)
	if err := m.Reset(); err != nil {
		t.Fatalf("Reset: %v", err)
	}

	offset, err := m.Resume()
	if err != nil {
		t.Fatalf("Resume after reset: %v", err)
	}
	if offset != 0 {
		t.Errorf("expected 0 after reset, got %d", offset)
	}
}

func TestReset_NoCheckpointIsNoOp(t *testing.T) {
	store := checkpoint.New(tempPath(t))
	m := checkpoint.NewManager(store, "/var/log/app.log")

	if err := m.Reset(); err != nil {
		t.Errorf("Reset on missing checkpoint should not error: %v", err)
	}
}
