package reporter_test

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/reporter"
	"github.com/user/portwatch/internal/scanner"
)

func diffs() []scanner.Diff {
	return []scanner.Diff{
		{Proto: "tcp", Port: 8080, Closed: false},
		{Proto: "udp", Port: 53, Closed: true},
	}
}

func TestReport_TextContainsPortAndState(t *testing.T) {
	var buf strings.Builder
	r := reporter.New(&buf, reporter.FormatText)
	if err := r.Report(diffs()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "OPENED") {
		t.Errorf("expected OPENED in output, got: %s", out)
	}
	if !strings.Contains(out, "CLOSED") {
		t.Errorf("expected CLOSED in output, got: %s", out)
	}
	if !strings.Contains(out, "tcp/8080") {
		t.Errorf("expected tcp/8080 in output, got: %s", out)
	}
	if !strings.Contains(out, "udp/53") {
		t.Errorf("expected udp/53 in output, got: %s", out)
	}
}

func TestReport_CSVFormat(t *testing.T) {
	var buf strings.Builder
	r := reporter.New(&buf, reporter.FormatCSV)
	if err := r.Report(diffs()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if !strings.Contains(lines[0], ",opened,tcp,8080") {
		t.Errorf("unexpected CSV line: %s", lines[0])
	}
}

func TestReport_EmptyDiffsWritesNothing(t *testing.T) {
	var buf strings.Builder
	r := reporter.New(&buf, reporter.FormatText)
	if err := r.Report(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty diffs, got: %s", buf.String())
	}
}

func TestNew_DefaultsToStdoutAndText(t *testing.T) {
	r := reporter.New(nil, "")
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}
