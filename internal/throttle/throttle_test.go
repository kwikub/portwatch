package throttle_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/throttle"
)

func TestAllow_FirstCallAlwaysAllowed(t *testing.T) {
	th := throttle.New(time.Second)
	if !th.Allow("tcp", 8080) {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallWithinCooldownBlocked(t *testing.T) {
	th := throttle.New(time.Second)
	th.Allow("tcp", 8080)
	if th.Allow("tcp", 8080) {
		t.Fatal("expected second call within cooldown to be blocked")
	}
}

func TestAllow_CallAfterCooldownAllowed(t *testing.T) {
	th := throttle.New(10 * time.Millisecond)
	th.Allow("tcp", 9000)
	time.Sleep(20 * time.Millisecond)
	if !th.Allow("tcp", 9000) {
		t.Fatal("expected call after cooldown to be allowed")
	}
}

func TestAllow_DifferentPortsAreIndependent(t *testing.T) {
	th := throttle.New(time.Second)
	th.Allow("tcp", 80)
	if !th.Allow("tcp", 443) {
		t.Fatal("expected different port to be allowed independently")
	}
}

func TestReset_AllowsImmediateRetry(t *testing.T) {
	th := throttle.New(time.Second)
	th.Allow("udp", 53)
	th.Reset("udp", 53)
	if !th.Allow("udp", 53) {
		t.Fatal("expected allow after reset")
	}
}

func TestNew_PanicsOnZeroCooldown(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on zero cooldown")
		}
	}()
	throttle.New(0)
}

func TestLen_TracksKeys(t *testing.T) {
	th := throttle.New(time.Second)
	th.Allow("tcp", 80)
	th.Allow("tcp", 443)
	if th.Len() != 2 {
		t.Fatalf("expected 2 tracked keys, got %d", th.Len())
	}
}
