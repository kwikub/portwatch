package rotation

import (
	"os"
	"path/filepath"
	"testing"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "rotation-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func TestNew_CreatesFile(t *testing.T) {
	dir := tempDir(t)
	r, err := New(Options{Dir: dir, Prefix: "pw-"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer r.Close()
	matches, _ := filepath.Glob(filepath.Join(dir, "pw-*.log"))
	if len(matches) != 1 {
		t.Fatalf("expected 1 file, got %d", len(matches))
	}
}

func TestWrite_PersistsData(t *testing.T) {
	dir := tempDir(t)
	r, _ := New(Options{Dir: dir, Prefix: "pw-"})
	defer r.Close()
	_, err := r.Write([]byte("hello\n"))
	if err != nil {
		t.Fatalf("write error: %v", err)
	}
}

func TestRotate_CreatesNewFile(t *testing.T) {
	dir := tempDir(t)
	r, _ := New(Options{Dir: dir, Prefix: "pw-", MaxFiles: 5})
	defer r.Close()
	if err := r.Rotate(); err != nil {
		t.Fatalf("rotate error: %v", err)
	}
	matches, _ := filepath.Glob(filepath.Join(dir, "pw-*.log"))
	if len(matches) < 1 {
		t.Fatal("expected at least one file after rotate")
	}
}

func TestRotate_PrunesOldFiles(t *testing.T) {
	dir := tempDir(t)
	max := 3
	r, _ := New(Options{Dir: dir, Prefix: "pw-", MaxFiles: max})
	defer r.Close()
	for i := 0; i < max+2; i++ {
		if err := r.Rotate(); err != nil {
			t.Fatalf("rotate %d: %v", i, err)
		}
	}
	matches, _ := filepath.Glob(filepath.Join(dir, "pw-*.log"))
	if len(matches) > max {
		t.Fatalf("expected at most %d files, got %d", max, len(matches))
	}
}

func TestNew_DefaultMaxFiles(t *testing.T) {
	dir := tempDir(t)
	r, err := New(Options{Dir: dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer r.Close()
	if r.opts.MaxFiles != 5 {
		t.Fatalf("expected default MaxFiles=5, got %d", r.opts.MaxFiles)
	}
}
