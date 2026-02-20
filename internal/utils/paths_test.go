package utils

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolvePaths(t *testing.T) {
	root := t.TempDir()
	abs1 := filepath.Join(root, "abs", "path1")
	abs2 := filepath.Join(root, "abs", "path2")
	baseDir := filepath.Join(root, "base")
	baseSub := filepath.Join(root, "base", "sub")

	tests := []struct {
		name     string
		paths    []string
		baseDir  string
		expected []string
	}{
		{
			name:     "empty list",
			paths:    []string{},
			baseDir:  baseDir,
			expected: nil,
		},
		{
			name:     "nil list",
			paths:    nil,
			baseDir:  baseDir,
			expected: nil,
		},
		{
			name:     "absolute paths unchanged",
			paths:    []string{abs1, abs2},
			baseDir:  baseDir,
			expected: []string{abs1, abs2},
		},
		{
			name:     "relative paths resolved",
			paths:    []string{"rel1", "rel2/sub"},
			baseDir:  baseDir,
			expected: []string{filepath.Join(baseDir, "rel1"), filepath.Join(baseDir, "rel2", "sub")},
		},
		{
			name:     "mixed paths",
			paths:    []string{abs1, "rel", "../parent"},
			baseDir:  baseSub,
			expected: []string{abs1, filepath.Join(baseSub, "rel"), filepath.Join(root, "base", "parent")},
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
