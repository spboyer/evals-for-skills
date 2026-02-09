package orchestration

import (
	"fmt"
	"path/filepath"

	"github.com/spboyer/waza/internal/models"
)

// FilterTestCases returns the subset of testCases whose DisplayName or TestID
// matches at least one of the given glob patterns. An empty patterns slice
// returns all test cases unchanged.
func FilterTestCases(testCases []*models.TestCase, patterns []string) ([]*models.TestCase, error) {
	if len(patterns) == 0 {
		return testCases, nil
	}

	var matched []*models.TestCase
	for _, tc := range testCases {
		ok, err := matchesAny(tc, patterns)
		if err != nil {
			return nil, err
		}
		if ok {
			matched = append(matched, tc)
		}
	}
	return matched, nil
}

// matchesAny reports whether a test case's DisplayName or TestID matches any pattern.
func matchesAny(tc *models.TestCase, patterns []string) (bool, error) {
	for _, p := range patterns {
		nameMatch, err := filepath.Match(p, tc.DisplayName)
		if err != nil {
			return false, fmt.Errorf("invalid task filter pattern %q: %w", p, err)
		}
		if nameMatch {
			return true, nil
		}
		idMatch, err := filepath.Match(p, tc.TestID)
		if err != nil {
			return false, fmt.Errorf("invalid task filter pattern %q: %w", p, err)
		}
		if idMatch {
			return true, nil
		}
	}
	return false, nil
}
