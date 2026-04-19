package portexpiry

import "time"

// Builder constructs a Tracker with a fluent API.
type Builder struct {
	ttl time.Duration
}

// NewBuilder returns a Builder with a default TTL of 24 hours.
func NewBuilder() *Builder {
	return &Builder{ttl: 24 * time.Hour}
}

// WithTTL sets the expiry duration.
func (b *Builder) WithTTL(d time.Duration) *Builder {
	b.ttl = d
	return b
}

// Build returns the configured Tracker.
func (b *Builder) Build() *Tracker {
	return New(b.ttl)
}
