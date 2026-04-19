package export

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuilder_DefaultsToStdoutJSON(t *testing.T) {
	ex, closer, err := NewBuilder().Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if closer != nil {
		t.Error("expected nil closer for stdout")
	}
	if ex == nil {
		t.Fatal("expected non-nil exporter")
	}
	if ex.format != FormatJSON {
		t.Errorf("want JSON, got %s", ex.format)
	}
}

func TestBuilder_WritesToFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "out.json")
	ex, closer, err := NewBuilder().WithPath(p).Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if closer == nil {
		t.Fatal("expected non-nil closer for file")
	}
	defer closer.Close()
	if ex == nil {
		t.Fatal("expected non-nil exporter")
	}
	if _, err := os.Stat(p); err != nil {
		t.Errorf("file not created: %v", err)
	}
}

func TestBuilder_CSVFormat(t *testing.T) {
	ex, _, err := NewBuilder().WithFormat(FormatCSV).Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ex.format != FormatCSV {
		t.Errorf("want CSV, got %s", ex.format)
	}
}
