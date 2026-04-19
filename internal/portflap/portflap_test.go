package portflap

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func diff(port int, proto, state string) scanner.Diff {
	return scanner.Diff{Port: port, Proto: proto, State: state}
}

func TestRecord_BelowThresholdNotFlapping(t *testing.T) {
	d := New(time.Minute, 3)
	if d.Record(diff(80, "tcp", "opened")) {
		t.Fatal("expected not flapping on first event")
	}
	if d.Record(diff(80, "tcp", "closed")) {
		t.Fatal("expected not flapping on second event")
	}
}

func TestRecord_AtThresholdIsFlapping(t *testing.T) {
	d := New(time.Minute, 3)
	d.Record(diff(443, "tcp", "opened"))
	d.Record(diff(443, "tcp", "closed"))
	if !d.Record(diff(443, "tcp", "opened")) {
		t.Fatal("expected flapping at threshold")
	}
}

func TestRecord_EventsOutsideWindowEvicted(t *testing.T) {
	now := time.Now()
	d := New(time.Second, 3)
	d.now = func() time.Time { return now }

	d.Record(diff(22, "tcp", "opened"))
	d.Record(diff(22, "tcp", "closed"))

	// Advance past window
	d.now = func() time.Time { return now.Add(2 * time.Second) }

	if d.Record(diff(22, "tcp", "opened")) {
		t.Fatal("old events should have been evicted; should not be flapping")
	}
}

func TestRecord_DifferentPortsAreIndependent(t *testing.T) {
	d := New(time.Minute, 2)
	d.Record(diff(80, "tcp", "opened"))
	if d.Record(diff(443, "tcp", "opened")) {
		t.Fatal("port 443 should not be affected by port 80 events")
	}
}

func TestReset_ClearsHistory(t *testing.T) {
	d := New(time.Minute, 2)
	d.Record(diff(8080, "tcp", "opened"))
	d.Reset(diff(8080, "tcp", ""))
	if d.Record(diff(8080, "tcp", "closed")) {
		t.Fatal("after reset first event should not trigger flapping")
	}
}

func TestNew_PanicsOnZeroWindow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on zero window")
		}
	}()
	New(0, 3)
}

func TestNew_PanicsOnLowThresh(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on thresh < 2")
		}
	}()
	New(time.Minute, 1)
}
