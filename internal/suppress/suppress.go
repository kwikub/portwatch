// Package suppress provides a mechanism to suppress repeated notifications
// for the same port event within a configurable time window.
package suppress

import (
	"sync"
	"time"
)

// Suppressor tracks seen port events and suppresses duplicates within a window.
type Suppressor struct {
	mu     sync.Mutex
	seen   map[string]time.Time
	window time.Duration
	now    func() time.Time
}

// New returns a Suppressor that suppresses repeated events within window.
func New(window time.Duration) *Suppressor {
	return &Suppressor{
		seen:   make(map[string]time.Time),
		window: window,
		now:    time.Now,
	}
}

// Allow returns true if the event identified by key should be allowed through.
// Repeated calls with the same key within the window return false.
func (s *Suppressor) Allow(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.now()
	if last, ok := s.seen[key]; ok && now.Sub(last) < s.window {
		return false
	}
	s.seen[key] = now
	return true
}

// Reset clears the suppression record for key, allowing the next call through.
func (s *Suppressor) Reset(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.seen, key)
}

// Flush removes all entries that have expired beyond the window.
func (s *Suppressor) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.now()
	for k, t := range s.seen {
		if now.Sub(t) >= s.window {
			delete(s.seen, k)
		}
	}
}
