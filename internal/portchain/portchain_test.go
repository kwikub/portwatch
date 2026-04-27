package portchain_test

import (
	"testing"

	"github.com/user/portwatch/internal/portchain"
	"github.com/user/portwatch/internal/scanner"
)

func diffs() []scanner.Diff {
	return []scanner.Diff{
		{Port: 80, Proto: "tcp", State: scanner.StateOpened},
		{Port: 443, Proto: "tcp", State: scanner.StateOpened},
		{Port: 22, Proto: "tcp", State: scanner.StateClosed},
	}
}

func TestRun_EmptyChainReturnsDiffsUnchanged(t *testing.T) {
	c := portchain.New()
	out := c.Run(diffs())
	if len(out) != 3 {
		t.Fatalf("expected 3 diffs, got %d", len(out))
	}
}

func TestRun_StageCanFilterDiffs(t *testing.T) {
	c := portchain.New()
	_ = c.Add("keep-opened", func(ds []scanner.Diff) []scanner.Diff {
		var out []scanner.Diff
		for _, d := range ds {
			if d.State == scanner.StateOpened {
				out = append(out, d)
			}
		}
		return out
	})
	out := c.Run(diffs())
	if len(out) != 2 {
		t.Fatalf("expected 2 diffs after filter, got %d", len(out))
	}
}

func TestRun_ShortCircuitsOnEmptyResult(t *testing.T) {
	c := portchain.New()
	called := false
	_ = c.Add("drop-all", func(ds []scanner.Diff) []scanner.Diff { return nil })
	_ = c.Add("should-not-run", func(ds []scanner.Diff) []scanner.Diff {
		called = true
		return ds
	})
	c.Run(diffs())
	if called {
		t.Fatal("second stage should not have been called after empty result")
	}
}

func TestAdd_EmptyNameReturnsError(t *testing.T) {
	c := portchain.New()
	if err := c.Add("", func(ds []scanner.Diff) []scanner.Diff { return ds }); err == nil {
		t.Fatal("expected error for empty stage name")
	}
}

func TestAdd_NilHandlerReturnsError(t *testing.T) {
	c := portchain.New()
	if err := c.Add("stage", nil); err == nil {
		t.Fatal("expected error for nil handler")
	}
}

func TestStages_ReturnsNamesInOrder(t *testing.T) {
	c := portchain.New()
	pass := func(ds []scanner.Diff) []scanner.Diff { return ds }
	_ = c.Add("alpha", pass)
	_ = c.Add("beta", pass)
	_ = c.Add("gamma", pass)
	names := c.Stages()
	expected := []string{"alpha", "beta", "gamma"}
	for i, n := range expected {
		if names[i] != n {
			t.Errorf("stage[%d]: want %q, got %q", i, n, names[i])
		}
	}
}

func TestLen_ReflectsAddedStages(t *testing.T) {
	c := portchain.New()
	if c.Len() != 0 {
		t.Fatalf("expected 0, got %d", c.Len())
	}
	_ = c.Add("s1", func(ds []scanner.Diff) []scanner.Diff { return ds })
	if c.Len() != 1 {
		t.Fatalf("expected 1, got %d", c.Len())
	}
}
