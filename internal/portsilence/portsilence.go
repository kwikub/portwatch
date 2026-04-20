// Package portsilence allows specific ports to be silenced (suppressed from
// output) for a configurable duration, useful for known-noisy ports.
package portsilence

import (
	"fmt"
	"sync"
	"time"
)

// Silencer tracks ports that should be silenced until their silence window expires.
type Silencer struct {
	mu      sync.Mutex
	entries map[string]time.Time
	dur     time.Duration
}

// New creates a Silencer with the given silence duration.
func New(dur time.Duration) *Silencer {
	if dur <= 0 {
		panic("portsilence: duration must be positive")
	}
	return &Silencer{
		entries: make(map[string]time.Time),
		dur:     dur,
	}
}

// Silence marks the given port+protocol as silenced starting now.
func (s *Silencer) Silence(port uint16, proto string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[portKey(port, proto)] = time.Now().Add(s.dur)
}

// IsSilenced reports whether the given port+protocol is currently silenced.
func (s *Silencer) IsSilenced(port uint16, proto string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	expiry, ok := s.entries[portKey(port, proto)]
	if !ok {
		return false
	}
	if time.Now().After(expiry) {
		delete(s.entries, portKey(port, proto))
		return false
	}
	return true
}

// Lift removes a silence entry immediately, regardless of expiry.
func (s *Silencer) Lift(port uint16, proto string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, portKey(port, proto))
}

// Active returns the number of currently active silence entries.
func (s *Silencer) Active() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	count := 0
	for k, expiry := range s.entries {
		if now.After(expiry) {
			delete(s.entries, k)
		} else {
			count++
		}
	}
	return count
}

func portKey(port uint16, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}
