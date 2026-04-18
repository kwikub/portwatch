package snapshot

import (
	"github.com/example/portwatch/internal/scanner"
)

// DiffEntry pairs a time-ordered diff result with its source entries.
type DiffEntry struct {
	From Entry
	To   Entry
	Diff []scanner.Diff
}

// Diffs computes diffs between consecutive entries in the history.
// Returns an empty slice when fewer than two snapshots are stored.
func (h *History) Diffs() []DiffEntry {
	entries := h.All()
	if len(entries) < 2 {
		return nil
	}
	out := make([]DiffEntry, 0, len(entries)-1)
	for i := 1; i < len(entries); i++ {
		d := scanner.Diff(entries[i-1].Snapshot, entries[i].Snapshot)
		if len(d) == 0 {
			continue
		}
		out = append(out, DiffEntry{
			From: entries[i-1],
			To:   entries[i],
			Diff: d,
		})
	}
	return out
}
