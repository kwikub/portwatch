package baseline

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "baseline.json")
}

func snap(ports []scanner.Port) *scanner.Snapshot {
	return &scanner.Snapshot{Ports: ports, At: time.Now()}
}

func TestNew_EmptyWhenMissing(t *testing.T) {
	b, err := New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(b.ports) != 0 {
		t.Errorf("expected empty baseline, got %d ports", len(b.ports))
	}
}

func TestRecord_PersistsAndReloads(t *testing.T) {
	path := tempPath(t)
	b, _ := New(path)
	ports := []scanner.Port{{Proto: "tcp", Addr: "0.0.0.0:80"}, {Proto: "udp", Addr: "0.0.0.0:53"}}
	if err := b.Record(snap(ports)); err != nil {
		t.Fatalf("Record: %v", err)
	}
	b2, err := New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if len(b2.ports) != 2 {
		t.Errorf("expected 2 ports after reload, got %d", len(b2.ports))
	}
}

func TestDeviations_AddedPort(t *testing.T) {
	b, _ := New(tempPath(t))
	base := []scanner.Port{{Proto: "tcp", Addr: "0.0.0.0:80"}}
	_ = b.Record(snap(base))

	newSnap := snap([]scanner.Port{
		{Proto: "tcp", Addr: "0.0.0.0:80"},
		{Proto: "tcp", Addr: "0.0.0.0:443"},
	})
	added, removed := b.Deviations(newSnap)
	if len(added) != 1 {
		t.Errorf("expected 1 added, got %d", len(added))
	}
	if len(removed) != 0 {
		t.Errorf("expected 0 removed, got %d", len(removed))
	}
}

func TestDeviations_NoDeviations(t *testing.T) {
	b, _ := New(tempPath(t))
	ports := []scanner.Port{{Proto: "tcp", Addr: "0.0.0.0:22"}}
	_ = b.Record(snap(ports))
	added, removed := b.Deviations(snap(ports))
	if len(added) != 0 || len(removed) != 0 {
		t.Errorf("expected no deviations, got added=%d removed=%d", len(added), len(removed))
	}
}

func TestNew_InvalidJSON(t *testing.T) {
	path := tempPath(t)
	_ = os.WriteFile(path, []byte("not-json"), 0o644)
	_, err := New(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
