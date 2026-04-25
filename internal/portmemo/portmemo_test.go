package portmemo

import (
	"sync"
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	m := New()
	m.Set(80, "tcp", "owner", "team-a")

	v, ok := m.Get(80, "tcp", "owner")
	if !ok {
		t.Fatal("expected key to be present")
	}
	if v != "team-a" {
		t.Fatalf("expected 'team-a', got %q", v)
	}
}

func TestGet_Missing(t *testing.T) {
	m := New()
	_, ok := m.Get(443, "tcp", "owner")
	if ok {
		t.Fatal("expected key to be absent")
	}
}

func TestSet_EmptyValueRemovesKey(t *testing.T) {
	m := New()
	m.Set(22, "tcp", "note", "ssh")
	m.Set(22, "tcp", "note", "")

	_, ok := m.Get(22, "tcp", "note")
	if ok {
		t.Fatal("expected key to be removed after empty-value set")
	}
}

func TestSet_EmptyValueCleansUpPort(t *testing.T) {
	m := New()
	m.Set(22, "tcp", "note", "ssh")
	m.Set(22, "tcp", "note", "")

	if ann := m.All(22, "tcp"); ann != nil {
		t.Fatalf("expected nil after all keys removed, got %v", ann)
	}
}

func TestAll_ReturnsAllAnnotations(t *testing.T) {
	m := New()
	m.Set(8080, "tcp", "env", "prod")
	m.Set(8080, "tcp", "team", "ops")

	ann := m.All(8080, "tcp")
	if len(ann) != 2 {
		t.Fatalf("expected 2 annotations, got %d", len(ann))
	}
	if ann["env"] != "prod" || ann["team"] != "ops" {
		t.Fatalf("unexpected annotations: %v", ann)
	}
}

func TestAll_IsCopy(t *testing.T) {
	m := New()
	m.Set(9090, "udp", "k", "v")
	ann := m.All(9090, "udp")
	ann["k"] = "mutated"

	v, _ := m.Get(9090, "udp", "k")
	if v != "v" {
		t.Fatal("All() should return a copy, not the internal map")
	}
}

func TestClear_RemovesAllAnnotations(t *testing.T) {
	m := New()
	m.Set(53, "udp", "role", "dns")
	m.Clear(53, "udp")

	if ann := m.All(53, "udp"); ann != nil {
		t.Fatalf("expected nil after Clear, got %v", ann)
	}
}

func TestProtocolDistinct(t *testing.T) {
	m := New()
	m.Set(53, "tcp", "role", "dns-tcp")
	m.Set(53, "udp", "role", "dns-udp")

	v, _ := m.Get(53, "tcp", "role")
	if v != "dns-tcp" {
		t.Fatalf("tcp annotation wrong: %q", v)
	}
	v, _ = m.Get(53, "udp", "role")
	if v != "dns-udp" {
		t.Fatalf("udp annotation wrong: %q", v)
	}
}

func TestSet_Concurrent(t *testing.T) {
	m := New()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			m.Set(port, "tcp", "key", "val")
			m.Get(port, "tcp", "key")
		}(i)
	}
	wg.Wait()
}
