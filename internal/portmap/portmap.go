// Package portmap maintains a runtime mapping of port+protocol pairs to
// human-readable service descriptions composed from multiple enrichment sources.
package portmap

import (
	"fmt"
	"sync"
)

// Entry holds the composed description for a single port.
type Entry struct {
	Port     int
	Protocol string
	Name     string
	Group    string
	Tag      string
}

// String returns a short human-readable label for the entry.
func (e Entry) String() string {
	if e.Name != "" {
		return fmt.Sprintf("%d/%s (%s)", e.Port, e.Protocol, e.Name)
	}
	return fmt.Sprintf("%d/%s", e.Port, e.Protocol)
}

// Map stores enriched entries keyed by port+protocol.
type Map struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// New returns an empty Map.
func New() *Map {
	return &Map{entries: make(map[string]Entry)}
}

func key(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// Set stores or replaces the entry for the given port and protocol.
func (m *Map) Set(e Entry) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[key(e.Port, e.Protocol)] = e
}

// Get retrieves the entry for the given port and protocol.
// The second return value is false when no entry exists.
func (m *Map) Get(port int, proto string) (Entry, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.entries[key(port, proto)]
	return e, ok
}

// Delete removes the entry for the given port and protocol.
func (m *Map) Delete(port int, proto string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, key(port, proto))
}

// All returns a snapshot of all current entries.
func (m *Map) All() []Entry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Entry, 0, len(m.entries))
	for _, e := range m.entries {
		out = append(out, e)
	}
	return out
}
