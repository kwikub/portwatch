package portcap_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/portcap"
	"github.com/user/portwatch/internal/scanner"
)

func snap(ports []scanner.Port, at time.Time) *scanner.Snapshot {
	return &scanner.Snapshot{Ports: ports, At: at}
}

func TestNew_PanicsOnZeroWindow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero window")
		}
	}()
	portcap.New(0)
}

func TestNew_PanicsOnNegativeWindow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for negative window")
		}
	}()
	portcap.New(-time.Second)
}

func TestPeak_ZeroBeforeRecord(t *testing.T) {
	c := portcap.New(time.Minute)
	if got := c.Peak(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestPeak_ReturnMaxCount(t *testing.T) {
	c := portcap.New(time.Minute)
	now := time.Now()
	c.Record(snap(make([]scanner.Port, 3), now))
	c.Record(snap(make([]scanner.Port, 7), now.Add(time.Second)))
	c.Record(snap(make([]scanner.Port, 2), now.Add(2*time.Second)))
	if got := c.Peak(); got != 7 {
		t.Fatalf("expected peak 7, got %d", got)
	}
}

func TestPeak_EvictsOldObservations(t *testing.T) {
	c := portcap.New(100 * time.Millisecond)
	old := time.Now().Add(-200 * time.Millisecond)
	c.Record(snap(make([]scanner.Port, 10), old))
	c.Record(snap(make([]scanner.Port, 2), time.Now()))
	if got := c.Peak(); got != 2 {
		t.Fatalf("expected peak 2 after eviction, got %d", got)
	}
}

func TestRecord_NilSnapshotIsNoop(t *testing.T) {
	c := portcap.New(time.Minute)
	c.Record(nil)
	if got := c.Len(); got != 0 {
		t.Fatalf("expected 0 entries, got %d", got)
	}
}

func TestLen_CountsRetainedObservations(t *testing.T) {
	c := portcap.New(time.Minute)
	now := time.Now()
	c.Record(snap(nil, now))
	c.Record(snap(nil, now.Add(time.Second)))
	if got := c.Len(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}
