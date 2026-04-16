package scanner

// State represents whether a port was opened or closed.
type State int

const (
	StateOpened State = iota
	StateClosed
)

// Diff describes a single port state change.
type Diff struct {
	Protocol string
	Port     int
	State    State
}

// Diff computes the changes between two port snapshots.
func Diff(prev, curr []PortInfo) []Diff {
	prevIdx := index(prev)
	currIdx := index(curr)

	var diffs []Diff

	for k, p := range currIdx {
		if _, existed := prevIdx[k]; !existed {
			diffs = append(diffs, Diff{Protocol: p.Protocol, Port: p.Port, State: StateOpened})
		}
	}

	for k, p := range prevIdx {
		if _, exists := currIdx[k]; !exists {
			diffs = append(diffs, Diff{Protocol: p.Protocol, Port: p.Port, State: StateClosed})
		}
	}

	return diffs
}

func index(ports []PortInfo) map[string]PortInfo {
	m := make(map[string]PortInfo, len(ports))
	for _, p := range ports {
		m[portKey(p)] = p
	}
	return m
}

func portKey(p PortInfo) string {
	return p.Protocol + ":" + itoa(p.Port)
}

func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}
