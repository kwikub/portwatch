package enrich_test

import (
	"testing"

	"github.com/user/portwatch/internal/enrich"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/tags"
)

func makeRegistry(t *testing.T, proto string, port int, label string) *tags.Registry {
	t.Helper()
	r := tags.New()
	if err := r.Add(proto, port, label); err != nil {
		t.Fatalf("tags.Add: %v", err)
	}
	return r
}

func TestEnrich_AttachesTag(t *testing.T) {
	reg := makeRegistry(t, "tcp", 80, "http")
	e := enrich.New(reg, nil)

	diffs := []scanner.Diff{{Proto: "tcp", Port: 80, State: "opened"}}
	got := e.Enrich(diffs)

	if len(got) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(got))
	}
	if got[0].Tag != "http" {
		t.Errorf("expected tag 'http', got %q", got[0].Tag)
	}
}

func TestEnrich_UnknownTagIsEmpty(t *testing.T) {
	reg := tags.New()
	e := enrich.New(reg, nil)

	diffs := []scanner.Diff{{Proto: "tcp", Port: 9999, State: "opened"}}
	got := e.Enrich(diffs)

	if got[0].Tag != "" {
		t.Errorf("expected empty tag, got %q", got[0].Tag)
	}
}

func TestEnrich_DeviationWhenNotInBaseline(t *testing.T) {
	reg := tags.New()
	baseline := map[string]bool{"tcp:80": true}
	e := enrich.New(reg, baseline)

	diffs := []scanner.Diff{
		{Proto: "tcp", Port: 80, State: "opened"},
		{Proto: "tcp", Port: 443, State: "opened"},
	}
	got := e.Enrich(diffs)

	if got[0].Deviation {
		t.Error("port 80 should not be a deviation")
	}
	if !got[1].Deviation {
		t.Error("port 443 should be a deviation")
	}
}

func TestEnrich_EmptyDiffsReturnsEmpty(t *testing.T) {
	e := enrich.New(tags.New(), nil)
	got := e.Enrich(nil)
	if len(got) != 0 {
		t.Errorf("expected empty, got %d entries", len(got))
	}
}

func TestEnrich_NilBaselineTreatedAsEmpty(t *testing.T) {
	e := enrich.New(tags.New(), nil)
	diffs := []scanner.Diff{{Proto: "udp", Port: 53, State: "opened"}}
	got := e.Enrich(diffs)
	if !got[0].Deviation {
		t.Error("expected deviation when baseline is nil")
	}
}
