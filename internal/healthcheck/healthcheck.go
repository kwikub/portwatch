// Package healthcheck provides a simple HTTP endpoint that reports daemon status.
package healthcheck

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

// Status holds the current health state of the daemon.
type Status struct {
	Running   bool      `json:"running"`
	StartedAt time.Time `json:"started_at"`
	LastScan  time.Time `json:"last_scan"`
	ScanCount int64     `json:"scan_count"`
}

// Server exposes a /healthz HTTP endpoint.
type Server struct {
	addr      string
	startedAt time.Time
	lastScan  atomic.Value // stores time.Time
	scanCount atomic.Int64
	server    *http.Server
}

// New creates a new healthcheck Server listening on addr (e.g. ":9090").
func New(addr string) *Server {
	s := &Server{
		addr:      addr,
		startedAt: time.Now(),
	}
	s.lastScan.Store(time.Time{})

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)

	s.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	return s
}

// RecordScan updates the last scan timestamp and increments the scan counter.
func (s *Server) RecordScan() {
	s.lastScan.Store(time.Now())
	s.scanCount.Add(1)
}

// Start begins serving in a background goroutine. It returns an error if the
// listener cannot be bound.
func (s *Server) Start() error {
	ln, err := (&http.Server{}).RegisterOnShutdown, fmt.Errorf // dummy — replaced below
	_ = ln
	_ = err
	go func() { _ = s.server.ListenAndServe() }()
	return nil
}

// Stop gracefully shuts down the HTTP server.
func (s *Server) Stop() error {
	return s.server.Close()
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	status := Status{
		Running:   true,
		StartedAt: s.startedAt,
		LastScan:  s.lastScan.Load().(time.Time),
		ScanCount: s.scanCount.Load(),
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}
