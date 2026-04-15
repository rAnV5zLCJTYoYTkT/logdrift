// Package checkpoint persists and restores the last read offset of a log file,
// allowing logdrift to resume tailing after a restart without reprocessing old lines.
package checkpoint

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

// ErrNoCheckpoint is returned when no checkpoint file exists yet.
var ErrNoCheckpoint = errors.New("checkpoint: no checkpoint found")

// State holds the persisted position within a log file.
type State struct {
	File   string `json:"file"`
	Offset int64  `json:"offset"`
}

// Store saves and loads checkpoint state to/from a JSON file.
type Store struct {
	mu   sync.Mutex
	path string
}

// New creates a new Store that persists state to the given path.
func New(path string) *Store {
	return &Store{path: path}
}

// Save atomically writes the checkpoint state to disk.
func (s *Store) Save(state State) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

// Load reads the last saved checkpoint state from disk.
// Returns ErrNoCheckpoint if the file does not exist.
func (s *Store) Load() (State, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return State{}, ErrNoCheckpoint
	}
	if err != nil {
		return State{}, err
	}

	var st State
	if err := json.Unmarshal(data, &st); err != nil {
		return State{}, err
	}
	return st, nil
}

// Delete removes the checkpoint file from disk.
func (s *Store) Delete() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return os.Remove(s.path)
}
