package config_test

import (
	"testing"

	"github.com/user/portwatch/internal/config"
)

func withArgs(args ...string) (*config.Config, error) {
	return config.Parse(args)
}

func TestParse_Defaults(t *testing.T) {
	cfg, err := withArgs()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.PortStart != 1 || cfg.PortEnd != 1024 {
		t.Errorf("unexpected range: %d-%d", cfg.PortStart, cfg.PortEnd)
	}
	if cfg.Protocol != "tcp" {
		t.Errorf("expected tcp, got %s", cfg.Protocol)
	}
	if cfg.AlertThreshold != 10 {
		t.Errorf("expected threshold 10, got %d", cfg.AlertThreshold)
	}
}

func TestParse_CustomValues(t *testing.T) {
	cfg, err := withArgs("-start=80", "-end=443", "-proto=tcp", "-alert-threshold=3")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.PortStart != 80 || cfg.PortEnd != 443 {
		t.Errorf("unexpected range: %d-%d", cfg.PortStart, cfg.PortEnd)
	}
	if cfg.AlertThreshold != 3 {
		t.Errorf("expected 3, got %d", cfg.AlertThreshold)
	}
}

func TestValidate_InvalidRange(t *testing.T) {
	_, err := withArgs("-start=500", "-end=100")
	if err == nil {
		t.Fatal("expected error for invalid range")
	}
}

func TestValidate_BadProtocol(t *testing.T) {
	_, err := withArgs("-proto=icmp")
	if err == nil {
		t.Fatal("expected error for bad protocol")
	}
}

func TestValidate_BadThreshold(t *testing.T) {
	_, err := withArgs("-alert-threshold=0")
	if err == nil {
		t.Fatal("expected error for zero threshold")
	}
}
