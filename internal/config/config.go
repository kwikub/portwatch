// Package config parses and validates portwatch CLI configuration.
package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

// Config holds all runtime configuration for portwatch.
type Config struct {
	PortStart      int
	PortEnd        int
	Protocol       string
	Interval       int    // seconds
	StatePath      string
	LogPath        string
	AlertThreshold int
	ExcludePorts   string // comma-separated
}

// Parse reads configuration from command-line arguments.
func Parse(args []string) (*Config, error) {
	fs := flag.NewFlagSet("portwatch", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	cfg := &Config{}
	fs.IntVar(&cfg.PortStart, "start", 1, "start of port range")
	fs.IntVar(&cfg.PortEnd, "end", 1024, "end of port range")
	fs.StringVar(&cfg.Protocol, "proto", "tcp", "protocol to scan (tcp|udp)")
	fs.IntVar(&cfg.Interval, "interval", 60, "scan interval in seconds")
	fs.StringVar(&cfg.StatePath, "state", "/tmp/portwatch.state", "path to state file")
	fs.StringVar(&cfg.LogPath, "log", "", "path to log file (default stdout)")
	fs.IntVar(&cfg.AlertThreshold, "alert-threshold", 10, "change count that triggers ALERT level")
	fs.StringVar(&cfg.ExcludePorts, "exclude", "", "comma-separated ports to exclude")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	return cfg, validate(cfg)
}

func validate(cfg *Config) error {
	if cfg.PortStart < 1 || cfg.PortEnd > 65535 || cfg.PortStart > cfg.PortEnd {
		return fmt.Errorf("invalid port range %d-%d", cfg.PortStart, cfg.PortEnd)
	}
	if cfg.Protocol != "tcp" && cfg.Protocol != "udp" {
		return errors.New("protocol must be tcp or udp")
	}
	if cfg.Interval < 1 {
		return errors.New("interval must be at least 1 second")
	}
	if cfg.AlertThreshold < 1 {
		return errors.New("alert-threshold must be at least 1")
	}
	return nil
}
