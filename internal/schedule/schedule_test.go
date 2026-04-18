package schedule_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/schedule"
)

func TestNew_FiresAtInterval(t *testing.T) {
	tick := schedule.New(20 * time.Millisecond)
	defer tick.Stop()

	select {
	case <-tick.C:
		// received a tick as expected
	case <-time.After(200 * time.Millisecond):
		t.Fatal("expected tick within 200ms")
	}
}

func TestStop_PreventsAdditionalTicks(t *testing.T) {
	tick := schedule.New(500 * time.Millisecond)
	tick.Stop()

	select {
	case <-tick.C:
		t.Fatal("received tick after Stop")
	case <-time.After(60 * time.Millisecond):
		// correctly silent
	}
}

func TestNew_PanicsOnZeroInterval(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero interval")
		}
	}()
	schedule.New(0)
}

func TestNew_PanicsOnNegativeInterval(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for negative interval")
		}
	}()
	schedule.New(-1 * time.Second)
}

func TestDone_ReceivesAfterStop(t *testing.T) {
	tick := schedule.New(1 * time.Second)
	go tick.Stop()

	select {
	case <-tick.Done():
		// stop signal received
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Done channel not signalled after Stop")
	}
}
