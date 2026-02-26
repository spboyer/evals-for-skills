package tokens

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Suggest workspace-aware tests (Issue #442)
// These mirror the check workspace tests for the suggest command.
// Tests marked [NEEDS-IMPL] require Linus's workspace-aware implementation.
// ---------------------------------------------------------------------------

// [NEEDS-IMPL] Skill name resolution should respect paths.skills from .waza.yaml.
func TestWorkspace_SuggestSkillWithConfiguredSkillsDir(t *testing.T) {
	dir := t.TempDir()

	// .waza.yaml with custom skills dir
	writeWazaYaml(t, dir, "paths:\n  skills: \"my-custom-skills\"\n")

	// Create a skill with content that triggers suggestions (excess emojis)
	skillDir := filepath.Join(dir, "my-custom-skills", "emoji-skill")
	require.NoError(t, os.MkdirAll(skillDir, 0o755))
	content := "---\nname: emoji-skill\ndescription: \"emoji heavy\"\n---\n# 🎉🎊🎈🎆🎇🌟💫✨🌈\nLots of emojis.\n"
	require.NoError(t, os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644))

	t.Chdir(dir)

	out := new(bytes.Buffer)
	cmd := newSuggestCmd()
	cmd.SetOut(out)
	cmd.SetArgs([]string{"emoji-skill"})
	err := cmd.Execute()
	require.NoError(t, err, "suggest should work with custom skills dir from .waza.yaml")

	result := out.String()
	// The emoji-heavy content should trigger the emojis suggestion
	assert.Contains(t, result, "emojis", "should detect excessive emojis in custom skills dir")
}

// Skill name resolution should work with default skills/ dir (no config).
func TestWorkspace_SuggestSkillWithDefaultSkillsDir(t *testing.T) {
	dir := t.TempDir()

	skillDir := filepath.Join(dir, "skills", "clean-skill")
	require.NoError(t, os.MkdirAll(skillDir, 0o755))
	content := "---\nname: clean-skill\ndescription: \"A clean skill\"\n---\n# Clean Content\nNo issues here.\n"
	require.NoError(t, os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644))

	t.Chdir(dir)

	out := new(bytes.Buffer)
	cmd := newSuggestCmd()
	cmd.SetOut(out)
	cmd.SetArgs([]string{"clean-skill"})
	require.NoError(t, cmd.Execute())

	result := out.String()
	assert.Contains(t, result, "No optimization suggestions found", "clean skill should have no suggestions")
}

// [NEEDS-IMPL] No args + multi-skill workspace should analyze all skills.
func TestWorkspace_SuggestNoArgs_MultiSkill(t *testing.T) {
	dir := t.TempDir()

	writeSkill(t, filepath.Join(dir, "skills", "skill-alpha"), "skill-alpha")
	writeSkill(t, filepath.Join(dir, "skills", "skill-beta"), "skill-beta")

	t.Chdir(dir)

	out := new(bytes.Buffer)
	cmd := newSuggestCmd()
	cmd.SetOut(out)
	// No args — should detect multi-skill workspace and suggest for all
	require.NoError(t, cmd.Execute())

	result := out.String()
	// Both skills should be referenced in output
	assert.Contains(t, result, "skill-alpha", "output should include skill-alpha")
	assert.Contains(t, result, "skill-beta", "output should include skill-beta")
}

// [NEEDS-IMPL] No args + multi-skill with custom skills dir from .waza.yaml.
func TestWorkspace_SuggestNoArgs_MultiSkillCustomDir(t *testing.T) {
	dir := t.TempDir()

	writeWazaYaml(t, dir, "paths:\n  skills: \"extensions\"\n")
	writeSkill(t, filepath.Join(dir, "extensions", "ext-a"), "ext-a")
	writeSkill(t, filepath.Join(dir, "extensions", "ext-b"), "ext-b")

	t.Chdir(dir)

	out := new(bytes.Buffer)
	cmd := newSuggestCmd()
	cmd.SetOut(out)
	require.NoError(t, cmd.Execute())

	result := out.String()
	assert.Contains(t, result, "ext-a", "output should reference ext-a")
	assert.Contains(t, result, "ext-b", "output should reference ext-b")
}

// [NEEDS-IMPL] Malformed .waza.yaml should not crash suggest.
func TestWorkspace_SuggestSkillWithMalformedConfig(t *testing.T) {
	dir := t.TempDir()

	writeWazaYaml(t, dir, "{{{{invalid yaml")
	writeSkill(t, filepath.Join(dir, "skills", "safe-skill"), "safe-skill")

	t.Chdir(dir)

	out := new(bytes.Buffer)
	cmd := newSuggestCmd()
	cmd.SetOut(out)
	cmd.SetArgs([]string{"safe-skill"})
	err := cmd.Execute()
	require.NoError(t, err, "malformed .waza.yaml should not cause suggest to fail")
}

// No args + no workspace → falls back to CWD scan (existing behavior).
func TestWorkspace_SuggestNoArgs_NoWorkspace(t *testing.T) {
	dir := t.TempDir()

	// Just a plain markdown file, no skill structure
	require.NoError(t, os.WriteFile(filepath.Join(dir, "notes.md"), []byte("# Notes\nSome notes.\n"), 0o644))

	t.Chdir(dir)

	out := new(bytes.Buffer)
	cmd := newSuggestCmd()
	cmd.SetOut(out)
	require.NoError(t, cmd.Execute())

	// Should complete without error — CWD scan is the fallback
	result := out.String()
	assert.NotEmpty(t, result, "should produce some output even without workspace")
}
