package ratelimit_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/ratelimit"
)

func TestAllow_FirstCallAlwaysAllowed(t *testing.T) {
	l := ratelimit.New(time.Second)
	if !l.AllowAt("k", time.Unix(0, 0)) {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallWithinIntervalBlocked(t *testing.T) {
	l := ratelimit.New(time.Second)
	t0 := time.Unix(100, 0)
	l.AllowAt("k", t0)

	if l.AllowAt("k", t0.Add(500*time.Millisecond)) {
		t.Fatal("expected call within interval to be blocked")
	}
}

func TestAllow_CallAfterIntervalAllowed(t *testing.T) {
	l := ratelimit.New(time.Second)
	t0 := time.Unix(100, 0)
	l.AllowAt("k", t0)

	if !l.AllowAt("k", t0.Add(time.Second)) {
		t.Fatal("expected call after interval to be allowed")
	}
}

func TestAllow_DifferentKeysAreIndependent(t *testing.T) {
	l := ratelimit.New(time.Second)
	t0 := time.Unix(100, 0)
	l.AllowAt("a", t0)

	if !l.AllowAt("b", t0) {
		t.Fatal("expected different key to be allowed independently")
	}
}

func TestReset_AllowsImmediateRetry(t *testing.T) {
	l := ratelimit.New(time.Second)
	t0 := time.Unix(100, 0)
	l.AllowAt("k", t0)
	l.Reset("k")

	if !l.AllowAt("k", t0.Add(10*time.Millisecond)) {
		t.Fatal("expected reset key to be allowed immediately")
	}
}
