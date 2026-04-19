// Package portage tracks how long a port has been in its current state.
package portage

import (
	"fmt"
	"sync"
	"time"
)

// Tracker records when a port was first seen open and reports its age.
type Tracker struct {
	mu    sync.Mutex
	first map[string]time.Time
	now   func() time.Time
}

// New returns a new Tracker.
func New() *Tracker {
	return &Tracker{
		first: make(map[string]time.Time),
		now:   time.Now,
	}
}

func key(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// Opened records the time a port was first seen open.
// If the port is already tracked, the existing timestamp is kept.
func (t *Tracker) Opened(port int, proto string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	k := key(port, proto)
	if _, ok := t.first[k]; !ok {
		t.first[k] = t.now()
	}
}

// Closed removes a port from tracking.
func (t *Tracker) Closed(port int, proto string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.first, key(port, proto))
}

// Age returns how long a port has been open and whether it is tracked.
func (t *Tracker) Age(port int, proto string) (time.Duration, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if ts, ok := t.first[key(port, proto)]; ok {
		return t.now().Sub(ts), true
	}
	return 0, false
}

// Reset clears all tracking state.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.first = make(map[string]time.Time)
}
