package summary_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/summary"
)

func diffs() []scanner.Diff {
	return []scanner.Diff{
		{Port: 80, Proto: "tcp", State: "opened"},
		{Port: 443, Proto: "tcp", State: "opened"},
		{Port: 8080, Proto: "tcp", State: "closed"},
	}
}

func TestBuild_CountsOpenedAndClosed(t *testing.T) {
	r := summary.Build(5, diffs())
	if len(r.Opened) != 2 {
		t.Fatalf("expected 2 opened, got %d", len(r.Opened))
	}
	if len(r.Closed) != 1 {
		t.Fatalf("expected 1 closed, got %d", len(r.Closed))
	}
	if r.TotalScans != 5 {
		t.Fatalf("expected 5 scans, got %d", r.TotalScans)
	}
}

func TestBuild_EmptyDiffs(t *testing.T) {
	r := summary.Build(3, nil)
	if len(r.Opened) != 0 || len(r.Closed) != 0 {
		t.Fatal("expected empty slices for nil diffs")
	}
}

func TestWrite_ContainsScanCount(t *testing.T) {
	var buf bytes.Buffer
	s := summary.New(&buf)
	r := summary.Build(7, diffs())
	if err := s.Write(r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "Scans: 7") {
		t.Errorf("output missing scan count: %s", buf.String())
	}
}

func TestWrite_ListsPorts(t *testing.T) {
	var buf bytes.Buffer
	s := summary.New(&buf)
	s.Write(summary.Build(1, diffs()))
	out := buf.String()
	if !strings.Contains(out, "tcp/80") {
		t.Errorf("expected tcp/80 in output: %s", out)
	}
	if !strings.Contains(out, "tcp/8080") {
		t.Errorf("expected tcp/8080 in output: %s", out)
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	s := summary.New(nil)
	if s == nil {
		t.Fatal("expected non-nil summarizer")
	}
}
