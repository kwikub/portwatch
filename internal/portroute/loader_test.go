package portroute

import (
	"os"
	"path/filepath"
	"testing"
)

func writeRouteFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "routes.conf")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeRouteFile: %v", err)
	}
	return p
}

func TestLoad_ValidFile(t *testing.T) {
	p := writeRouteFile(t, "tcp 443 api-gateway web\nudp 53 dns-server\n")
	r := New()
	if err := Load(r, p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.All()) != 2 {
		t.Errorf("expected 2 routes, got %d", len(r.All()))
	}
}

func TestLoad_MissingFileIsNoop(t *testing.T) {
	r := New()
	if err := Load(r, "/nonexistent/routes.conf"); err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
}

func TestLoad_CommentsAndBlankLines(t *testing.T) {
	p := writeRouteFile(t, "# comment\n\ntcp 80 http web\n")
	r := New()
	_ = Load(r, p)
	if len(r.All()) != 1 {
		t.Errorf("expected 1 route, got %d", len(r.All()))
	}
}

func TestLoad_MalformedLineReturnsError(t *testing.T) {
	p := writeRouteFile(t, "tcp\n")
	r := New()
	if err := Load(r, p); err == nil {
		t.Fatal("expected error for malformed line")
	}
}

func TestLoad_InvalidPortReturnsError(t *testing.T) {
	p := writeRouteFile(t, "tcp notaport target\n")
	r := New()
	if err := Load(r, p); err == nil {
		t.Fatal("expected error for invalid port")
	}
}
