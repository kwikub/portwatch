package healthcheck_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/healthcheck"
)

func TestRecordScan_IncrementsCount(t *testing.T) {
	s := healthcheck.New(":0")
	s.RecordScan()
	s.RecordScan()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	// exercise via exported Status by calling the handler indirectly through a
	// test server.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.RecordScan() // third call
		_ = rec
		_ = req
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"ok": 1})
	}))
	defer ts.Close()

	// Direct struct inspection via a fresh recorder.
	if got := scanCountViaHandler(t, s); got < 2 {
		t.Fatalf("expected scan_count >= 2, got %d", got)
	}
}

func TestStatus_RunningIsTrue(t *testing.T) {
	s := healthcheck.New(":0")
	status := fetchStatus(t, s)
	if !status.Running {
		t.Fatal("expected running to be true")
	}
}

func TestStatus_StartedAtIsRecent(t *testing.T) {
	before := time.Now()
	s := healthcheck.New(":0")
	status := fetchStatus(t, s)
	if status.StartedAt.Before(before) {
		t.Fatal("started_at should be after test start")
	}
}

func TestStatus_LastScanZeroInitially(t *testing.T) {
	s := healthcheck.New(":0")
	status := fetchStatus(t, s)
	if !status.LastScan.IsZero() {
		t.Fatal("last_scan should be zero before any scan")
	}
}

func TestStatus_LastScanUpdatesAfterRecord(t *testing.T) {
	s := healthcheck.New(":0")
	before := time.Now()
	s.RecordScan()
	status := fetchStatus(t, s)
	if status.LastScan.Before(before) {
		t.Fatal("last_scan should be updated after RecordScan")
	}
}

// helpers

func handlerFor(s *healthcheck.Server) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		// Re-use exported behaviour by hitting the real server via httptest.
		// Since Start() binds a port, we instead spin a test server wrapping
		// the same mux that New() would create — we replicate via fetchStatus.
		_ = s
		w.WriteHeader(http.StatusOK)
	})
	return mux
}

func fetchStatus(t *testing.T, s *healthcheck.Server) healthcheck.Status {
	t.Helper()
	ts := httptest.NewServer(testHandler(s))
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/healthz")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	var st healthcheck.Status
	if err := json.NewDecoder(resp.Body).Decode(&st); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return st
}

func scanCountViaHandler(t *testing.T, s *healthcheck.Server) int64 {
	t.Helper()
	return fetchStatus(t, s).ScanCount
}

func testHandler(s *healthcheck.Server) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.ServeHealth)
	return mux
}
