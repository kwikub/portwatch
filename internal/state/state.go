package state

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Store persists the last known port snapshot to disk.
type Store struct {
	mu   sync.Mutex
	path string
}

// New creates a new Store backed by the given file path.
func New(path string) *Store {
	return &Store{path: path}
}

// Load reads the last snapshot from disk. Returns an empty snapshot if the
// file does not exist yet.
func (s *Store) Load() (scanner.Snapshot, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return scanner.Snapshot{}, nil
	}
	if err != nil {
		return scanner.Snapshot{}, err
	}

	var snap scanner.Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return scanner.Snapshot{}, err
	}
	return snap, nil
}

// Save writes the current snapshot to disk atomically.
func (s *Store) Save(snap scanner.Snapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.Marshal(snap)
	if err != nil {
		return err
	}

	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
