package portfence_test

import (
	"testing"

	"github.com/user/portwatch/internal/portfence"
	"github.com/user/portwatch/internal/scanner"
)

func TestPermitted_EmptyAllowlistAlwaysTrue(t *testing.T) {
	f := portfence.New()
	if !f.Permitted(8080, "tcp") {
		t.Fatal("expected empty fence to permit every port")
	}
}

func TestAdd_InvalidPortReturnsError(t *testing.T) {
	f := portfence.New()
	if err := f.Add(0, "tcp"); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := f.Add(65536, "tcp"); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestAdd_InvalidProtocolReturnsError(t *testing.T) {
	f := portfence.New()
	if err := f.Add(80, "icmp"); err == nil {
		t.Fatal("expected error for unknown protocol")
	}
}

func TestPermitted_AllowedPortReturnsTrue(t *testing.T) {
	f := portfence.New()
	if err := f.Add(443, "tcp"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if !f.Permitted(443, "tcp") {
		t.Fatal("expected 443/tcp to be permitted")
	}
}

func TestPermitted_UnlistedPortReturnsFalse(t *testing.T) {
	f := portfence.New()
	_ = f.Add(443, "tcp")
	if f.Permitted(80, "tcp") {
		t.Fatal("expected 80/tcp to be blocked")
	}
}

func TestPermitted_ProtocolDistinct(t *testing.T) {
	f := portfence.New()
	_ = f.Add(53, "udp")
	if f.Permitted(53, "tcp") {
		t.Fatal("expected 53/tcp to be blocked when only 53/udp is allowed")
	}
	if !f.Permitted(53, "udp") {
		t.Fatal("expected 53/udp to be permitted")
	}
}

func TestFilter_RemovesUnpermittedDiffs(t *testing.T) {
	f := portfence.New()
	_ = f.Add(22, "tcp")

	diffs := []scanner.Diff{
		{Port: 22, Proto: "tcp", State: "opened"},
		{Port: 80, Proto: "tcp", State: "opened"},
		{Port: 53, Proto: "udp", State: "opened"},
	}

	got := f.Filter(diffs)
	if len(got) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(got))
	}
	if got[0].Port != 22 {
		t.Fatalf("expected port 22, got %d", got[0].Port)
	}
}

func TestFilter_EmptyAllowlistReturnsAll(t *testing.T) {
	f := portfence.New()
	diffs := []scanner.Diff{
		{Port: 80, Proto: "tcp"},
		{Port: 443, Proto: "tcp"},
	}
	got := f.Filter(diffs)
	if len(got) != 2 {
		t.Fatalf("expected 2 diffs, got %d", len(got))
	}
}
