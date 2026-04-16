package state_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "state.json")
}

func TestLoad_EmptyWhenMissing(t *testing.T) {
	s := state.New(tempPath(t))
	snap, err := s.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(snap.Ports) != 0 {
		t.Errorf("expected empty snapshot, got %d ports", len(snap.Ports))
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	s := state.New(tempPath(t))

	orig := scanner.Snapshot{
		Ports: []scanner.Port{
			{Number: 80, Protocol: "tcp"},
			{Number: 443, Protocol: "tcp"},
		},
	}

	if err := s.Save(orig); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := s.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(loaded.Ports) != len(orig.Ports) {
		t.Fatalf("expected %d ports, got %d", len(orig.Ports), len(loaded.Ports))
	}
	for i, p := range loaded.Ports {
		if p != orig.Ports[i] {
			t.Errorf("port mismatch at %d: got %+v, want %+v", i, p, orig.Ports[i])
		}
	}
}

func TestSave_LeavesNoTempFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")
	s := state.New(path)

	if err := s.Save(scanner.Snapshot{}); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		t.Errorf("expected 1 file in dir, got %d", len(entries))
	}
}
