package portname

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLookup_WellKnownPort(t *testing.T) {
	r := New()
	if got := r.Lookup(80, "tcp"); got != "http" {
		t.Fatalf("expected http, got %q", got)
	}
}

func TestLookup_UnknownReturnsEmpty(t *testing.T) {
	r := New()
	if got := r.Lookup(9999, "tcp"); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}

func TestLookup_ProtocolDistinct(t *testing.T) {
	r := New()
	if r.Lookup(53, "tcp") == "" || r.Lookup(53, "udp") == "" {
		t.Fatal("expected dns for both tcp and udp")
	}
}

func TestRegister_AddsCustomEntry(t *testing.T) {
	r := New()
	r.Register(9200, "tcp", "elasticsearch")
	if got := r.Lookup(9200, "tcp"); got != "elasticsearch" {
		t.Fatalf("expected elasticsearch, got %q", got)
	}
}

func TestLabel_FallsBackToPortProto(t *testing.T) {
	r := New()
	if got := r.Label(9999, "tcp"); got != "9999/tcp" {
		t.Fatalf("expected 9999/tcp, got %q", got)
	}
}

func TestLabel_UsesServiceName(t *testing.T) {
	r := New()
	if got := r.Label(443, "tcp"); got != "https" {
		t.Fatalf("expected https, got %q", got)
	}
}

func TestLoad_ReadsFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "ports.txt")
	content := "# comment\n9200/tcp elasticsearch\n9300/tcp elasticsearch-cluster\n"
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	r := New()
	if err := Load(r, p); err != nil {
		t.Fatal(err)
	}
	if got := r.Lookup(9200, "tcp"); got != "elasticsearch" {
		t.Fatalf("expected elasticsearch, got %q", got)
	}
}

func TestLoad_MissingFileReturnsError(t *testing.T) {
	r := New()
	if err := Load(r, "/nonexistent/ports.txt"); err == nil {
		t.Fatal("expected error for missing file")
	}
}
