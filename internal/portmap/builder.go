package portmap

import "os"

// FromFile loads a Map from a file at path, returning an empty Map if the
// file does not exist.
func FromFile(path string) (*Map, error) {
	m := New()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return m, nil
	}
	if err := Load(m, path); err != nil {
		return nil, err
	}
	return m, nil
}
