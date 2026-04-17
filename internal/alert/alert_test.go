package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func diffs(n int) []scanner.Diff {
	out := make([]scanner.Diff, n)
	for i := range out {
		out[i] = scanner.Diff{Port: 8000 + i, Proto: "tcp", State: "opened"}
	}
	return out
}

func TestEvaluate_InfoBelowThreshold(t *testing.T) {
	var buf bytes.Buffer
	a := alert.New(5, &buf)
	events := a.Evaluate(diffs(2))
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	for _, ev := range events {
		if ev.Level != alert.LevelInfo {
			t.Errorf("expected INFO, got %s", ev.Level)
		}
	}
}

func TestEvaluate_AlertAtThreshold(t *testing.T) {
	var buf bytes.Buffer
	a := alert.New(3, &buf)
	events := a.Evaluate(diffs(3))
	for _, ev := range events {
		if ev.Level != alert.LevelAlert {
			t.Errorf("expected ALERT, got %s", ev.Level)
		}
	}
}

func TestEvaluate_EmptyDiffs(t *testing.T) {
	var buf bytes.Buffer
	a := alert.New(1, &buf)
	events := a.Evaluate(nil)
	if len(events) != 0 {
		t.Fatalf("expected no events, got %d", len(events))
	}
	if buf.Len() != 0 {
		t.Error("expected no output for empty diffs")
	}
}

func TestEvaluate_OutputContainsLevel(t *testing.T) {
	var buf bytes.Buffer
	a := alert.New(1, &buf)
	a.Evaluate(diffs(1))
	if !strings.Contains(buf.String(), "[ALERT]") {
		t.Errorf("expected [ALERT] in output, got: %s", buf.String())
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	a := alert.New(5, nil)
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
}
