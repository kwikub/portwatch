// Package healthcheck exposes a lightweight HTTP /healthz endpoint so that
// process supervisors and monitoring tools can verify that the portwatch daemon
// is alive and report basic runtime statistics such as scan count and last scan
// timestamp.
package healthcheck
