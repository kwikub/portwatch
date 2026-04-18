// Package dedupe provides deduplication of port change events
// within a sliding time window to prevent duplicate notifications.
package dedupe

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Deduplicator filters out duplicate diffs seen within a window.
type Deduplicator struct {
	mu     sync.Mutex
	seen   map[string]time.Time
	window time.Duration
	now    func() time.Time
}

// New returns a Deduplicator that suppresses repeated diffs within window.
func New(window time.Duration) *Deduplicator {
	return &Deduplicator{
		seen:   make(map[string]time.Time),
		window: window,
		now:    time.Now,
	}
}

// Filter returns only diffs not seen within the current window.
func (d *Deduplicator) Filter(diffs []scanner.Diff) []scanner.Diff {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	d.evict(now)

	var out []scanner.Diff
	for _, diff := range diffs {
		k := key(diff)
		if _, seen := d.seen[k]; !seen {
			d.seen[k] = now
			out = append(out, diff)
		}
	}
	return out
}

// Reset clears all seen entries.
func (d *Deduplicator) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]time.Time)
}

func (d *Deduplicator) evict(now time.Time) {
	for k, t := range d.seen {
		if now.Sub(t) >= d.window {
			delete(d.seen, k)
		}
	}
}

func key(diff scanner.Diff) string {
	return fmt.Sprintf("%s:%d:%s", diff.Proto, diff.Port, diff.State)
}
