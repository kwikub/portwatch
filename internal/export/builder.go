package export

import (
	"io"
	"os"
)

// Builder constructs an Exporter from common configuration values.
type Builder struct {
	path   string
	format Format
}

// NewBuilder returns a Builder with sensible defaults (stdout, JSON).
func NewBuilder() *Builder {
	return &Builder{format: FormatJSON}
}

// WithFormat sets the output format.
func (b *Builder) WithFormat(f Format) *Builder {
	b.format = f
	return b
}

// WithPath sets a file path to write to instead of stdout.
func (b *Builder) WithPath(p string) *Builder {
	b.path = p
	return b
}

// Build creates the Exporter, opening the output file when a path is set.
// The caller is responsible for closing the returned io.Closer (may be nil
// when writing to stdout).
func (b *Builder) Build() (*Exporter, io.Closer, error) {
	if b.path == "" {
		return New(os.Stdout, b.format), nil, nil
	}
	f, err := os.OpenFile(b.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, nil, err
	}
	return New(f, b.format), f, nil
}
