package dev

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunScaffoldTriggers_WritesFile(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "my-skill")
	require.NoError(t, os.MkdirAll(skillDir, 0o755))

	skillMD := `---
name: my-skill
description: |
  A test skill.
  USE FOR: "do task A", "do task B".
  DO NOT USE FOR: unrelated work (use other-skill).
---

# My Skill
`
	require.NoError(t, os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(skillMD), 0o644))

	cmd := NewCommand()
	cmd.SetArgs([]string{skillDir, "--scaffold-triggers"})
	buf := &captureWriter{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	require.NoError(t, err)

	// Verify file was created
	outPath := filepath.Join(skillDir, "tests", "trigger_tests.yaml")
	data, err := os.ReadFile(outPath)
	require.NoError(t, err)

	content := string(data)
	assert.Contains(t, content, "skill: my-skill")
	assert.Contains(t, content, "should_trigger_prompts:")
	assert.Contains(t, content, `"do task A"`)
	assert.Contains(t, content, `"do task B"`)
	assert.Contains(t, content, "should_not_trigger_prompts:")
	assert.Contains(t, content, `"unrelated work"`)

	// Verify output message
	assert.Contains(t, buf.String(), "Scaffolded trigger tests")
	assert.Contains(t, buf.String(), "2 should-trigger")
	assert.Contains(t, buf.String(), "1 should-not-trigger")
}

func TestRunScaffoldTriggers_NoDescription(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "empty-skill")
	require.NoError(t, os.MkdirAll(skillDir, 0o755))

	skillMD := `---
name: empty-skill
---

# Empty Skill
`
	require.NoError(t, os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(skillMD), 0o644))

	cmd := NewCommand()
	cmd.SetArgs([]string{skillDir, "--scaffold-triggers"})
	buf := &captureWriter{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no description")
}

func TestRunScaffoldTriggers_NoPhrases(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "no-phrases")
	require.NoError(t, os.MkdirAll(skillDir, 0o755))

	skillMD := `---
name: no-phrases
description: "A simple description without trigger phrases."
---

# No Phrases
`
	require.NoError(t, os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(skillMD), 0o644))

	cmd := NewCommand()
	cmd.SetArgs([]string{skillDir, "--scaffold-triggers"})
	buf := &captureWriter{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no USE FOR or DO NOT USE FOR")
}

func TestRunScaffoldTriggers_NoSkillMD(t *testing.T) {
	dir := t.TempDir()

	cmd := NewCommand()
	cmd.SetArgs([]string{dir, "--scaffold-triggers"})
	buf := &captureWriter{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reading SKILL.md")
}

func TestRunScaffoldTriggers_RequiresArg(t *testing.T) {
	cmd := NewCommand()
	cmd.SetArgs([]string{"--scaffold-triggers"})
	buf := &captureWriter{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "skill name or path required")
}

// captureWriter collects output for test assertions.
type captureWriter struct {
	data []byte
}

func (w *captureWriter) Write(p []byte) (int, error) {
	w.data = append(w.data, p...)
	return len(p), nil
}

func (w *captureWriter) String() string {
	return string(w.data)
}
