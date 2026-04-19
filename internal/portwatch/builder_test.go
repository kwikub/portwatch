package portwatch_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/portwatch"
)

func TestBuilder_DefaultsBuild(t *testing.T) {
	dir := t.TempDir()
	_, err := portwatch.NewBuilder().
		WithStatePath(filepath.Join(dir, "state.json")).
		Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
}

func TestBuilder_WithOutput(t *testing.T) {
	dir := t.TempDir()
	var buf bytes.Buffer
	_, err := portwatch.NewBuilder().
		WithStatePath(filepath.Join(dir, "state.json")).
		WithOutput(&buf).
		Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
}

func TestBuilder_WithRange(t *testing.T) {
	dir := t.TempDir()
	_, err := portwatch.NewBuilder().
		WithStatePath(filepath.Join(dir, "state.json")).
		WithRange(8000, 9000).
		Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
}

func TestBuilder_InvalidStateDir(t *testing.T) {
	// Use a path inside a non-existent directory to force state.New to fail.
	_, err := portwatch.NewBuilder().
		WithStatePath(filepath.Join(os.DevNull, "no", "such", "dir", "state.json")).
		Build()
	if err == nil {
		t.Fatal("expected error for invalid state path, got nil")
	}
}
