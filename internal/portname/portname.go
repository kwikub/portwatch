package portname

import "fmt"

// Registry maps port/protocol pairs to human-readable service names.
type Registry struct {
	entries map[string]string
}

// New returns a Registry pre-loaded with common well-known port names.
func New() *Registry {
	r := &Registry{entries: make(map[string]string)}
	for _, e := range builtins {
		r.entries[key(e.port, e.proto)] = e.name
	}
	return r
}

// Lookup returns the service name for the given port and protocol.
// An empty string is returned when no entry exists.
func (r *Registry) Lookup(port int, proto string) string {
	return r.entries[key(port, proto)]
}

// Register adds or overwrites a service name entry.
func (r *Registry) Register(port int, proto, name string) {
	r.entries[key(port, proto)] = name
}

// Label returns the service name if known, otherwise a formatted port string.
func (r *Registry) Label(port int, proto string) string {
	if name := r.Lookup(port, proto); name != "" {
		return name
	}
	return fmt.Sprintf("%d/%s", port, proto)
}

func key(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}
