// Package rotation provides log file rotation for audit and report outputs.
package rotation

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Options configures rotation behaviour.
type Options struct {
	Dir        string
	Prefix     string
	MaxFiles   int
}

// Rotator manages a rolling set of timestamped log files.
type Rotator struct {
	mu      sync.Mutex
	opts    Options
	current *os.File
}

// New creates a Rotator and opens the first file.
func New(opts Options) (*Rotator, error) {
	if opts.Dir == "" {
		opts.Dir = "."
	}
	if opts.MaxFiles <= 0 {
		opts.MaxFiles = 5
	}
	if err := os.MkdirAll(opts.Dir, 0o755); err != nil {
		return nil, fmt.Errorf("rotation: mkdir: %w", err)
	}
	r := &Rotator{opts: opts}
	if err := r.rotate(); err != nil {
		return nil, err
	}
	return r, nil
}

// Write implements io.Writer, delegating to the current file.
func (r *Rotator) Write(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.current.Write(p)
}

// Rotate closes the current file and opens a new one, pruning old files.
func (r *Rotator) Rotate() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.rotate()
}

func (r *Rotator) rotate() error {
	if r.current != nil {
		_ = r.current.Close()
	}
	name := filepath.Join(r.opts.Dir, fmt.Sprintf("%s%s.log", r.opts.Prefix, time.Now().Format("20060102T150405")))
	f, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("rotation: create: %w", err)
	}
	r.current = f
	return r.prune()
}

func (r *Rotator) prune() error {
	pattern := filepath.Join(r.opts.Dir, r.opts.Prefix+"*.log")
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) <= r.opts.MaxFiles {
		return err
	}
	for _, old := range matches[:len(matches)-r.opts.MaxFiles] {
		_ = os.Remove(old)
	}
	return nil
}

// Close closes the current underlying file.
func (r *Rotator) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.current != nil {
		return r.current.Close()
	}
	return nil
}
