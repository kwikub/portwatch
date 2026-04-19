package export

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func diffs() []scanner.Diff {
	at := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	return []scanner.Diff{
		{Port: 80, Protocol: "tcp", Opened: true, At: at},
		{Port: 443, Protocol: "tcp", Opened: false, At: at},
	}
}

func TestWrite_JSONContainsFields(t *testing.T) {
	var buf bytes.Buffer
	e := New(&buf, FormatJSON)
	if err := e.Write(diffs()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var records []Record
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("want 2 records, got %d", len(records))
	}
	if records[0].Port != 80 || records[0].State != "opened" {
		t.Errorf("unexpected first record: %+v", records[0])
	}
	if records[1].Port != 443 || records[1].State != "closed" {
		t.Errorf("unexpected second record: %+v", records[1])
	}
}

func TestWrite_CSVContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	e := New(&buf, FormatCSV)
	if err := e.Write(diffs()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if lines[0] != "timestamp,port,protocol,state" {
		t.Errorf("unexpected header: %s", lines[0])
	}
	if len(lines) != 3 {
		t.Errorf("want 3 lines (header+2), got %d", len(lines))
	}
}

func TestWrite_EmptyDiffsWritesNothing(t *testing.T) {
	var buf bytes.Buffer
	e := New(&buf, FormatJSON)
	if err := e.Write(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got %q", buf.String())
	}
}

func TestWrite_UnknownFormatReturnsError(t *testing.T) {
	var buf bytes.Buffer
	e := New(&buf, Format("xml"))
	if err := e.Write(diffs()); err == nil {
		t.Fatal("expected error for unknown format")
	}
}
