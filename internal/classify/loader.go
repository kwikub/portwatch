package classify

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Load reads classification rules from a file.
// Each-blank, non-comment line has format:
//n//	<min <protocol> <level>
//023 tcp critical
//	1024-49151 tcp warning
func Load(path string) ([]Rule, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var rules []Rule
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 3 {
			return nil, fmt.Errorf("classify: malformed line: %q", line)
		}
		rangeParts := strings.SplitN(parts[0], "-", 2)
		if len(rangeParts) != 2 {
			return nil, fmt.Errorf("classify: invalid range: %q", parts[0])
		}
		min, err := strconv.Atoi(rangeParts[0])
		if err != nil {
			return nil, fmt.Errorf("classify: invalid min port: %w", err)
		}
		max, err := strconv.Atoi(rangeParts[1])
		if err != nil {
			return nil, fmt.Errorf("classify: invalid max port: %w", err)
		}
		rules = append(rules, Rule{
			MinPort:  min,
			MaxPort:  max,
			Protocol: parts[1],
			Level:    Level(parts[2]),
		})
	}
	return rules, sc.Err()
}
