// Package portrank assigns a numeric priority rank to ports based on
// configurable rules, allowing downstream consumers to sort or filter
// diffs by importance.
package portrank

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Ranker maps ports to a priority rank. Lower numbers indicate higher priority.
type Ranker struct {
	mu    sync.RWMutex
	rules []rule
}

type rule struct {
	port     int
	protocol string
	rank     int
}

// New returns an empty Ranker.
func New() *Ranker {
	return &Ranker{}
}

// Add registers a port/protocol pair with the given rank.
// Protocol must be "tcp" or "udp". Rank must be >= 0.
func (r *Ranker) Add(port int, protocol string, rank int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("portrank: invalid port %d", port)
	}
	if protocol != "tcp" && protocol != "udp" {
		return fmt.Errorf("portrank: invalid protocol %q", protocol)
	}
	if rank < 0 {
		return fmt.Errorf("portrank: rank must be >= 0, got %d", rank)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, existing := range r.rules {
		if existing.port == port && existing.protocol == protocol {
			r.rules[i].rank = rank
			return nil
		}
	}
	r.rules = append(r.rules, rule{port: port, protocol: protocol, rank: rank})
	return nil
}

// Rank returns the priority rank for the given diff. If no rule matches,
// defaultRank (100) is returned.
func (r *Ranker) Rank(d scanner.Diff) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, rule := range r.rules {
		if rule.port == d.Port && rule.protocol == d.Protocol {
			return rule.rank
		}
	}
	return 100
}

// Sort returns a copy of diffs ordered by ascending rank (highest priority first).
// Diffs with equal rank preserve their original relative order.
func (r *Ranker) Sort(diffs []scanner.Diff) []scanner.Diff {
	out := make([]scanner.Diff, len(diffs))
	copy(out, diffs)
	n := len(out)
	for i := 1; i < n; i++ {
		for j := i; j > 0 && r.Rank(out[j]) < r.Rank(out[j-1]); j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}
