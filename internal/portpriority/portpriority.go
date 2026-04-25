// Package portpriority assigns a numeric priority level to port diff events
// based on configurable rules. Higher priority values indicate more important
// events that should be processed or surfaced first.
package portpriority

import (
	"errors"
	"fmt"

	"github.com/user/portwatch/internal/scanner"
)

// DefaultPriority is returned when no rule matches.
const DefaultPriority = 0

// Rule associates a port/protocol pair with a priority level.
type Rule struct {
	Port     int
	Protocol string
	Priority int
}

// Prioritizer assigns priority levels to diffs.
type Prioritizer struct {
	rules    []Rule
	defaultP int
}

// New returns a Prioritizer with the given default priority.
// Panics if defaultPriority is negative.
func New(defaultPriority int) *Prioritizer {
	if defaultPriority < 0 {
		panic("portpriority: default priority must be non-negative")
	}
	return &Prioritizer{defaultP: defaultPriority}
}

// Add registers a priority rule. Returns an error for invalid input.
func (p *Prioritizer) Add(port int, protocol string, priority int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("portpriority: port %d out of range", port)
	}
	if protocol != "tcp" && protocol != "udp" {
		return errors.New("portpriority: protocol must be tcp or udp")
	}
	if priority < 0 {
		return errors.New("portpriority: priority must be non-negative")
	}
	p.rules = append(p.rules, Rule{Port: port, Protocol: protocol, Priority: priority})
	return nil
}

// Priority returns the priority for the given diff. The first matching rule
// wins; if no rule matches the configured default is returned.
func (p *Prioritizer) Priority(d scanner.Diff) int {
	for _, r := range p.rules {
		if r.Port == d.Port && r.Protocol == d.Protocol {
			return r.Priority
		}
	}
	return p.defaultP
}

// Sort returns diffs ordered from highest to lowest priority.
// The original slice is not modified.
func (p *Prioritizer) Sort(diffs []scanner.Diff) []scanner.Diff {
	out := make([]scanner.Diff, len(diffs))
	copy(out, diffs)
	// simple insertion sort — diff slices are typically small
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && p.Priority(out[j]) > p.Priority(out[j-1]); j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}
