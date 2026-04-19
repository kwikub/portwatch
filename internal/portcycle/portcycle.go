// Package portcycle tracks how many times a port has transitioned
// between open and closed states within the lifetime of the daemon.
package portcycle

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Tracker counts open/close cycles per port+protocol pair.
type Tracker struct {
	mu     sync.Mutex
	counts map[string]int
}

// New returns an initialised Tracker.
func New() *Tracker {
	return &Tracker{counts: make(map[string]int)}
}

func key(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// Record increments the cycle counter for every diff entry whose state
// transitions to either opened or closed.
func (t *Tracker) Record(diffs []scanner.Diff) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, d := range diffs {
		t.counts[key(d.Port, d.Proto)]++
	}
}

// Count returns the number of transitions recorded for the given port and
// protocol. Returns 0 for unknown pairs.
func (t *Tracker) Count(port int, proto string) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.counts[key(port, proto)]
}

// Reset clears the cycle counter for a specific port+protocol pair.
func (t *Tracker) Reset(port int, proto string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.counts, key(port, proto))
}

// Top returns the port+protocol key with the highest cycle count, and that
// count. Returns an empty string and 0 if no data has been recorded.
func (t *Tracker) Top() (string, int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	var topKey string
	var topVal int
	for k, v := range t.counts {
		if v > topVal {
			topKey = k
			topVal = v
		}
	}
	return topKey, topVal
}
