package classify_test

import (
	"testing"

	"github.com/user/portwatch/internal/classify"
	"github.com/user/portwatch/internal/scanner"
)

func rules() []classify.Rule {
	return []classify.Rule{
		{MinPort: 0, MaxPort: 1023, Protocol: "tcp", Level: classify.LevelCritical},
		{MinPort: 1024, MaxPort: 49151, Protocol: "tcp", Level: classify.LevelWarning},
	}
}

func TestClassify_CriticalRange(t *testing.T) {
	c := classify.New(rules())
	d := scanner.Diff{Port: 80, Protocol: "tcp", State: scanner.StateOpened}
	if got := c.Classify(d); got != classify.LevelCritical {
		t.Fatalf("expected critical, got %s", got)
	}
}

func TestClassify_WarningRange(t *testing.T) {
	c := classify.New(rules())
	d := scanner.Diff{Port: 8080, Protocol: "tcp", State: scanner.StateOpened}
	if got := c.Classify(d); got != classify.LevelWarning {
		t.Fatalf("expected warning, got %s", got)
	}
}

func TestClassify_NoMatchDefaultsToInfo(t *testing.T) {
	c := classify.New(rules())
	d := scanner.Diff{Port: 55000, Protocol: "tcp"}
	if got := c.Classify(d); got != classify.LevelInfo {
		t.Fatalf("expected info, got %s", got)
	}
}

func TestClassify_ProtocolMismatchSkipsRule(t *testing.T) {
	c := classify.New(rules())
	d := scanner.Diff{Port: 80, Protocol: "udp"}
	if got := c.Classify(d); got != classify.LevelInfo {
		t.Fatalf("expected info for udp:80, got %s", got)
	}
}

func TestClassifyAll_ReturnsMapForEachDiff(t *testing.T) {
	c := classify.New(rules())
	diffs := []scanner.Diff{
		{Port: 22, Protocol: "tcp"},
		{Port: 8080, Protocol: "tcp"},
		{Port: 55000, Protocol: "tcp"},
	}
	m := c.ClassifyAll(diffs)
	if len(m) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(m))
	}
	if m[0] != classify.LevelCritical {
		t.Errorf("index 0: want critical, got %s", m[0])
	}
	if m[1] != classify.LevelWarning {
		t.Errorf("index 1: want warning, got %s", m[1])
	}
	if m[2] != classify.LevelInfo {
		t.Errorf("index 2: want info, got %s", m[2])
	}
}

func TestNew_EmptyRulesAlwaysInfo(t *testing.T) {
	c := classify.New(nil)
	d := scanner.Diff{Port: 443, Protocol: "tcp"}
	if got := c.Classify(d); got != classify.LevelInfo {
		t.Fatalf("expected info, got %s", got)
	}
}
