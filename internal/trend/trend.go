// Package trend tracks port change frequency over a sliding window.
package trend

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry records a single change event.
type Entry struct {
	At    time.Time
	State string // "opened" | "closed"
}

// Tracker accumulates change events per port key and reports frequency.
type Tracker struct {
	mu     sync.Mutex
	window time.Duration
	events map[string][]Entry
}

// New returns a Tracker with the given sliding window duration.
func New(window time.Duration) *Tracker {
	if window <= 0 {
		panic("trend: window must be positive")
	}
	return &Tracker{
		window: window,
		events: make(map[string][]Entry),
	}
}

// Record adds a diff's changes to the tracker.
func (t *Tracker) Record(diffs []scanner.Diff) {
	now := time.Now()
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, d := range diffs {
		k := key(d)
		t.events[k] = append(t.events[k], Entry{At: now, State: d.State})
	}
}

// Count returns the number of events for a port within the window.
func (t *Tracker) Count(d scanner.Diff) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	cutoff := time.Now().Add(-t.window)
	k := key(d)
	var n int
	for _, e := range t.events[k] {
		if e.At.After(cutoff) {
			n++
		}
	}
	return n
}

// Flush removes all events outside the current window.
func (t *Tracker) Flush() {
	t.mu.Lock()
	defer t.mu.Unlock()
	cutoff := time.Now().Add(-t.window)
	for k, entries := range t.events {
		var kept []Entry
		for _, e := range entries {
			if e.At.After(cutoff) {
				kept = append(kept, e)
			}
		}
		if len(kept) == 0 {
			delete(t.events, k)
		} else {
			t.events[k] = kept
		}
	}
}

func key(d scanner.Diff) string {
	return fmt.Sprintf("%s:%d", d.Proto, d.Port)
}
