package portmap

import (
	"os"
	"path/filepath"
	"testing"
)

func writeMapFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "portmap.conf")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestFromFile_MissingFileReturnsEmpty(t *testing.T) {
	m, err := FromFile("/nonexistent/portmap.conf")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.All()) != 0 {
		t.Errorf("expected empty map, got %d entries", len(m.All()))
	}
}

func TestFromFile_LoadsEntries(t *testing.T) {
	path := writeMapFile(t, "8080 tcp internal-http\n443 tcp secure-web\n")
	m, err := FromFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.All()) != 2 {
		t.Errorf("expected 2 entries, got %d", len(m.All()))
	}
}

func TestFromFile_InvalidContentReturnsError(t *testing.T) {
	path := writeMapFile(t, "notaport tcp label\n")
	_, err := FromFile(path)
	if err == nil {
		t.Error("expected error for invalid port, got nil")
	}
}
