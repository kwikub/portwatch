package portquota

import "time"

// Builder constructs a Quota with a fluent API.
type Builder struct {
	window    time.Duration
	threshold int
}

// NewBuilder returns a Builder with sensible defaults:
// a 1-minute window and a threshold of 5 open events.
func NewBuilder() *Builder {
	return &Builder{
		window:    time.Minute,
		threshold: 5,
	}
}

// WithWindow sets the rolling window duration.
func (b *Builder) WithWindow(d time.Duration) *Builder {
	b.window = d
	return b
}

// WithThreshold sets the maximum number of open events allowed within the window.
func (b *Builder) WithThreshold(n int) *Builder {
	b.threshold = n
	return b
}

// Build returns the configured Quota.
func (b *Builder) Build() *Quota {
	return New(b.window, b.threshold)
}
