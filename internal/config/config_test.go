package config

import (
	"os"
	"testing"
)

func withArgs(args []string, fn func()) {
	old := os.Args
	os.Args = append([]string{"portwatch"}, args...)
	defer func() { os.Args = old }()
	fn()
}

func TestParse_Defaults(t *testing.T) {
	withArgs([]string{}, func() {
		cfg, err := Parse()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.StartPort != DefaultStartPort {
			t.Errorf("expected start %d, got %d", DefaultStartPort, cfg.StartPort)
		}
		if cfg.EndPort != DefaultEndPort {
			t.Errorf("expected end %d, got %d", DefaultEndPort, cfg.EndPort)
		}
		if cfg.Interval != DefaultInterval {
			t.Errorf("expected interval %d, got %d", DefaultInterval, cfg.Interval)
		}
		if cfg.Protocol != DefaultProtocol {
			t.Errorf("expected protocol %q, got %q", DefaultProtocol, cfg.Protocol)
		}
	})
}

func TestParse_CustomValues(t *testing.T) {
	withArgs([]string{"-start=8000", "-end=9000", "-interval=10", "-proto=tcp", "-log=/tmp/pw.log"}, func() {
		cfg, err := Parse()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.StartPort != 8000 || cfg.EndPort != 9000 {
			t.Errorf("port range mismatch")
		}
		if cfg.LogFile != "/tmp/pw.log" {
			t.Errorf("log file mismatch")
		}
	})
}

func TestValidate_InvalidRange(t *testing.T) {
	cfg := &Config{StartPort: 9000, EndPort: 8000, Interval: 5, Protocol: "tcp"}
	if err := cfg.validate(); err == nil {
		t.Error("expected error for inverted port range")
	}
}

func TestValidate_BadProtocol(t *testing.T) {
	cfg := &Config{StartPort: 1, EndPort: 100, Interval: 5, Protocol: "icmp"}
	if err := cfg.validate(); err == nil {
		t.Error("expected error for unsupported protocol")
	}
}

func TestValidate_BadInterval(t *testing.T) {
	cfg := &Config{StartPort: 1, EndPort: 100, Interval: 0, Protocol: "tcp"}
	if err := cfg.validate(); err == nil {
		t.Error("expected error for zero interval")
	}
}
