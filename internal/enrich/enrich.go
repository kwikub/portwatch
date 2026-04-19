// Package enrich attaches metadata (tags, baseline flags) to scanner diffs.
package enrich

import (
	"fmt"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/tags"
)

// Diff wraps scanner.Diff with enrichment metadata.
type Diff struct {
	scanner.Diff
	Tag       string // human label from tags registry
	Deviation bool   // true when port is absent from baseline
}

// Enricher decorates diffs with tag and baseline information.
type Enricher struct {
	tags     *tags.Registry
	baseline map[string]bool
}

// New constructs an Enricher.
// baseline is a set of "proto:port" keys considered normal (non-deviating).
func New(t *tags.Registry, baseline map[string]bool) *Enricher {
	if baseline == nil {
		baseline = make(map[string]bool)
	}
	return &Enricher{tags: t, baseline: baseline}
}

// Enrich returns a slice of decorated diffs.
func (e *Enricher) Enrich(diffs []scanner.Diff) []Diff {
	out := make([]Diff, 0, len(diffs))
	for _, d := range diffs {
		label, _ := e.tags.Lookup(d.Proto, d.Port)
		key := fmt.Sprintf("%s:%d", d.Proto, d.Port)
		out = append(out, Diff{
			Diff:      d,
			Tag:       label,
			Deviation: !e.baseline[key],
		})
	}
	return out
}
