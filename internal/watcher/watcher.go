package watcher

import (
	"time"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/logger"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

// Watcher orchestrates periodic port scanning and diff logging.
type Watcher struct {
	cfg    *config.Config
	log    *logger.Logger
	state  *state.State
	stop   chan struct{}
}

// New creates a new Watcher.
func New(cfg *config.Config, log *logger.Logger, st *state.State) *Watcher {
	return &Watcher{
		cfg:   cfg,
		log:   log,
		state: st,
		stop:  make(chan struct{}),
	}
}

// Start begins the watch loop, blocking until Stop is called.
func (w *Watcher) Start() error {
	if err := w.tick(); err != nil {
		return err
	}
	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := w.tick(); err != nil {
				return err
			}
		case <-w.stop:
			return nil
		}
	}
}

// Stop signals the watch loop to exit.
func (w *Watcher) Stop() {
	close(w.stop)
}

func (w *Watcher) tick() error {
	sc := scanner.New(w.cfg)
	current, err := sc.Scan()
	if err != nil {
		return err
	}
	prev := w.state.Last()
	diffs := scanner.Diff(prev, current)
	w.log.Log(diffs)
	return w.state.Save(current)
}
