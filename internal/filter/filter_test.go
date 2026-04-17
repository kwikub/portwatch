package filter_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/scanner"
)

func snap(ports []uint16) scanner.Snapshot {
	return scanner.Snapshot{Timestamp: time.Now(), Protocol: "tcp", Ports: ports}
}

func TestFilter_NoRulesReturnsAll(t *testing.T) {
	f := filter.New(nil)
	s := snap([]uint16{80, 443, 8080})
	out := f.Apply(s)
	if len(out.Ports) != 3 {
		t.Fatalf("expected 3 ports, got %d", len(out.Ports))
	}
}

func TestFilter_AllowsMatchingPorts(t *testing.T) {
	f := filter.New([]filter.Rule{
		{Ports: []uint16{80, 443}, Protocol: "tcp"},
	})
	out := f.Apply(snap([]uint16{80, 443, 8080}))
	if len(out.Ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(out.Ports))
	}
}

func TestFilter_ProtocolMismatchExcludesAll(t *testing.T) {
	f := filter.New([]filter.Rule{
		{Ports: []uint16{80}, Protocol: "udp"},
	})
	out := f.Apply(snap([]uint16{80, 443}))
	if len(out.Ports) != 0 {
		t.Fatalf("expected 0 ports, got %d", len(out.Ports))
	}
}
