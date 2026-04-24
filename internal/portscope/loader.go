package portscope

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Load reads a scope definition file and populates s.
//
// Each non-blank, non-comment line must have the form:
//
//	<protocol> <lo> <hi>
//
// Example:
//
//	tcp 1 1024
//	udp 53 53
func Load(s *Scope, path string) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 3 {
			return fmt.Errorf("portscope: line %d: expected '<proto> <lo> <hi>', got %q", lineNo, line)
		}
		lo, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("portscope: line %d: invalid lo port %q", lineNo, parts[1])
		}
		hi, err := strconv.Atoi(parts[2])
		if err != nil {
			return fmt.Errorf("portscope: line %d: invalid hi port %q", lineNo, parts[2])
		}
		if err := s.Add(parts[0], lo, hi); err != nil {
			return fmt.Errorf("portscope: line %d: %w", lineNo, err)
		}
	}
	return scanner.Err()
}
