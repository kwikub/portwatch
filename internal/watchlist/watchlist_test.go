package watchlist_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/watchlist"
)

func snap(ports []scanner.Port) *scanner.Snapshot {
	return &scanner.Snapshot{Ports: ports, At: time.Now()}
}

func TestAdd_ValidEntry(t *testing.T) {
	w := watchlist.New()
	if err := w.Add(80, "tcp"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", w.Len())
	}
}

func TestAdd_InvalidProtocol(t *testing.T) {
	w := watchlist.New()
	if err := w.Add(80, "icmp"); err == nil {
		t.Fatal("expected error for unsupported protocol")
	}
}

func TestAdd_InvalidPort(t *testing.T) {
	w := watchlist.New()
	if err := w.Add(0, "tcp"); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := w.Add(65536, "tcp"); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestMissingFrom_AllPresent(t *testing.T) {
	w := watchlist.New()
	_ = w.Add(80, "tcp")
	s := snap([]scanner.Port{{Port: 80, Protocol: "tcp"}})
	if got := w.MissingFrom(s); len(got) != 0 {
		t.Fatalf("expected no missing entries, got %v", got)
	}
}

func TestMissingFrom_SomeMissing(t *testing.T) {
	w := watchlist.New()
	_ = w.Add(80, "tcp")
	_ = w.Add(443, "tcp")
	s := snap([]scanner.Port{{Port: 80, Protocol: "tcp"}})
	missing := w.MissingFrom(s)
	if len(missing) != 1 {
		t.Fatalf("expected 1 missing entry, got %d", len(missing))
	}
	if missing[0].Port != 443 {
		t.Errorf("expected port 443 missing, got %d", missing[0].Port)
	}
}

func TestRemove_DeregistersEntry(t *testing.T) {
	w := watchlist.New()
	_ = w.Add(22, "tcp")
	w.Remove(22, "tcp")
	if w.Len() != 0 {
		t.Fatal("expected empty watchlist after remove")
	}
}

func TestMissingFrom_EmptySnapshot(t *testing.T) {
	w := watchlist.New()
	_ = w.Add(8080, "tcp")
	s := snap([]scanner.Port{})
	if got := w.MissingFrom(s); len(got) != 1 {
		t.Fatalf("expected 1 missing, got %d", len(got))
	}
}
