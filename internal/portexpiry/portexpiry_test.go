package portexpiry

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func closed(port int, proto string) scanner.Diff {
	return scanner.Diff{Port: port, Protocol: proto, State: scanner.StateClosed}
}

func opened(port int, proto string) scanner.Diff {
	return scanner.Diff{Port: port, Protocol: proto, State: scanner.StateOpen}
}

func TestExpired_EmptyInitially(t *testing.T) {
	tr := New(time.Minute)
	if len(tr.Expired()) != 0 {
		t.Fatal("expected no expired entries initially")
	}
}

func TestExpired_BelowTTLNotReturned(t *testing.T) {
	tr := New(time.Hour)
	tr.Record([]scanner.Diff{closed(8080, "tcp")})
	if len(tr.Expired()) != 0 {
		t.Fatal("expected no expired entries before TTL")
	}
}

func TestExpired_AfterTTLReturned(t *testing.T) {
	tr := New(time.Minute)
	now := time.Now()
	tr.now = func() time.Time { return now }
	tr.Record([]scanner.Diff{closed(8080, "tcp")})
	tr.now = func() time.Time { return now.Add(2 * time.Minute) }
	exp := tr.Expired()
	if len(exp) != 1 {
		t.Fatalf("expected 1 expired entry, got %d", len(exp))
	}
	if exp[0].Port != 8080 || exp[0].Protocol != "tcp" {
		t.Errorf("unexpected entry: %+v", exp[0])
	}
}

func TestRecord_OpenedRemovesEntry(t *testing.T) {
	tr := New(time.Minute)
	now := time.Now()
	tr.now = func() time.Time { return now }
	tr.Record([]scanner.Diff{closed(443, "tcp")})
	tr.Record([]scanner.Diff{opened(443, "tcp")})
	tr.now = func() time.Time { return now.Add(2 * time.Minute) }
	if len(tr.Expired()) != 0 {
		t.Fatal("expected entry removed after port reopened")
	}
}

func TestEvict_RemovesEntry(t *testing.T) {
	tr := New(time.Minute)
	now := time.Now()
	tr.now = func() time.Time { return now }
	tr.Record([]scanner.Diff{closed(22, "tcp")})
	tr.now = func() time.Time { return now.Add(2 * time.Minute) }
	tr.Evict(22, "tcp")
	if len(tr.Expired()) != 0 {
		t.Fatal("expected entry evicted")
	}
}

func TestRecord_ProtocolDistinct(t *testing.T) {
	tr := New(time.Minute)
	now := time.Now()
	tr.now = func() time.Time { return now }
	tr.Record([]scanner.Diff{closed(53, "tcp"), closed(53, "udp")})
	tr.now = func() time.Time { return now.Add(2 * time.Minute) }
	if len(tr.Expired()) != 2 {
		t.Fatal("expected tcp and udp tracked separately")
	}
}
