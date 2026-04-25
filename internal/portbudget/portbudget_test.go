package portbudget

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeSnap(ports []scanner.Port) *scanner.Snapshot {
	return &scanner.Snapshot{
		At:    time.Now(),
		Ports: ports,
	}
}

func TestNew_PanicsOnZeroMax(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero max")
		}
	}()
	New(0)
}

func TestNew_PanicsOnNegativeMax(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for negative max")
		}
	}()
	New(-5)
}

func TestExceeded_FalseWhenUnderBudget(t *testing.T) {
	b := New(10)
	ports := []scanner.Port{
		{Port: 80, Proto: "tcp"},
		{Port: 443, Proto: "tcp"},
	}
	b.Record(makeSnap(ports))
	if b.Exceeded() {
		t.Errorf("expected not exceeded with 2/%d ports", b.Max())
	}
}

func TestExceeded_FalseAtExactBudget(t *testing.T) {
	b := New(2)
	ports := []scanner.Port{
		{Port: 80, Proto: "tcp"},
		{Port: 443, Proto: "tcp"},
	}
	b.Record(makeSnap(ports))
	if b.Exceeded() {
		t.Error("expected not exceeded when count equals max")
	}
}

func TestExceeded_TrueWhenOverBudget(t *testing.T) {
	b := New(1)
	ports := []scanner.Port{
		{Port: 80, Proto: "tcp"},
		{Port: 443, Proto: "tcp"},
	}
	b.Record(makeSnap(ports))
	if !b.Exceeded() {
		t.Error("expected exceeded when count > max")
	}
}

func TestCurrent_ReflectsLatestSnapshot(t *testing.T) {
	b := New(100)
	b.Record(makeSnap([]scanner.Port{{Port: 22, Proto: "tcp"}}))
	if got := b.Current(); got != 1 {
		t.Errorf("expected current=1, got %d", got)
	}
	b.Record(makeSnap([]scanner.Port{
		{Port: 22, Proto: "tcp"},
		{Port: 80, Proto: "tcp"},
		{Port: 443, Proto: "tcp"},
	}))
	if got := b.Current(); got != 3 {
		t.Errorf("expected current=3, got %d", got)
	}
}

func TestReset_SetCurrentToZero(t *testing.T) {
	b := New(5)
	b.Record(makeSnap([]scanner.Port{{Port: 80, Proto: "tcp"}}))
	b.Reset()
	if got := b.Current(); got != 0 {
		t.Errorf("expected current=0 after reset, got %d", got)
	}
	if b.Exceeded() {
		t.Error("expected not exceeded after reset")
	}
}
