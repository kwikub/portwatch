package portquota

import (
	"testing"
	"time"
)

var (
	now   = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	win   = 10 * time.Second
	thr   = 3
)

func TestExceeded_BelowThreshold(t *testing.T) {
	q := New(win, thr)
	q.Record(80, "tcp", now)
	q.Record(80, "tcp", now.Add(1*time.Second))
	if q.Exceeded(80, "tcp", now.Add(2*time.Second)) {
		t.Fatal("expected not exceeded")
	}
}

func TestExceeded_AtThreshold(t *testing.T) {
	q := New(win, thr)
	for i := 0; i < thr; i++ {
		q.Record(80, "tcp", now.Add(time.Duration(i)*time.Second))
	}
	if !q.Exceeded(80, "tcp", now.Add(time.Duration(thr)*time.Second)) {
		t.Fatal("expected exceeded at threshold")
	}
}

func TestExceeded_EventsOutsideWindowEvicted(t *testing.T) {
	q := New(win, thr)
	for i := 0; i < thr; i++ {
		q.Record(80, "tcp", now.Add(time.Duration(i)*time.Second))
	}
	// advance past window so all events expire
	later := now.Add(win + time.Second)
	if q.Exceeded(80, "tcp", later) {
		t.Fatal("expected not exceeded after window expiry")
	}
}

func TestExceeded_ProtocolDistinct(t *testing.T) {
	q := New(win, thr)
	for i := 0; i < thr; i++ {
		q.Record(80, "tcp", now.Add(time.Duration(i)*time.Second))
	}
	if q.Exceeded(80, "udp", now.Add(time.Duration(thr)*time.Second)) {
		t.Fatal("udp should be independent of tcp")
	}
}

func TestReset_ClearsCount(t *testing.T) {
	q := New(win, thr)
	for i := 0; i < thr; i++ {
		q.Record(443, "tcp", now.Add(time.Duration(i)*time.Second))
	}
	q.Reset(443, "tcp")
	if q.Exceeded(443, "tcp", now.Add(time.Duration(thr)*time.Second)) {
		t.Fatal("expected count cleared after reset")
	}
}

func TestNew_PanicsOnZeroWindow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on zero window")
		}
	}()
	New(0, thr)
}

func TestNew_PanicsOnZeroThreshold(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on zero threshold")
		}
	}()
	New(win, 0)
}
