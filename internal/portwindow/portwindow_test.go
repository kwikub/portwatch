package portwindow

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func at(hour, min int) time.Time {
	return time.Date(2024, 1, 15, hour, min, 0, 0, time.UTC)
}

func TestNew_InvalidStart(t *testing.T) {
	_, err := New("25:00", "18:00")
	if err == nil {
		t.Fatal("expected error for invalid start")
	}
}

func TestNew_InvalidEnd(t *testing.T) {
	_, err := New("09:00", "99:99")
	if err == nil {
		t.Fatal("expected error for invalid end")
	}
}

func TestInWindow_InsideRange(t *testing.T) {
	c, _ := New("09:00", "17:00")
	c.now = func() time.Time { return at(12, 0) }
	if !c.InWindow() {
		t.Error("expected 12:00 to be inside 09:00-17:00")
	}
}

func TestInWindow_OutsideRange(t *testing.T) {
	c, _ := New("09:00", "17:00")
	c.now = func() time.Time { return at(18, 30) }
	if c.InWindow() {
		t.Error("expected 18:30 to be outside 09:00-17:00")
	}
}

func TestInWindow_OvernightInsideRange(t *testing.T) {
	c, _ := New("22:00", "06:00")
	c.now = func() time.Time { return at(23, 30) }
	if !c.InWindow() {
		t.Error("expected 23:30 to be inside overnight window")
	}
}

func TestInWindow_OvernightOutsideRange(t *testing.T) {
	c, _ := New("22:00", "06:00")
	c.now = func() time.Time { return at(10, 0) }
	if c.InWindow() {
		t.Error("expected 10:00 to be outside overnight window")
	}
}

func TestFilter_ReturnsOnlyInsideDiffs(t *testing.T) {
	c, _ := New("09:00", "17:00")
	c.now = func() time.Time { return at(12, 0) }

	diffs := []scanner.Diff{
		{Port: 80, Proto: "tcp", State: scanner.Opened, Time: at(10, 0)},
		{Port: 443, Proto: "tcp", State: scanner.Opened, Time: at(20, 0)},
		{Port: 22, Proto: "tcp", State: scanner.Closed, Time: at(14, 59)},
	}

	got := c.Filter(diffs)
	if len(got) != 2 {
		t.Fatalf("expected 2 diffs, got %d", len(got))
	}
	if got[0].Port != 80 || got[1].Port != 22 {
		t.Errorf("unexpected ports: %v", got)
	}
}

func TestFilter_ZeroTimeUsesNow(t *testing.T) {
	c, _ := New("09:00", "17:00")
	c.now = func() time.Time { return at(12, 0) } // inside window

	diffs := []scanner.Diff{
		{Port: 8080, Proto: "tcp", State: scanner.Opened}, // zero Time
	}
	got := c.Filter(diffs)
	if len(got) != 1 {
		t.Fatalf("expected 1 diff when zero time uses now (inside window), got %d", len(got))
	}
}

func TestFilter_EmptyDiffsReturnsNil(t *testing.T) {
	c, _ := New("09:00", "17:00")
	got := c.Filter(nil)
	if got != nil {
		t.Errorf("expected nil for empty input, got %v", got)
	}
}
