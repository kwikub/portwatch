package portmatch

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func diffs() []scanner.Diff {
	return []scanner.Diff{
		{Port: 80, Protocol: "tcp", State: scanner.StateOpened},
		{Port: 443, Protocol: "tcp", State: scanner.StateOpened},
		{Port: 8080, Protocol: "tcp", State: scanner.StateOpened},
		{Port: 53, Protocol: "udp", State: scanner.StateOpened},
	}
}

func TestMatch_NoRulesReturnsAll(t *testing.T) {
	m := New()
	got := m.Match(diffs())
	if len(got) != 4 {
		t.Fatalf("expected 4 diffs, got %d", len(got))
	}
}

func TestMatch_SinglePortTCP(t *testing.T) {
	m := New()
	_ = m.Add("80", "tcp")
	got := m.Match(diffs())
	if len(got) != 1 || got[0].Port != 80 {
		t.Fatalf("expected only port 80, got %v", got)
	}
}

func TestMatch_PortRange(t *testing.T) {
	m := New()
	_ = m.Add("8000-8999", "tcp")
	got := m.Match(diffs())
	if len(got) != 1 || got[0].Port != 8080 {
		t.Fatalf("expected port 8080, got %v", got)
	}
}

func TestMatch_WildcardProtocol(t *testing.T) {
	m := New()
	_ = m.Add("53", "*")
	got := m.Match(diffs())
	if len(got) != 1 || got[0].Port != 53 {
		t.Fatalf("expected port 53, got %v", got)
	}
}

func TestMatch_ProtocolMismatchExcludes(t *testing.T) {
	m := New()
	_ = m.Add("53", "tcp")
	got := m.Match(diffs())
	if len(got) != 0 {
		t.Fatalf("expected no matches, got %v", got)
	}
}

func TestAdd_InvalidPort(t *testing.T) {
	m := New()
	if err := m.Add("99999", "tcp"); err == nil {
		t.Fatal("expected error for out-of-range port")
	}
}

func TestAdd_InvalidRange(t *testing.T) {
	m := New()
	if err := m.Add("900-100", "tcp"); err == nil {
		t.Fatal("expected error for inverted range")
	}
}

func TestAdd_InvalidProtocol(t *testing.T) {
	m := New()
	if err := m.Add("80", "icmp"); err == nil {
		t.Fatal("expected error for unsupported protocol")
	}
}

func TestLoad_ValidFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "rules.txt")
	content := "# comment\n80 tcp\n8000-8999 *\n\n53 udp\n"
	_ = os.WriteFile(p, []byte(content), 0o644)

	m := New()
	if err := Load(m, p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := m.Match(diffs())
	// should match 80/tcp, 8080/tcp, 53/udp
	if len(got) != 3 {
		t.Fatalf("expected 3 matches, got %d", len(got))
	}
}

func TestLoad_MissingFileIsNoop(t *testing.T) {
	m := New()
	if err := Load(m, "/no/such/file.txt"); err != nil {
		t.Fatalf("expected nil for missing file, got %v", err)
	}
}

func TestLoad_MalformedLineReturnsError(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "bad.txt")
	_ = os.WriteFile(p, []byte("80\n"), 0o644)

	m := New()
	if err := Load(m, p); err == nil {
		t.Fatal("expected error for malformed line")
	}
}
