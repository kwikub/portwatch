package portage

import (
	"testing"
	"time"
)

func TestOpened_TracksFirstSeen(t *testing.T) {
	tr := New()
	tr.Opened(80, "tcp")
	age, ok := tr.Age(80, "tcp")
	if !ok {
		t.Fatal("expected port to be tracked")
	}
	if age < 0 {
		t.Fatalf("unexpected negative age: %v", age)
	}
}

func TestOpened_DoesNotOverwriteExisting(t *testing.T) {
	now := time.Now()
	tr := New()
	tr.now = func() time.Time { return now }
	tr.Opened(443, "tcp")
	tr.now = func() time.Time { return now.Add(5 * time.Second) }
	tr.Opened(443, "tcp") // should not overwrite
	age, ok := tr.Age(443, "tcp")
	if !ok {
		t.Fatal("expected port to be tracked")
	}
	if age < 5*time.Second {
		t.Fatalf("expected age >= 5s, got %v", age)
	}
}

func TestClosed_RemovesPort(t *testing.T) {
	tr := New()
	tr.Opened(22, "tcp")
	tr.Closed(22, "tcp")
	_, ok := tr.Age(22, "tcp")
	if ok {
		t.Fatal("expected port to be removed after close")
	}
}

func TestAge_UnknownPortReturnsFalse(t *testing.T) {
	tr := New()
	_, ok := tr.Age(9999, "udp")
	if ok {
		t.Fatal("expected false for untracked port")
	}
}

func TestProtocolDistinct(t *testing.T) {
	tr := New()
	tr.Opened(53, "tcp")
	_, ok := tr.Age(53, "udp")
	if ok {
		t.Fatal("tcp and udp should be tracked independently")
	}
}

func TestReset_ClearsAll(t *testing.T) {
	tr := New()
	tr.Opened(80, "tcp")
	tr.Opened(443, "tcp")
	tr.Reset()
	_, ok1 := tr.Age(80, "tcp")
	_, ok2 := tr.Age(443, "tcp")
	if ok1 || ok2 {
		t.Fatal("expected all ports cleared after reset")
	}
}
