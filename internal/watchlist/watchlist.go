// Package watchlist manages a set of ports that should always be monitored
// and trigger alerts when they transition to a closed state.
package watchlist

import (
	"fmt"
	"strings"

	"github.com/user/portwatch/internal/scanner"
)

// Entry represents a single watched port+protocol pair.
type Entry struct {
	Port     int
	Protocol string
}

// Watchlist holds a set of critical ports.
type Watchlist struct {
	entries map[string]Entry
}

// New creates an empty Watchlist.
func New() *Watchlist {
	return &Watchlist{entries: make(map[string]Entry)}
}

// Add registers a port/protocol pair for watching.
func (w *Watchlist) Add(port int, protocol string) error {
	protocol = strings.ToLower(protocol)
	if protocol != "tcp" && protocol != "udp" {
		return fmt.Errorf("watchlist: unsupported protocol %q", protocol)
	}
	if port < 1 || port > 65535 {
		return fmt.Errorf("watchlist: port %d out of range", port)
	}
	k := key(port, protocol)
	w.entries[k] = Entry{Port: port, Protocol: protocol}
	return nil
}

// Remove deregisters a port/protocol pair.
func (w *Watchlist) Remove(port int, protocol string) {
	delete(w.entries, key(port, strings.ToLower(protocol)))
}

// MissingFrom returns entries that are absent from the given snapshot.
func (w *Watchlist) MissingFrom(snap *scanner.Snapshot) []Entry {
	var missing []Entry
	for k, e := range w.entries {
		found := false
		for _, p := range snap.Ports {
			if key(p.Port, p.Protocol) == k {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, e)
		}
	}
	return missing
}

// Len returns the number of watched entries.
func (w *Watchlist) Len() int { return len(w.entries) }

func key(port int, protocol string) string {
	return fmt.Sprintf("%s:%d", protocol, port)
}
