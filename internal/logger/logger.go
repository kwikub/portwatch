package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Logger writes port change events to an output.
type Logger struct {
	out io.Writer
}

// New creates a Logger writing to the given writer.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}
	return &Logger{out: w}
}

// Log writes a human-readable line for each diff entry.
func (l *Logger) Log(diffs []scanner.DiffEntry) {
	for _, d := range diffs {
		timestamp := time.Now().Format(time.RFC3339)
		fmt.Fprintf(l.out, "%s  %-8s  %s:%d\n",
			timestamp, label(d.State), d.Entry.Host, d.Entry.Port)
	}
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
