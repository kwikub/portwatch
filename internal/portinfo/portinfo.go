// Package portinfo aggregates enriched metadata for a port diff entry,
// combining port name, group membership, and geo data into a single struct.
package portinfo

import (
	"fmt"

	"github.com/user/portwatch/internal/portname"
	"github.com/user/portwatch/internal/portgroup"
	"github.com/user/portwatch/internal/scanner"
)

// Info holds aggregated metadata for a single port.
type Info struct {
	Port     int
	Proto    string
	Name     string
	Groups   []string
	Label    string
}

// Resolver builds Info values from registered sources.
type Resolver struct {
	names  *portname.Registry
	groups *portgroup.Registry
}

// New returns a Resolver backed by the given registries.
func New(names *portname.Registry, groups *portgroup.Registry) *Resolver {
	return &Resolver{names: names, groups: groups}
}

// Resolve returns an Info for the given port diff entry.
func (r *Resolver) Resolve(d scanner.Diff) Info {
	name := ""
	if r.names != nil {
		name = r.names.Lookup(d.Port, d.Proto)
	}

	var matched []string
	if r.groups != nil {
		for _, g := range r.groups.All() {
			if r.groups.Contains(g, d.Port, d.Proto) {
				matched = append(matched, g)
			}
		}
	}

	label := name
	if label == "" {
		label = fmt.Sprintf("%d/%s", d.Port, d.Proto)
	}

	return Info{
		Port:   d.Port,
		Proto:  d.Proto,
		Name:   name,
		Groups: matched,
		Label:  label,
	}
}
