package portevent_test

import (
	"testing"

	"github.com/user/portwatch/internal/portevent"
	"github.com/user/portwatch/internal/scanner"
)

func diff(port int, proto string, open bool) scanner.Diff {
	return scanner.Diff{Port: port, Proto: proto, Open: open}
}

func TestPublish_DeliversToSubscriber(t *testing.T) {
	b := portevent.New()
	var got []portevent.Event
	b.Subscribe(func(e portevent.Event) { got = append(got, e) })

	b.Publish(portevent.Event{Type: portevent.EventOpened, Diff: diff(80, "tcp", true)})

	if len(got) != 1 {
		t.Fatalf("expected 1 event, got %d", len(got))
	}
	if got[0].Type != portevent.EventOpened {
		t.Errorf("expected EventOpened, got %s", got[0].Type)
	}
}

func TestPublish_MultipleSubscribers(t *testing.T) {
	b := portevent.New()
	count := 0
	b.Subscribe(func(portevent.Event) { count++ })
	b.Subscribe(func(portevent.Event) { count++ })

	b.Publish(portevent.Event{Type: portevent.EventClosed, Diff: diff(443, "tcp", false)})

	if count != 2 {
		t.Errorf("expected 2 handler calls, got %d", count)
	}
}

func TestPublishDiffs_SetsCorrectEventType(t *testing.T) {
	b := portevent.New()
	var events []portevent.Event
	b.Subscribe(func(e portevent.Event) { events = append(events, e) })

	diffs := []scanner.Diff{
		diff(22, "tcp", true),
		diff(8080, "tcp", false),
	}
	b.PublishDiffs(diffs)

	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].Type != portevent.EventOpened {
		t.Errorf("port 22 should be EventOpened")
	}
	if events[1].Type != portevent.EventClosed {
		t.Errorf("port 8080 should be EventClosed")
	}
}

func TestPublishDiffs_EmptyIsNoop(t *testing.T) {
	b := portevent.New()
	called := false
	b.Subscribe(func(portevent.Event) { called = true })
	b.PublishDiffs(nil)
	if called {
		t.Error("handler should not be called for empty diffs")
	}
}
