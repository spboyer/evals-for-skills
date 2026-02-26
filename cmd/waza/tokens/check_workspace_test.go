package tokens

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheck_SkillFlag(t *testing.T) {
	dir := t.TempDir()
	skillsDir := filepath.Join(dir, "skills")

	// Create two skills with different content
	for _, name := range []string{"skill-x", "skill-y"} {
		skillDir := filepath.Join(skillsDir, name)
		require.NoError(t, os.MkdirAll(skillDir, 0o755))
		content := "---\nname: " + name + "\ndescription: \"desc for " + name + "\"\n---\n# Body of " + name + "\n"
		require.NoError(t, os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644))
	}

	t.Chdir(dir)

	out := new(bytes.Buffer)
	cmd := newCheckCmd()
	cmd.SetOut(out)
	cmd.SetArgs([]string{"skill-x"})
	require.NoError(t, cmd.Execute())

	result := out.String()
	assert.Contains(t, result, "SKILL.md")
	// Should only see one file (the skill-x SKILL.md)
	assert.Contains(t, result, "1/1 files within limits")
}

func TestCheck_SkillFlagNotFound(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	cmd := newCheckCmd()
	cmd.SetOut(new(bytes.Buffer))
	cmd.SetArgs([]string{"nonexistent"})
	err := cmd.Execute()
	require.Error(t, err)
}

// ---------------------------------------------------------------------------
// Workspace-aware tests (Issue #442)
// These test the configDetectOptions behavior and no-args workspace detection.
// Tests that depend on Linus's implementation are marked with [NEEDS-IMPL].
// ---------------------------------------------------------------------------

// writeSkill is a test helper that creates a SKILL.md file with frontmatter.
func writeSkill(t *testing.T, dir, name string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(dir, 0o755))
	content := "---\nname: " + name + "\ndescription: \"desc for " + name + "\"\n---\n# Body of " + name + "\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(content), 0o644))
}

// writeWazaYaml creates a .waza.yaml file in the given directory.
func writeWazaYaml(t *testing.T, dir, content string) {
	t.Helper()
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".waza.yaml"), []byte(content), 0o644))
}

// --- configDetectOptions behavior: .waza.yaml with custom paths.skills ---

// [NEEDS-IMPL] Skill name resolution should respect paths.skills from .waza.yaml.
func TestWorkspace_CheckSkillWithConfiguredSkillsDir(t *testing.T) {
	dir := t.TempDir()

	// .waza.yaml points skills to a non-default directory
	writeWazaYaml(t, dir, "paths:\n  skills: \"my-custom-skills\"\n")

	// Create skill in the custom directory
	writeSkill(t, filepath.Join(dir, "my-custom-skills", "alpha"), "alpha")

	// Also create a skill in the default skills/ dir (should NOT be found
	// if config is respected)
	writeSkill(t, filepath.Join(dir, "skills", "decoy"), "decoy")

	t.Chdir(dir)

	out := new(bytes.Buffer)
	cmd := newCheckCmd()
	cmd.SetOut(out)
	cmd.SetArgs([]string{"alpha"})
	err := cmd.Execute()
	require.NoError(t, err, "check alpha should succeed with custom skills dir from .waza.yaml")

	result := out.String()
	assert.Contains(t, result, "SKILL.md")
	assert.Contains(t, result, "1/1 files within limits")
}

// --- configDetectOptions behavior: no .waza.yaml → defaults ---

func TestWorkspace_CheckSkillWithDefaultSkillsDir(t *testing.T) {
	dir := t.TempDir()

	// No .waza.yaml — workspace detection should use default "skills/" dir
	writeSkill(t, filepath.Join(dir, "skills", "beta"), "beta")

	t.Chdir(dir)

	out := new(bytes.Buffer)
	cmd := newCheckCmd()
	cmd.SetOut(out)
	cmd.SetArgs([]string{"beta"})
	require.NoError(t, cmd.Execute())

	result := out.String()
	assert.Contains(t, result, "SKILL.md")
	assert.Contains(t, result, "1/1 files within limits")
}

// --- configDetectOptions behavior: malformed .waza.yaml → warn, don't crash ---

