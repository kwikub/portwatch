package portping_test

import (
	"net"
	"testing"
	"time"

	"portwatch/internal/portping"
)

func startTCP(t *testing.T) (port int, stop func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return ln.Addr().(*net.TCPAddr).Port, func() { ln.Close() }
}

func TestProbe_OpenPort(t *testing.T) {
	port, stop := startTCP(t)
	defer stop()

	p := portping.New("127.0.0.1", time.Second)
	r := p.Probe(port, "tcp")

	if !r.Open {
		t.Fatalf("expected open, got closed")
	}
	if r.Latency <= 0 {
		t.Fatalf("expected positive latency")
	}
}

func TestProbe_ClosedPort(t *testing.T) {
	p := portping.New("127.0.0.1", 200*time.Millisecond)
	r := p.Probe(19999, "tcp")
	if r.Open {
		t.Fatalf("expected closed")
	}
}

func TestProbeAll_ReturnsAllResults(t *testing.T) {
	port, stop := startTCP(t)
	defer stop()

	p := portping.New("127.0.0.1", time.Second)
	results := p.ProbeAll([]int{port, 19998}, "tcp")

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if !results[0].Open {
		t.Errorf("first port should be open")
	}
	if results[1].Open {
		t.Errorf("second port should be closed")
	}
}

func TestBuilder_Defaults(t *testing.T) {
	p := portping.NewBuilder().Build()
	if p == nil {
		t.Fatal("expected non-nil prober")
	}
}

func TestBuilder_CustomValues(t *testing.T) {
	port, stop := startTCP(t)
	defer stop()

	p := portping.NewBuilder().
		WithHost("127.0.0.1").
		WithTimeout(time.Second).
		Build()

	r := p.Probe(port, "tcp")
	if !r.Open {
		t.Errorf("expected open port")
	}
}
