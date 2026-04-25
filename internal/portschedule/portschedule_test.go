package portschedule

import (
	"testing"
	"time"
)

var weekdays = []time.Weekday{
	time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday,
}

func at(weekday time.Weekday, hour, minute int) time.Time {
	// Find the next occurrence of the given weekday from a fixed base.
	base := time.Date(2024, 1, 1, hour, minute, 0, 0, time.UTC) // Monday
	offset := int(weekday) - int(base.Weekday())
	if offset < 0 {
		offset += 7
	}
	return base.AddDate(0, 0, offset)
}

func TestActive_NoRulesAlwaysTrue(t *testing.T) {
	s := New()
	if !s.Active(80, "tcp", at(time.Saturday, 3, 0)) {
		t.Fatal("expected true when no rules defined")
	}
}

func TestActive_InsideWindowReturnsTrue(t *testing.T) {
	s := New()
	if err := s.Add(80, "tcp", "09:00", "17:00", weekdays); err != nil {
		t.Fatal(err)
	}
	if !s.Active(80, "tcp", at(time.Wednesday, 12, 0)) {
		t.Fatal("expected active inside window")
	}
}

func TestActive_OutsideWindowReturnsFalse(t *testing.T) {
	s := New()
	if err := s.Add(80, "tcp", "09:00", "17:00", weekdays); err != nil {
		t.Fatal(err)
	}
	if s.Active(80, "tcp", at(time.Wednesday, 18, 0)) {
		t.Fatal("expected inactive outside window")
	}
}

func TestActive_WeekendExcluded(t *testing.T) {
	s := New()
	if err := s.Add(443, "tcp", "08:00", "20:00", weekdays); err != nil {
		t.Fatal(err)
	}
	if s.Active(443, "tcp", at(time.Saturday, 10, 0)) {
		t.Fatal("expected inactive on weekend")
	}
}

func TestActive_ProtocolDistinct(t *testing.T) {
	s := New()
	if err := s.Add(53, "udp", "06:00", "22:00", weekdays); err != nil {
		t.Fatal(err)
	}
	// tcp/53 has no rule → should be active
	if !s.Active(53, "tcp", at(time.Monday, 23, 0)) {
		t.Fatal("expected tcp rule to be independent of udp rule")
	}
	// udp/53 outside window → inactive
	if s.Active(53, "udp", at(time.Monday, 23, 0)) {
		t.Fatal("expected udp/53 inactive outside window")
	}
}

func TestAdd_InvalidPortReturnsError(t *testing.T) {
	s := New()
	if err := s.Add(0, "tcp", "09:00", "17:00", weekdays); err == nil {
		t.Fatal("expected error for port 0")
	}
}

func TestAdd_InvalidProtocolReturnsError(t *testing.T) {
	s := New()
	if err := s.Add(80, "icmp", "09:00", "17:00", weekdays); err == nil {
		t.Fatal("expected error for unknown protocol")
	}
}

func TestAdd_EndBeforeStartReturnsError(t *testing.T) {
	s := New()
	if err := s.Add(80, "tcp", "17:00", "09:00", weekdays); err == nil {
		t.Fatal("expected error when end <= start")
	}
}

func TestAdd_EmptyWeekdaysReturnsError(t *testing.T) {
	s := New()
	if err := s.Add(80, "tcp", "09:00", "17:00", nil); err == nil {
		t.Fatal("expected error for empty weekdays")
	}
}
