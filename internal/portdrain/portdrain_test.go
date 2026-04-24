package portdrain

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeSnap(ports ...scanner.Port) *scanner.Snapshot {
	return &scanner.Snapshot{Ports: ports, At: time.Now()}
}

func TestNew_PanicsOnZeroWindow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on zero window")
		}
	}()
	New(0)
}

func TestNew_PanicsOnNegativeWindow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on negative window")
		}
	}()
	New(-time.Second)
}

func TestStale_EmptyBeforeRecord(t *testing.T) {
	d := New(time.Minute)
	if got := d.Stale(time.Now()); len(got) != 0 {
		t.Fatalf("expected no stale ports, got %d", len(got))
	}
}

func TestStale_BelowWindowNotReturned(t *testing.T) {
	d := New(time.Hour)
	snap := makeSnap(scanner.Port{Number: 80, Protocol: "tcp"})
	d.Record(snap)
	if got := d.Stale(time.Now()); len(got) != 0 {
		t.Fatalf("expected no stale ports, got %d", len(got))
	}
}

func TestStale_AfterWindowReturned(t *testing.T) {
	d := New(time.Millisecond)
	snap := makeSnap(scanner.Port{Number: 443, Protocol: "tcp"})
	d.Record(snap)
	time.Sleep(5 * time.Millisecond)
	got := d.Stale(time.Now())
	if len(got) != 1 {
		t.Fatalf("expected 1 stale port, got %d", len(got))
	}
	if got[0].Key != "443/tcp" {
		t.Errorf("unexpected key: %s", got[0].Key)
	}
	if got[0].Age < time.Millisecond {
		t.Errorf("age too small: %v", got[0].Age)
	}
}

func TestRecord_EvictsClosedPort(t *testing.T) {
	d := New(time.Millisecond)
	snap := makeSnap(scanner.Port{Number: 22, Protocol: "tcp"})
	d.Record(snap)
	time.Sleep(5 * time.Millisecond)
	// Record a snapshot without port 22.
	d.Record(makeSnap())
	if got := d.Stale(time.Now()); len(got) != 0 {
		t.Fatalf("expected evicted port to be absent, got %d stale", len(got))
	}
}

func TestRecord_DoesNotResetExistingTimer(t *testing.T) {
	d := New(time.Millisecond)
	port := scanner.Port{Number: 8080, Protocol: "tcp"}
	d.Record(makeSnap(port))
	time.Sleep(5 * time.Millisecond)
	// Record again — firstSeen should not be updated.
	d.Record(makeSnap(port))
	got := d.Stale(time.Now())
	if len(got) != 1 {
		t.Fatalf("expected 1 stale port after re-record, got %d", len(got))
	}
}
