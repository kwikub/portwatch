// Package audit provides a persistent audit trail of port change events.
package audit

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry represents a single audit log entry.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Proto     string    `json:"proto"`
	Port      int       `json:"port"`
	State     string    `json:"state"`
}

// Auditor appends port change events to a JSON-lines file.
type Auditor struct {
	mu   sync.Mutex
	path string
}

// New returns an Auditor that writes to path.
func New(path string) *Auditor {
	return &Auditor{path: path}
}

// Record appends one entry per diff to the audit file.
func (a *Auditor) Record(diffs []scanner.Diff) error {
	if len(diffs) == 0 {
		return nil
	}
	a.mu.Lock()
	defer a.mu.Unlock()

	f, err := os.OpenFile(a.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	for _, d := range diffs {
		state := "opened"
		if d.Closed {
			state = "closed"
		}
		e := Entry{
			Timestamp: time.Now().UTC(),
			Proto:     d.Proto,
			Port:      d.Port,
			State:     state,
		}
		if err := enc.Encode(e); err != nil {
			return err
		}
	}
	return nil
}
