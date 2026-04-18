package metrics

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Collector tracks runtime scan metrics for portwatch.
type Collector struct {
	mu          sync.Mutex
	scansTotal  int
	openEvents  int
	closeEvents int
	started     time.Time
	out         io.Writer
}

// New returns a new Collector writing summaries to out.
// If out is nil, os.Stdout is used.
func New(out io.Writer) *Collector {
	if out == nil {
		out = os.Stdout
	}
	return &Collector{started: time.Now(), out: out}
}

// RecordScan increments the total scan counter and tallies
// opened/closed port events from a diff count.
func (c *Collector) RecordScan(opened, closed int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.scansTotal++
	c.openEvents += opened
	c.closeEvents += closed
}

// Summary returns a snapshot of current metrics.
func (c *Collector) Summary() Summary {
	c.mu.Lock()
	defer c.mu.Unlock()
	return Summary{
		ScansTotal:  c.scansTotal,
		OpenEvents:  c.openEvents,
		CloseEvents: c.closeEvents,
		Uptime:      time.Since(c.started).Round(time.Second),
	}
}

// Print writes a human-readable summary line to the configured writer.
func (c *Collector) Print() {
	s := c.Summary()
	fmt.Fprintf(c.out, "[metrics] uptime=%s scans=%d opened=%d closed=%d\n",
		s.Uptime, s.ScansTotal, s.OpenEvents, s.CloseEvents)
}

// Summary holds a point-in-time snapshot of collected metrics.
type Summary struct {
	ScansTotal  int
	OpenEvents  int
	CloseEvents int
	Uptime      time.Duration
}
