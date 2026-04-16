package config

import (
	"errors"
	"flag"
	"fmt"
	"time"
)

// Config holds all runtime configuration for portwatch.
type Config struct {
	Host      string
	Protocol  string
	StartPort int
	EndPort   int
	Interval  time.Duration
	Timeout   time.Duration
	StatePath string
	LogFile   string
}

// Parse reads configuration from command-line arguments.
func Parse(args []string) (*Config, error) {
	fs := flag.NewFlagSet("portwatch", flag.ContinueOnError)

	host := fs.String("host", "127.0.0.1", "host to scan")
	proto := fs.String("proto", "tcp", "protocol: tcp or udp")
	start := fs.Int("start", 1, "start port")
	end := fs.Int("end", 1024, "end port")
	interval := fs.Duration("interval", 30*time.Second, "scan interval")
	timeout := fs.Duration("timeout", 500*time.Millisecond, "per-port dial timeout")
	statePath := fs.String("state", "/tmp/portwatch_state.json", "path to state file")
	logFile := fs.String("log", "", "log file path (default stdout)")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	cfg := &Config{
		Host:      *host,
		Protocol:  *proto,
		StartPort: *start,
		EndPort:   *end,
		Interval:  *interval,
		Timeout:   *timeout,
		StatePath: *statePath,
		LogFile:   *logFile,
	}
	return cfg, validate(cfg)
}

func validate(c *Config) error {
	if c.Protocol != "tcp" && c.Protocol != "udp" {
		return fmt.Errorf("invalid protocol %q: must be tcp or udp", c.Protocol)
	}
	if c.StartPort < 1 || c.EndPort > 65535 || c.StartPort > c.EndPort {
		return errors.New("invalid port range")
	}
	return nil
}
