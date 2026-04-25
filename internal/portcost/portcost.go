// Package portcost assigns a numeric cost to port diff events based on
// configurable per-port and per-protocol rules. Higher cost indicates
// a more significant or sensitive change.
package portcost

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

const defaultCost = 1

type entry struct {
	port     int
	protocol string
	cost     int
}

// Registry holds cost rules for ports.
type Registry struct {
	mu      sync.RWMutex
	rules   []entry
	default_ int
}

// New returns a Registry with the given default cost.
// Panics if defaultCost is negative.
func New(def int) *Registry {
	if def < 0 {
		panic("portcost: default cost must be non-negative")
	}
	return &Registry{default_: def}
}

// Add registers a cost for the given port and protocol.
// Protocol must be "tcp" or "udp". Port must be in [1, 65535].
func (r *Registry) Add(port int, protocol string, cost int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("portcost: invalid port %d", port)
	}
	if protocol != "tcp" && protocol != "udp" {
		return fmt.Errorf("portcost: invalid protocol %q", protocol)
	}
	if cost < 0 {
		return fmt.Errorf("portcost: cost must be non-negative")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, e := range r.rules {
		if e.port == port && e.protocol == protocol {
			r.rules[i].cost = cost
			return nil
		}
	}
	r.rules = append(r.rules, entry{port: port, protocol: protocol, cost: cost})
	return nil
}

// Cost returns the cost for the given diff. Falls back to the default.
func (r *Registry) Cost(d scanner.Diff) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, e := range r.rules {
		if e.port == d.Port && e.protocol == d.Protocol {
			return e.cost
		}
	}
	return r.default_
}

// Total sums the costs of all provided diffs.
func (r *Registry) Total(diffs []scanner.Diff) int {
	total := 0
	for _, d := range diffs {
		total += r.Cost(d)
	}
	return total
}
