package logger

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func TestLog_OpenedEntry(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf)

	l.Log([]scanner.DiffEntry{
		{
			State: scanner.StateOpened,
			Entry: scanner.Entry{Host: "127.0.0.1", Port: 8080},
		},
	})

	out := buf.String()
	if !strings.Contains(out, "OPENED") {
		t.Errorf("expected OPENED in output, got: %s", out)
	}
	if !strings.Contains(out, "127.0.0.1:8080") {
		t.Errorf("expected host:port in output, got: %s", out)
	}
}

func TestLog_ClosedEntry(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf)

	l.Log([]scanner.DiffEntry{
		{
			State: scanner.StateClosed,
			Entry: scanner.Entry{Host: "localhost", Port: 443},
		},
	})

	out := buf.String()
	if !strings.Contains(out, "CLOSED") {
		t.Errorf("expected CLOSED in output, got: %s", out)
	}
	if !strings.Contains(out, "localhost:443") {
		t.Errorf("expected host:port in output, got: %s", out)
	}
}

func TestLog_EmptyDiffs(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf)
	l.Log(nil)
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty diffs, got: %s", buf.String())
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	l := New(nil)
	if l.out == nil {
		t.Error("expected non-nil writer when nil passed to New")
	}
}
