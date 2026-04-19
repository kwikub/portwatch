// Package export writes port diff snapshots to external formats (JSON, CSV).
package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Format selects the output encoding.
type Format string

const (
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
)

// Record is a serialisable representation of a single diff entry.
type Record struct {
	Timestamp time.Time `json:"timestamp"`
	Port      int       `json:"port"`
	Protocol  string    `json:"protocol"`
	State     string    `json:"state"`
}

// Exporter writes diffs to an io.Writer in the configured format.
type Exporter struct {
	format Format
	w      io.Writer
}

// New returns an Exporter that writes to w using the given format.
func New(w io.Writer, f Format) *Exporter {
	return &Exporter{format: f, w: w}
}

// Write encodes diffs and writes them to the underlying writer.
func (e *Exporter) Write(diffs []scanner.Diff) error {
	if len(diffs) == 0 {
		return nil
	}
	records := toRecords(diffs)
	switch e.format {
	case FormatJSON:
		return writeJSON(e.w, records)
	case FormatCSV:
		return writeCSV(e.w, records)
	default:
		return fmt.Errorf("export: unknown format %q", e.format)
	}
}

func toRecords(diffs []scanner.Diff) []Record {
	out := make([]Record, len(diffs))
	for i, d := range diffs {
		state := "opened"
		if !d.Opened {
			state = "closed"
		}
		out[i] = Record{
			Timestamp: d.At,
			Port:      d.Port,
			Protocol:  d.Protocol,
			State:     state,
		}
	}
	return out
}

func writeJSON(w io.Writer, records []Record) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}

func writeCSV(w io.Writer, records []Record) error {
	cw := csv.NewWriter(w)
	_ = cw.Write([]string{"timestamp", "port", "protocol", "state"})
	for _, r := range records {
		_ = cw.Write([]string{
			r.Timestamp.Format(time.RFC3339),
			fmt.Sprintf("%d", r.Port),
			r.Protocol,
			r.State,
		})
	}
	cw.Flush()
	return cw.Error()
}
