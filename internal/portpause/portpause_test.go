package portpause

import (
	"testing"
	"time"
)

func TestIsPaused_FalseWhenEmpty(t *testing.T) {
	p := New()
	if p.IsPaused(80, "tcp") {
		t.Fatal("expected false for unregistered port")
	}
}

func TestPause_MakesPortPaused(t *testing.T) {
	p := New()
	p.Pause(443, "tcp", 5*time.Minute)
	if !p.IsPaused(443, "tcp") {
		t.Fatal("expected port to be paused")
	}
}

func TestIsPaused_FalseAfterExpiry(t *testing.T) {
	now := time.Now()
	p := New()
	p.now = func() time.Time { return now }
	p.Pause(8080, "tcp", 1*time.Second)

	// advance clock past expiry
	p.now = func() time.Time { return now.Add(2 * time.Second) }
	if p.IsPaused(8080, "tcp") {
		t.Fatal("expected pause to have expired")
	}
}

func TestLift_RemovesPauseImmediately(t *testing.T) {
	p := New()
	p.Pause(22, "tcp", 10*time.Minute)
	p.Lift(22, "tcp")
	if p.IsPaused(22, "tcp") {
		t.Fatal("expected pause to be lifted")
	}
}

func TestProtocolDistinct(t *testing.T) {
	p := New()
	p.Pause(53, "tcp", 5*time.Minute)
	if p.IsPaused(53, "udp") {
		t.Fatal("udp should not be paused when only tcp was paused")
	}
}

func TestPause_ExtendsExistingDeadline(t *testing.T) {
	now := time.Now()
	p := New()
	p.now = func() time.Time { return now }
	p.Pause(3306, "tcp", 1*time.Minute)
	p.Pause(3306, "tcp", 10*time.Minute)

	// advance to 5 minutes — first pause would have expired, second should not
	p.now = func() time.Time { return now.Add(5 * time.Minute) }
	if !p.IsPaused(3306, "tcp") {
		t.Fatal("expected extended pause to still be active")
	}
}

func TestActive_ReturnsOnlyCurrentlyPaused(t *testing.T) {
	now := time.Now()
	p := New()
	p.now = func() time.Time { return now }
	p.Pause(80, "tcp", 5*time.Minute)
	p.Pause(81, "tcp", 1*time.Second)

	// advance past the short pause only
	p.now = func() time.Time { return now.Add(2 * time.Second) }
	actives := p.Active()
	if len(actives) != 1 {
		t.Fatalf("expected 1 active pause, got %d", len(actives))
	}
	if actives[0] != "tcp:80" {
		t.Fatalf("unexpected active key: %s", actives[0])
	}
}
