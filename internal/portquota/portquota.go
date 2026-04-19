// Package portquota tracks how many times a port has been seen opened
// within a rolling time window and flags ports that exceed a threshold.
package portquota

import (
	"fmt"
	"sync"
	"time"
)

// Quota holds per-port open-event counts within a sliding window.
type Quota struct {
	mu        sync.Mutex
	window    time.Duration
	threshold int
	events    map[string][]time.Time
}

// New returns a Quota that flags a port when it opens more than threshold
// times within window.
func New(window time.Duration, threshold int) *Quota {
	if window <= 0 {
		panic("portquota: window must be positive")
	}
	if threshold <= 0 {
		panic("portquota: threshold must be positive")
	}
	return &Quota{
		window:    window,
		threshold: threshold,
		events:    make(map[string][]time.Time),
	}
}

func portKey(port int, proto string) string {
	return fmt.Sprintf("%s:%d", proto, port)
}

// Record registers an open event for the given port/proto at now.
func (q *Quota) Record(port int, proto string, now time.Time) {
	q.mu.Lock()
	defer q.mu.Unlock()
	k := portKey(port, proto)
	q.events[k] = append(q.evict(k, now), now)
}

// Exceeded reports whether the port has surpassed the threshold within the window.
func (q *Quota) Exceeded(port int, proto string, now time.Time) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	k := portKey(port, proto)
	return len(q.evict(k, now)) >= q.threshold
}

// Reset clears all recorded events for the given port/proto.
func (q *Quota) Reset(port int, proto string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.events, portKey(port, proto))
}

// evict removes events outside the window and returns the remaining slice.
// Caller must hold q.mu.
func (q *Quota) evict(k string, now time.Time) []time.Time {
	cutoff := now.Add(-q.window)
	slice := q.events[k]
	i := 0
	for i < len(slice) && slice[i].Before(cutoff) {
		i++
	}
	slice = slice[i:]
	q.events[k] = slice
	return slice
}
