package audit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/audit"
	"github.com/user/portwatch/internal/scanner"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "audit.jsonl")
}

func diffs() []scanner.Diff {
	return []scanner.Diff{
		{Proto: "tcp", Port: 8080, Closed: false},
		{Proto: "udp", Port: 53, Closed: true},
	}
}

func TestRecord_WritesEntries(t *testing.T) {
	p := tempPath(t)
	a := audit.New(p)
	if err := a.Record(diffs()); err != nil {
		t.Fatalf("Record: %v", err)
	}
	entries, err := audit.ReadAll(p)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Port != 8080 || entries[0].State != "opened" {
		t.Errorf("unexpected first entry: %+v", entries[0])
	}
	if entries[1].Port != 53 || entries[1].State != "closed" {
		t.Errorf("unexpected second entry: %+v", entries[1])
	}
}

func TestRecord_EmptyDiffsWritesNothing(t *testing.T) {
	p := tempPath(t)
	a := audit.New(p)
	if err := a.Record(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		t.Error("expected no file to be created for empty diffs")
	}
}

func TestReadAll_MissingFileReturnsEmpty(t *testing.T) {
	entries, err := audit.ReadAll("/tmp/portwatch_nonexistent_audit.jsonl")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(entries))
	}
}

func TestRecord_Appends(t *testing.T) {
	p := tempPath(t)
	a := audit.New(p)
	_ = a.Record(diffs())
	_ = a.Record(diffs())
	entries, _ := audit.ReadAll(p)
	if len(entries) != 4 {
		t.Fatalf("expected 4 entries after two records, got %d", len(entries))
	}
}
