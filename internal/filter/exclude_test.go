package filter_test

import (
	"testing"

	"github.com/user/portwatch/internal/filter"
)

func TestExcluder_NoRulesReturnsAll(t *testing.T) {
	e := filter.NewExcluder(nil)
	s := snap([]uint16{22, 80, 443})
	out := e.Apply(s)
	if len(out.Ports) != 3 {
		t.Fatalf("expected 3 ports, got %d", len(out.Ports))
	}
}

func TestExcluder_RemovesBlockedPorts(t *testing.T) {
	e := filter.NewExcluder([]filter.ExcludeRule{
		{Ports: []uint16{22}, Protocol: "tcp"},
	})
	out := e.Apply(snap([]uint16{22, 80, 443}))
	if len(out.Ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(out.Ports))
	}
	for _, p := range out.Ports {
		if p == 22 {
			t.Fatal("port 22 should have been excluded")
		}
	}
}

func TestExcluder_ProtocolMismatchKeepsAll(t *testing.T) {
	e := filter.NewExcluder([]filter.ExcludeRule{
		{Ports: []uint16{22}, Protocol: "udp"},
	})
	out := e.Apply(snap([]uint16{22, 80}))
	if len(out.Ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(out.Ports))
	}
}
