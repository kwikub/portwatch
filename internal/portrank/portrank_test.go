package portrank

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func d(port int, proto string) scanner.Diff {
	return scanner.Diff{Port: port, Protocol: proto}
}

func TestAdd_ValidEntry(t *testing.T) {
	r := New()
	if err := r.Add(80, "tcp", 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAdd_InvalidPort(t *testing.T) {
	r := New()
	if err := r.Add(0, "tcp", 1); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := r.Add(70000, "tcp", 1); err == nil {
		t.Fatal("expected error for port 70000")
	}
}

func TestAdd_InvalidProtocol(t *testing.T) {
	r := New()
	if err := r.Add(443, "icmp", 1); err == nil {
		t.Fatal("expected error for unknown protocol")
	}
}

func TestAdd_NegativeRankReturnsError(t *testing.T) {
	r := New()
	if err := r.Add(22, "tcp", -1); err == nil {
		t.Fatal("expected error for negative rank")
	}
}

func TestRank_DefaultWhenNoRule(t *testing.T) {
	r := New()
	got := r.Rank(d(9999, "tcp"))
	if got != 100 {
		t.Fatalf("want 100, got %d", got)
	}
}

func TestRank_ReturnsRegisteredRank(t *testing.T) {
	r := New()
	_ = r.Add(22, "tcp", 5)
	if got := r.Rank(d(22, "tcp")); got != 5 {
		t.Fatalf("want 5, got %d", got)
	}
}

func TestRank_ProtocolDistinct(t *testing.T) {
	r := New()
	_ = r.Add(53, "tcp", 10)
	_ = r.Add(53, "udp", 2)
	if got := r.Rank(d(53, "tcp")); got != 10 {
		t.Fatalf("tcp want 10, got %d", got)
	}
	if got := r.Rank(d(53, "udp")); got != 2 {
		t.Fatalf("udp want 2, got %d", got)
	}
}

func TestAdd_OverwritesExisting(t *testing.T) {
	r := New()
	_ = r.Add(80, "tcp", 50)
	_ = r.Add(80, "tcp", 3)
	if got := r.Rank(d(80, "tcp")); got != 3 {
		t.Fatalf("want 3 after overwrite, got %d", got)
	}
}

func TestSort_OrdersByRankAscending(t *testing.T) {
	r := New()
	_ = r.Add(22, "tcp", 1)
	_ = r.Add(80, "tcp", 50)
	_ = r.Add(443, "tcp", 10)

	input := []scanner.Diff{d(80, "tcp"), d(443, "tcp"), d(22, "tcp")}
	sorted := r.Sort(input)

	want := []int{22, 443, 80}
	for i, w := range want {
		if sorted[i].Port != w {
			t.Fatalf("position %d: want port %d, got %d", i, w, sorted[i].Port)
		}
	}
}

func TestSort_UnknownPortsGetDefaultRank(t *testing.T) {
	r := New()
	_ = r.Add(22, "tcp", 1)

	input := []scanner.Diff{d(9999, "tcp"), d(22, "tcp")}
	sorted := r.Sort(input)

	if sorted[0].Port != 22 {
		t.Fatalf("want port 22 first, got %d", sorted[0].Port)
	}
}

func TestSort_DoesNotMutateInput(t *testing.T) {
	r := New()
	_ = r.Add(22, "tcp", 1)
	input := []scanner.Diff{d(80, "tcp"), d(22, "tcp")}
	_ = r.Sort(input)
	if input[0].Port != 80 {
		t.Fatal("original slice was mutated")
	}
}
