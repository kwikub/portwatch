package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Logger writes port diff events to an output sink.
type Logger struct {
	out    io.Writer
	closer io.Closer // non-nil when we own the file
}

// New creates a Logger. If path is empty, output goes to stdout.
func New(path string) (*Logger, error) {
	if path == "" {
		return &Logger{out: os.Stdout}, nil
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	return &Logger{out: f, closer: f}, nil
}

// Close releases any underlying file resource.
func (l *Logger) Close() error {
	if l.closer != nil {
		return l.closer.Close()
	}
	return nil
}

// Log writes one line per diff entry.
func (l *Logger) Log(diffs []scanner.Diff) error {
	for _, d := range diffs {
		ts := time.Now().UTC().Format(time.RFC3339)
		_, err := fmt.Fprintf(l.out, "%s  %-6s  %s/%d\n", ts, label(d.State), d.Protocol, d.Port)
		if err != nil {
			return err
		}
	}
	return nil
}

func label(state scanner.State) string {
	switch state {
	case scanner.StateOpened:
		return "OPENED"
	case scanner.StateClosed:
		return "CLOSED"
	default:
		return "UNKNOWN"
	}
}
