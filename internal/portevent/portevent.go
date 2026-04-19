// Package portevent provides a typed event bus for port change notifications.
package portevent

import (
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// EventType classifies a port change.
type EventType string

const (
	EventOpened EventType = "opened"
	EventClosed EventType = "closed"
)

// Event carries a single port change notification.
type Event struct {
	Type  EventType
	Diff  scanner.Diff
}

// Handler is a function that receives a port event.
type Handler func(Event)

// Bus dispatches port events to registered handlers.
type Bus struct {
	mu       sync.RWMutex
	handlers []Handler
}

// New returns an initialised Bus.
func New() *Bus {
	return &Bus{}
}

// Subscribe registers a handler to receive all future events.
func (b *Bus) Subscribe(h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers = append(b.handlers, h)
}

// Publish sends an event to every registered handler sequentially.
func (b *Bus) Publish(e Event) {
	b.mu.RLock()
	handlers := make([]Handler, len(b.handlers))
	copy(handlers, b.handlers)
	b.mu.RUnlock()

	for _, h := range handlers {
		h(e)
	}
}

// PublishDiffs converts a slice of diffs into events and publishes each one.
func (b *Bus) PublishDiffs(diffs []scanner.Diff) {
	for _, d := range diffs {
		et := EventOpened
		if !d.Open {
			et = EventClosed
		}
		b.Publish(Event{Type: et, Diff: d})
	}
}
