package watcher_test

import (
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/logger"
	"github.com/user/portwatch/internal/state"
	"github.com/user/portwatch/internal/watcher"
)

func tempState(t *testing.T) *state.State {
	t.Helper()
	dir := t.TempDir()
	st, err := state.New(dir + "/state.json")
	if err != nil {
		t.Fatalf("state.New: %v", err)
	}
	return st
}

func TestWatcher_DetectsOpenPort(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()
	addr := ln.Addr().(*net.TCPAddr)

	cfg := &config.Config{
		Protocol:  "tcp",
		StartPort: addr.Port,
		EndPort:   addr.Port,
		Host:      "127.0.0.1",
		Interval:  100 * time.Millisecond,
		Timeout:   200 * time.Millisecond,
	}

	var buf safeBuffer
	log := logger.New(&buf)
	st := tempState(t)
	w := watcher.New(cfg, log, st)

	done := make(chan error, 1)
	go func() { done <- w.Start() }()
	time.Sleep(250 * time.Millisecond)
	w.Stop()
	if err := <-done; err != nil {
		t.Fatalf("Start: %v", err)
	}
	if buf.String() == "" {
		t.Error("expected log output, got none")
	}
}

// safeBuffer is a thread-safe bytes.Buffer.
import "bytes"
import "sync"

type safeBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (s *safeBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}

func (s *safeBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.String()
}
