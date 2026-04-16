package scanner_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func TestNewSnapshot_SetsTimestamp(t *testing.T) {
	before := time.Now()
	snap := scanner.NewSnapshot(nil)
	after := time.Now()

	if snap.Timestamp.Before(before) || snap.Timestamp.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", snap.Timestamp, before, after)
	}
}

func TestNewSnapshot_StoresPorts(t *testing.T) {
	ports := []scanner.Port{
		{Number: 22, Protocol: "tcp"},
		{Number: 8080, Protocol: "tcp"},
	}

	snap := scanner.NewSnapshot(ports)

	if len(snap.Ports) != len(ports) {
		t.Fatalf("expected %d ports, got %d", len(ports), len(snap.Ports))
	}
	for i, p := range snap.Ports {
		if p != ports[i] {
			t.Errorf("port %d: got %+v, want %+v", i, p, ports[i])
		}
	}
}

func TestNewSnapshot_NilPortsIsEmpty(t *testing.T) {
	snap := scanner.NewSnapshot(nil)
	if snap.Ports != nil {
		t.Errorf("expected nil ports slice, got %v", snap.Ports)
	}
}
