// Package portmemo stores arbitrary key-value annotations against a port+protocol
// pair, allowing other subsystems to attach metadata without coupling.
package portmemo

import (
	"fmt"
	"sync"
)

// Memo holds annotations keyed by port+protocol.
type Memo struct {
	mu      sync.RWMutex
	entries map[string]map[string]string
}

// New returns an empty Memo.
func New() *Memo {
	return &Memo{
		entries: make(map[string]map[string]string),
	}
}

func portKey(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// Set stores a value under the given annotation key for a port+protocol pair.
// An empty value removes the annotation key.
func (m *Memo) Set(port int, proto, key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	pk := portKey(port, proto)
	if _, ok := m.entries[pk]; !ok {
		m.entries[pk] = make(map[string]string)
	}
	if value == "" {
		delete(m.entries[pk], key)
		if len(m.entries[pk]) == 0 {
			delete(m.entries, pk)
		}
		return
	}
	m.entries[pk][key] = value
}

// Get retrieves an annotation value. The second return value reports whether
// the key was present.
func (m *Memo) Get(port int, proto, key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	pk := portKey(port, proto)
	ann, ok := m.entries[pk]
	if !ok {
		return "", false
	}
	v, ok := ann[key]
	return v, ok
}

// All returns a copy of every annotation for the given port+protocol pair.
// Returns nil if no annotations exist.
func (m *Memo) All(port int, proto string) map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	pk := portKey(port, proto)
	ann, ok := m.entries[pk]
	if !ok {
		return nil
	}
	out := make(map[string]string, len(ann))
	for k, v := range ann {
		out[k] = v
	}
	return out
}

// Clear removes all annotations for the given port+protocol pair.
func (m *Memo) Clear(port int, proto string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, portKey(port, proto))
}
