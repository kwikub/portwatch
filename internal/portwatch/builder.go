package portwatch

import (
	"io"

	"github.com/user/portwatch/internal/baseline"
	"github.com/user/portwatch/internal/enrich"
	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
	"github.com/user/portwatch/internal/tags"
)

// Builder provides a fluent API for constructing a Coordinator.
type Builder struct {
	host      string
	portStart int
	portEnd   int
	protocol  string
	statePath string
	output    io.Writer
	tagReg    *tags.Registry
	baselineP *baseline.Baseline
}

// NewBuilder returns a Builder with sensible defaults.
func NewBuilder() *Builder {
	return &Builder{
		host:      "127.0.0.1",
		portStart: 1,
		portEnd:   65535,
		protocol:  "tcp",
	}
}

func (b *Builder) WithHost(h string) *Builder        { b.host = h; return b }
func (b *Builder) WithRange(lo, hi int) *Builder     { b.portStart = lo; b.portEnd = hi; return b }
func (b *Builder) WithProtocol(p string) *Builder    { b.protocol = p; return b }
func (b *Builder) WithStatePath(p string) *Builder   { b.statePath = p; return b }
func (b *Builder) WithOutput(w io.Writer) *Builder   { b.output = w; return b }
func (b *Builder) WithTags(r *tags.Registry) *Builder { b.tagReg = r; return b }
func (b *Builder) WithBaseline(bl *baseline.Baseline) *Builder { b.baselineP = bl; return b }

// Build assembles and returns a ready-to-use Coordinator.
func (b *Builder) Build() (*Coordinator, error) {
	sc := scanner.New(scanner.Options{
		Host:      b.host,
		PortStart: b.portStart,
		PortEnd:   b.portEnd,
		Protocol:  b.protocol,
	})

	st, err := state.New(b.statePath)
	if err != nil {
		return nil, err
	}

	enOpts := enrich.Options{}
	if b.tagReg != nil {
		enOpts.Tags = b.tagReg
	}
	if b.baselineP != nil {
		enOpts.Baseline = b.baselineP
	}
	en := enrich.New(enOpts)

	plOpts := pipeline.Options{}
	if b.output != nil {
		plOpts.Output = b.output
	}
	pl := pipeline.New(plOpts)

	return New(Config{
		Scanner:  sc,
		State:    st,
		Enricher: en,
		Pipeline: pl,
	}), nil
}
