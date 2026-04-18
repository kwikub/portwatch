// Package throttle limits how frequently notifications are sent for a given port.
package throttle

import (
	"fmt"
	"sync"
	"time"
)

// Throttle suppresses repeated notifications for the same port within a cooldown window.
type Throttle struct {
	mu       sync.Mutex
	cooldown time.Duration
	last     map[string]time.Time
}

// New returns a Throttle with the given cooldown duration.
// Panics if cooldown is zero or negative.
func New(cooldown time.Duration) *Throttle {
	if cooldown <= 0 {
		panic("throttle: cooldown must be positive")
	}
	return &Throttle{
		cooldown: cooldown,
		last:     make(map[string]time.Time),
	}
}

// Allow returns true if the given proto/port combination is outside the cooldown window.
// If allowed, the timestamp for that key is updated.
func (t *Throttle) Allow(proto string, port int) bool {
	key := fmt.Sprintf("%s:%d", proto, port)
	t.mu.Lock()
	defer t.mu.Unlock()
	if ts, ok := t.last[key]; ok && time.Since(ts) < t.cooldown {
		return false
	}
	t.last[key] = time.Now()
	return true
}

// Reset clears the cooldown record for the given proto/port, allowing immediate retry.
func (t *Throttle) Reset(proto string, port int) {
	key := fmt.Sprintf("%s:%d", proto, port)
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.last, key)
}

// Len returns the number of tracked keys.
func (t *Throttle) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.last)
}
