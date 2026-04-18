package summary_test

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/summary"
)

func TestAccumulator_RecordIncrementsScans(t *testing.T) {
	a := summary.NewAccumulator()
	a.Record(nil)
	a.Record(nil)
	if a.Scans() != 2 {
		t.Fatalf("expected 2 scans, got %d", a.Scans())
	}
}

func TestAccumulator_FlushReturnsDiffs(t *testing.T) {
	a := summary.NewAccumulator()
	a.Record([]scanner.Diff{{Port: 22, Proto: "tcp", State: "opened"}})
	a.Record([]scanner.Diff{{Port: 80, Proto: "tcp", State: "closed"}})
	r := a.Flush()
	if len(r.Opened) != 1 {
		t.Fatalf("expected 1 opened, got %d", len(r.Opened))
	}
	if len(r.Closed) != 1 {
		t.Fatalf("expected 1 closed, got %d", len(r.Closed))
	}
	if r.TotalScans != 2 {
		t.Fatalf("expected 2 total scans, got %d", r.TotalScans)
	}
}

func TestAccumulator_FlushResetsState(t *testing.T) {
	a := summary.NewAccumulator()
	a.Record([]scanner.Diff{{Port: 9090, Proto: "tcp", State: "opened"}})
	a.Flush()
	if a.Scans() != 0 {
		t.Fatalf("expected 0 scans after flush, got %d", a.Scans())
	}
	r := a.Flush()
	if len(r.Opened)+len(r.Closed) != 0 {
		t.Fatal("expected empty report after second flush")
	}
}

func TestAccumulator_ConcurrentRecord(t *testing.T) {
	a := summary.NewAccumulator()
	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			a.Record([]scanner.Diff{{Port: 1234, Proto: "tcp", State: "opened"}})
			done <- struct{}{}
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}
	if a.Scans() != 10 {
		t.Fatalf("expected 10 scans, got %d", a.Scans())
	}
}
