package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spboyer/waza/internal/workspace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLooksLikePath(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"eval.yaml", true},
		{"./skills/foo", true},
		{"skills/foo", true},
		{".", true},
		{"code-explainer", false},
		{"my-skill", false},
		{"foo.yaml", true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, workspace.LooksLikePath(tt.input))
		})
	}
}

func TestResolveWorkspace_ExplicitPath(t *testing.T) {
	ctx, err := resolveWorkspace([]string{"eval.yaml"})
	assert.NoError(t, err)
	assert.NotNil(t, ctx)

	ctx, err = resolveWorkspace([]string{"./my-skill"})
	assert.NoError(t, err)
	assert.NotNil(t, ctx)

	ctx, err = resolveWorkspace([]string{"."})
	assert.NoError(t, err)
	assert.NotNil(t, ctx)
}

func TestResolveWorkspace_SingleSkillWorkspace(t *testing.T) {
	dir := t.TempDir()
	skillContent := `---
name: test-skill
description: "A test skill."
---
# Test
`
	require.NoError(t, os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(skillContent), 0o644))
	t.Chdir(dir)

	ctx, err := resolveWorkspace(nil)
	require.NoError(t, err)
	require.Len(t, ctx.Skills, 1)
	assert.Equal(t, "test-skill", ctx.Skills[0].Name)
}

func TestResolveWorkspace_MultiSkillWorkspace(t *testing.T) {
	dir := t.TempDir()
	skillsDir := filepath.Join(dir, "skills")

	for _, name := range []string{"skill-a", "skill-b"} {
		skillDir := filepath.Join(skillsDir, name)
		require.NoError(t, os.MkdirAll(skillDir, 0o755))
		content := "---\nname: " + name + "\ndescription: \"desc\"\n---\n# Body\n"
		require.NoError(t, os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644))
	}
	t.Chdir(dir)

	ctx, err := resolveWorkspace(nil)
	require.NoError(t, err)
	require.Len(t, ctx.Skills, 2)

	names := []string{ctx.Skills[0].Name, ctx.Skills[1].Name}
	assert.Contains(t, names, "skill-a")
	assert.Contains(t, names, "skill-b")
}

func TestResolveWorkspace_SkillName(t *testing.T) {
	dir := t.TempDir()
	skillsDir := filepath.Join(dir, "skills")

	for _, name := range []string{"alpha", "beta"} {
		skillDir := filepath.Join(skillsDir, name)
		require.NoError(t, os.MkdirAll(skillDir, 0o755))
		content := "---\nname: " + name + "\ndescription: \"desc\"\n---\n# Body\n"
		require.NoError(t, os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644))
	}
	t.Chdir(dir)

	ctx, err := resolveWorkspace([]string{"alpha"})
	require.NoError(t, err)
	require.Len(t, ctx.Skills, 1)
	assert.Equal(t, "alpha", ctx.Skills[0].Name)
}

func TestResolveWorkspace_NoWorkspace(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	_, err := resolveWorkspace(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no skills detected")
}

func TestResolveEvalPath_WithColocatedEval(t *testing.T) {
	dir := t.TempDir()
	skillContent := "---\nname: test-skill\ndescription: \"desc\"\n---\n# Body\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(skillContent), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "eval.yaml"), []byte("name: test\n"), 0o644))
	t.Chdir(dir)

	si := workspace.SkillInfo{Name: "test-skill", Dir: dir, SkillPath: filepath.Join(dir, "SKILL.md")}
	evalPath, err := resolveEvalPath(&si)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(dir, "eval.yaml"), evalPath)
}
