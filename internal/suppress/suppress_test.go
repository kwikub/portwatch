package suppress

import (
	"testing"
	"time"
)

func TestAllow_FirstCallAlwaysAllowed(t *testing.T) {
	s := New(5 * time.Second)
	if !s.Allow("tcp:8080:opened") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallWithinWindowBlocked(t *testing.T) {
	s := New(5 * time.Second)
	s.Allow("tcp:8080:opened")
	if s.Allow("tcp:8080:opened") {
		t.Fatal("expected second call within window to be blocked")
	}
}

func TestAllow_CallAfterWindowAllowed(t *testing.T) {
	now := time.Now()
	s := New(5 * time.Second)
	s.now = func() time.Time { return now }
	s.Allow("tcp:9090:closed")

	s.now = func() time.Time { return now.Add(6 * time.Second) }
	if !s.Allow("tcp:9090:closed") {
		t.Fatal("expected call after window to be allowed")
	}
}

func TestAllow_DifferentKeysAreIndependent(t *testing.T) {
	s := New(5 * time.Second)
	s.Allow("tcp:80:opened")
	if !s.Allow("tcp:443:opened") {
		t.Fatal("expected different key to be allowed independently")
	}
}

func TestReset_AllowsImmediateRetry(t *testing.T) {
	s := New(5 * time.Second)
	s.Allow("tcp:22:opened")
	s.Reset("tcp:22:opened")
	if !s.Allow("tcp:22:opened") {
		t.Fatal("expected allow after reset")
	}
}

func TestFlush_RemovesExpiredEntries(t *testing.T) {
	now := time.Now()
	s := New(5 * time.Second)
	s.now = func() time.Time { return now }
	s.Allow("tcp:3000:opened")

	s.now = func() time.Time { return now.Add(10 * time.Second) }
	s.Flush()

	// after flush, key should be gone and next Allow should succeed
	if !s.Allow("tcp:3000:opened") {
		t.Fatal("expected allow after flush removed expired entry")
	}
}
