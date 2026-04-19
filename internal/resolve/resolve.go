// Package resolve maps port numbers to well-known service names.
package resolve

import (
	"fmt"
	"strconv"
)

// Resolver maps port/protocol pairs to service names.
type Resolver struct {
	table map[string]string
}

// New returns a Resolver seeded with common well-known services.
func New() *Resolver {
	r := &Resolver{table: make(map[string]string)}
	defaults := map[string]string{
		"tcp/21":  "ftp",
		"tcp/22":  "ssh",
		"tcp/23":  "telnet",
		"tcp/25":  "smtp",
		"tcp/53":  "dns",
		"udp/53":  "dns",
		"tcp/80":  "http",
		"tcp/110": "pop3",
		"tcp/143": "imap",
		"tcp/443": "https",
		"tcp/3306": "mysql",
		"tcp/5432": "postgres",
		"tcp/6379": "redis",
		"tcp/8080": "http-alt",
	}
	for k, v := range defaults {
		r.table[k] = v
	}
	return r
}

// Lookup returns the service name for the given port and protocol.
// If unknown, it returns an empty string.
func (r *Resolver) Lookup(port int, proto string) string {
	return r.table[key(port, proto)]
}

// Register adds or overwrites a service name for a port/protocol pair.
func (r *Resolver) Register(port int, proto, name string) {
	r.table[key(port, proto)] = name
}

// LookupOrPort returns the service name, falling back to the port number as a string.
func (r *Resolver) LookupOrPort(port int, proto string) string {
	if name := r.Lookup(port, proto); name != "" {
		return name
	}
	return strconv.Itoa(port)
}

func key(port int, proto string) string {
	return fmt.Sprintf("%s/%d", proto, port)
}
