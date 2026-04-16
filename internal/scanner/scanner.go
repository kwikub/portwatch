package scanner

import (
	"fmt"
	"net"
	"time"
)

// PortState represents the state of a single port.
type PortState struct {
	Port     int
	Protocol string
	Open     bool
	SeenAt   time.Time
}

// Scanner scans a range of ports on a given host.
type Scanner struct {
	Host    string
	Timeout time.Duration
}

// New creates a new Scanner for the given host.
func New(host string, timeout time.Duration) *Scanner {
	return &Scanner{Host: host, Timeout: timeout}
}

// Scan checks each port in [startPort, endPort] over the given protocol
// and returns a slice of PortState results.
func (s *Scanner) Scan(startPort, endPort int, protocol string) ([]PortState, error) {
	if startPort < 1 || endPort > 65535 || startPort > endPort {
		return nil, fmt.Errorf("invalid port range: %d-%d", startPort, endPort)
	}

	results := make([]PortState, 0, endPort-startPort+1)
	for port := startPort; port <= endPort; port++ {
		address := fmt.Sprintf("%s:%d", s.Host, port)
		conn, err := net.DialTimeout(protocol, address, s.Timeout)
		open := err == nil
		if open {
			conn.Close()
		}
		results = append(results, PortState{
			Port:     port,
			Protocol: protocol,
			Open:     open,
			SeenAt:   time.Now(),
		})
	}
	return results, nil
}
