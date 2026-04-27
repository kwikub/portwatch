package portroute

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func TestAdd_ValidEntry(t *testing.T) {
	r := New()
	if err := r.Add(443, "tcp", "api-gateway", "web"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAdd_InvalidPort(t *testing.T) {
	r := New()
	if err := r.Add(0, "tcp", "svc", ""); err == nil {
		t.Fatal("expected error for port 0")
	}
}

func TestAdd_InvalidProtocol(t *testing.T) {
	r := New()
	if err := r.Add(80, "sctp", "svc", ""); err == nil {
		t.Fatal("expected error for unknown protocol")
	}
}

func TestAdd_EmptyTargetReturnsError(t *testing.T) {
	r := New()
	if err := r.Add(80, "tcp", "", ""); err == nil {
		t.Fatal("expected error for empty target")
	}
}

func TestLookup_Found(t *testing.T) {
	r := New()
	_ = r.Add(8080, "tcp", "backend", "app")
	d := scanner.Diff{Port: 8080, Protocol: "tcp"}
	route, ok := r.Lookup(d)
	if !ok {
		t.Fatal("expected route to be found")
	}
	if route.Target != "backend" {
		t.Errorf("got target %q, want %q", route.Target, "backend")
	}
	if route.Group != "app" {
		t.Errorf("got group %q, want %q", route.Group, "app")
	}
}

func TestLookup_NotFound(t *testing.T) {
	r := New()
	d := scanner.Diff{Port: 9999, Protocol: "tcp"}
	_, ok := r.Lookup(d)
	if ok {
		t.Fatal("expected no route for unregistered port")
	}
}

func TestLookup_ProtocolDistinct(t *testing.T) {
	r := New()
	_ = r.Add(53, "tcp", "dns-tcp", "")
	_ = r.Add(53, "udp", "dns-udp", "")

	tcpRoute, _ := r.Lookup(scanner.Diff{Port: 53, Protocol: "tcp"})
	udpRoute, _ := r.Lookup(scanner.Diff{Port: 53, Protocol: "udp"})

	if tcpRoute.Target != "dns-tcp" {
		t.Errorf("tcp target: got %q, want dns-tcp", tcpRoute.Target)
	}
	if udpRoute.Target != "dns-udp" {
		t.Errorf("udp target: got %q, want dns-udp", udpRoute.Target)
	}
}

func TestAll_ReturnsRegisteredRoutes(t *testing.T) {
	r := New()
	_ = r.Add(80, "tcp", "http", "web")
	_ = r.Add(443, "tcp", "https", "web")
	if len(r.All()) != 2 {
		t.Errorf("expected 2 routes, got %d", len(r.All()))
	}
}