// [NEEDS-IMPL] Malformed .waza.yaml should not prevent check from running.
// The command should fall back to defaults when config parsing fails.
func TestWorkspace_CheckSkillWithMalformedConfig(t *testing.T) {
	dir := t.TempDir()

	// Write invalid YAML that will fail to parse
	writeWazaYaml(t, dir, "{{{{invalid yaml content that cannot be parsed")

	// Create skill in the default skills/ dir (fallback)
	writeSkill(t, filepath.Join(dir, "skills", "gamma"), "gamma")

	t.Chdir(dir)

	out := new(bytes.Buffer)
	cmd := newCheckCmd()
	cmd.SetOut(out)
	cmd.SetArgs([]string{"gamma"})
	// Should not crash — should gracefully fall back to default skills dir
	err := cmd.Execute()
	require.NoError(t, err, "malformed .waza.yaml should not cause check to fail; should fall back to defaults")

	result := out.String()
	assert.Contains(t, result, "SKILL.md")
}

// --- No-args workspace detection ---

// [NEEDS-IMPL] No args + multi-skill workspace should check all skills with per-skill output.
func TestWorkspace_CheckNoArgs_MultiSkill(t *testing.T) {
	dir := t.TempDir()

	// Create a multi-skill workspace with default "skills/" dir
	writeSkill(t, filepath.Join(dir, "skills", "skill-a"), "skill-a")
	writeSkill(t, filepath.Join(dir, "skills", "skill-b"), "skill-b")

	t.Chdir(dir)

	out := new(bytes.Buffer)
	cmd := newCheckCmd()
	cmd.SetOut(out)
	// No args — should detect multi-skill workspace and check all skills
	require.NoError(t, cmd.Execute())

	result := out.String()
	// Both skills should appear in output
	assert.Contains(t, result, "skill-a", "output should mention skill-a")
	assert.Contains(t, result, "skill-b", "output should mention skill-b")
}

// [NEEDS-IMPL] No args + multi-skill workspace with custom skills dir.
func TestWorkspace_CheckNoArgs_MultiSkillCustomDir(t *testing.T) {
	dir := t.TempDir()

	writeWazaYaml(t, dir, "paths:\n  skills: \"custom-skills\"\n")
	writeSkill(t, filepath.Join(dir, "custom-skills", "one"), "one")
	writeSkill(t, filepath.Join(dir, "custom-skills", "two"), "two")

	t.Chdir(dir)

	out := new(bytes.Buffer)
	cmd := newCheckCmd()
	cmd.SetOut(out)
	// No args — should use .waza.yaml paths.skills and find both skills
	require.NoError(t, cmd.Execute())

	result := out.String()
	assert.Contains(t, result, "one", "output should reference skill 'one'")
	assert.Contains(t, result, "two", "output should reference skill 'two'")
}

// [NEEDS-IMPL] No args + single-skill workspace should check that skill's directory.
func TestWorkspace_CheckNoArgs_SingleSkill(t *testing.T) {
	dir := t.TempDir()

	// Single skill in CWD — SKILL.md directly in dir
	writeSkill(t, dir, "solo-skill")

	t.Chdir(dir)

	out := new(bytes.Buffer)
	cmd := newCheckCmd()
	cmd.SetOut(out)
	require.NoError(t, cmd.Execute())

	result := out.String()
	assert.Contains(t, result, "SKILL.md", "single-skill workspace should check SKILL.md")
}

// No args + no workspace → falls back to CWD scan (existing behavior).
func TestWorkspace_CheckNoArgs_NoWorkspace(t *testing.T) {
	dir := t.TempDir()

	// No SKILL.md, no skills/ dir — just a markdown file
	require.NoError(t, os.WriteFile(filepath.Join(dir, "README.md"), []byte("# Hello\n"), 0o644))

	t.Chdir(dir)

	out := new(bytes.Buffer)
	cmd := newCheckCmd()
	cmd.SetOut(out)
	require.NoError(t, cmd.Execute())

	result := out.String()
	assert.Contains(t, result, "README.md", "no-workspace fallback should scan CWD for markdown files")
}

// --- Edge cases ---

