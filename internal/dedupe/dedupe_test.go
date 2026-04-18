package dedupe

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeDiff(proto string, port int, state string) scanner.Diff {
	return scanner.Diff{Proto: proto, Port: port, State: state}
}

func TestFilter_FirstCallPassesThrough(t *testing.T) {
	d := New(time.Minute)
	diffs := []scanner.Diff{makeDiff("tcp", 80, "opened")}
	out := d.Filter(diffs)
	if len(out) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(out))
	}
}

func TestFilter_DuplicateWithinWindowBlocked(t *testing.T) {
	d := New(time.Minute)
	diff := []scanner.Diff{makeDiff("tcp", 80, "opened")}
	d.Filter(diff)
	out := d.Filter(diff)
	if len(out) != 0 {
		t.Fatalf("expected 0 diffs, got %d", len(out))
	}
}

func TestFilter_DuplicateAfterWindowAllowed(t *testing.T) {
	d := New(time.Millisecond * 10)
	fixed := time.Now()
	d.now = func() time.Time { return fixed }

	diff := []scanner.Diff{makeDiff("tcp", 443, "opened")}
	d.Filter(diff)

	d.now = func() time.Time { return fixed.Add(time.Millisecond * 20) }
	out := d.Filter(diff)
	if len(out) != 1 {
		t.Fatalf("expected 1 diff after window, got %d", len(out))
	}
}

func TestFilter_DifferentPortsAreIndependent(t *testing.T) {
	d := New(time.Minute)
	out := d.Filter([]scanner.Diff{
		makeDiff("tcp", 80, "opened"),
		makeDiff("tcp", 443, "opened"),
	})
	if len(out) != 2 {
		t.Fatalf("expected 2 diffs, got %d", len(out))
	}
}

func TestReset_AllowsRepeat(t *testing.T) {
	d := New(time.Minute)
	diff := []scanner.Diff{makeDiff("udp", 53, "opened")}
	d.Filter(diff)
	d.Reset()
	out := d.Filter(diff)
	if len(out) != 1 {
		t.Fatalf("expected 1 diff after reset, got %d", len(out))
	}
}
