package portscorer_test

import (
	"testing"

	"github.com/user/portwatch/internal/portscorer"
	"github.com/user/portwatch/internal/scanner"
)

func d(port int, proto string) scanner.Diff {
	return scanner.Diff{Port: port, Proto: proto, State: scanner.StateOpened}
}

func TestScore_DefaultNoModifiers(t *testing.T) {
	s := portscorer.New()
	got := s.Score(d(80, "tcp"))
	if got != 1.0 {
		t.Fatalf("expected 1.0, got %f", got)
	}
}

func TestScore_CustomDefault(t *testing.T) {
	s := portscorer.New(portscorer.WithDefault(2.0))
	got := s.Score(d(80, "tcp"))
	if got != 2.0 {
		t.Fatalf("expected 2.0, got %f", got)
	}
}

func TestScore_WeightMultiplies(t *testing.T) {
	s := portscorer.New(portscorer.WithWeight(443, "tcp", 3.0))
	got := s.Score(d(443, "tcp"))
	if got != 3.0 {
		t.Fatalf("expected 3.0, got %f", got)
	}
}

func TestScore_CostAddsToBase(t *testing.T) {
	s := portscorer.New(portscorer.WithCost(22, "tcp", 4.0))
	got := s.Score(d(22, "tcp"))
	if got != 5.0 {
		t.Fatalf("expected 5.0 (1+4), got %f", got)
	}
}

func TestScore_PriorityScalesResult(t *testing.T) {
	// score = 1.0 * 1.0 * (1 + 10/10) = 2.0
	s := portscorer.New(portscorer.WithPriority(8080, "tcp", 10))
	got := s.Score(d(8080, "tcp"))
	if got != 2.0 {
		t.Fatalf("expected 2.0, got %f", got)
	}
}

func TestScore_AllModifiersCombined(t *testing.T) {
	// base=1+2=3, weight=2, prio=5 => 3*2*(1+0.5)=9
	s := portscorer.New(
		portscorer.WithCost(9000, "udp", 2.0),
		portscorer.WithWeight(9000, "udp", 2.0),
		portscorer.WithPriority(9000, "udp", 5),
	)
	got := s.Score(d(9000, "udp"))
	if got != 9.0 {
		t.Fatalf("expected 9.0, got %f", got)
	}
}

func TestScore_ProtocolDistinct(t *testing.T) {
	s := portscorer.New(portscorer.WithWeight(53, "tcp", 5.0))
	tcp := s.Score(d(53, "tcp"))
	udp := s.Score(d(53, "udp"))
	if tcp == udp {
		t.Fatal("expected tcp and udp scores to differ")
	}
}

func TestScore_UnknownPortUsesDefaults(t *testing.T) {
	s := portscorer.New(
		portscorer.WithWeight(80, "tcp", 10.0),
		portscorer.WithDefault(3.0),
	)
	got := s.Score(d(9999, "tcp"))
	if got != 3.0 {
		t.Fatalf("expected 3.0 for unknown port, got %f", got)
	}
}
