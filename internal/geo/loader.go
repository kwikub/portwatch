package geo

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Load reads a TSV file of the form:
//
//	<ip>\t<country>\t<city>\t<asn>
//
// Lines beginning with '#' are treated as comments.
func Load(path string, r *Registry) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("geo: open %s: %w", path, err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	line := 0
	for sc.Scan() {
		line++
		text := strings.TrimSpace(sc.Text())
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}
		parts := strings.SplitN(text, "\t", 4)
		if len(parts) != 4 {
			return fmt.Errorf("geo: %s line %d: expected 4 fields, got %d", path, line, len(parts))
		}
		loc := Location{
			Country: parts[1],
			City:    parts[2],
			ASN:     parts[3],
		}
		if err := r.Register(parts[0], loc); err != nil {
			return fmt.Errorf("geo: %s line %d: %w", path, line, err)
		}
	}
	return sc.Err()
}
