package portping

import "time"

// Builder constructs a Prober with optional configuration.
type Builder struct {
	host    string
	timeout time.Duration
}

// NewBuilder returns a Builder with defaults.
func NewBuilder() *Builder {
	return &Builder{
		host:    "127.0.0.1",
		timeout: 2 * time.Second,
	}
}

// WithHost sets the target host.
func (b *Builder) WithHost(host string) *Builder {
	b.host = host
	return b
}

// WithTimeout sets the dial timeout.
func (b *Builder) WithTimeout(d time.Duration) *Builder {
	b.timeout = d
	return b
}

// Build returns the configured Prober.
func (b *Builder) Build() *Prober {
	return New(b.host, b.timeout)
}
