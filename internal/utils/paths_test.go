package utils

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testAbsRoot() string {
	if runtime.GOOS == "windows" {
		return `C:\`
	}
	return "/"
}

func TestResolvePaths(t *testing.T) {
	root := testAbsRoot()
	tests := []struct {
		name     string
		paths    []string
		baseDir  string
		expected []string
	}{
		{
			name:     "empty list",
			paths:    []string{},
			baseDir:  filepath.Join(root, "base"),
			expected: nil,
		},
		{
			name:     "nil list",
			paths:    nil,
			baseDir:  filepath.Join(root, "base"),
			expected: nil,
		},
		{
			name:     "absolute paths unchanged",
			paths:    []string{filepath.Join(root, "abs", "path1"), filepath.Join(root, "abs", "path2")},
			baseDir:  filepath.Join(root, "base"),
			expected: []string{filepath.Join(root, "abs", "path1"), filepath.Join(root, "abs", "path2")},
		},
		{
			name:     "relative paths resolved",
			paths:    []string{"rel1", "rel2/sub"},
			baseDir:  filepath.Join(root, "base"),
			expected: []string{filepath.Join(root, "base", "rel1"), filepath.Join(root, "base", "rel2", "sub")},
		},
		{
			name:     "mixed paths",
			paths:    []string{filepath.Join(root, "abs"), "rel", "../parent"},
			baseDir:  filepath.Join(root, "base", "sub"),
			expected: []string{filepath.Join(root, "abs"), filepath.Join(root, "base", "sub", "rel"), filepath.Join(root, "base", "parent")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ResolvePaths(tt.paths, tt.baseDir)

			// Clean paths for comparison (normalize separators and . .. references)
			if tt.expected != nil {
				cleanExpected := make([]string, len(tt.expected))
				for i, p := range tt.expected {
					cleanExpected[i] = filepath.Clean(p)
				}
				cleanResult := make([]string, len(result))
				for i, p := range result {
					cleanResult[i] = filepath.Clean(p)
				}
				assert.Equal(t, cleanExpected, cleanResult)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}
