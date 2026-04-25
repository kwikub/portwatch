// Package porttoken assigns and validates short opaque tokens to port+protocol
// pairs, allowing external systems to reference a port without exposing its
// numeric value directly.
package porttoken

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
)

// Registry maps port/protocol pairs to tokens and back.
type Registry struct {
	mu       sync.RWMutex
	toToken  map[string]string
	toPort   map[string]string
	tokenLen int
}

// New returns a Registry that generates tokens of tokenBytes random bytes
// (hex-encoded). tokenBytes must be > 0.
func New(tokenBytes int) *Registry {
	if tokenBytes <= 0 {
		panic("porttoken: tokenBytes must be > 0")
	}
	return &Registry{
		toToken:  make(map[string]string),
		toPort:   make(map[string]string),
		tokenLen: tokenBytes,
	}
}

func portKey(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// Assign returns the token for the given port/protocol, creating one if it
// does not already exist.
func (r *Registry) Assign(port int, proto string) (string, error) {
	k := portKey(port, proto)

	r.mu.Lock()
	defer r.mu.Unlock()

	if tok, ok := r.toToken[k]; ok {
		return tok, nil
	}

	buf := make([]byte, r.tokenLen)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("porttoken: generate token: %w", err)
	}
	tok := hex.EncodeToString(buf)

	r.toToken[k] = tok
	r.toPort[tok] = k
	return tok, nil
}

// Lookup returns the port/protocol key for the given token, or false if
// the token is unknown.
func (r *Registry) Lookup(token string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	k, ok := r.toPort[token]
	return k, ok
}

// Revoke removes the token associated with the given port/protocol pair.
func (r *Registry) Revoke(port int, proto string) {
	k := portKey(port, proto)

	r.mu.Lock()
	defer r.mu.Unlock()

	if tok, ok := r.toToken[k]; ok {
		delete(r.toPort, tok)
		delete(r.toToken, k)
	}
}
