// Package portroute maps ports to logical service routes or destinations.
package portroute

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Route describes a logical destination associated with a port.
type Route struct {
	Port     int
	Protocol string
	Target   string // e.g. "api-gateway", "db-primary"
	Group    string // optional grouping label
}

func (r Route) String() string {
	return fmt.Sprintf("%s/%d -> %s", r.Protocol, r.Port, r.Target)
}

type entry struct {
	target string
	group  string
}

// Registry maps port+protocol pairs to route metadata.
type Registry struct {
	mu      sync.RWMutex
	routes  map[string]entry
}

// New returns an empty Registry.
func New() *Registry {
	return &Registry{routes: make(map[string]entry)}
}

func portKey(port int, proto string) string {
	return fmt.Sprintf("%s:%d", proto, port)
}

// Add registers a route for the given port and protocol.
func (r *Registry) Add(port int, proto, target, group string) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("portroute: invalid port %d", port)
	}
	if proto != "tcp" && proto != "udp" {
		return fmt.Errorf("portroute: invalid protocol %q", proto)
	}
	if target == "" {
		return fmt.Errorf("portroute: target must not be empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.routes[portKey(port, proto)] = entry{target: target, group: group}
	return nil
}

// Lookup returns the Route for the given diff's port and protocol.
// ok is false if no route is registered.
func (r *Registry) Lookup(d scanner.Diff) (Route, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.routes[portKey(d.Port, d.Protocol)]
	if !ok {
		return Route{}, false
	}
	return Route{Port: d.Port, Protocol: d.Protocol, Target: e.target, Group: e.group}, true
}

// All returns every registered Route.
func (r *Registry) All() []Route {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Route, 0, len(r.routes))
	for _, e := range r.routes {
		out = append(out, Route{Target: e.target, Group: e.group})
	}
	return out
}
