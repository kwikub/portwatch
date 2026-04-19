// Package portflap detects ports that open and close repeatedly within a time window.
package portflap

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Detector tracks open/close transitions per port and flags flapping behaviour.
type Detector struct {
	mu      sync.Mutex
	window  time.Duration
	thresh  int
	events  map[string][]time.Time
	now     func() time.Time
}

// New returns a Detector that considers a port flapping when it transitions
// at least thresh times within window.
func New(window time.Duration, thresh int) *Detector {
	if window <= 0 {
		panic("portflap: window must be positive")
	}
	if thresh < 2 {
		panic("portflap: thresh must be at least 2")
	}
	return &Detector{
		window: window,
		thresh: thresh,
		events: make(map[string][]time.Time),
		now:    time.Now,
	}
}

// Record registers a state-change event for the given diff and reports
// whether the port is currently considered flapping.
func (d *Detector) Record(diff scanner.Diff) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	k := key(diff)
	now := d.now()
	cutoff := now.Add(-d.window)

	prev := d.events[k]
	filtered := prev[:0]
	for _, t := range prev {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	filtered = append(filtered, now)
	d.events[k] = filtered

	return len(filtered) >= d.thresh
}

// Reset clears the event history for a port.
func (d *Detector) Reset(diff scanner.Diff) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.events, key(diff))
}

func key(diff scanner.Diff) string {
	return fmt.Sprintf("%d/%s", diff.Port, diff.Proto)
}
