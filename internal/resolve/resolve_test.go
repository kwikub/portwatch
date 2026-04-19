package resolve_test

import (
	"testing"

	"github.com/user/portwatch/internal/resolve"
)

func TestLookup_WellKnownPort(t *testing.T) {
	r := resolve.New()
	name := r.Lookup(22, "tcp")
	if name != "ssh" {
		t.Fatalf("expected ssh, got %q", name)
	}
}

func TestLookup_UnknownPortReturnsEmpty(t *testing.T) {
	r := resolve.New()
	name := r.Lookup(9999, "tcp")
	if name != "" {
		t.Fatalf("expected empty, got %q", name)
	}
}

func TestLookup_ProtocolDistinct(t *testing.T) {
	r := resolve.New()
	// port 53 exists for both tcp and udp; ensure they resolve independently
	if r.Lookup(53, "tcp") != "dns" {
		t.Fatal("expected dns for tcp/53")
	}
	if r.Lookup(53, "udp") != "dns" {
		t.Fatal("expected dns for udp/53")
	}
}

func TestRegister_AddsCustomEntry(t *testing.T) {
	r := resolve.New()
	r.Register(9200, "tcp", "elasticsearch")
	if got := r.Lookup(9200, "tcp"); got != "elasticsearch" {
		t.Fatalf("expected elasticsearch, got %q", got)
	}
}

func TestRegister_OverwritesExisting(t *testing.T) {
	r := resolve.New()
	r.Register(80, "tcp", "custom-http")
	if got := r.Lookup(80, "tcp"); got != "custom-http" {
		t.Fatalf("expected custom-http, got %q", got)
	}
}

func TestLookupOrPort_ReturnsNameWhenKnown(t *testing.T) {
	r := resolve.New()
	if got := r.LookupOrPort(443, "tcp"); got != "https" {
		t.Fatalf("expected https, got %q", got)
	}
}

func TestLookupOrPort_FallsBackToPortString(t *testing.T) {
	r := resolve.New()
	if got := r.LookupOrPort(12345, "tcp"); got != "12345" {
		t.Fatalf("expected 12345, got %q", got)
	}
}
