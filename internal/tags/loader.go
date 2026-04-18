package tags

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Load reads a tag file into a Registry.
// Each non-blank, non-comment line must have the form:
//
//	<port>/<protocol>  <label>
//
// Lines starting with '#' are ignored.
func Load(path string) (*Registry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("tags: open %s: %w", path, err)
	}
	defer f.Close()
	return parse(f)
}

func parse(r io.Reader) (*Registry, error) {
	reg := New()
	sc := bufio.NewScanner(r)
	lineNo := 0
	for sc.Scan() {
		lineNo++
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("tags: line %d: expected '<port>/<proto> <label>'", lineNo)
		}
		parts := strings.SplitN(fields[0], "/", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("tags: line %d: bad port/proto %q", lineNo, fields[0])
		}
		port, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("tags: line %d: bad port %q", lineNo, parts[0])
		}
		label := strings.Join(fields[1:], " ")
		if err := reg.Add(port, parts[1], label); err != nil {
			return nil, fmt.Errorf("tags: line %d: %w", lineNo, err)
		}
	}
	return reg, sc.Err()
}
