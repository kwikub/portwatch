// Package portping probes individual ports to verify reachability.
package portping

import (
	"fmt"
	"net"
	"time"
)

// Result holds the outcome of a single probe.
type Result struct {
	Port     int
	Protocol string
	Open     bool
	Latency  time.Duration
}

// Prober probes ports on demand.
type Prober struct {
	timeout time.Duration
	host    string
}

// New returns a Prober targeting host with the given dial timeout.
func New(host string, timeout time.Duration) *Prober {
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	return &Prober{host: host, timeout: timeout}
}

// Probe attempts to connect to the given port/protocol and returns a Result.
func (p *Prober) Probe(port int, protocol string) Result {
	addr := fmt.Sprintf("%s:%d", p.host, port)
	start := time.Now()
	conn, err := net.DialTimeout(protocol, addr, p.timeout)
	latency := time.Since(start)
	if err != nil {
		return Result{Port: port, Protocol: protocol, Open: false, Latency: latency}
	}
	conn.Close()
	return Result{Port: port, Protocol: protocol, Open: true, Latency: latency}
}

// ProbeAll probes each port/protocol pair and returns all results.
func (p *Prober) ProbeAll(ports []int, protocol string) []Result {
	results := make([]Result, 0, len(ports))
	for _, port := range ports {
		results = append(results, p.Probe(port, protocol))
	}
	return results
}
