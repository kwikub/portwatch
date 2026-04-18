// Package snapshot provides rolling history of port snapshots.
package snapshot

import (
	"sync"
	"time"

	"github.com/example/portwatch/internal/scanner"
)

// Entry holds a snapshot paired with its capture time.
type Entry struct {
	CapturedAt time.Time
	Snapshot   *scanner.Snapshot
}

// History keeps the last N snapshots in memory.
type History struct {
	mu      sync.Mutex
	entries []Entry
	maxSize int
}

// New creates a History that retains at most maxSize snapshots.
func New(maxSize int) *History {
	if maxSize < 1 {
		maxSize = 1
	}
	return &History{maxSize: maxSize}
}

// Add appends a snapshot to the history, evicting the oldest if full.
func (h *History) Add(snap *scanner.Snapshot) {
	h.mu.Lock()
	defer h.mu.Unlock()
	e := Entry{CapturedAt: time.Now(), Snapshot: snap}
	h.entries = append(h.entries, e)
	if len(h.entries) > h.maxSize {
		h.entries = h.entries[len(h.entries)-h.maxSize:]
	}
}

// All returns a copy of all stored entries, oldest first.
func (h *History) All() []Entry {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]Entry, len(h.entries))
	copy(out, h.entries)
	return out
}

// Latest returns the most recent entry and true, or false if empty.
func (h *History) Latest() (Entry, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.entries) == 0 {
		return Entry{}, false
	}
	return h.entries[len(h.entries)-1], true
}

// Len returns the number of stored entries.
func (h *History) Len() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.entries)
}
