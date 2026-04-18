package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"portwatch/internal/notify"
)

func TestWebhookChannel_SendsJSON(t *testing.T) {
	var received map[string]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			http.Error(w, "bad json", 400)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	ch := notify.NewWebhookChannel(srv.URL, time.Second)
	err := ch.Send(notify.Message{
		Level: notify.LevelAlert,
		Title: "port opened",
		Body:  "443/tcp",
		Timestamp: time.Now(),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["level"] != "ALERT" {
		t.Errorf("expected level ALERT, got %s", received["level"])
	}
	if received["body"] != "443/tcp" {
		t.Errorf("expected body 443/tcp, got %s", received["body"])
	}
}

func TestWebhookChannel_ErrorOnNon2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	ch := notify.NewWebhookChannel(srv.URL, time.Second)
	err := ch.Send(notify.Message{Level: notify.LevelInfo, Title: "t", Body: "b", Timestamp: time.Now()})
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}
