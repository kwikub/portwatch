// Package portcap tracks the maximum number of simultaneously open ports
// observed within a sliding time window.
package portcap

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Cap records peak open-port counts over a rolling window.
type Cap struct {
	mu      sync.Mutex
	window  time.Duration
	buckets []bucket
}

type bucket struct {
	at    time.Time
	count int
}

// New returns a Cap that retains observations for the given window.
// It panics if window is zero or negative.
func New(window time.Duration) *Cap {
	if window <= 0 {
		panic("portcap: window must be positive")
	}
	return &Cap{window: window}
}

// Record stores the number of open ports observed in the supplied snapshot.
func (c *Cap) Record(snap *scanner.Snapshot) {
	if snap == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.evict()
	c.buckets = append(c.buckets, bucket{
		at:    snap.At,
		count: len(snap.Ports),
	})
}

// Peak returns the highest open-port count seen within the current window.
// Returns 0 if no observations have been recorded.
func (c *Cap) Peak() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.evict()
	max := 0
	for _, b := range c.buckets {
		if b.count > max {
			max = b.count
		}
	}
	return max
}

// Len returns the number of observations currently retained in the window.
func (c *Cap) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.evict()
	return len(c.buckets)
}

// evict removes observations older than the window. Must be called with mu held.
func (c *Cap) evict() {
	cutoff := time.Now().Add(-c.window)
	i := 0
	for i < len(c.buckets) && c.buckets[i].at.Before(cutoff) {
		i++
	}
	c.buckets = c.buckets[i:]
}
