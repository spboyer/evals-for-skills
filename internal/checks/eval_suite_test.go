package checks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spboyer/waza/internal/skill"
	"github.com/stretchr/testify/require"
)

var _ ComplianceChecker = (*EvalSuiteChecker)(nil)

func TestEvalSuiteChecker_Found(t *testing.T) {
	dir := t.TempDir()

	// Create eval.yaml
	require.NoError(t, os.WriteFile(filepath.Join(dir, "eval.yaml"), []byte("name: test\n"), 0o644))

	checker := &EvalSuiteChecker{}
	result, err := checker.Check(skill.Skill{Path: filepath.Join(dir, "SKILL.md")})
	require.NoError(t, err)
	require.True(t, result.Passed) // always true (recommendation)
	require.Equal(t, "eval-suite", result.Name)

	data, ok := result.Data.(*EvalSuiteData)
	require.True(t, ok)
	require.True(t, data.Found)
	require.Contains(t, result.Summary, "Found")
}

func TestEvalSuiteChecker_NotFound(t *testing.T) {
	dir := t.TempDir()

	checker := &EvalSuiteChecker{}
	result, err := checker.Check(skill.Skill{Path: filepath.Join(dir, "SKILL.md")})
	require.NoError(t, err)
	require.True(t, result.Passed) // always true (recommendation)

	data, ok := result.Data.(*EvalSuiteData)
	require.True(t, ok)
	require.False(t, data.Found)
	require.Contains(t, result.Summary, "Not Found")
}

func TestEvalSuiteChecker_WithSkillName(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "eval.yaml"), []byte("name: test\n"), 0o644))

	sk := skill.Skill{
		Path:        filepath.Join(dir, "SKILL.md"),
		Frontmatter: skill.Frontmatter{Name: "my-skill"},
	}

	checker := &EvalSuiteChecker{}
	result, err := checker.Check(sk)
	require.NoError(t, err)

	data, ok := result.Data.(*EvalSuiteData)
	require.True(t, ok)
	require.True(t, data.Found)
}
