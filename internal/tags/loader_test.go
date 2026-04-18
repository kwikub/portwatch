package tags

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "tags.conf")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
	return p
}

func TestLoad_ValidFile(t *testing.T) {
	p := writeFile(t, "# comment\n22/tcp SSH\n80/tcp HTTP\n")
	r, err := Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.All()) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(r.All()))
	}
	if label, _ := r.Lookup(22, "tcp"); label != "SSH" {
		t.Fatalf("expected SSH, got %s", label)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/tags.conf")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_MalformedLine(t *testing.T) {
	p := writeFile(t, "badline\n")
	_, err := Load(p)
	if err == nil {
		t.Fatal("expected error for malformed line")
	}
}

func TestLoad_MultiWordLabel(t *testing.T) {
	p := writeFile(t, "443/tcp HTTPS Web\n")
	r, err := Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	label, _ := r.Lookup(443, "tcp")
	if !strings.Contains(label, "Web") {
		t.Fatalf("expected multi-word label, got %q", label)
	}
}

func TestLoad_BlankLinesIgnored(t *testing.T) {
	p := writeFile(t, "\n\n53/udp DNS\n\n")
	r, err := Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.All()) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.All()))
	}
}
