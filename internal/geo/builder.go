package geo

import "os"

// FromFile loads a Registry from path if it exists.
// If path is empty or the file does not exist, an empty Registry is returned.
func FromFile(path string) (*Registry, error) {
	r := New()
	if path == "" {
		return r, nil
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return r, nil
	}
	if err := Load(path, r); err != nil {
		return nil, err
	}
	return r, nil
}
