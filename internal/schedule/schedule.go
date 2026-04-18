// Package schedule provides tick-based interval control for the watcher loop.
package schedule

import (
	"time"
)

// Ticker wraps a time.Ticker and exposes a channel-based API
// so the watcher can be driven at a configurable interval.
type Ticker struct {
	ticker   *time.Ticker
	C        <-chan time.Time
	stopCh   chan struct{}
}

// New creates a Ticker that fires at the given interval.
// Panics if interval is zero or negative.
func New(interval time.Duration) *Ticker {
	if interval <= 0 {
		panic("schedule: interval must be positive")
	}
	t := time.NewTicker(interval)
	return &Ticker{
		ticker: t,
		C:      t.C,
		stopCh: make(chan struct{}),
	}
}

// Stop halts the ticker and releases resources.
func (t *Ticker) Stop() {
	t.ticker.Stop()
	select {
	case t.stopCh <- struct{}{}:
	default:
	}
}

// Done returns a channel that is closed after Stop is called.
func (t *Ticker) Done() <-chan struct{} {
	return t.stopCh
}
