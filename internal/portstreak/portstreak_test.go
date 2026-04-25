package portstreak

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeSnap(ports ...scanner.Port) *scanner.Snapshot {
	return &scanner.Snapshot{Ports: ports, At: time.Now()}
}

func TestGet_UnknownPortReturnsFalse(t *testing.T) {
	tr := New()
	_, ok := tr.Get(80, "tcp")
	if ok {
		t.Fatal("expected false for unrecorded port")
	}
}

func TestRecord_OpenPortStartsAtOne(t *testing.T) {
	tr := New()
	tr.Record(makeSnap(scanner.Port{Port: 80, Protocol: "tcp"}))
	s, ok := tr.Get(80, "tcp")
	if !ok {
		t.Fatal("expected streak to exist")
	}
	if s.State != "open" || s.Count != 1 {
		t.Fatalf("want open/1, got %s/%d", s.State, s.Count)
	}
}

func TestRecord_ConsecutiveOpenIncrementsCount(t *testing.T) {
	tr := New()
	p := scanner.Port{Port: 443, Protocol: "tcp"}
	tr.Record(makeSnap(p))
	tr.Record(makeSnap(p))
	tr.Record(makeSnap(p))
	s, _ := tr.Get(443, "tcp")
	if s.Count != 3 {
		t.Fatalf("want 3, got %d", s.Count)
	}
}

func TestRecord_PortDisappearsSwitchesToClosed(t *testing.T) {
	tr := New()
	p := scanner.Port{Port: 22, Protocol: "tcp"}
	tr.Record(makeSnap(p))
	tr.Record(makeSnap()) // port gone
	s, _ := tr.Get(22, "tcp")
	if s.State != "closed" || s.Count != 1 {
		t.Fatalf("want closed/1, got %s/%d", s.State, s.Count)
	}
}

func TestRecord_ClosedStreakIncrements(t *testing.T) {
	tr := New()
	p := scanner.Port{Port: 22, Protocol: "tcp"}
	tr.Record(makeSnap(p))
	tr.Record(makeSnap())
	tr.Record(makeSnap())
	s, _ := tr.Get(22, "tcp")
	if s.Count != 2 {
		t.Fatalf("want 2, got %d", s.Count)
	}
}

func TestRecord_ProtocolDistinct(t *testing.T) {
	tr := New()
	tr.Record(makeSnap(scanner.Port{Port: 53, Protocol: "tcp"}))
	tr.Record(makeSnap(scanner.Port{Port: 53, Protocol: "tcp"}))
	// udp/53 was never seen open, so it starts closed
	tr.Record(makeSnap(scanner.Port{Port: 53, Protocol: "tcp"}))

	tcp, _ := tr.Get(53, "tcp")
	_, udpOk := tr.Get(53, "udp")
	if tcp.Count != 3 {
		t.Fatalf("tcp want 3, got %d", tcp.Count)
	}
	if udpOk {
		t.Fatal("udp/53 should not exist")
	}
}

func TestReset_ClearsAllStreaks(t *testing.T) {
	tr := New()
	tr.Record(makeSnap(scanner.Port{Port: 80, Protocol: "tcp"}))
	tr.Reset()
	_, ok := tr.Get(80, "tcp")
	if ok {
		t.Fatal("expected streak to be cleared after Reset")
	}
}
