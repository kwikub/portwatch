package watchlist

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Load populates a Watchlist from a reader containing lines of "port/protocol".
// Lines starting with '#' and blank lines are ignored.
func Load(r io.Reader) (*Watchlist, error) {
	w := New()
	scanner := bufio.NewScanner(r)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "/", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("watchlist: line %d: expected port/protocol, got %q", lineNum, line)
		}
		port, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, fmt.Errorf("watchlist: line %d: invalid port %q", lineNum, parts[0])
		}
		protocol := strings.TrimSpace(parts[1])
		if err := w.Add(port, protocol); err != nil {
			return nil, fmt.Errorf("watchlist: line %d: %w", lineNum, err)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return w, nil
}
