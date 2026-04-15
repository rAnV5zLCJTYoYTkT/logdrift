package checkpoint

import (
	"errors"
	"os"
)

// Manager wraps a Store and provides higher-level helpers used by the watcher
// to track and advance the read offset of a specific log file.
type Manager struct {
	store *Store
	file  string
}

// NewManager creates a Manager for the given log file backed by store.
func NewManager(store *Store, file string) *Manager {
	return &Manager{store: store, file: file}
}

// Resume returns the last saved offset for the managed file.
// If no checkpoint exists or the file has been replaced (inode changed),
// it returns 0 so the caller starts from the beginning.
func (m *Manager) Resume() (int64, error) {
	st, err := m.store.Load()
	if errors.Is(err, ErrNoCheckpoint) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	if st.File != m.file {
		return 0, nil
	}
	return st.Offset, nil
}

// Commit saves the current offset for the managed file.
func (m *Manager) Commit(offset int64) error {
	return m.store.Save(State{File: m.file, Offset: offset})
}

// Reset deletes the stored checkpoint, ignoring a missing-file error.
func (m *Manager) Reset() error {
	err := m.store.Delete()
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
