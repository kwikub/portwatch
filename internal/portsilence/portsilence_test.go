package portsilence

import (
	"testing"
	"time"
)

func TestIsSilenced_FalseBeforeSilence(t *testing.T) {
	s := New(time.Second)
	if s.IsSilenced(80, "tcp") {
		t.Fatal("expected port to not be silenced initially")
	}
}

func TestIsSilenced_TrueAfterSilence(t *testing.T) {
	s := New(time.Second)
	s.Silence(80, "tcp")
	if !s.IsSilenced(80, "tcp") {
		t.Fatal("expected port to be silenced")
	}
}

func TestIsSilenced_FalseAfterExpiry(t *testing.T) {
	s := New(20 * time.Millisecond)
	s.Silence(443, "tcp")
	time.Sleep(40 * time.Millisecond)
	if s.IsSilenced(443, "tcp") {
		t.Fatal("expected silence to have expired")
	}
}

func TestLift_RemovesSilenceImmediately(t *testing.T) {
	s := New(time.Hour)
	s.Silence(22, "tcp")
	s.Lift(22, "tcp")
	if s.IsSilenced(22, "tcp") {
		t.Fatal("expected silence to be lifted")
	}
}

func TestProtocolDistinct(t *testing.T) {
	s := New(time.Second)
	s.Silence(53, "tcp")
	if s.IsSilenced(53, "udp") {
		t.Fatal("expected udp to be unaffected by tcp silence")
	}
}

func TestActive_CountsLiveEntries(t *testing.T) {
	s := New(time.Second)
	s.Silence(80, "tcp")
	s.Silence(443, "tcp")
	if got := s.Active(); got != 2 {
		t.Fatalf("expected 2 active, got %d", got)
	}
}

func TestActive_ExpiryReducesCount(t *testing.T) {
	s := New(20 * time.Millisecond)
	s.Silence(8080, "tcp")
	time.Sleep(40 * time.Millisecond)
	if got := s.Active(); got != 0 {
		t.Fatalf("expected 0 active after expiry, got %d", got)
	}
}

func TestNew_PanicsOnZeroDuration(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on zero duration")
		}
	}()
	New(0)
}
