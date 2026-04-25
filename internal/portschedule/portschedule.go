// Package portschedule allows ports to be monitored only during defined
// time-of-week windows (e.g. weekdays 09:00–17:00).
package portschedule

import (
	"fmt"
	"strings"
	"time"
)

// Entry represents a single scheduled window for a port/protocol pair.
type Entry struct {
	Port     int
	Protocol string
	Weekdays []time.Weekday
	Start    int // minutes since midnight
	End      int // minutes since midnight
}

// Schedule holds active-window rules keyed by port+protocol.
type Schedule struct {
	entries []Entry
}

// New returns an empty Schedule.
func New() *Schedule {
	return &Schedule{}
}

// Add registers a monitoring window for the given port and protocol.
// start and end are "HH:MM" strings; weekdays is a slice of time.Weekday values.
func (s *Schedule) Add(port int, protocol, start, end string, weekdays []time.Weekday) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("portschedule: invalid port %d", port)
	}
	proto := strings.ToLower(protocol)
	if proto != "tcp" && proto != "udp" {
		return fmt.Errorf("portschedule: invalid protocol %q", protocol)
	}
	if len(weekdays) == 0 {
		return fmt.Errorf("portschedule: weekdays must not be empty")
	}
	st, err := parseClock(start)
	if err != nil {
		return fmt.Errorf("portschedule: invalid start %q: %w", start, err)
	}
	en, err := parseClock(end)
	if err != nil {
		return fmt.Errorf("portschedule: invalid end %q: %w", end, err)
	}
	if en <= st {
		return fmt.Errorf("portschedule: end must be after start")
	}
	s.entries = append(s.entries, Entry{
		Port:     port,
		Protocol: proto,
		Weekdays: weekdays,
		Start:    st,
		End:      en,
	})
	return nil
}

// Active returns true if the given port/protocol should be monitored at t.
// If no rule exists for the port/protocol, Active returns true (open by default).
func (s *Schedule) Active(port int, protocol string, t time.Time) bool {
	proto := strings.ToLower(protocol)
	minutes := t.Hour()*60 + t.Minute()
	weekday := t.Weekday()

	matched := false
	for _, e := range s.entries {
		if e.Port != port || e.Protocol != proto {
			continue
		}
		matched = true
		for _, wd := range e.Weekdays {
			if wd == weekday && minutes >= e.Start && minutes < e.End {
				return true
			}
		}
	}
	if !matched {
		return true
	}
	return false
}

func parseClock(s string) (int, error) {
	var h, m int
	if _, err := fmt.Sscanf(s, "%d:%d", &h, &m); err != nil {
		return 0, fmt.Errorf("expected HH:MM")
	}
	if h < 0 || h > 23 || m < 0 || m > 59 {
		return 0, fmt.Errorf("out of range")
	}
	return h*60 + m, nil
}
