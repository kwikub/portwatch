// Package notify provides a fan-out notification system for portwatch.
//
// Notifications can be delivered via pluggable Channel implementations.
// Built-in channels include:
//
//   - LogChannel: writes human-readable messages to an io.Writer
//   - WebhookChannel: sends a JSON payload via HTTP POST to a configured URL
//
// A Notifier fans out each event to all registered channels concurrently,
// collecting and returning any errors that occur during delivery.
package notify
