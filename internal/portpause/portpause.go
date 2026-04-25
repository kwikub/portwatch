// Package portpause allows individual ports to be temporarily paused
// from generating diff events for a configurable duration.
package portpause

import (
	"fmt"
	"sync"
	"time"
)

// Pauser tracks ports that have been paused and suppresses diffs for them
// until the pause duration expires.
type Pauser struct {
	mu      sync.Mutex
	paused  map[string]time.Time
	now     func() time.Time
}

// New returns a new Pauser.
func New() *Pauser {
	return &Pauser{
		paused: make(map[string]time.Time),
		now:    time.Now,
	}
}

func portKey(port int, proto string) string {
	return fmt.Sprintf("%s:%d", proto, port)
}

// Pause suppresses events for the given port/protocol pair for the
// specified duration. Calling Pause again before expiry extends the deadline.
func (p *Pauser) Pause(port int, proto string, d time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.paused[portKey(port, proto)] = p.now().Add(d)
}

// Lift removes a pause immediately, regardless of the remaining duration.
func (p *Pauser) Lift(port int, proto string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.paused, portKey(port, proto))
}

// IsPaused reports whether the given port/protocol pair is currently paused.
func (p *Pauser) IsPaused(port int, proto string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	expiry, ok := p.paused[portKey(port, proto)]
	if !ok {
		return false
	}
	if p.now().After(expiry) {
		delete(p.paused, portKey(port, proto))
		return false
	}
	return true
}

// Active returns all port keys that are currently paused (expiry not yet reached).
func (p *Pauser) Active() []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	now := p.now()
	out := make([]string, 0, len(p.paused))
	for k, expiry := range p.paused {
		if now.Before(expiry) {
			out = append(out, k)
		}
	}
	return out
}
