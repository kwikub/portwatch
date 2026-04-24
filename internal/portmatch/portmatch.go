// Package portmatch provides pattern-based port matching, allowing rules
// expressed as port ranges or wildcards to be tested against observed diffs.
package portmatch

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/user/portwatch/internal/scanner"
)

// Rule represents a single match pattern: a port (or range) and a protocol.
type Rule struct {
	Low      int
	High     int
	Protocol string // "tcp", "udp", or "*" for any
}

// Matcher holds a set of rules and tests diffs against them.
type Matcher struct {
	rules []Rule
}

// New returns an empty Matcher.
func New() *Matcher {
	return &Matcher{}
}

// Add parses and appends a rule. port may be a single value ("80") or a
// range ("8000-8999"). protocol may be "tcp", "udp", or "*".
func (m *Matcher) Add(port, protocol string) error {
	protocol = strings.ToLower(strings.TrimSpace(protocol))
	if protocol != "tcp" && protocol != "udp" && protocol != "*" {
		return fmt.Errorf("portmatch: unsupported protocol %q", protocol)
	}

	var low, high int
	if idx := strings.IndexByte(port, '-'); idx >= 0 {
		lo, err1 := strconv.Atoi(port[:idx])
		hi, err2 := strconv.Atoi(port[idx+1:])
		if err1 != nil || err2 != nil || lo < 1 || hi > 65535 || lo > hi {
			return fmt.Errorf("portmatch: invalid port range %q", port)
		}
		low, high = lo, hi
	} else {
		v, err := strconv.Atoi(port)
		if err != nil || v < 1 || v > 65535 {
			return fmt.Errorf("portmatch: invalid port %q", port)
		}
		low, high = v, v
	}

	m.rules = append(m.rules, Rule{Low: low, High: high, Protocol: protocol})
	return nil
}

// Match returns all diffs whose port and protocol satisfy at least one rule.
func (m *Matcher) Match(diffs []scanner.Diff) []scanner.Diff {
	if len(m.rules) == 0 {
		return diffs
	}
	out := diffs[:0:0]
	for _, d := range diffs {
		if m.matches(d) {
			out = append(out, d)
		}
	}
	return out
}

func (m *Matcher) matches(d scanner.Diff) bool {
	for _, r := range m.rules {
		if (r.Protocol == "*" || r.Protocol == strings.ToLower(d.Protocol)) &&
			d.Port >= r.Low && d.Port <= r.High {
			return true
		}
	}
	return false
}
