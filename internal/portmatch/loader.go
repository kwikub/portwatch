package portmatch

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Load reads a rule file into m. Each non-blank, non-comment line must have
// the form:
//
//	<port-or-range> <protocol>
//
// Example:
//
//	80       tcp
//	8000-8999 *
//	443      tcp
func Load(m *Matcher, path string) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("portmatch: open %s: %w", path, err)
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
		if len(fields) != 2 {
			return fmt.Errorf("portmatch: %s:%d: expected '<port> <protocol>', got %q", path, lineNo, line)
		}
		if err := m.Add(fields[0], fields[1]); err != nil {
			return fmt.Errorf("portmatch: %s:%d: %w", path, lineNo, err)
		}
	}
	return sc.Err()
}
