// Package pipeline wires scanner diffs through filter, dedupe, throttle,
// alert, reporter and notify stages in a single reusable chain.
package pipeline

import (
	"context"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/dedupe"
	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/reporter"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/throttle"
)

// Stage is a function that filters or transforms a diff slice.
type Stage func([]scanner.Diff) []scanner.Diff

// Pipeline runs diffs through an ordered set of stages then hands the
// survivors to reporter and notify.
type Pipeline struct {
	stages   []Stage
	rep      *reporter.Reporter
	notifier *notify.Notifier
	alert    *alert.Evaluator
	dedupe   *dedupe.Filter
	throttle *throttle.Throttle
}

// Config holds dependencies injected at construction time.
type Config struct {
	Reporter *reporter.Reporter
	Notifier *notify.Notifier
	Alert    *alert.Evaluator
	Dedupe   *dedupe.Filter
	Throttle *throttle.Throttle
	Extra    []Stage
}

// New constructs a Pipeline from cfg.
func New(cfg Config) *Pipeline {
	return &Pipeline{
		rep:      cfg.Reporter,
		notifier: cfg.Notifier,
		alert:    cfg.Alert,
		dedupe:   cfg.Dedupe,
		throttle: cfg.Throttle,
		stages:   cfg.Extra,
	}
}

// Run processes diffs through every stage and dispatches results.
func (p *Pipeline) Run(ctx context.Context, diffs []scanner.Diff) error {
	if len(diffs) == 0 {
		return nil
	}

	// built-in stages
	if p.dedupe != nil {
		diffs = p.dedupe.Filter(diffs)
	}
	if p.throttle != nil {
		diffs = p.throttle.Allow(diffs)
	}
	for _, s := range p.stages {
		diffs = s(diffs)
		if len(diffs) == 0 {
			return nil
		}
	}

	if p.rep != nil {
		p.rep.Report(diffs)
	}
	if p.alert != nil {
		p.alert.Evaluate(diffs)
	}
	if p.notifier != nil {
		return p.notifier.Dispatch(ctx, diffs)
	}
	return nil
}
