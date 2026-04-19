package trend_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/trend"
)

func diff(proto string, port int, state string) scanner.Diff {
	return scanner.Diff{Proto: proto, Port: port, State: state}
}

func TestNew_PanicsOnZeroWindow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	trend.New(0)
}

func TestCount_ZeroBeforeRecord(t *testing.T) {
	tr := trend.New(time.Minute)
	if n := tr.Count(diff("tcp", 80, "opened")); n != 0 {
		t.Fatalf("expected 0, got %d", n)
	}
}

func TestRecord_IncrementsCount(t *testing.T) {
	tr := trend.New(time.Minute)
	d := diff("tcp", 443, "opened")
	tr.Record([]scanner.Diff{d})
	tr.Record([]scanner.Diff{d})
	if n := tr.Count(d); n != 2 {
		t.Fatalf("expected 2, got %d", n)
	}
}

func TestCount_ProtocolDistinct(t *testing.T) {
	tr := trend.New(time.Minute)
	tcp := diff("tcp", 53, "opened")
	udp := diff("udp", 53, "opened")
	tr.Record([]scanner.Diff{tcp, tcp, udp})
	if n := tr.Count(tcp); n != 2 {
		t.Fatalf("expected 2 for tcp, got %d", n)
	}
	if n := tr.Count(udp); n != 1 {
		t.Fatalf("expected 1 for udp, got %d", n)
	}
}

func TestFlush_RemovesExpiredEvents(t *testing.T) {
	tr := trend.New(50 * time.Millisecond)
	d := diff("tcp", 8080, "closed")
	tr.Record([]scanner.Diff{d})
	time.Sleep(80 * time.Millisecond)
	tr.Flush()
	if n := tr.Count(d); n != 0 {
		t.Fatalf("expected 0 after flush, got %d", n)
	}
}

func TestCount_ExcludesExpiredWithoutFlush(t *testing.T) {
	tr := trend.New(50 * time.Millisecond)
	d := diff("tcp", 22, "opened")
	tr.Record([]scanner.Diff{d})
	time.Sleep(80 * time.Millisecond)
	if n := tr.Count(d); n != 0 {
		t.Fatalf("expected 0 for expired event, got %d", n)
	}
}
