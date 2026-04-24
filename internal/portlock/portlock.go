// Package portlock tracks ports that have been administratively locked,
// preventing any state changes from being reported for them.
package portlock

import (
	"fmt"
	"sync"
	"time"
)

// Lock holds metadata about a locked port.
type Lock struct {
	Port      int
	Protocol  string
	Reason    string
	LockedAt  time.Time
	ExpiresAt time.Time // zero means no expiry
}

// Expired reports whether the lock has passed its expiry time.
func (l Lock) Expired(now time.Time) bool {
	return !l.ExpiresAt.IsZero() && now.After(l.ExpiresAt)
}

type registry struct {
	mu    sync.RWMutex
	locks map[string]Lock
}

// New returns an empty lock registry.
func New() *registry {
	return &registry{locks: make(map[string]Lock)}
}

func portKey(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// Lock adds or replaces a lock for the given port/protocol pair.
// A zero ttl means the lock never expires.
func (r *registry) Lock(port int, proto, reason string, ttl time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var exp time.Time
	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}
	r.locks[portKey(port, proto)] = Lock{
		Port:      port,
		Protocol:  proto,
		Reason:    reason,
		LockedAt:  time.Now(),
		ExpiresAt: exp,
	}
}

// Unlock removes a lock. It is a no-op if the port is not locked.
func (r *registry) Unlock(port int, proto string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.locks, portKey(port, proto))
}

// IsLocked reports whether the port is currently locked (and not expired).
func (r *registry) IsLocked(port int, proto string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	l, ok := r.locks[portKey(port, proto)]
	if !ok {
		return false
	}
	if l.Expired(time.Now()) {
		return false
	}
	return true
}

// Get returns the Lock for a port, and whether it exists and is active.
func (r *registry) Get(port int, proto string) (Lock, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	l, ok := r.locks[portKey(port, proto)]
	if !ok || l.Expired(time.Now()) {
		return Lock{}, false
	}
	return l, true
}

// All returns a snapshot of all active (non-expired) locks.
func (r *registry) All() []Lock {
	r.mu.RLock()
	defer r.mu.RUnlock()
	now := time.Now()
	out := make([]Lock, 0, len(r.locks))
	for _, l := range r.locks {
		if !l.Expired(now) {
			out = append(out, l)
		}
	}
	return out
}
