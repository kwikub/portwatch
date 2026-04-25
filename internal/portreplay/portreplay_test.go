package portreplay_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portreplay"
	"github.com/user/portwatch/internal/scanner"
)

func entries() []portreplay.Entry {
	now := time.Now()
	return []portreplay.Entry{
		{
			At:    now,
			Diffs: []scanner.Diff{{Port: 80, Proto: "tcp", State: "opened"}},
		},
		{
			At:    now.Add(2 * time.Second),
			Diffs: []scanner.Diff{{Port: 443, Proto: "tcp", State: "opened"}},
		},
	}
}

func TestRun_EmptyEntriesIsNoop(t *testing.T) {
	var buf bytes.Buffer
	r := portreplay.New(nil, 0, &buf)
	if err := r.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output, got %q", buf.String())
	}
}

func TestRun_WritesEachEntry(t *testing.T) {
	var buf bytes.Buffer
	r := portreplay.New(entries(), 0, &buf)
	if err := r.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Count(out, "[replay]") != 2 {
		t.Errorf("expected 2 replay lines, got:\n%s", out)
	}
}

func TestRun_HandlerReceivesAllEntries(t *testing.T) {
	var buf bytes.Buffer
	r := portreplay.New(entries(), 0, &buf)
	var received []portreplay.Entry
	r.OnEntry(func(e portreplay.Entry) {
		received = append(received, e)
	})
	_ = r.Run()
	if len(received) != 2 {
		t.Errorf("expected 2 handler calls, got %d", len(received))
	}
}

func TestLen_ReturnsEntryCount(t *testing.T) {
	var buf bytes.Buffer
	r := portreplay.New(entries(), 0, &buf)
	if r.Len() != 2 {
		t.Errorf("expected Len 2, got %d", r.Len())
	}
}

func TestNew_NegativeSpeedClampedToZero(t *testing.T) {
	var buf bytes.Buffer
	r := portreplay.New(entries(), -5, &buf)
	// Should complete without sleeping (speed clamped to 0).
	start := time.Now()
	_ = r.Run()
	if time.Since(start) > time.Second {
		t.Error("negative speed caused unexpected delay")
	}
}
