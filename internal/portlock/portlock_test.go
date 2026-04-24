package portlock

import (
	"testing"
	"time"
)

func TestIsLocked_FalseWhenEmpty(t *testing.T) {
	r := New()
	if r.IsLocked(80, "tcp") {
		t.Fatal("expected port to be unlocked")
	}
}

func TestLock_MakesPortLocked(t *testing.T) {
	r := New()
	r.Lock(80, "tcp", "maintenance", 0)
	if !r.IsLocked(80, "tcp") {
		t.Fatal("expected port 80/tcp to be locked")
	}
}

func TestUnlock_RemovesLock(t *testing.T) {
	r := New()
	r.Lock(443, "tcp", "test", 0)
	r.Unlock(443, "tcp")
	if r.IsLocked(443, "tcp") {
		t.Fatal("expected port to be unlocked after Unlock")
	}
}

func TestIsLocked_ProtocolDistinct(t *testing.T) {
	r := New()
	r.Lock(53, "tcp", "test", 0)
	if r.IsLocked(53, "udp") {
		t.Fatal("locking tcp should not affect udp")
	}
}

func TestIsLocked_ExpiredLockReturnsFalse(t *testing.T) {
	r := New()
	r.Lock(8080, "tcp", "short", 1*time.Millisecond)
	time.Sleep(10 * time.Millisecond)
	if r.IsLocked(8080, "tcp") {
		t.Fatal("expired lock should report as unlocked")
	}
}

func TestGet_ReturnsLockDetails(t *testing.T) {
	r := New()
	r.Lock(22, "tcp", "ssh-maintenance", 0)
	l, ok := r.Get(22, "tcp")
	if !ok {
		t.Fatal("expected lock to be found")
	}
	if l.Reason != "ssh-maintenance" {
		t.Fatalf("expected reason %q, got %q", "ssh-maintenance", l.Reason)
	}
	if l.Port != 22 || l.Protocol != "tcp" {
		t.Fatalf("unexpected lock fields: %+v", l)
	}
}

func TestGet_MissingReturnsFalse(t *testing.T) {
	r := New()
	_, ok := r.Get(9999, "tcp")
	if ok {
		t.Fatal("expected no lock for unregistered port")
	}
}

func TestAll_ReturnsOnlyActiveLocks(t *testing.T) {
	r := New()
	r.Lock(80, "tcp", "a", 0)
	r.Lock(81, "tcp", "b", 1*time.Millisecond)
	time.Sleep(10 * time.Millisecond)
	all := r.All()
	if len(all) != 1 {
		t.Fatalf("expected 1 active lock, got %d", len(all))
	}
	if all[0].Port != 80 {
		t.Fatalf("expected port 80, got %d", all[0].Port)
	}
}

func TestLock_ZeroTTLNeverExpires(t *testing.T) {
	r := New()
	r.Lock(3306, "tcp", "db", 0)
	l, ok := r.Get(3306, "tcp")
	if !ok {
		t.Fatal("expected lock to be present")
	}
	if !l.ExpiresAt.IsZero() {
		t.Fatal("expected zero ExpiresAt for no-TTL lock")
	}
}
