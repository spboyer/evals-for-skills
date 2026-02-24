package checks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spboyer/waza/internal/skill"
	"github.com/stretchr/testify/require"
)

var _ ComplianceChecker = (*TokenLimitsChecker)(nil)

func skillInDir(dir string) skill.Skill {
	return skill.Skill{Path: filepath.Join(dir, "SKILL.md")}
}

func TestTokenLimitsChecker_AllPass(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "README.md"), []byte("# Hello"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# Skill"), 0o644))

	checker := &TokenLimitsChecker{}
	result, err := checker.Check(skillInDir(dir))
	require.NoError(t, err)
	require.True(t, result.Passed)
	require.Equal(t, "token-limits", result.Name)

	data, ok := result.Data.(*TokenLimitsData)
	require.True(t, ok)
	require.Equal(t, 2, data.TotalFiles)
	require.Equal(t, 0, data.ExceededCount)
}

func TestTokenLimitsChecker_SomeExceed(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "README.md"), []byte("# Hello"), 0o644))

	// Use a config with a very low limit
	checker := &TokenLimitsChecker{
		Config: TokenLimitsConfig{
			Defaults:  map[string]int{"*.md": 1},
			Overrides: map[string]int{},
		},
	}
	result, err := checker.Check(skillInDir(dir))
	require.NoError(t, err)
	require.False(t, result.Passed)

	data, ok := result.Data.(*TokenLimitsData)
	require.True(t, ok)
	require.Equal(t, 1, data.ExceededCount)
}

func TestTokenLimitsChecker_SpecificPaths(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "one.md"), []byte("# One"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "two.md"), []byte("# Two"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "three.md"), []byte("# Three"), 0o644))

	checker := &TokenLimitsChecker{
		Paths: []string{filepath.Join(dir, "one.md")},
	}
	result, err := checker.Check(skillInDir(dir))
	require.NoError(t, err)

	data, ok := result.Data.(*TokenLimitsData)
	require.True(t, ok)
	require.Equal(t, 1, data.TotalFiles)
}

func TestTokenLimitsChecker_EmptyDir(t *testing.T) {
	dir := t.TempDir()

	checker := &TokenLimitsChecker{}
	result, err := checker.Check(skillInDir(dir))
	require.NoError(t, err)
	require.True(t, result.Passed)

	data, ok := result.Data.(*TokenLimitsData)
	require.True(t, ok)
	require.Equal(t, 0, data.TotalFiles)
}

func TestTokenLimitsChecker_ExcludedDirs(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "node_modules"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "node_modules", "dep.md"), []byte("# Dep"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "main.md"), []byte("# Main"), 0o644))

	checker := &TokenLimitsChecker{}
	result, err := checker.Check(skillInDir(dir))
	require.NoError(t, err)

	data, ok := result.Data.(*TokenLimitsData)
	require.True(t, ok)
	require.Equal(t, 1, data.TotalFiles) // only main.md, not the one in node_modules
}
