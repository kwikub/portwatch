package portgroup

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Load reads a portgroup definition file into the registry.
 non-blank non-comment line has the format:
//n//	g port/protocol
//
// Example 80/tcp
//	w	dns 53/udp
func Load(path string, r *Registry) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	accum := map[string][]Entry{}
	sc := bufio.NewScanner(f)
	line := 0
	for sc.Scan() {
		line++
		text := strings.TrimSpace(sc.Text())
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}
		parts := strings.Fields(text)
		if len(parts) != 2 {
			return fmt.Errorf("portgroup: line %d: expected 'name port/proto', got %q", line, text)
		}
		name := parts[0]
		pp := strings.SplitN(parts[1], "/", 2)
		if len(pp) != 2 {
			return fmt.Errorf("portgroup: line %d: invalid port/proto %q", line, parts[1])
		}
		port, err := strconv.Atoi(pp[0])
		if err != nil {
			return fmt.Errorf("portgroup: line %d: invalid port %q", line, pp[0])
		}
		accum[name] = append(accum[name], Entry{Port: port, Protocol: pp[1]})
	}
	if err := sc.Err(); err != nil {
		return err
	}
	for name, entries := range accum {
		if err := r.Add(name, entries); err != nil {
			return err
		}
	}
	return nil
}
