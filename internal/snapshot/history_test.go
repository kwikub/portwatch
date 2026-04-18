package snapshot_test

import (
	"testing"
	"time"

	"github.com/example/portwatch/internal/scanner"
	"github.com/example/portwatch/internal/snapshot"
)

func makeSnap(ports []scanner.Port) *scanner.Snapshot {
	return scanner.NewSnapshot(ports)
}

func TestAdd_StoresEntry(t *testing.T) {
	h := snapshot.New(5)
	h.Add(makeSnap(nil))
	if h.Len() != 1 {
		t.Fatalf("expected 1, got %d", h.Len())
	}
}

func TestAdd_EvictsOldestWhenFull(t *testing.T) {
	h := snapshot.New(3)
	for i := 0; i < 5; i++ {
		h.Add(makeSnap(nil))
	}
	if h.Len() != 3 {
		t.Fatalf("expected 3, got %d", h.Len())
	}
}

func TestLatest_EmptyReturnsFalse(t *testing.T) {
	h := snapshot.New(5)
	_, ok := h.Latest()
	if ok {
		t.Fatal("expected false for empty history")
	}
}

func TestLatest_ReturnsNewest(t *testing.T) {
	h := snapshot.New(5)
	h.Add(makeSnap(nil))
	time.Sleep(time.Millisecond)
	want := makeSnap([]scanner.Port{{Port: 9090, Protocol: "tcp"}})
	h.Add(want)
	e, ok := h.Latest()
	if !ok {
		t.Fatal("expected entry")
	}
	if len(e.Snapshot.Ports) != 1 {
		t.Fatalf("expected 1 port, got %d", len(e.Snapshot.Ports))
	}
}

func TestAll_OrderOldestFirst(t *testing.T) {
	h := snapshot.New(5)
	for i := 0; i < 3; i++ {
		h.Add(makeSnap(nil))
		time.Sleep(time.Millisecond)
	}
	entries := h.All()
	for i := 1; i < len(entries); i++ {
		if entries[i].CapturedAt.Before(entries[i-1].CapturedAt) {
			t.Fatal("entries not in ascending order")
		}
	}
}

func TestNew_MinSizeIsOne(t *testing.T) {
	h := snapshot.New(0)
	h.Add(makeSnap(nil))
	h.Add(makeSnap(nil))
	if h.Len() != 1 {
		t.Fatalf("expected 1, got %d", h.Len())
	}
}
