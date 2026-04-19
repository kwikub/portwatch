package portmap

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Load populates m from a TSV file with columns: port, protocol, name, groupn// Lines beginning with '# blank lines are ignored.
 trailing columns default to empty string.
func Load *Map, path)
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
		parts := strings.Split(line, "\t")
		if len(parts) < 2 {
			return fmt.Errorf("portmap: line %d: expected at least port and protocol", lineNo)
		}
		port, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return fmt.Errorf("portmap: line %d: invalid port %q", lineNo, parts[0])
		}
		col := func(i int) string {
			if i < len(parts) {
				return strings.TrimSpace(parts[i])
			}
			return ""
		}
		m.Set(Entry{
			Port:     port,
			Protocol: col(1),
			Name:     col(2),
			Group:    col(3),
			Tag:      col(4),
		})
	}
	return sc.Err()
}
