// Package portwindow tracks whether a port event occurred within a
// configured time-of-day window (e.g. business hours).
package portwindow

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Window defines a start and end time-of-day in HH:MM format.
type Window struct {
	start time.Duration // offset from midnight
	end   time.Duration // offset from midnight
}

// Checker evaluates whether diffs fall inside or outside a Window.
type Checker struct {
	window Window
	now    func() time.Time
}

// New creates a Checker for the given time window.
// start and end are "HH:MM" strings (24-hour clock).
// Returns an error if either value cannot be parsed.
func New(start, end string) (*Checker, error) {
	s, err := parseClock(start)
	if err != nil {
		return nil, fmt.Errorf("portwindow: invalid start %q: %w", start, err)
	}
	e, err := parseClock(end)
	if err != nil {
		return nil, fmt.Errorf("portwindow: invalid end %q: %w", end, err)
	}
	return &Checker{
		window: Window{start: s, end: e},
		now:    time.Now,
	}, nil
}

// InWindow reports whether the current time falls within the configured window.
// Overnight windows (start > end) are supported.
func (c *Checker) InWindow() bool {
	now := c.now()
	offset := sinceMidnight(now)
	s, e := c.window.start, c.window.end
	if s <= e {
		return offset >= s && offset < e
	}
	// overnight: e.g. 22:00 – 06:00
	return offset >= s || offset < e
}

// Filter returns only those diffs whose timestamp falls within the window.
// If d.Time is zero the current wall clock is used instead.
func (c *Checker) Filter(diffs []scanner.Diff) []scanner.Diff {
	var out []scanner.Diff
	for _, d := range diffs {
		t := d.Time
		if t.IsZero() {
			t = c.now()
		}
		offset := sinceMidnight(t)
		s, e := c.window.start, c.window.end
		var inside bool
		if s <= e {
			inside = offset >= s && offset < e
		} else {
			inside = offset >= s || offset < e
		}
		if inside {
			out = append(out, d)
		}
	}
	return out
}

func sinceMidnight(t time.Time) time.Duration {
	h, m, _ := t.Clock()
	return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute
}

func parseClock(s string) (time.Duration, error) {
	var h, m int
	_, err := fmt.Sscanf(s, "%d:%d", &h, &m)
	if err != nil || h < 0 || h > 23 || m < 0 || m > 59 {
		return 0, fmt.Errorf("expected HH:MM")
	}
	return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute, nil
}
