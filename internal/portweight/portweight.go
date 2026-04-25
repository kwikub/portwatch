// Package portweight assigns a numeric weight to port diff events based on
// configurable rules. Higher weights indicate higher operational significance.
package portweight

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Rule associates a port+protocol pair with a weight.
type Rule struct {
	Port     int
	Protocol string
	Weight   int
}

// Weigher holds a set of rules and returns weights for diff events.
type Weigher struct {
	mu      sync.RWMutex
	rules   map[string]int
	default_ int
}

// New returns a Weigher with the given default weight used when no rule matches.
// Panics if defaultWeight is negative.
func New(defaultWeight int) *Weigher {
	if defaultWeight < 0 {
		panic("portweight: defaultWeight must be non-negative")
	}
	return &Weigher{
		rules:    make(map[string]int),
		default_: defaultWeight,
	}
}

// Add registers a weight rule. Returns an error if port or weight is invalid.
func (w *Weigher) Add(r Rule) error {
	if r.Port < 1 || r.Port > 65535 {
		return fmt.Errorf("portweight: port %d out of range", r.Port)
	}
	if r.Protocol != "tcp" && r.Protocol != "udp" {
		return fmt.Errorf("portweight: unknown protocol %q", r.Protocol)
	}
	if r.Weight < 0 {
		return fmt.Errorf("portweight: weight must be non-negative")
	}
	w.mu.Lock()
	w.rules[key(r.Port, r.Protocol)] = r.Weight
	w.mu.Unlock()
	return nil
}

// Weight returns the weight for the given diff. Falls back to the default
// weight when no matching rule is found.
func (w *Weigher) Weight(d scanner.Diff) int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	if v, ok := w.rules[key(d.Port, d.Protocol)]; ok {
		return v
	}
	return w.default_
}

// WeighDiffs returns a map of each diff to its computed weight.
func (w *Weigher) WeighDiffs(diffs []scanner.Diff) map[scanner.Diff]int {
	out := make(map[scanner.Diff]int, len(diffs))
	for _, d := range diffs {
		out[d] = w.Weight(d)
	}
	return out
}

func key(port int, proto string) string {
	return fmt.Sprintf("%s:%d", proto, port)
}
