// Package portexpiry tracks ports that have been closed for a configurable
// duration and emits expiry events when the deadline is exceeded.
package portexpiry

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry records when a port was first seen closed.
type Entry struct {
	ClosedAt time.Time
	Port     int
	Protocol string
}

// Tracker watches for ports that remain closed beyond a TTL.
type Tracker struct {
	mu      sync.Mutex
	closed  map[string]Entry
	ttl     time.Duration
	now     func() time.Time
}

// New creates a Tracker with the given TTL.
func New(ttl time.Duration) *Tracker {
	return &Tracker{
		closed: make(map[string]Entry),
		ttl:    ttl,
		now:    time.Now,
	}
}

func key(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// Record updates internal state from a batch of diffs.
func (t *Tracker) Record(diffs []scanner.Diff) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, d := range diffs {
		k := key(d.Port, d.Protocol)
		if d.State == scanner.StateClosed {
			if _, exists := t.closed[k]; !exists {
				t.closed[k] = Entry{ClosedAt: t.now(), Port: d.Port, Protocol: d.Protocol}
			}
		} else {
			delete(t.closed, k)
		}
	}
}

// Expired returns all entries whose closed duration exceeds the TTL.
func (t *Tracker) Expired() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.now()
	var out []Entry
	for _, e := range t.closed {
		if now.Sub(e.ClosedAt) >= t.ttl {
			out = append(out, e)
		}
	}
	return out
}

// Evict removes an entry so it will not be reported again.
func (t *Tracker) Evict(port int, proto string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.closed, key(port, proto))
}
