// Package portchain provides a composable middleware chain for processing
// port diff events through an ordered sequence of named stages.
package portchain

import (
	"fmt"

	"github.com/user/portwatch/internal/scanner"
)

// Handler is a function that processes a slice of diffs and returns a
// (possibly filtered or modified) slice for the next stage.
type Handler func(diffs []scanner.Diff) []scanner.Diff

// Stage is a named processing step in the chain.
type Stage struct {
	Name    string
	Handler Handler
}

// Chain executes an ordered sequence of stages against a set of diffs.
type Chain struct {
	stages []Stage
}

// New returns an empty Chain.
func New() *Chain {
	return &Chain{}
}

// Add appends a named stage to the chain. It returns an error if name is
// empty or handler is nil.
func (c *Chain) Add(name string, h Handler) error {
	if name == "" {
		return fmt.Errorf("portchain: stage name must not be empty")
	}
	if h == nil {
		return fmt.Errorf("portchain: handler for stage %q must not be nil")
	}
	c.stages = append(c.stages, Stage{Name: name, Handler: h})
	return nil
}

// Run passes diffs through each stage in order. If any stage returns a nil
// or empty slice the chain short-circuits and returns that result immediately.
func (c *Chain) Run(diffs []scanner.Diff) []scanner.Diff {
	current := diffs
	for _, s := range c.stages {
		current = s.Handler(current)
		if len(current) == 0 {
			return current
		}
	}
	return current
}

// Stages returns a copy of the registered stage names in order.
func (c *Chain) Stages() []string {
	names := make([]string, len(c.stages))
	for i, s := range c.stages {
		names[i] = s.Name
	}
	return names
}

// Len returns the number of stages in the chain.
func (c *Chain) Len() int { return len(c.stages) }
