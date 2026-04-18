// Package notify delivers alert notifications via pluggable channels.
package notify

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of a notification.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Message holds the data for a single notification.
type Message struct {
	Level     Level
	Title     string
	Body      string
	Timestamp time.Time
}

// Channel is anything that can send a notification.
type Channel interface {
	Send(msg Message) error
}

// Notifier fans a message out to one or more channels.
type Notifier struct {
	channels []Channel
}

// New returns a Notifier that writes to the supplied channels.
// If no channels are provided it falls back to a LogChannel writing to stdout.
func New(channels ...Channel) *Notifier {
	if len(channels) == 0 {
		channels = []Channel{NewLogChannel(os.Stdout)}
	}
	return &Notifier{channels: channels}
}

// Dispatch sends msg to every registered channel, collecting errors.
func (n *Notifier) Dispatch(msg Message) []error {
	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}
	var errs []error
	for _, ch := range n.channels {
		if err := ch.Send(msg); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// LogChannel writes notifications as plain-text lines to a writer.
type LogChannel struct {
	w io.Writer
}

// NewLogChannel returns a LogChannel that writes to w.
func NewLogChannel(w io.Writer) *LogChannel {
	return &LogChannel{w: w}
}

// Send implements Channel.
func (l *LogChannel) Send(msg Message) error {
	_, err := fmt.Fprintf(l.w, "%s [%s] %s: %s\n",
		msg.Timestamp.Format(time.RFC3339), msg.Level, msg.Title, msg.Body)
	return err
}
