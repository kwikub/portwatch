package portgroup

import (
	"testing"
)

func TestAdd_ValidGroup(t *testing.T) {
	r := New()
	err := r.Add("web", []Entry{{80, "tcp"}, {443, "tcp"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	g, ok := r.Get("web")
	if !ok {
		t.Fatal("expected group to exist")
	}
	if len(g.Ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(g.Ports))
	}
}

func TestAdd_EmptyNameReturnsError(t *testing.T) {
	r := New()
	if err := r.Add("", []Entry{{80, "tcp"}}); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestAdd_InvalidPortReturnsError(t *testing.T) {
	r := New()
	if err := r.Add("bad", []Entry{{0, "tcp"}}); err == nil {
		t.Fatal("expected error for port 0")
	}
}

func TestAdd_InvalidProtocolReturnsError(t *testing.T) {
	r := New()
	if err := r.Add("bad", []Entry{{80, "sctp"}}); err == nil {
		t.Fatal("expected error for unknown protocol")
	}
}

func TestContains_Found(t *testing.T) {
	r := New()
	_ = r.Add("dns", []Entry{{53, "udp"}})
	if !r.Contains("dns", 53, "udp") {
		t.Fatal("expected Contains to return true")
	}
}

func TestContains_WrongProtocol(t *testing.T) {
	r := New()
	_ = r.Add("dns", []Entry{{53, "udp"}})
	if r.Contains("dns", 53, "tcp") {
		t.Fatal("expected Contains to return false for wrong protocol")
	}
}

func TestContains_MissingGroup(t *testing.T) {
	r := New()
	if r.Contains("nope", 80, "tcp") {
		t.Fatal("expected false for missing group")
	}
}

func TestNames_ReturnsAll(t *testing.T) {
	r := New()
	_ = r.Add("a", []Entry{{80, "tcp"}})
	_ = r.Add("b", []Entry{{443, "tcp"}})
	names := r.Names()
	if len(names) != 2 {
		t.Fatalf("expected 2 names, got %d", len(names))
	}
}
