package enrich

import (
	"github.com/user/portwatch/internal/baseline"
	"github.com/user/portwatch/internal/tags"
)

// FromBaseline constructs an Enricher by converting a baseline.Baseline
// snapshot into the flat map expected by New.
func FromBaseline(t *tags.Registry, b *baseline.Baseline) *Enricher {
	snap, ok := b.Snapshot()
	if !ok {
		return New(t, nil)
	}
	flat := make(map[string]bool, len(snap.Ports))
	for _, p := range snap.Ports {
		key := p.Proto + ":" + itoa(p.Port)
		flat[key] = true
	}
	return New(t, flat)
}
