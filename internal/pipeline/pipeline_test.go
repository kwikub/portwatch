package pipeline_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/dedupe"
	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/reporter"
	"github.com/user/portwatch/internal/scanner"
)

func diffs() []scanner.Diff {
	return []scanner.Diff{
		{Port: 80, Proto: "tcp", State: scanner.Opened},
	}
}

func TestRun_EmptyDiffsIsNoop(t *testing.T) {
	var buf bytes.Buffer
	rep := reporter.New(reporter.Config{Writer: &buf})
	p := pipeline.New(pipeline.Config{Reporter: rep})
	if err := p.Run(context.Background(), nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Fatalf("expected no output, got %q", buf.String())
	}
}

func TestRun_ReporterReceivesDiffs(t *testing.T) {
	var buf bytes.Buffer
	rep := reporter.New(reporter.Config{Writer: &buf})
	p := pipeline.New(pipeline.Config{Reporter: rep})
	if err := p.Run(context.Background(), diffs()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Fatal("expected reporter output")
	}
}

func TestRun_DedupeBlocksDuplicate(t *testing.T) {
	var buf bytes.Buffer
	rep := reporter.New(reporter.Config{Writer: &buf})
	dd := dedupe.New(dedupe.Config{Window: time.Minute})
	p := pipeline.New(pipeline.Config{Reporter: rep, Dedupe: dd})

	_ = p.Run(context.Background(), diffs())
	buf.Reset()
	_ = p.Run(context.Background(), diffs())

	if buf.Len() != 0 {
		t.Fatalf("expected duplicate to be suppressed, got %q", buf.String())
	}
}

func TestRun_ExtraStageCanDropAll(t *testing.T) {
	var buf bytes.Buffer
	rep := reporter.New(reporter.Config{Writer: &buf})
	dropAll := func(d []scanner.Diff) []scanner.Diff { return nil }
	p := pipeline.New(pipeline.Config{Reporter: rep, Extra: []pipeline.Stage{dropAll}})

	_ = p.Run(context.Background(), diffs())
	if buf.Len() != 0 {
		t.Fatalf("expected stage to drop all diffs")
	}
}

func TestRun_MultipleExtraStages(t *testing.T) {
	var buf bytes.Buffer
	rep := reporter.New(reporter.Config{Writer: &buf})
	called := 0
	count := func(d []scanner.Diff) []scanner.Diff { called++; return d }
	p := pipeline.New(pipeline.Config{Reporter: rep, Extra: []pipeline.Stage{count, count}})

	_ = p.Run(context.Background(), diffs())
	if called != 2 {
		t.Fatalf("expected 2 stage calls, got %d", called)
	}
}
