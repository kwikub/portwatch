// Package watcher provides the core watch loop for portwatch.
// It periodically invokes the scanner, computes diffs against the
// previously saved snapshot, logs any changes, and persists the
// latest snapshot via the state store.
package watcher
