package scanner

// ChangeType describes how a port's state changed between scans.
type ChangeType string

const (
	Opened ChangeType = "opened"
	Closed ChangeType = "closed"
)

// Change represents a detected change in port state.
type Change struct {
	PortState
	Change ChangeType
}

// Diff compares two snapshots (previous, current) and returns any changes.
// Both slices must be sorted by Port and cover the same range.
func Diff(previous, current []PortState) []Change {
	prev := index(previous)
	var changes []Change

	for _, cur := range current {
		key := portKey(cur.Port, cur.Protocol)
		old, exists := prev[key]
		switch {
		case !exists && cur.Open:
			changes = append(changes, Change{PortState: cur, Change: Opened})
		case exists && old.Open && !cur.Open:
			changes = append(changes, Change{PortState: cur, Change: Closed})
		case exists && !old.Open && cur.Open:
			changes = append(changes, Change{PortState: cur, Change: Opened})
		}
	}
	return changes
}

func index(states []PortState) map[string]PortState {
	m := make(map[string]PortState, len(states))
	for _, s := range states {
		m[portKey(s.Port, s.Protocol)] = s
	}
	return m
}

func portKey(port int, protocol string) string {
	return protocol + ":" + string(rune(port))
}
