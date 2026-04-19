// Package portwatch provides a high-level coordinator that wires together
// scanning, diffing, enrichment, and pipeline execution for a single watch cycle.
package portwatch

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/enrich"
	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

// Coordinator runs one watch cycle: scan → diff → enrich → pipeline.
type Coordinator struct {
	scanner  *scanner.Scanner
	state    *state.State
	enricher *enrich.Enricher
	pipeline *pipeline.Pipeline
}

// Config holds the dependencies needed to build a Coordinator.
type Config struct {
	Scanner  *scanner.Scanner
	State    *state.State
	Enricher *enrich.Enricher
	Pipeline *pipeline.Pipeline
}

// New creates a Coordinator from the provided Config.
func New(cfg Config) *Coordinator {
	return &Coordinator{
		scanner:  cfg.Scanner,
		state:    cfg.State,
		enricher: cfg.Enricher,
		pipeline: cfg.Pipeline,
	}
}

// Run executes a single watch cycle. It scans open ports, computes the diff
// against the previously saved state, enriches the diffs, then runs the
// pipeline. The new snapshot is persisted before returning.
func (c *Coordinator) Run(ctx context.Context) error {
	snap, err := c.scanner.Scan(ctx)
	if err != nil {
		return err
	}

	prev, _ := c.state.Load()

	diffs := scanner.ComputeDiff(prev, snap)

	enriched := c.enricher.Enrich(diffs)

	if err := c.pipeline.Run(ctx, enriched); err != nil {
		return err
	}

	return c.state.Save(snap)
}

// RunForever calls Run repeatedly on the given ticker channel until ctx is
// cancelled.
func (c *Coordinator) RunForever(ctx context.Context, tick <-chan time.Time) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case _, ok := <-tick:
			if !ok {
				return nil
			}
			if err := c.Run(ctx); err != nil {
				return err
			}
		}
	}
}
