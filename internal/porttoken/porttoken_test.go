package porttoken

import (
	"strings"
	"testing"
)

func TestAssign_ReturnsToken(t *testing.T) {
	r := New(8)
	tok, err := r.Assign(443, "tcp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tok) != 16 { // 8 bytes → 16 hex chars
		t.Fatalf("expected 16 hex chars, got %d", len(tok))
	}
}

func TestAssign_SamePortReturnsSameToken(t *testing.T) {
	r := New(8)
	t1, _ := r.Assign(80, "tcp")
	t2, _ := r.Assign(80, "tcp")
	if t1 != t2 {
		t.Fatalf("expected same token, got %q and %q", t1, t2)
	}
}

func TestAssign_DifferentPortsGetDifferentTokens(t *testing.T) {
	r := New(8)
	t1, _ := r.Assign(80, "tcp")
	t2, _ := r.Assign(443, "tcp")
	if t1 == t2 {
		t.Fatal("expected different tokens for different ports")
	}
}

func TestAssign_ProtocolDistinct(t *testing.T) {
	r := New(8)
	t1, _ := r.Assign(53, "tcp")
	t2, _ := r.Assign(53, "udp")
	if t1 == t2 {
		t.Fatal("expected different tokens for different protocols")
	}
}

func TestLookup_Found(t *testing.T) {
	r := New(8)
	tok, _ := r.Assign(22, "tcp")
	key, ok := r.Lookup(tok)
	if !ok {
		t.Fatal("expected lookup to succeed")
	}
	if !strings.Contains(key, "22") || !strings.Contains(key, "tcp") {
		t.Fatalf("unexpected key: %q", key)
	}
}

func TestLookup_NotFound(t *testing.T) {
	r := New(8)
	_, ok := r.Lookup("deadbeef")
	if ok {
		t.Fatal("expected lookup to fail for unknown token")
	}
}

func TestRevoke_RemovesToken(t *testing.T) {
	r := New(8)
	tok, _ := r.Assign(8080, "tcp")
	r.Revoke(8080, "tcp")

	if _, ok := r.Lookup(tok); ok {
		t.Fatal("expected token to be revoked")
	}
	// Re-assigning should produce a fresh (possibly different) token.
	tok2, err := r.Assign(8080, "tcp")
	if err != nil {
		t.Fatalf("unexpected error after revoke: %v", err)
	}
	_ = tok2
}

func TestNew_PanicsOnZeroBytes(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero tokenBytes")
		}
	}()
	New(0)
}
