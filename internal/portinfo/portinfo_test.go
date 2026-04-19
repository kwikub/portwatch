package portinfo_test

import (
	"testing"

	"github.com/user/portwatch/internal/portinfo"
	"github.com/user/portwatch/internal/portname"
	"github.com/user/portwatch/internal/portgroup"
	"github.com/user/portwatch/internal/scanner"
)

func diff(port int, proto string) scanner.Diff {
	return scanner.Diff{Port: port, Proto: proto, State: scanner.Opened}
}

func TestResolve_KnownPortHasName(t *testing.T) {
	names := portname.New()
	names.Register(80, "tcp", "http")
	r := portinfo.New(names, nil)
	info := r.Resolve(diff(80, "tcp"))
	if info.Name != "http" {
		t.Fatalf("expected http, got %q", info.Name)
	}
	if info.Label != "http" {
		t.Fatalf("expected label http, got %q", info.Label)
	}
}

func TestResolve_UnknownPortFallsBackToPortProto(t *testing.T) {
	r := portinfo.New(nil, nil)
	info := r.Resolve(diff(9999, "udp"))
	if info.Label != "9999/udp" {
		t.Fatalf("unexpected label %q", info.Label)
	}
	if info.Name != "" {
		t.Fatalf("expected empty name, got %q", info.Name)
	}
}

func TestResolve_GroupMembership(t *testing.T) {
	groups := portgroup.New()
	_ = groups.Add("web", 80, "tcp")
	_ = groups.Add("web", 443, "tcp")
	r := portinfo.New(nil, groups)
	info := r.Resolve(diff(80, "tcp"))
	if len(info.Groups) != 1 || info.Groups[0] != "web" {
		t.Fatalf("expected [web], got %v", info.Groups)
	}
}

func TestResolve_NoGroupMatch(t *testing.T) {
	groups := portgroup.New()
	_ = groups.Add("web", 80, "tcp")
	r := portinfo.New(nil, groups)
	info := r.Resolve(diff(22, "tcp"))
	if len(info.Groups) != 0 {
		t.Fatalf("expected no groups, got %v", info.Groups)
	}
}
