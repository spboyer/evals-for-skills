package trigger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Discover searches dir for a trigger test file (trigger_tests.yaml)
// and returns the parsed spec. Returns nil, nil when no test file is found.
func Discover(dir string) (*TestSpec, error) {
	p := filepath.Join(dir, "trigger_tests.yaml")
	data, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	spec, err := ParseSpec(data)
	if err != nil {
		return nil, fmt.Errorf("loading trigger tests from %s: %w", p, err)
	}
	return spec, nil
}
