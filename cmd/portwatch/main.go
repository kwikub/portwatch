package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/logger"
	"github.com/user/portwatch/internal/scanner"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "portwatch: %v\n", err)
		os.Exit(1)
	}

	log,LogFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, to open log: %v\n", err)
		os.Exit(1)
	}
	defer log.Close()

	sc := scanner.New(cfg.Protocol, cfg.StartPort, cfg.EndPort)

	// Initial scan to establish baseline — no diff logged.
	prev, err := sc.Scan()
	if err != nil {
		fmt.Fprintf(os.Stderr, "portwatch: initial scan failed: %v\n", err)
		os.Exit(1)
	}

	ticker := time.NewTicker(time.Duration(cfg.Interval) * time.Second)
	defer ticker.Stop()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("portwatch: monitoring %s ports %d-%d every %ds\n",
		cfg.Protocol, cfg.StartPort, cfg.EndPort, cfg.Interval)

	for {
		select {
		case <-ticker.C:
			curr, err := sc.Scan()
			if err != nil {
				fmt.Fprintf(os.Stderr, "portwatch: scan error: %v\n", err)
				continue
			}
			diffs := scanner.Diff(prev, curr)
			if err := log.Log(diffs); err != nil {
				fmt.Fprintf(os.Stderr, "portwatch: log error: %v\n", err)
			}
			prev = curr
		case <-sigs:
			fmt.Println("\nportwatch: shutting down")
			return
		}
	}
}
