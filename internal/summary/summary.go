package summary

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Report holds aggregated scan summary data.
type Report struct {
	GeneratedAt time.Time
	TotalScans  int
	Opened      []scanner.Diff
	Closed      []scanner.Diff
}

// Summarizer builds and writes periodic summary reports.
type Summarizer struct {
	out io.Writer
}

// New returns a Summarizer writing to w. If w is nil, os.Stdout is used.
func New(w io.Writer) *Summarizer {
	if w == nil {
		w = os.Stdout
	}
	return &Summarizer{out: w}
}

// Build constructs a Report from accumulated diffs and scan count.
func Build(scans int, diffs []scanner.Diff) Report {
	r := Report{
		GeneratedAt: time.Now(),
		TotalScans:  scans,
	}
	for _, d := range diffs {
		if d.State == "opened" {
			r.Opened = append(r.Opened, d)
		} else {
			r.Closed = append(r.Closed, d)
		}
	}
	return r
}

// Write prints the summary report to the configured writer.
func (s *Summarizer) Write(r Report) error {
	_, err := fmt.Fprintf(s.out,
		"=== Port Summary [%s] ===\nScans: %d | Opened: %d | Closed: %d\n",
		r.GeneratedAt.Format(time.RFC3339),
		r.TotalScans,
		len(r.Opened),
		len(r.Closed),
	)
	if err != nil {
		return err
	}
	for _, d := range r.Opened {
		fmt.Fprintf(s.out, "  + %s/%d\n", d.Proto, d.Port)
	}
	for _, d := range r.Closed {
		fmt.Fprintf(s.out, "  - %s/%d\n", d.Proto, d.Port)
	}
	return nil
}
