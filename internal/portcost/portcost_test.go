package portcost_test

import (
	"testing"

	"github.com/user/portwatch/internal/portcost"
	"github.com/user/portwatch/internal/scanner"
)

func diff(port int, proto, state string) scanner.Diff {
	return scanner.Diff{Port: port, Protocol: proto, State: state}
}

func TestNew_PanicsOnNegativeDefault(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	portcost.New(-1)
}

func TestCost_ReturnsDefaultWhenNoRule(t *testing.T) {
	r := portcost.New(5)
	got := r.Cost(diff(8080, "tcp", "opened"))
	if got != 5 {
		t.Fatalf("want 5, got %d", got)
	}
}

func TestAdd_ValidRule(t *testing.T) {
	r := portcost.New(1)
	if err := r.Add(443, "tcp", 10); err != nil {
		t.Fatal(err)
	}
	got := r.Cost(diff(443, "tcp", "opened"))
	if got != 10 {
		t.Fatalf("want 10, got %d", got)
	}
}

func TestAdd_InvalidPort(t *testing.T) {
	r := portcost.New(1)
	if err := r.Add(0, "tcp", 5); err == nil {
		t.Fatal("expected error for port 0")
	}
}

func TestAdd_InvalidProtocol(t *testing.T) {
	r := portcost.New(1)
	if err := r.Add(80, "sctp", 5); err == nil {
		t.Fatal("expected error for unknown protocol")
	}
}

func TestAdd_NegativeCostReturnsError(t *testing.T) {
	r := portcost.New(1)
	if err := r.Add(80, "tcp", -1); err == nil {
		t.Fatal("expected error for negative cost")
	}
}

func TestAdd_OverwritesExistingRule(t *testing.T) {
	r := portcost.New(1)
	_ = r.Add(22, "tcp", 3)
	_ = r.Add(22, "tcp", 7)
	got := r.Cost(diff(22, "tcp", "opened"))
	if got != 7 {
		t.Fatalf("want 7, got %d", got)
	}
}

func TestCost_ProtocolDistinct(t *testing.T) {
	r := portcost.New(1)
	_ = r.Add(53, "tcp", 2)
	_ = r.Add(53, "udp", 9)
	if got := r.Cost(diff(53, "tcp", "opened")); got != 2 {
		t.Fatalf("tcp: want 2, got %d", got)
	}
	if got := r.Cost(diff(53, "udp", "opened")); got != 9 {
		t.Fatalf("udp: want 9, got %d", got)
	}
}

func TestTotal_SumsCosts(t *testing.T) {
	r := portcost.New(1)
	_ = r.Add(80, "tcp", 3)
	_ = r.Add(443, "tcp", 5)
	diffs := []scanner.Diff{
		diff(80, "tcp", "opened"),
		diff(443, "tcp", "opened"),
		diff(8080, "tcp", "opened"), // uses default cost 1
	}
	got := r.Total(diffs)
	if got != 9 {
		t.Fatalf("want 9, got %d", got)
	}
}

func TestTotal_EmptyDiffsIsZero(t *testing.T) {
	r := portcost.New(5)
	if got := r.Total(nil); got != 0 {
		t.Fatalf("want 0, got %d", got)
	}
}