// [NEEDS-IMPL] Custom paths.skills pointing to a non-existent directory.
func TestWorkspace_CheckCustomSkillsDirNonExistent(t *testing.T) {
	dir := t.TempDir()

	writeWazaYaml(t, dir, "paths:\n  skills: \"nonexistent-dir\"\n")

	t.Chdir(dir)

	out := new(bytes.Buffer)
	cmd := newCheckCmd()
	cmd.SetOut(out)
	cmd.SetArgs([]string{"phantom"})
	err := cmd.Execute()
	// Should produce an error — the skill can't be found
	require.Error(t, err, "skill resolution should fail when skills dir does not exist")
}

// [NEEDS-IMPL] Empty skills directory (no SKILL.md files found).
func TestWorkspace_CheckEmptySkillsDir(t *testing.T) {
	dir := t.TempDir()

	// Create skills dir but don't put any SKILL.md files in it
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "skills", "empty-skill"), 0o755))
	// empty-skill/ dir exists but has no SKILL.md

	t.Chdir(dir)

	out := new(bytes.Buffer)
	cmd := newCheckCmd()
	cmd.SetOut(out)
	cmd.SetArgs([]string{"empty-skill"})
	err := cmd.Execute()
	// Should fail — no SKILL.md means no skill to find
	require.Error(t, err, "should error when skill dir exists but has no SKILL.md")
}

// ---------------------------------------------------------------------------
// Workspace-root-relative pattern tests (Issue #444)
// ---------------------------------------------------------------------------

// TestWorkspace_CheckPrefixedPatterns verifies that workspace-root-relative
// patterns in .waza.yaml (e.g. "custom-skills/**/SKILL.md") resolve
// correctly when tokens check runs against a specific skill.
func TestWorkspace_CheckPrefixedPatterns(t *testing.T) {
	dir := t.TempDir()

	writeWazaYaml(t, dir, `paths:
  skills: custom-skills/
tokens:
  limits:
    defaults:
      "custom-skills/**/SKILL.md": 1000
      "custom-skills/**/references/*.md": 800
`)

	skillDir := filepath.Join(dir, "custom-skills", "my-skill")
	writeSkill(t, skillDir, "my-skill")

	refsDir := filepath.Join(skillDir, "references")
	require.NoError(t, os.MkdirAll(refsDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(refsDir, "doc.md"), []byte("# Doc\n"), 0o644))

	t.Chdir(dir)

	out := new(bytes.Buffer)
	cmd := newCheckCmd()
	cmd.SetOut(out)
	cmd.SetArgs([]string{"my-skill", "--format", "json"})
	require.NoError(t, cmd.Execute())

	var report checkReport
	require.NoError(t, json.Unmarshal(out.Bytes(), &report))
	require.Equal(t, 2, report.TotalFiles, "should find SKILL.md and references/doc.md")

	for _, r := range report.Results {
		switch {
		case strings.HasSuffix(r.File, "SKILL.md"):
			assert.Equal(t, 1000, r.Limit, "SKILL.md should match prefixed pattern")
			assert.Equal(t, "custom-skills/**/SKILL.md", r.Pattern)
		case strings.HasSuffix(r.File, "doc.md"):
			assert.Equal(t, 800, r.Limit, "doc.md should match prefixed references pattern")
			assert.Equal(t, "custom-skills/**/references/*.md", r.Pattern)
		default:
			t.Errorf("unexpected file in results: %s", r.File)
		}
	}
}

// TestWorkspace_CheckPrefixedPatterns_BatchMode verifies workspace-root-relative
// patterns work when no args are given (multi-skill batch mode).
func TestWorkspace_CheckPrefixedPatterns_BatchMode(t *testing.T) {
	dir := t.TempDir()

	writeWazaYaml(t, dir, `paths:
  skills: custom-skills/
tokens:
  limits:
    defaults:
      "custom-skills/**/SKILL.md": 750
`)

	writeSkill(t, filepath.Join(dir, "custom-skills", "sk-a"), "sk-a")
	writeSkill(t, filepath.Join(dir, "custom-skills", "sk-b"), "sk-b")

	t.Chdir(dir)

	out := new(bytes.Buffer)
	cmd := newCheckCmd()
	cmd.SetOut(out)
	// No args — batch mode
	require.NoError(t, cmd.Execute())

	result := out.String()
	assert.Contains(t, result, "sk-a")
	assert.Contains(t, result, "sk-b")
}
