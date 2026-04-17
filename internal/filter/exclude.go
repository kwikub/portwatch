package filter

import "github.com/user/portwatch/internal/scanner"

// ExcludeRule defines ports that should always be suppressed.
type ExcludeRule struct {
	Ports    []uint16
	Protocol string
}

// Excluder removes specific ports from a snapshot.
type Excluder struct {
	rules []ExcludeRule
}

// NewExcluder creates an Excluder from the provided rules.
func NewExcluder(rules []ExcludeRule) *Excluder {
	return &Excluder{rules: rules}
}

// Apply strips excluded ports from snap.
func (e *Excluder) Apply(snap scanner.Snapshot) scanner.Snapshot {
	if len(e.rules) == 0 {
		return snap
	}

	blocked := make(map[uint16]bool)
	for _, r := range e.rules {
		if r.Protocol == "" || r.Protocol == snap.Protocol {
			for _, p := range r.Ports {
				blocked[p] = true
			}
		}
	}

	kept := make([]uint16, 0, len(snap.Ports))
	for _, p := range snap.Ports {
		if !blocked[p] {
			kept = append(kept, p)
		}
	}

	return scanner.Snapshot{
		Timestamp: snap.Timestamp,
		Protocol:  snap.Protocol,
		Ports:     kept,
	}
}
