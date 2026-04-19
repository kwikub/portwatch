package portmap

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	m := New()
	e := Entry{Port: 80, Protocol: "tcp", Name: "http", Group: "web", Tag: "public"}
	m.Set(e)
	got, ok := m.Get(80, "tcp")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if got.Name != "http" {
		t.Errorf("name: got %q want %q", got.Name, "http")
	}
}

func TestGet_Missing(t *testing.T) {
	m := New()
	_, ok := m.Get(9999, "tcp")
	if ok {
		t.Fatal("expected no entry")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	m := New()
	m.Set(Entry{Port: 443, Protocol: "tcp", Name: "https"})
	m.Delete(443, "tcp")
	_, ok := m.Get(443, "tcp")
	if ok {
		t.Fatal("expected entry to be deleted")
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	m := New()
	m.Set(Entry{Port: 22, Protocol: "tcp", Name: "ssh"})
	m.Set(Entry{Port: 53, Protocol: "udp", Name: "dns"})
	if len(m.All()) != 2 {
		t.Errorf("expected 2 entries, got %d", len(m.All()))
	}
}

func TestEntry_String(t *testing.T) {
	e := Entry{Port: 80, Protocol: "tcp", Name: "http"}
	if e.String() != "80/tcp (http)" {
		t.Errorf("unexpected string: %s", e.String())
	}
	e2 := Entry{Port: 9999, Protocol: "tcp"}
	if e2.String() != "9999/tcp" {
		t.Errorf("unexpected string: %s", e2.String())
	}
}

func writeFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "portmap.tsv")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoad_ValidFile(t *testing.T) {
	path := writeFile(t, "# comment\n80\ttcp\thttp\tweb\tpublic\n53\tudp\tdns\t\t\n")
	m := New()
	if err := Load(m, path); err != nil {
		t.Fatal(err)
	}
	e, ok := m.Get(80, "tcp")
	if !ok || e.Name != "http" || e.Group != "web" || e.Tag != "public" {
		t.Errorf("unexpected entry: %+v", e)
	}
	e2, ok := m.Get(53, "udp")
	if !ok || e2.Name != "dns" {
		t.Errorf("unexpected entry: %+v", e2)
	}
}

func TestLoad_MissingFileIsNoop(t *testing.T) {
	m := New()
	if err := Load(m, "/nonexistent/portmap.tsv"); err != nil {
		t.Fatal(err)
	}
	if len(m.All()) != 0 {
		t.Error("expected empty map")
	}
}

func TestLoad_MalformedPortReturnsError(t *testing.T) {
	path := writeFile(t, "bad\ttcp\n")
	m := New()
	if err := Load(m, path); err == nil {
		t.Fatal("expected error for malformed port")
	}
}
