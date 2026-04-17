package filter

import "github.com/user/portwatch/internal/scanner"

// Rule defines inclusion/exclusion criteria for ports.
type Rule struct {
	Ports    []uint16
	Protocol string // "tcp", "udp", or "" for both
}

// Filter applies rules to a snapshot, returning a filtered copy.
type Filter struct {
	rules []Rule
}

// New creates a Filter from the provided rules.
func New(rules []Rule) *Filter {
	return &Filter{rules: rules}
}

// Apply returns only the ports from snap that match at least one rule.
// If no rules are defined all ports are returned unchanged.
func (f *Filter) Apply(snap scanner.Snapshot) scanner.Snapshot {
	if len(f.rules) == 0 {
		return snap
	}

	allowed := make(map[uint16]bool)
	for _, r := range f.rules {
		for _, p := range r.Ports {
			if r.Protocol == "" || r.Protocol == snap.Protocol {
				allowed[p] = true
			}
		}
	}

	filtered := make([]uint16, 0, len(snap.Ports))
	for _, p := range snap.Ports {
		if allowed[p] {
			filtered = append(filtered, p)
		}
	}

	return scanner.Snapshot{
		Timestamp: snap.Timestamp,
		Protocol:  snap.Protocol,
		Ports:     filtered,
	}
}
