package scanner

import "time"

// Port represents a single open port observed during a scan.
type Port struct {
	Number   int    `json:"number"`
	Protocol string `json:"protocol"`
}

// Snapshot captures the set of open ports at a point in time.
type Snapshot struct {
	Ports     []Port    `json:"ports"`
	Timestamp time.Time `json:"timestamp"`
}

// NewSnapshot constructs a Snapshot from a slice of Ports, stamped now.
func NewSnapshot(ports []Port) Snapshot {
	return Snapshot{
		Ports:     ports,
		Timestamp: time.Now(),
	}
}
