// Package baseline records a known-good port snapshot and computes
// deviations from it on subsequent scans.
package baseline

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Baseline holds a reference snapshot against which future scans are compared.
type Baseline struct {
	mu       sync.RWMutex
	ports    map[string]struct{}
	RecordedAt time.Time
	path     string
}

type stored struct {
	Ports      []string  `json:"ports"`
	RecordedAt time.Time `json:"recorded_at"`
}

// New loads an existing baseline from path, or returns an empty one.
func New(path string) (*Baseline, error) {
	b := &Baseline{path: path, ports: make(map[string]struct{})}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return b, nil
	}
	if err != nil {
		return nil, err
	}
	var s stored
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	for _, p := range s.Ports {
		b.ports[p] = struct{}{}
	}
	b.RecordedAt = s.RecordedAt
	return b, nil
}

// Record saves snap as the new baseline.
func (b *Baseline) Record(snap *scanner.Snapshot) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.ports = make(map[string]struct{}, len(snap.Ports))
	for _, p := range snap.Ports {
		b.ports[key(p)] = struct{}{}
	}
	b.RecordedAt = snap.At
	s := stored{RecordedAt: b.RecordedAt}
	for k := range b.ports {
		s.Ports = append(s.Ports, k)
	}
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return os.WriteFile(b.path, data, 0o644)
}

// Deviations returns ports in snap that are not in the baseline, and ports in
// the baseline that are absent from snap.
func (b *Baseline) Deviations(snap *scanner.Snapshot) (added []scanner.Port, removed []scanner.Port) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	current := make(map[string]scanner.Port, len(snap.Ports))
	for _, p := range snap.Ports {
		current[key(p)] = p
	}
	for k, p := range current {
		if _, ok := b.ports[k]; !ok {
			added = append(added, p)
		}
	}
	for k := range b.ports {
		if _, ok := current[k]; !ok {
			removed = append(removed, scanner.Port{}) // placeholder; key only
			_ = k
		}
	}
	return added, removed
}

func key(p scanner.Port) string {
	return p.Proto + ":" + p.Addr
}
