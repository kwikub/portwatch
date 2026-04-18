// Package tags provides port tagging — associating human-readable labels
// with specific port/protocol pairs for richer log output.
package tags

import (
	"fmt"
	"strings"
)

// Tag associates a label with a port/protocol pair.
type Tag struct {
	Port     int
	Protocol string
	Label    string
}

// Registry holds a set of port tags.
type Registry struct {
	entries map[string]string
}

// New returns an empty Registry.
func New() *Registry {
	return &Registry{entries: make(map[string]string)}
}

// Add registers a label for the given port and protocol.
// Protocol is normalised to lowercase.
func (r *Registry) Add(port int, protocol, label string) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("tags: invalid port %d", port)
	}
	proto := strings.ToLower(strings.TrimSpace(protocol))
	if proto != "tcp" && proto != "udp" {
		return fmt.Errorf("tags: unsupported protocol %q", protocol)
	}
	if strings.TrimSpace(label) == "" {
		return fmt.Errorf("tags: label must not be empty")
	}
	r.entries[key(port, proto)] = strings.TrimSpace(label)
	return nil
}

// Lookup returns the label for a port/protocol pair and whether it was found.
func (r *Registry) Lookup(port int, protocol string) (string, bool) {
	v, ok := r.entries[key(port, strings.ToLower(protocol))]
	return v, ok
}

// All returns a slice of every registered Tag.
func (r *Registry) All() []Tag {
	out := make([]Tag, 0, len(r.entries))
	for k, label := range r.entries {
		var port int
		var proto string
		fmt.Sscanf(k, "%d/%s", &port, &proto)
		out = append(out, Tag{Port: port, Protocol: proto, Label: label})
	}
	return out
}

func key(port int, protocol string) string {
	return fmt.Sprintf("%d/%s", port, protocol)
}
