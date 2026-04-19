package portgroup

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "groups.conf")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoad_ValidFile(t *testing.T) {
	p := writeFile(t, "web 80/tcp\nweb 443/tcp\ndns 53/udp\n")
	r := New()
	if err := Load(p, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Contains("web", 80, "tcp") {
		t.Error("expected web group to contain 80/tcp")
	}
	if !r.Contains("dns", 53, "udp") {
		t.Error("expected dns group to contain 53/udp")
	}
}

func TestLoad_MissingFileIsNoop(t *testing.T) {
	r := New()
	if err := Load("/nonexistent/path.conf", r); err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
}

func TestLoad_CommentsAndBlankLines(t *testing.T) {
	p := writeFile(t, "# comment\n\nweb 80/tcp\n")
	r := New()
	if err := Load(p, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Contains("web", 80, "tcp") {
		t.Error("expected web group")
	}
}

func TestLoad_MalformedLineReturnsError(t *testing.T) {
	p := writeFile(t, "web\n")
	r := New()
	if err := Load(p, r); err == nil {
		t.Fatal("expected error for malformed line")
	}
}

func TestLoad_InvalidPortReturnsError(t *testing.T) {
	p := writeFile(t, "web abc/tcp\n")
	r := New()
	if err := Load(p, r); err == nil {
		t.Fatal("expected error for non-numeric port")
	}
}
