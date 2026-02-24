package dev

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spboyer/waza/internal/scoring"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// makeSkillDir creates a skill directory with SKILL.md and returns the path.
func makeSkillDir(t *testing.T, parent, name, content string) string {
	t.Helper()
	dir := filepath.Join(parent, name)
	require.NoError(t, os.MkdirAll(dir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(content), 0o644))
	return dir
}

func TestBatch_MultipleSkillsByName(t *testing.T) {
	dir := t.TempDir()
	skillsDir := filepath.Join(dir, "skills")

	makeSkillDir(t, skillsDir, "alpha", `---
name: alpha
description: "Short"
---
# Alpha
`)
	makeSkillDir(t, skillsDir, "beta", `---
name: beta
description: "Short"
---
# Beta
`)

	t.Chdir(dir)

	out := new(bytes.Buffer)
	cmd := NewCommand()
	cmd.SetOut(out)
	cmd.SetErr(new(bytes.Buffer))
	cmd.SetIn(strings.NewReader(""))
	cmd.SetArgs([]string{"alpha", "beta", "--auto", "--max-iterations", "1"})

	err := cmd.Execute()
	require.NoError(t, err)

	output := out.String()
	assert.Contains(t, output, "Batch processing 2 skill(s)")
	assert.Contains(t, output, "[1/2] alpha")
	assert.Contains(t, output, "[2/2] beta")
	assert.Contains(t, output, "BATCH SUMMARY")
	assert.Contains(t, output, "alpha")
	assert.Contains(t, output, "beta")
}

func TestBatch_AllFlag(t *testing.T) {
	dir := t.TempDir()
	skillsDir := filepath.Join(dir, "skills")

	makeSkillDir(t, skillsDir, "one", `---
name: one
description: "Short"
---
# One
`)
	makeSkillDir(t, skillsDir, "two", `---
name: two
description: "Short"
---
# Two
`)

	t.Chdir(dir)

	out := new(bytes.Buffer)
	cmd := NewCommand()
	cmd.SetOut(out)
	cmd.SetErr(new(bytes.Buffer))
	cmd.SetIn(strings.NewReader(""))
	cmd.SetArgs([]string{"--all", "--auto", "--max-iterations", "1"})

	err := cmd.Execute()
	require.NoError(t, err)

	output := out.String()
	assert.Contains(t, output, "Batch processing 2 skill(s)")
	assert.Contains(t, output, "BATCH SUMMARY")
}

func TestBatch_FilterByLevel(t *testing.T) {
	dir := t.TempDir()
	skillsDir := filepath.Join(dir, "skills")

	// Low adherence skill
	makeSkillDir(t, skillsDir, "low-one", `---
name: low-one
description: "Short"
---
# Low
`)

	// High adherence skill
	makeSkillDir(t, skillsDir, "high-one", `---
name: high-one
description: |
  **WORKFLOW SKILL** - Handle complex workflow operations involving multiple steps.
  USE FOR: "run workflow", "execute pipeline", "orchestrate tasks", "automate process".
  DO NOT USE FOR: simple one-off commands (use cli-runner), debugging (use debugger).
  INVOKES: task-runner for execution, config-parser for setup.
  FOR SINGLE OPERATIONS: Use task-runner directly for individual tasks.
---
# High
`)

	t.Chdir(dir)

	out := new(bytes.Buffer)
	cmd := NewCommand()
	cmd.SetOut(out)
	cmd.SetErr(new(bytes.Buffer))
	cmd.SetIn(strings.NewReader(""))
	cmd.SetArgs([]string{"--all", "--filter", "low", "--auto", "--max-iterations", "1"})

	err := cmd.Execute()
	require.NoError(t, err)

	output := out.String()
	// Should only process the low skill, not the high one
	assert.Contains(t, output, "Batch processing 1 skill(s)")
	assert.Contains(t, output, "low-one")
	assert.NotContains(t, output, "[2/2]") // only 1 skill
}

func TestBatch_FilterNoMatch(t *testing.T) {
	dir := t.TempDir()
	skillsDir := filepath.Join(dir, "skills")

	makeSkillDir(t, skillsDir, "high-only", `---
name: high-only
description: |
  **WORKFLOW SKILL** - Handle complex workflow operations involving multiple steps.
  USE FOR: "run workflow", "execute pipeline", "orchestrate tasks", "automate process".
  DO NOT USE FOR: simple one-off commands (use cli-runner), debugging (use debugger).
  INVOKES: task-runner for execution, config-parser for setup.
  FOR SINGLE OPERATIONS: Use task-runner directly for individual tasks.
---
# High
`)

	t.Chdir(dir)

	out := new(bytes.Buffer)
	cmd := NewCommand()
	cmd.SetOut(out)
	cmd.SetErr(new(bytes.Buffer))
	cmd.SetIn(strings.NewReader(""))
	cmd.SetArgs([]string{"--all", "--filter", "low", "--auto", "--max-iterations", "1"})

	err := cmd.Execute()
	require.NoError(t, err)

	output := out.String()
	assert.Contains(t, output, "No skills at Low adherence level.")
}

func TestBatch_FilterRequiresAll(t *testing.T) {
	cmd := NewCommand()
	cmd.SetOut(new(bytes.Buffer))
	cmd.SetErr(new(bytes.Buffer))
	cmd.SetArgs([]string{"some-skill", "--filter", "low"})

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--filter requires --all")
}

func TestBatch_NoWorkspace(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	cmd := NewCommand()
	cmd.SetOut(new(bytes.Buffer))
	cmd.SetErr(new(bytes.Buffer))
	cmd.SetArgs([]string{"--all"})

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no skills detected")
}

func TestBatch_SummaryShowsImprovement(t *testing.T) {
	dir := t.TempDir()
	skillsDir := filepath.Join(dir, "skills")

	makeSkillDir(t, skillsDir, "improve-me", `---
name: improve-me
description: "Short"
---
# Improve Me
## Data Processing
## File Management
`)

	t.Chdir(dir)

	out := new(bytes.Buffer)
	cmd := NewCommand()
	cmd.SetOut(out)
	cmd.SetErr(new(bytes.Buffer))
	cmd.SetIn(strings.NewReader(""))
	cmd.SetArgs([]string{"--all", "--auto", "--max-iterations", "2", "--target", "medium"})

	err := cmd.Execute()
	require.NoError(t, err)

	output := out.String()
	assert.Contains(t, output, "BATCH SUMMARY")
	assert.Contains(t, output, "improve-me")
	// The skill should have been improved (📈 status)
	assert.Contains(t, output, "📈")
}

func TestDisplayBatchSummary(t *testing.T) {
	var buf bytes.Buffer
	results := []batchSkillResult{
		{Name: "skill-a", BeforeLevel: scoring.AdherenceLow, AfterLevel: scoring.AdherenceMedium, BeforeTokens: 10, AfterTokens: 50},
		{Name: "skill-b", BeforeLevel: scoring.AdherenceHigh, AfterLevel: scoring.AdherenceHigh, BeforeTokens: 100, AfterTokens: 100},
		{Name: "skill-c", BeforeLevel: scoring.AdherenceLow, AfterLevel: scoring.AdherenceLow, Err: assert.AnError},
	}

	DisplayBatchSummary(&buf, results)

	output := buf.String()
	assert.Contains(t, output, "BATCH SUMMARY")
	assert.Contains(t, output, "skill-a")
	assert.Contains(t, output, "📈") // skill-a improved
	assert.Contains(t, output, "✅") // skill-b unchanged
	assert.Contains(t, output, "❌") // skill-c errored
}
