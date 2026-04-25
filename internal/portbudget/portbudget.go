// Package portbudget tracks the total number of concurrently open ports and
// enforces a configurable upper bound (budget). When the budget is exceeded,
// Exceeded returns true so callers can emit alerts or suppress further actions.
package portbudget

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Budget tracks open-port counts against a maximum allowance.
type Budget struct {
	mu      sync.Mutex
	max     int
	current int
}

// New creates a Budget with the given maximum. Panics if max is <= 0.
func New(max int) *Budget {
	if max <= 0 {
		panic(fmt.Sprintf("portbudget: max must be positive, got %d", max))
	}
	return &Budget{max: max}
}

// Record updates the current open-port count from the provided snapshot.
// It counts all ports present in the snapshot regardless of protocol.
func (b *Budget) Record(snap *scanner.Snapshot) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.current = len(snap.Ports)
}

// Exceeded reports whether the current open-port count is strictly greater
// than the configured maximum.
func (b *Budget) Exceeded() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.current > b.max
}

// Current returns the most recently recorded open-port count.
func (b *Budget) Current() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.current
}

// Max returns the configured budget ceiling.
func (b *Budget) Max() int {
	return b.max
}

// Reset sets the current count back to zero.
func (b *Budget) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.current = 0
}
