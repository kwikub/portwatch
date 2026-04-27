// Package portfence enforces an allowlist of ports per protocol.
// Any port not explicitly permitted is considered blocked.
package portfence

import (
	"errors"
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Fence holds the set of allowed port/protocol combinations.
type Fence struct {
	mu      sync.RWMutex
	allowed map[string]struct{}
}

// New returns an empty Fence. Use Add to populate the allowlist.
func New() *Fence {
	return &Fence{allowed: make(map[string]struct{})}
}

func portKey(port int, proto string) string {
	return fmt.Sprintf("%s:%d", proto, port)
}

// Add registers a port/protocol pair as permitted.
// Returns an error if port is out of range or protocol is unknown.
func (f *Fence) Add(port int, proto string) error {
	if port < 1 || port > 65535 {
		return errors.New("portfence: port must be between 1 and 65535")
	}
	if proto != "tcp" && proto != "udp" {
		return fmt.Errorf("portfence: unknown protocol %q", proto)
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.allowed[portKey(port, proto)] = struct{}{}
	return nil
}

// Permitted reports whether the given port/protocol pair is on the allowlist.
// When the allowlist is empty every port is considered permitted.
func (f *Fence) Permitted(port int, proto string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if len(f.allowed) == 0 {
		return true
	}
	_, ok := f.allowed[portKey(port, proto)]
	return ok
}

// Filter returns only the diffs whose port/protocol pair is permitted.
func (f *Fence) Filter(diffs []scanner.Diff) []scanner.Diff {
	if len(f.allowed) == 0 {
		return diffs
	}
	out := diffs[:0:0]
	for _, d := range diffs {
		if f.Permitted(d.Port, d.Proto) {
			out = append(out, d)
		}
	}
	return out
}
