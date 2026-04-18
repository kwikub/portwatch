package notify_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"portwatch/internal/notify"
)

// errorChannel always returns an error from Send.
type errorChannel struct{}

func (e *errorChannel) Send(_ notify.Message) error {
	return errors.New("send failed")
}

func TestDispatch_WritesToLogChannel(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(notify.NewLogChannel(&buf))
	msg := notify.Message{Level: notify.LevelAlert, Title: "port change", Body: "22/tcp opened"}
	errs := n.Dispatch(msg)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	out := buf.String()
	if !strings.Contains(out, "ALERT") {
		t.Errorf("expected ALERT in output, got: %s", out)
	}
	if !strings.Contains(out, "22/tcp opened") {
		t.Errorf("expected body in output, got: %s", out)
	}
}

func TestDispatch_SetsTimestampWhenZero(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(notify.NewLogChannel(&buf))
	before := time.Now()
	n.Dispatch(notify.Message{Level: notify.LevelInfo, Title: "t", Body: "b"})
	after := time.Now()
	_ = before
	_ = after
	// If timestamp were zero the RFC3339 string would be "0001-…"; check year.
	if strings.Contains(buf.String(), "0001") {
		t.Error("timestamp was not set by Dispatch")
	}
}

func TestDispatch_CollectsErrors(t *testing.T) {
	n := notify.New(&errorChannel{}, &errorChannel{})
	errs := n.Dispatch(notify.Message{Level: notify.LevelWarn, Title: "x", Body: "y"})
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(errs))
	}
}

func TestNew_DefaultsToLogChannel(t *testing.T) {
	// Should not panic and should return a usable notifier.
	n := notify.New()
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
