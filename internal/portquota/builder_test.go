package portquota

import (
	"testing"
	"time"
)

func TestBuilder_Defaults(t *testing.T) {
	q := NewBuilder().Build()
	if q.window != time.Minute {
		t.Fatalf("expected 1m window, got %v", q.window)
	}
	if q.threshold != 5 {
		t.Fatalf("expected threshold 5, got %d", q.threshold)
	}
}

func TestBuilder_CustomValues(t *testing.T) {
	q := NewBuilder().
		WithWindow(30 * time.Second).
		WithThreshold(10).
		Build()
	if q.window != 30*time.Second {
		t.Fatalf("unexpected window: %v", q.window)
	}
	if q.threshold != 10 {
		t.Fatalf("unexpected threshold: %d", q.threshold)
	}
}

func TestBuilder_FunctionalQuota(t *testing.T) {
	q := NewBuilder().
		WithWindow(5 * time.Second).
		WithThreshold(2).
		Build()
	base := time.Now()
	q.Record(22, "tcp", base)
	q.Record(22, "tcp", base.Add(time.Second))
	if !q.Exceeded(22, "tcp", base.Add(2*time.Second)) {
		t.Fatal("expected quota exceeded")
	}
}
