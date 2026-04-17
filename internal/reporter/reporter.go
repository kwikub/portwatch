// Package reporter formats and emits port change summaries.
package reporter

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Format controls the output style.
type Format string

const (
	FormatText Format = "text"
	FormatCSV  Format = "csv"
)

// Reporter writes port change summaries to a writer.
type Reporter struct {
	w      io.Writer
	format Format
}

// New creates a Reporter. If w is nil, os.Stdout is used.
func New(w io.Writer, format Format) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	if format == "" {
		format = FormatText
	}
	return &Reporter{w: w, format: format}
}

// Report writes a summary of diffs to the underlying writer.
func (r *Reporter) Report(diffs []scanner.Diff) error {
	if len(diffs) == 0 {
		return nil
	}
	switch r.format {
	case FormatCSV:
		return r.writeCSV(diffs)
	default:
		return r.writeText(diffs)
	}
}

func (r *Reporter) writeText(diffs []scanner.Diff) error {
	ts := time.Now().UTC().Format(time.RFC3339)
	for _, d := range diffs {
		state := "opened"
		if d.Closed {
			state = "closed"
		}
		_, err := fmt.Fprintf(r.w, "[%s] %s %s/%d\n", ts, strings.ToUpper(state), d.Proto, d.Port)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Reporter) writeCSV(diffs []scanner.Diff) error {
	ts := time.Now().UTC().Format(time.RFC3339)
	for _, d := range diffs {
		state := "opened"
		if d.Closed {
			state = "closed"
		}
		_, err := fmt.Fprintf(r.w, "%s,%s,%s,%d\n", ts, state, d.Proto, d.Port)
		if err != nil {
			return err
		}
	}
	return nil
}
