package pipeline_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/pipeline"
)

func TestBuilder_NilOutputNoReporter(t *testing.T) {
	p := pipeline.Builder(pipeline.BuilderConfig{})
	if err := p.Run(context.Background(), diffs()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuilder_WritesToOutput(t *testing.T) {
	var buf bytes.Buffer
	p := pipeline.Builder(pipeline.BuilderConfig{Output: &buf, Format: "text"})
	_ = p.Run(context.Background(), diffs())
	if buf.Len() == 0 {
		t.Fatal("expected output from builder pipeline")
	}
}

func TestBuilder_DedupeWindowApplied(t *testing.T) {
	var buf bytes.Buffer
	p := pipeline.Builder(pipeline.BuilderConfig{
		Output:       &buf,
		DedupeWindow: time.Minute,
	})
	_ = p.Run(context.Background(), diffs())
	buf.Reset()
	_ = p.Run(context.Background(), diffs())
	if buf.Len() != 0 {
		t.Fatal("expected dedupe to suppress second run")
	}
}

func TestBuilder_CSVFormat(t *testing.T) {
	var buf bytes.Buffer
	p := pipeline.Builder(pipeline.BuilderConfig{Output: &buf, Format: "csv"})
	_ = p.Run(context.Background(), diffs())
	if buf.Len() == 0 {
		t.Fatal("expected CSV output")
	}
}
