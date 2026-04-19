package geo

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRegister_ValidIP(t *testing.T) {
	r := New()
	err := r.Register("1.2.3.4", Location{Country: "US", City: "NYC", ASN: "AS1234"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRegister_InvalidIP(t *testing.T) {
	r := New()
	err := r.Register("not-an-ip", Location{})
	if err == nil {
		t.Fatal("expected error for invalid IP")
	}
}

func TestLookup_Found(t *testing.T) {
	r := New()
	_ = r.Register("10.0.0.1", Location{Country: "DE", City: "Berlin", ASN: "AS999"})
	loc, ok := r.Lookup("10.0.0.1")
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if loc.Country != "DE" {
		t.Errorf("country: got %q, want %q", loc.Country, "DE")
	}
}

func TestLookup_NotFound(t *testing.T) {
	r := New()
	_, ok := r.Lookup("9.9.9.9")
	if ok {
		t.Fatal("expected miss")
	}
}

func TestLoad_ValidFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "geo.tsv")
	content := "# comment\n192.168.1.1\tUS\tBoston\tAS5678\n"
	_ = os.WriteFile(p, []byte(content), 0o644)

	r := New()
	if err := Load(p, r); err != nil {
		t.Fatalf("Load: %v", err)
	}
	loc, ok := r.Lookup("192.168.1.1")
	if !ok {
		t.Fatal("expected loaded entry")
	}
	if loc.City != "Boston" {
		t.Errorf("city: got %q", loc.City)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	r := New()
	err := Load("/nonexistent/geo.tsv", r)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_MalformedLine(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "bad.tsv")
	_ = os.WriteFile(p, []byte("1.2.3.4\tUS\n"), 0o644)
	r := New()
	if err := Load(p, r); err == nil {
		t.Fatal("expected error for malformed line")
	}
}
