package config

import (
	"flag"
	"fmt"
	"os"
)

// Config holds all runtime configuration for portwatch.
type Config struct {
	StartPort int
	EndPort   int
	Interval  int    // seconds between scans
	LogFile   string // empty means stdout
	Protocol  string // "tcp" or "udp"
}

// Default values.
const (
	DefaultStartPort = 1
	DefaultEndPort   = 65535
	DefaultInterval  = 5
	DefaultProtocol  = "tcp"
)

// Parse reads flags from os.Args and returns a validated Config.
func Parse() (*Config, error) {
	fs := flag.NewFlagSet("portwatch", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	start := fs.Int("start", DefaultStartPort, "first port to scan")
	end := fs.Int("end", DefaultEndPort, "last port to scan")
	interval := fs.Int("interval", DefaultInterval, "seconds between scans")
	logFile := fs.String("log", "", "path to log file (default: stdout)")
	protocol := fs.String("proto", DefaultProtocol, "protocol to scan: tcp or udp")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	cfg := &Config{
		StartPort: *start,
		EndPort:   *end,
		Interval:  *interval,
		LogFile:   *logFile,
		Protocol:  *protocol,
	}

	return cfg, cfg.validate()
}

func (c *Config) validate() error {
	if c.StartPort < 1 || c.StartPort > 65535 {
		return fmt.Errorf("start port %d out of range [1, 65535]", c.StartPort)
	}
	if c.EndPort < 1 || c.EndPort > 65535 {
		return fmt.Errorf("end port %d out of range [1, 65535]", c.EndPort)
	}
	if c.StartPort > c.EndPort {
		return fmt.Errorf("start port %d must be <= end port %d", c.StartPort, c.EndPort)
	}
	if c.Interval < 1 {
		return fmt.Errorf("interval must be at least 1 second")
	}
	if c.Protocol != "tcp" && c.Protocol != "udp" {
		return fmt.Errorf("unsupported protocol %q: must be tcp or udp", c.Protocol)
	}
	return nil
}
