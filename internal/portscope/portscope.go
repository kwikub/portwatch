// Package portscope restricts scanner activity to a defined set of port ranges
// per protocol, acting as a pre-scan gate before ports are evaluated.
package portscope

import "fmt"

// Scope holds inclusive port ranges for a given protocol.
type Scope struct {
	ranges []portRange
}

type portRange struct {
	protocol string
	lo, hi   int
}

// New returns an empty Scope. Use Add to register ranges.
func New() *Scope {
	return &Scope{}
}

// Add registers an inclusive [lo, hi] port range for the given protocol.
// protocol must be "tcp" or "udp". lo must be <= hi and both must be in [1, 65535].
func (s *Scope) Add(protocol string, lo, hi int) error {
	if protocol != "tcp" && protocol != "udp" {
		return fmt.Errorf("portscope: unsupported protocol %q", protocol)
	}
	if lo < 1 || hi > 65535 || lo > hi {
		return fmt.Errorf("portscope: invalid range %d-%d", lo, hi)
	}
	s.ranges = append(s.ranges, portRange{protocol: protocol, lo: lo, hi: hi})
	return nil
}

// Contains reports whether the given port/protocol falls within any registered
// range. If no ranges have been registered, Contains always returns true so
// that an unconfigured Scope is a no-op pass-through.
func (s *Scope) Contains(protocol string, port int) bool {
	if len(s.ranges) == 0 {
		return true
	}
	for _, r := range s.ranges {
		if r.protocol == protocol && port >= r.lo && port <= r.hi {
			return true
		}
	}
	return false
}

// Size returns the total number of port/protocol pairs covered by all ranges.
func (s *Scope) Size() int {
	total := 0
	for _, r := range s.ranges {
		_ = r.protocol
		total += r.hi - r.lo + 1
	}
	return total
}
