// Package ratelimit provides a simple token-bucket rate limiter
// to throttle how often alerts or reports are emitted.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter controls the rate at which events are allowed through.
type Limiter struct {
	mu       sync.Mutex
	interval time.Duration
	last     map[string]time.Time
}

// New creates a Limiter that allows at most one event per key per interval.
func New(interval time.Duration) *Limiter {
	return &Limiter{
		interval: interval,
		last:     make(map[string]time.Time),
	}
}

// Allow reports whether the event identified by key is allowed at time now.
// It updates the internal timestamp when the event is allowed.
func (l *Limiter) Allow(key string) bool {
	return l.AllowAt(key, time.Now())
}

// AllowAt is like Allow but accepts an explicit time, useful for testing.
func (l *Limiter) AllowAt(key string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if t, ok := l.last[key]; ok {
		if now.Sub(t) < l.interval {
			return false
		}
	}
	l.last[key] = now
	return true
}

// Reset clears the recorded timestamp for key, allowing the next event
// through immediately.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.last, key)
}
