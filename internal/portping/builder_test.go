package portping_test

import (
	"testing"
	"time"

	"portwatch/internal/portping"
)

func TestNewBuilder_DefaultHost(t *testing.T) {
	b := portping.NewBuilder()
	p := b.Build()
	if p == nil {
		t.Fatal("expected prober")
	}
}

func TestNewBuilder_WithHostAndTimeout(t *testing.T) {
	p := portping.NewBuilder().
		WithHost("localhost").
		WithTimeout(500 * time.Millisecond).
		Build()
	if p == nil {
		t.Fatal("expected prober")
	}
}

func TestNewBuilder_ZeroTimeoutUsesDefault(t *testing.T) {
	p := portping.NewBuilder().WithTimeout(0).Build()
	if p == nil {
		t.Fatal("expected prober with default timeout")
	}
	// Probe a closed port — should not hang indefinitely
	r := p.Probe(19997, "tcp")
	if r.Open {
		t.Errorf("expected closed")
	}
}
