// Package portstreak tracks consecutive scan cycles in which a port
// has remained in the same state (open or closed).
package portstreak

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Streak holds the current consecutive-cycle count for a single port+protocol
// pair together with the state that produced the streak.
type Streak struct {
	State string // "open" or "closed"
	Count int
}

// Tracker counts how many consecutive scan cycles each port has maintained
// its current state.
type Tracker struct {
	mu      sync.Mutex
	streaks map[string]*Streak
}

// New returns an initialised Tracker.
func New() *Tracker {
	return &Tracker{streaks: make(map[string]*Streak)}
}

func portKey(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// Record updates the streak for every port present in the snapshot.
// A port that keeps the same state increments its counter; a port whose
// state differs from the last recorded state resets to 1.
func (t *Tracker) Record(snap *scanner.Snapshot) {
	t.mu.Lock()
	defer t.mu.Unlock()

	seen := make(map[string]bool, len(snap.Ports))
	for _, p := range snap.Ports {
		k := portKey(p.Port, p.Protocol)
		seen[k] = true
		if s, ok := t.streaks[k]; ok && s.State == "open" {
			s.Count++
		} else {
			t.streaks[k] = &Streak{State: "open", Count: 1}
		}
	}

	// Ports absent from the snapshot are considered closed.
	for k, s := range t.streaks {
		if seen[k] {
			continue
		}
		if s.State == "closed" {
			s.Count++
		} else {
			t.streaks[k] = &Streak{State: "closed", Count: 1}
		}
	}
}

// Get returns the current Streak for the given port and protocol.
// ok is false when the port has never been recorded.
func (t *Tracker) Get(port int, proto string) (Streak, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	s, ok := t.streaks[portKey(port, proto)]
	if !ok {
		return Streak{}, false
	}
	return *s, true
}

// Reset clears all streak data.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.streaks = make(map[string]*Streak)
}
