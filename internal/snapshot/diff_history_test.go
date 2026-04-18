package snapshot_test

import (
	"testing"

	"github.com/example/portwatch/internal/scanner"
	"github.com/example/portwatch/internal/snapshot"
)

func TestDiffs_FewerThanTwoSnapsReturnsNil(t *testing.T) {
	h := snapshot.New(5)
	h.Add(makeSnap(nil))
	if d := h.Diffs(); d != nil {
		t.Fatalf("expected nil, got %v", d)
	}
}

func TestDiffs_NoChangeProducesNoDiffEntry(t *testing.T) {
	h := snapshot.New(5)
	ports := []scanner.Port{{Port: 80, Protocol: "tcp"}}
	h.Add(makeSnap(ports))
	h.Add(makeSnap(ports))
	if len(h.Diffs()) != 0 {
		t.Fatal("expected no diff entries when ports unchanged")
	}
}

func TestDiffs_DetectsNewPort(t *testing.T) {
	h := snapshot.New(5)
	h.Add(makeSnap(nil))
	h.Add(makeSnap([]scanner.Port{{Port: 443, Protocol: "tcp"}}))
	diffs := h.Diffs()
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff entry, got %d", len(diffs))
	}
	if len(diffs[0].Diff) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs[0].Diff))
	}
}

func TestDiffs_MultipleConsecutiveChanges(t *testing.T) {
	h := snapshot.New(10)
	h.Add(makeSnap(nil))
	h.Add(makeSnap([]scanner.Port{{Port: 80, Protocol: "tcp"}}))
	h.Add(makeSnap([]scanner.Port{{Port: 80, Protocol: "tcp"}, {Port: 443, Protocol: "tcp"}}))
	diffs := h.Diffs()
	if len(diffs) != 2 {
		t.Fatalf("expected 2 diff entries, got %d", len(diffs))
	}
}
