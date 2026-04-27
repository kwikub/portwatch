package portroute

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Load reads a route file into the Registry.
//
// Each non-blank, non-comment line has the format:
//
//	<protocol> <port> <target> [group]
//
// Example:
//
//	tcp 443 api-gateway web
//	udp 53  dns-server
func Load(r *Registry, path string) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	lineNo := 0
	for sc.Scan() {
		lineNo++
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			return fmt.Errorf("portroute: line %d: expected at least 3 fields", lineNo)
		}
		proto := fields[0]
		port, err := strconv.Atoi(fields[1])
		if err != nil {
			return fmt.Errorf("portroute: line %d: invalid port %q", lineNo, fields[1])
		}
		target := fields[2]
		group := ""
		if len(fields) >= 4 {
			group = fields[3]
		}
		if err := r.Add(port, proto, target, group); err != nil {
			return fmt.Errorf("portroute: line %d: %w", lineNo, err)
		}
	}
	return sc.Err()
}
