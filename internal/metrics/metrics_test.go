package metrics_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

func TestRecordScan_IncrementsTotals(t *testing.T) {
	c := metrics.New(nil)
	c.RecordScan(3, 1)
	c.RecordScan(0, 2)
	s := c.Summary()
	if s.ScansTotal != 2 {
		t.Errorf("expected 2 scans, got %d", s.ScansTotal)
	}
	if s.OpenEvents != 3 {
		t.Errorf("expected 3 open events, got %d", s.OpenEvents)
	}
	if s.CloseEvents != 3 {
		t.Errorf("expected 3 close events, got %d", s.CloseEvents)
	}
}

func TestSummary_UptimeIsPositive(t *testing.T) {
	c := metrics.New(nil)
	time.Sleep(2 * time.Millisecond)
	if c.Summary().Uptime <= 0 {
		t.Error("expected positive uptime")
	}
}

func TestNew_ZeroValues(t *testing.T) {
	c := metrics.New(nil)
	s := c.Summary()
	if s.ScansTotal != 0 || s.OpenEvents != 0 || s.CloseEvents != 0 {
		t.Error("expected zero initial values")
	}
}

func TestPrint_ContainsKeyFields(t *testing.T) {
	var buf bytes.Buffer
	c := metrics.New(&buf)
	c.RecordScan(2, 1)
	c.Print()
	out := buf.String()
	for _, want := range []string{"scans=", "opened=", "closed=", "uptime="} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output: %s", want, out)
		}
	}
}

func TestRecordScan_Concurrent(t *testing.T) {
	c := metrics.New(nil)
	done := make(chan struct{})
	for i := 0; i < 50; i++ {
		go func() {
			c.RecordScan(1, 1)
			done <- struct{}{}
		}()
	}
	for i := 0; i < 50; i++ {
		<-done
	}
	if c.Summary().ScansTotal != 50 {
		t.Errorf("expected 50 scans under concurrency")
	}
}
