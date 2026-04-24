// Package portdrain detects ports that have been continuously open
// for longer than a configured drain window, signalling they may be
// stale services that were never properly shut down.
package portdrain

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Drain tracks how long each port has been continuously open and
// reports ports whose open duration exceeds the configured window.
type Drain struct {
	mu     sync.Mutex
	window time.Duration
	firstSeen map[string]time.Time
}

// New creates a Drain with the given drain window. Panics if window
// is zero or negative.
func New(window time.Duration) *Drain {
	if window <= 0 {
		panic("portdrain: window must be positive")
	}
	return &Drain{
		window:    window,
		firstSeen: make(map[string]time.Time),
	}
}

// Record updates the drain tracker with the current open ports from
// a snapshot. Ports no longer present are removed from tracking.
func (d *Drain) Record(snap *scanner.Snapshot) {
	d.mu.Lock()
	defer d.mu.Unlock()

	active := make(map[string]struct{}, len(snap.Ports))
	for _, p := range snap.Ports {
		k := portKey(p)
		active[k] = struct{}{}
		if _, exists := d.firstSeen[k]; !exists {
			d.firstSeen[k] = time.Now()
		}
	}

	// Evict ports no longer open.
	for k := range d.firstSeen {
		if _, ok := active[k]; !ok {
			delete(d.firstSeen, k)
		}
	}
}

// Stale returns the ports that have been continuously open for longer
// than the drain window, along with how long each has been open.
func (d *Drain) Stale(now time.Time) []StalePort {
	d.mu.Lock()
	defer d.mu.Unlock()

	var out []StalePort
	for k, t := range d.firstSeen {
		if age := now.Sub(t); age >= d.window {
			out = append(out, StalePort{Key: k, Age: age})
		}
	}
	return out
}

// StalePort describes a port that has exceeded the drain window.
type StalePort struct {
	Key string
	Age time.Duration
}

func portKey(p scanner.Port) string {
	return fmt.Sprintf("%d/%s", p.Number, p.Protocol)
}
