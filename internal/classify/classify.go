// Package classify assigns severity labels to port diff events
// based on configurable rules (port range, protocol, direction).
package classify

import "github.com/user/portwatch/internal/scanner"

// Level represents a severity classification.
type Level string

const (
	LevelInfo     Level = "info"
	LevelWarning  Level = "warning"
	LevelCritical Level = "critical"
)

// Rule maps a port range and protocol to a severity level.
type Rule struct {
	MinPort  int
	MaxPort  int
	Protocol string
	Level    Level
}

// Classifier assigns a Level to a scanner.Diff.
type Classifier struct {
	rules []Rule
}

// New returns a Classifier with the given rules.
// Rules are evaluated in order; the first match wins.
func New(rules []Rule) *Classifier {
	return &Classifier{rules: rules}
}

// Classify returns the Level for the given diff.
// If no rule matches, LevelInfo is returned.
func (c *Classifier) Classify(d scanner.Diff) Level {
	for _, r := range c.rules {
		if r.Protocol != "" && r.Protocol != d.Protocol {
			continue
		}
		if d.Port >= r.MinPort && d.Port <= r.MaxPort {
			return r.Level
		}
	}
	return LevelInfo
}

// ClassifyAll returns a map from diff index to Level for a slice of diffs.
func (c *Classifier) ClassifyAll(diffs []scanner.Diff) map[int]Level {
	out := make(map[int]Level, len(diffs))
	for i, d := range diffs {
		out[i] = c.Classify(d)
	}
	return out
}
