package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookChannel posts notifications as JSON to an HTTP endpoint.
type WebhookChannel struct {
	url    string
	client *http.Client
}

// NewWebhookChannel returns a WebhookChannel targeting url.
func NewWebhookChannel(url string, timeout time.Duration) *WebhookChannel {
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	return &WebhookChannel{
		url:    url,
		client: &http.Client{Timeout: timeout},
	}
}

type webhookPayload struct {
	Level     string `json:"level"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Timestamp string `json:"timestamp"`
}

// Send implements Channel.
func (w *WebhookChannel) Send(msg Message) error {
	payload := webhookPayload{
		Level:     string(msg.Level),
		Title:     msg.Title,
		Body:      msg.Body,
		Timestamp: msg.Timestamp.Format(time.RFC3339),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("notify/webhook: marshal: %w", err)
	}
	resp, err := w.client.Post(w.url, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("notify/webhook: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("notify/webhook: unexpected status %d", resp.StatusCode)
	}
	return nil
}
