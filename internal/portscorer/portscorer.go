// Package portscorer computes a composite risk score for a port diff
// by combining weight, priority, cost, and classification level.
package portscorer

import (
	"fmt"

	"github.com/user/portwatch/internal/scanner"
)

// Scorer computes a numeric risk score for a port diff.
type Scorer struct {
	weights    map[string]float64
	priorities map[string]int
	costs      map[string]float64
	defScore   float64
}

// Option configures a Scorer.
type Option func(*Scorer)

// WithWeight registers a weight multiplier for a port/protocol pair.
func WithWeight(port int, proto string, w float64) Option {
	return func(s *Scorer) {
		s.weights[portKey(port, proto)] = w
	}
}

// WithPriority registers a priority value for a port/protocol pair.
func WithPriority(port int, proto string, p int) Option {
	return func(s *Scorer) {
		s.priorities[portKey(port, proto)] = p
	}
}

// WithCost registers a cost value for a port/protocol pair.
func WithCost(port int, proto string, c float64) Option {
	return func(s *Scorer) {
		s.costs[portKey(port, proto)] = c
	}
}

// WithDefault sets the baseline score before modifiers are applied.
func WithDefault(score float64) Option {
	return func(s *Scorer) {
		s.defScore = score
	}
}

// New creates a Scorer with the given options.
func New(opts ...Option) *Scorer {
	s := &Scorer{
		weights:    make(map[string]float64),
		priorities: make(map[string]int),
		costs:      make(map[string]float64),
		defScore:   1.0,
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Score returns a composite risk score for the given diff.
// score = (defaultScore + cost) * weight * (1 + priority/10)
func (s *Scorer) Score(d scanner.Diff) float64 {
	k := portKey(d.Port, d.Proto)

	base := s.defScore
	if c, ok := s.costs[k]; ok {
		base += c
	}

	w := 1.0
	if wv, ok := s.weights[k]; ok {
		w = wv
	}

	prio := 0
	if pv, ok := s.priorities[k]; ok {
		prio = pv
	}

	return base * w * (1.0 + float64(prio)/10.0)
}

func portKey(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}
