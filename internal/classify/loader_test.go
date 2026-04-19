package classify_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/classify"
)

func writeFile(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "rules.conf")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoad_ValidFile(t *testing.T) {
	p := writeFile(t, "# comment\n0-1023 tcp critical\n1024-49151 udp warning\n")
	rules, err := classify.Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Level != classify.LevelCritical {
		t.Errorf("rule 0: want critical, got %s", rules[0].Level)
	}
	if rules[1].Protocol != "udp" {
		t.Errorf("rule 1: want udp, got %s", rules[1].Protocol)
	}
}

func TestLoad_MissingFileReturnsNil(t *testing.T) {
	rules, err := classify.Load("/nonexistent/path/rules.conf")
	if err != nil {
		t.Fatal(err)
	}
	if rules != nil {
		t.Fatal("expected nil rules for missing file")
	}
}

func TestLoad_MalformedLineReturnsError(t *testing.T) {
	p := writeFile(t, "badline\n")
	_, err := classify.Load(p)
	if err == nil {
		t.Fatal("expected error for malformed line")
	}
}

func TestLoad_BlankLinesIgnored(t *testing.T) {
	p := writeFile(t, "\n\n0-1023 tcp critical\n\n")
	rules, err := classify.Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
}
