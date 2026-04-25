package portweight

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func diff(port int, proto, state string) scanner.Diff {
	return scanner.Diff{Port: port, Protocol: proto, State: state}
}

func TestNew_PanicsOnNegativeDefault(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for negative default weight")
		}
	}()
	New(-1)
}

func TestWeight_ReturnsDefaultWhenNoRule(t *testing.T) {
	w := New(5)
	if got := w.Weight(diff(8080, "tcp", "opened")); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestAdd_ValidRule(t *testing.T) {
	w := New(1)
	if err := w.Add(Rule{Port: 443, Protocol: "tcp", Weight: 10}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := w.Weight(diff(443, "tcp", "opened")); got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}

func TestAdd_InvalidPort(t *testing.T) {
	w := New(0)
	if err := w.Add(Rule{Port: 0, Protocol: "tcp", Weight: 5}); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := w.Add(Rule{Port: 70000, Protocol: "tcp", Weight: 5}); err == nil {
		t.Fatal("expected error for port 70000")
	}
}

func TestAdd_InvalidProtocol(t *testing.T) {
	w := New(0)
	if err := w.Add(Rule{Port: 80, Protocol: "icmp", Weight: 3}); err == nil {
		t.Fatal("expected error for unknown protocol")
	}
}

func TestAdd_NegativeWeightReturnsError(t *testing.T) {
	w := New(0)
	if err := w.Add(Rule{Port: 80, Protocol: "tcp", Weight: -1}); err == nil {
		t.Fatal("expected error for negative weight")
	}
}

func TestWeight_ProtocolDistinct(t *testing.T) {
	w := New(0)
	_ = w.Add(Rule{Port: 53, Protocol: "tcp", Weight: 7})
	_ = w.Add(Rule{Port: 53, Protocol: "udp", Weight: 3})
	if got := w.Weight(diff(53, "tcp", "opened")); got != 7 {
		t.Fatalf("tcp: expected 7, got %d", got)
	}
	if got := w.Weight(diff(53, "udp", "opened")); got != 3 {
		t.Fatalf("udp: expected 3, got %d", got)
	}
}

func TestWeighDiffs_MapsEachDiff(t *testing.T) {
	w := New(1)
	_ = w.Add(Rule{Port: 22, Protocol: "tcp", Weight: 9})
	diffs := []scanner.Diff{
		diff(22, "tcp", "opened"),
		diff(9999, "udp", "closed"),
	}
	result := w.WeighDiffs(diffs)
	if result[diffs[0]] != 9 {
		t.Fatalf("port 22: expected 9, got %d", result[diffs[0]])
	}
	if result[diffs[1]] != 1 {
		t.Fatalf("port 9999: expected 1 (default), got %d", result[diffs[1]])
	}
}
