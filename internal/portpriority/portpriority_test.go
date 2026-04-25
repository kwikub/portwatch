package portpriority

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func d(port int, proto string) scanner.Diff {
	return scanner.Diff{Port: port, Protocol: proto}
}

func TestNew_PanicsOnNegativeDefault(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for negative default")
		}
	}()
	New(-1)
}

func TestPriority_ReturnsDefaultWhenNoRule(t *testing.T) {
	p := New(5)
	got := p.Priority(d(80, "tcp"))
	if got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestAdd_ValidRule(t *testing.T) {
	p := New(0)
	if err := p.Add(443, "tcp", 10); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := p.Priority(d(443, "tcp"))
	if got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}

func TestAdd_InvalidPort(t *testing.T) {
	p := New(0)
	if err := p.Add(0, "tcp", 1); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := p.Add(70000, "tcp", 1); err == nil {
		t.Fatal("expected error for port 70000")
	}
}

func TestAdd_InvalidProtocol(t *testing.T) {
	p := New(0)
	if err := p.Add(80, "icmp", 1); err == nil {
		t.Fatal("expected error for invalid protocol")
	}
}

func TestAdd_NegativePriorityReturnsError(t *testing.T) {
	p := New(0)
	if err := p.Add(80, "tcp", -1); err == nil {
		t.Fatal("expected error for negative priority")
	}
}

func TestPriority_ProtocolDistinct(t *testing.T) {
	p := New(0)
	_ = p.Add(53, "tcp", 3)
	_ = p.Add(53, "udp", 7)

	if got := p.Priority(d(53, "tcp")); got != 3 {
		t.Fatalf("tcp: expected 3, got %d", got)
	}
	if got := p.Priority(d(53, "udp")); got != 7 {
		t.Fatalf("udp: expected 7, got %d", got)
	}
}

func TestSort_OrdersByPriorityDescending(t *testing.T) {
	p := New(1)
	_ = p.Add(22, "tcp", 5)
	_ = p.Add(443, "tcp", 10)

	input := []scanner.Diff{d(22, "tcp"), d(9999, "udp"), d(443, "tcp")}
	sorted := p.Sort(input)

	if sorted[0].Port != 443 {
		t.Fatalf("expected 443 first, got %d", sorted[0].Port)
	}
	if sorted[1].Port != 22 {
		t.Fatalf("expected 22 second, got %d", sorted[1].Port)
	}
	if sorted[2].Port != 9999 {
		t.Fatalf("expected 9999 last, got %d", sorted[2].Port)
	}
}

func TestSort_DoesNotModifyOriginal(t *testing.T) {
	p := New(0)
	_ = p.Add(80, "tcp", 9)

	input := []scanner.Diff{d(9999, "tcp"), d(80, "tcp")}
	_ = p.Sort(input)

	if input[0].Port != 9999 {
		t.Fatal("original slice was modified")
	}
}
