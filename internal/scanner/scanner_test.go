package scanner

import (
	"net"
	"testing"
	"time"
)

// startTestServer opens a TCP listener on a random port and returns the port number and a stop func.
func startTestServer(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return port, func() { ln.Close() }
}

func TestScan_DetectsOpenPort(t *testing.T) {
	port, stop := startTestServer(t)
	defer stop()

	s := New("127.0.0.1", 500*time.Millisecond)
	states, err := s.Scan(port, port, "tcp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(states) != 1 {
		t.Fatalf("expected 1 result, got %d", len(states))
	}
	if !states[0].Open {
		t.Errorf("expected port %d to be open", port)
	}
}

func TestScan_DetectsClosedPort(t *testing.T) {
	s := New("127.0.0.1", 200*time.Millisecond)
	// Port 1 is almost certainly closed in test environments.
	states, err := s.Scan(1, 1, "tcp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if states[0].Open {
		t.Errorf("expected port 1 to be closed")
	}
}

func TestScan_InvalidRange(t *testing.T) {
	s := New("127.0.0.1", 200*time.Millisecond)
	_, err := s.Scan(100, 10, "tcp")
	if err == nil {
		t.Error("expected error for invalid range, got nil")
	}
}
