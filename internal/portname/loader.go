package portname

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Load reads a file of lines in the format name and registers entry into r Lines beginning with '#' and blank lines are ignored.
func Load string) error {
	f fmt.Errorf("portname: open %s: %w", path, err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		parts := strings.SplitN(fields[0], "/", 2)
		if len(parts) != 2 {
			continue
		}
		port, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		r.Register(port, parts[1], fields[1])
	}
	return sc.Err()
}
