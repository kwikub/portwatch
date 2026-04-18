package tags

import (
	"testing"
)

func TestAdd_ValidEntry(t *testing.T) {
	r := New()
	if err := r.Add(80, "tcp", "HTTP"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	label, ok := r.Lookup(80, "tcp")
	if !ok {
		t.Fatal("expected label to be found")
	}
	if label != "HTTP" {
		t.Fatalf("expected HTTP, got %s", label)
	}
}

func TestAdd_InvalidPort(t *testing.T) {
	r := New()
	if err := r.Add(0, "tcp", "bad"); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := r.Add(70000, "tcp", "bad"); err == nil {
		t.Fatal("expected error for port 70000")
	}
}

func TestAdd_BadProtocol(t *testing.T) {
	r := New()
	if err := r.Add(443, "sctp", "label"); err == nil {
		t.Fatal("expected error for unsupported protocol")
	}
}

func TestAdd_EmptyLabel(t *testing.T) {
	r := New()
	if err := r.Add(22, "tcp", "  "); err == nil {
		t.Fatal("expected error for blank label")
	}
}

func TestLookup_MissingReturnsNotFound(t *testing.T) {
	r := New()
	_, ok := r.Lookup(9999, "tcp")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestLookup_CaseInsensitiveProtocol(t *testing.T) {
	r := New()
	_ = r.Add(53, "UDP", "DNS")
	label, ok := r.Lookup(53, "udp")
	if !ok || label != "DNS" {
		t.Fatalf("expected DNS, got %q ok=%v", label, ok)
	}
}

func TestAll_ReturnsAllTags(t *testing.T) {
	r := New()
	_ = r.Add(22, "tcp", "SSH")
	_ = r.Add(80, "tcp", "HTTP")
	if len(r.All()) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(r.All()))
	}
}
