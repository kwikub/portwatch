// Package portlabel maps ports to human-readable labels combining name and classification.
package portlabel

import (
	"fmt"

	"github.com/user/portwatch/internal/classify"
	"github.com/user/portwatch/internal/portname"
	"github.com/user/portwatch/internal/scanner"
)

// Label holds the resolved display information for a port.
type Label struct {
	Port     int
	Proto    string
	Name     string
	Class    string
	Display  string
}

// Resolver combines portname and classify lookups into a single Label.
type Resolver struct {
	names   *portname.Registry
	classes *classify.Classifier
}

// New returns a Resolver backed by the provided registries.
func New(names *portname.Registry, classes *classify.Classifier) *Resolver {
	return &Resolver{names: names, classes: classes}
}

// Resolve returns a Label for the given port entry.
func (r *Resolver) Resolve(p scanner.Port) Label {
	name := r.names.Label(p.Port, p.Proto)
	class := "info"
	if r.classes != nil {
		class = r.classes.Classify(p.Port, p.Proto)
	}
	display := fmt.Sprintf("%s (%s/%d)", name, p.Proto, p.Port)
	return Label{
		Port:    p.Port,
		Proto:   p.Proto,
		Name:    name,
		Class:   class,
		Display: display,
	}
}
