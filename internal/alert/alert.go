// Package alert provides threshold-based alerting for port activity.
package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Level represents alert severity.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Event represents a single alert event.
type Event struct {
	Time    time.Time
	Level   Level
	Message string
	Port    int
	Proto   string
}

// Alerter emits alert events when port changes exceed a threshold.
type Alerter struct {
	threshold int
	out       io.Writer
}

// New creates an Alerter. threshold is the number of changes that triggers
// an ALERT level event; below that, INFO is used.
func New(threshold int, out io.Writer) *Alerter {
	if out == nil {
		out = os.Stdout
	}
	return &Alerter{threshold: threshold, out: out}
}

// Evaluate inspects a diff slice and writes alert events to the writer.
func (a *Alerter) Evaluate(diffs []scanner.Diff) []Event {
	var events []Event
	for _, d := range diffs {
		lvl := LevelInfo
		if len(diffs) >= a.threshold {
			lvl = LevelAlert
		}
		msg := fmt.Sprintf("port %d/%s %s", d.Port, d.Proto, d.State)
		ev := Event{
			Time:    time.Now(),
			Level:   lvl,
			Message: msg,
			Port:    d.Port,
			Proto:   d.Proto,
		}
		events = append(events, ev)
		fmt.Fprintf(a.out, "[%s] %s %s\n", ev.Level, ev.Time.Format(time.RFC3339), ev.Message)
	}
	return events
}
