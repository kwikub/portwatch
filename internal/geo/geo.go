// Package geo provides IP geolocation lookup for port activity events.
package geo

import (
	"net"
	"sync"
)

// Location holds geolocation metadata for an IP address.
type Location struct {
	IP      string
	Country string
	City    string
	ASN     string
}

// Registry maps IP addresses to Location records.
type Registry struct {
	mu      sync.RWMutex
	entries map[string]Location
}

// New returns an empty Registry.
func New() *Registry {
	return &Registry{entries: make(map[string]Location)}
}

// Register adds or replaces a Location entry for the given IP.
func (r *Registry) Register(ip string, loc Location) error {
	if net.ParseIP(ip) == nil {
		return &InvalidIPError{IP: ip}
	}
	loc.IP = ip
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[ip] = loc
	return nil
}

// Lookup returns the Location for ip, and whether it was found.
func (r *Registry) Lookup(ip string) (Location, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	loc, ok := r.entries[ip]
	return loc, ok
}

// InvalidIPError is returned when an IP address cannot be parsed.
type InvalidIPError struct {
	IP string
}

func (e *InvalidIPError) Error() string {
	return "geo: invalid IP address: " + e.IP
}
