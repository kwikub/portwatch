// Package portgroup provides named groupings of ports for easier monitoring rules.
package portgroup

import (
	"fmt"
	"sync"
)

// Group is a named collection of port+protocol pairs.
type Group struct {
	Name  string
	Ports []Entry
}

// Entry is a single port+protocol pair within a group.
type Entry struct {
	Port     int
	Protocol string
}

// Registry holds named port groups.
type Registry struct {
	mu     sync.RWMutex
	groups map[string]*Group
}

// New returns an empty Registry.
func New() *Registry {
	return &Registry{groups: make(map[string]*Group)}
}

// Add registers a named group. Overwrites any existing group with the same name.
func (r *Registry) Add(name string, entries []Entry) error {
	if name == "" {
		return fmt.Errorf("portgroup: name must not be empty")
	}
	for _, e := range entries {
		if e.Port < 1 || e.Port > 65535 {
			return fmt.Errorf("portgroup: invalid port %d in group %q", e.Port, name)
		}
		if e.Protocol != "tcp" && e.Protocol != "udp" {
			return fmt.Errorf("portgroup: invalid protocol %q in group %q", e.Protocol, name)
		}
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.groups[name] = &Group{Name: name, Ports: entries}
	return nil
}

// Get returns the group for the given name, or false if not found.
func (r *Registry) Get(name string) (*Group, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	g, ok := r.groups[name]
	return g, ok
}

// Names returns all registered group names.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.groups))
	for n := range r.groups {
		names = append(names, n)
	}
	return names
}

// Contains reports whether the group named name includes the given port+protocol.
func (r *Registry) Contains(name string, port int, proto string) bool {
	g, ok := r.Get(name)
	if !ok {
		return false
	}
	for _, e := range g.Ports {
		if e.Port == port && e.Protocol == proto {
			return true
		}
	}
	return false
}
