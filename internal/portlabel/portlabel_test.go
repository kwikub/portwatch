package portlabel_test

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/classify"
	"github.com/user/portwatch/internal/portlabel"
	"github.com/user/portwatch/internal/portname"
	"github.com/user/portwatch/internal/scanner"
)

func makeResolver(t *testing.T) *portlabel.Resolver {
	t.Helper()
	names := portname.New()
	classes := classify.New(nil)
	return portlabel.New(names, classes)
}

func TestResolve_KnownPortHasName(t *testing.T) {
	names := portname.New()
	names.Register(80, "tcp", "http")
	r := portlabel.New(names, classify.New(nil))
	lbl := r.Resolve(scanner.Port{Port: 80, Proto: "tcp"})
	if lbl.Name != "http" {
		t.Fatalf("expected http, got %q", lbl.Name)
	}
}

func TestResolve_UnknownPortFallsBack(t *testing.T) {
	r := makeResolver(t)
	lbl := r.Resolve(scanner.Port{Port: 9999, Proto: "tcp"})
	if lbl.Name == "" {
		t.Fatal("expected non-empty fallback label")
	}
}

func TestResolve_DisplayContainsProtoAndPort(t *testing.T) {
	r := makeResolver(t)
	lbl := r.Resolve(scanner.Port{Port: 443, Proto: "tcp"})
	if !strings.Contains(lbl.Display, "tcp") || !strings.Contains(lbl.Display, "443") {
		t.Fatalf("unexpected display: %q", lbl.Display)
	}
}

func TestResolve_ClassDefaultsToInfo(t *testing.T) {
	r := makeResolver(t)
	lbl := r.Resolve(scanner.Port{Port: 5000, Proto: "tcp"})
	if lbl.Class != "info" {
		t.Fatalf("expected info, got %q", lbl.Class)
	}
}

func TestResolve_NilClassifierUsesInfo(t *testing.T) {
	names := portname.New()
	r := portlabel.New(names, nil)
	lbl := r.Resolve(scanner.Port{Port: 22, Proto: "tcp"})
	if lbl.Class != "info" {
		t.Fatalf("expected info, got %q", lbl.Class)
	}
}
