package portwatch_test

import (
	"context"
	"net"
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/enrich"
	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/portwatch"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

func tempState(t *testing.T) *state.State {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "state-*.json")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	os.Remove(f.Name())
	s, err := state.New(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func openPort(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return port, func() { ln.Close() }
}

func TestRun_PersistsStateAfterScan(t *testing.T) {
	port, close := openPort(t)
	defer close()

	sc := scanner.New(scanner.Options{
		Host:      "127.0.0.1",
		PortStart: port,
		PortEnd:   port,
		Protocol:  "tcp",
	})
	st := tempState(t)
	en := enrich.New(enrich.Options{})
	pl := pipeline.New(pipeline.Options{})

	c := portwatch.New(portwatch.Config{
		Scanner:  sc,
		State:    st,
		Enricher: en,
		Pipeline: pl,
	})

	ctx := context.Background()
	if err := c.Run(ctx); err != nil {
		t.Fatalf("Run: %v", err)
	}

	snap, err := st.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(snap.Ports) == 0 {
		t.Fatal("expected at least one port persisted")
	}
}

func TestRunForever_StopsOnCancelledContext(t *testing.T) {
	st := tempState(t)
	sc := scanner.New(scanner.Options{
		Host:      "127.0.0.1",
		PortStart: 1,
		PortEnd:   1,
		Protocol:  "tcp",
	})
	en := enrich.New(enrich.Options{})
	pl := pipeline.New(pipeline.Options{})

	c := portwatch.New(portwatch.Config{
		Scanner:  sc,
		State:    st,
		Enricher: en,
		Pipeline: pl,
	})

	ctx, cancel := context.WithCancel(context.Background())
	tick := make(chan time.Time)

	done := make(chan error, 1)
	go func() { done <- c.RunForever(ctx, tick) }()

	cancel()
	select {
	case err := <-done:
		if err != context.Canceled {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("RunForever did not stop after context cancel")
	}
}

func TestRunForever_StopsWhenTickClosed(t *testing.T) {
	st := tempState(t)
	sc := scanner.New(scanner.Options{
		Host:      "127.0.0.1",
		PortStart: 1,
		PortEnd:   1,
		Protocol:  "tcp",
	})
	en := enrich.New(enrich.Options{})
	pl := pipeline.New(pipeline.Options{})

	c := portwatch.New(portwatch.Config{
		Scanner:  sc,
		State:    st,
		Enricher: en,
		Pipeline: pl,
	})

	tick := make(chan time.Time)
	close(tick)

	ctx := context.Background()
	if err := c.RunForever(ctx, tick); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
