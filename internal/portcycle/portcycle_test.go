package portcycle

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func diffs(port int, proto string, n int) []scanner.Diff {
	out := make([]scanner.Diff, n)
	for i := range out {
		out[i] = scanner.Diff{Port: port, Proto: proto, State: "opened"}
	}
	return out
}

func TestCount_ZeroBeforeRecord(t *testing.T) {
	tr := New()
	if c := tr.Count(80, "tcp"); c != 0 {
		t.Fatalf("expected 0, got %d", c)
	}
}

func TestRecord_IncrementsCount(t *testing.T) {
	tr := New()
	tr.Record(diffs(443, "tcp", 3))
	if c := tr.Count(443, "tcp"); c != 3 {
		t.Fatalf("expected 3, got %d", c)
	}
}

func TestRecord_ProtocolDistinct(t *testing.T) {
	tr := New()
	tr.Record(diffs(53, "tcp", 2))
	tr.Record(diffs(53, "udp", 5))
	if c := tr.Count(53, "tcp"); c != 2 {
		t.Fatalf("tcp: expected 2, got %d", c)
	}
	if c := tr.Count(53, "udp"); c != 5 {
		t.Fatalf("udp: expected 5, got %d", c)
	}
}

func TestReset_ClearsCount(t *testing.T) {
	tr := New()
	tr.Record(diffs(8080, "tcp", 4))
	tr.Reset(8080, "tcp")
	if c := tr.Count(8080, "tcp"); c != 0 {
		t.Fatalf("expected 0 after reset, got %d", c)
	}
}

func TestTop_ReturnsHighestCyclePort(t *testing.T) {
	tr := New()
	tr.Record(diffs(22, "tcp", 1))
	tr.Record(diffs(80, "tcp", 7))
	tr.Record(diffs(443, "tcp", 3))
	k, v := tr.Top()
	if k != "80/tcp" || v != 7 {
		t.Fatalf("expected 80/tcp:7, got %s:%d", k, v)
	}
}

func TestTop_EmptyTrackerReturnsZero(t *testing.T) {
	tr := New()
	k, v := tr.Top()
	if k != "" || v != 0 {
		t.Fatalf("expected empty result, got %s:%d", k, v)
	}
}
