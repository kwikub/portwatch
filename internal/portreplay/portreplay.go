// Package portreplay replays historical port diff sequences for debugging
// and analysis purposes.
package portreplay

import (
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry holds a single replay frame.
type Entry struct {
	At    time.Time
	Diffs []scanner.Diff
}

// Replayer reads a sequence of diff entries and emits them to a handler
// with optional time compression.
type Replayer struct {
	entries  []Entry
	speed    float64
	output   io.Writer
	handler  func(Entry)
}

// New returns a Replayer for the provided entries.
// speed controls playback rate (1.0 = real-time, 2.0 = double speed, 0 = instant).
func New(entries []Entry, speed float64, out io.Writer) *Replayer {
	if speed < 0 {
		speed = 0
	}
	return &Replayer{
		entries: entries,
		speed:   speed,
		output:  out,
	}
}

// OnEntry sets a callback invoked for each replayed entry.
func (r *Replayer) OnEntry(fn func(Entry)) {
	r.handler = fn
}

// Run replays all entries in order, sleeping between frames when speed > 0.
func (r *Replayer) Run() error {
	if len(r.entries) == 0 {
		return nil
	}
	for i, entry := range r.entries {
		if i > 0 && r.speed > 0 {
			gap := entry.At.Sub(r.entries[i-1].At)
			delay := time.Duration(float64(gap) / r.speed)
			if delay > 0 {
				time.Sleep(delay)
			}
		}
		fmt.Fprintf(r.output, "[replay] %s — %d diff(s)\n", entry.At.Format(time.RFC3339), len(entry.Diffs))
		if r.handler != nil {
			r.handler(entry)
		}
	}
	return nil
}

// Len returns the number of entries in the replay sequence.
func (r *Replayer) Len() int { return len(r.entries) }
