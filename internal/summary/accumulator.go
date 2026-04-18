package summary

import (
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Accumulator collects diffs across multiple scans for summary reporting.
type Accumulator struct {
	mu     sync.Mutex
	scans  int
	diffs  []scanner.Diff
}

// NewAccumulator returns an initialised Accumulator.
func NewAccumulator() *Accumulator {
	return &Accumulator{}
}

// Record adds diffs from a completed scan.
func (a *Accumulator) Record(diffs []scanner.Diff) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.scans++
	a.diffs = append(a.diffs, diffs...)
}

// Flush returns the current Report and resets internal state.
func (a *Accumulator) Flush() Report {
	a.mu.Lock()
	defer a.mu.Unlock()
	r := Build(a.scans, a.diffs)
	a.scans = 0
	a.diffs = nil
	return r
}

// Scans returns the number of scans recorded since last flush.
func (a *Accumulator) Scans() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.scans
}
