package portscope_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/portscope"
)

func writeFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "scope.conf")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestContains_NoRangesAlwaysTrue(t *testing.T) {
	s := portscope.New()
	if !s.Contains("tcp", 80) {
		t.Error("expected true for unconfigured scope")
	}
}

func TestAdd_InvalidProtocolReturnsError(t *testing.T) {
	s := portscope.New()
	if err := s.Add("icmp", 1, 100); err == nil {
		t.Error("expected error for unsupported protocol")
	}
}

func TestAdd_InvalidRangeReturnsError(t *testing.T) {
	s := portscope.New()
	if err := s.Add("tcp", 500, 100); err == nil {
		t.Error("expected error when lo > hi")
	}
	if err := s.Add("tcp", 0, 100); err == nil {
		t.Error("expected error when lo < 1")
	}
}

func TestContains_PortInsideRange(t *testing.T) {
	s := portscope.New()
	_ = s.Add("tcp", 80, 443)
	if !s.Contains("tcp", 80) {
		t.Error("expected port 80/tcp to be in scope")
	}
	if !s.Contains("tcp", 443) {
		t.Error("expected port 443/tcp to be in scope")
	}
}

func TestContains_PortOutsideRange(t *testing.T) {
	s := portscope.New()
	_ = s.Add("tcp", 80, 443)
	if s.Contains("tcp", 22) {
		t.Error("expected port 22/tcp to be out of scope")
	}
}

func TestContains_ProtocolMismatch(t *testing.T) {
	s := portscope.New()
	_ = s.Add("tcp", 53, 53)
	if s.Contains("udp", 53) {
		t.Error("expected udp/53 to be out of scope when only tcp/53 registered")
	}
}

func TestSize_ReflectsTotalPorts(t *testing.T) {
	s := portscope.New()
	_ = s.Add("tcp", 1, 1024)
	_ = s.Add("udp", 53, 53)
	if got := s.Size(); got != 1025 {
		t.Errorf("expected size 1025, got %d", got)
	}
}

func TestLoad_ValidFile(t *testing.T) {
	p := writeFile(t, "# comment\ntcp 1 1024\nudp 53 53\n")
	s := portscope.New()
	if err := portscope.Load(s, p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.Contains("tcp", 80) {
		t.Error("expected tcp/80 in scope after load")
	}
	if !s.Contains("udp", 53) {
		t.Error("expected udp/53 in scope after load")
	}
}

func TestLoad_MissingFileIsNoop(t *testing.T) {
	s := portscope.New()
	if err := portscope.Load(s, "/nonexistent/scope.conf"); err != nil {
		t.Errorf("expected nil for missing file, got %v", err)
	}
}

func TestLoad_MalformedLineReturnsError(t *testing.T) {
	p := writeFile(t, "tcp 80\n")
	s := portscope.New()
	if err := portscope.Load(s, p); err == nil {
		t.Error("expected error for malformed line")
	}
}
