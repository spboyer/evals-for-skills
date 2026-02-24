package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spboyer/waza/cmd/waza/dev"
	"github.com/spboyer/waza/internal/scaffold"
	"github.com/spboyer/waza/internal/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckCommand(t *testing.T) {
	// Create a test skill with proper YAML frontmatter
	tmpDir := t.TempDir()
	skillContent := `---
name: test-skill
description: This is a test skill for unit testing the check command functionality.
---

# Test Skill

This is the body of the test skill.
`
	skillPath := filepath.Join(tmpDir, "SKILL.md")
	err := os.WriteFile(skillPath, []byte(skillContent), 0644)
	require.NoError(t, err)

	// Run check command
	cmd := newCheckCommand()
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetErr(&output)
	cmd.SetArgs([]string{tmpDir})

	err = cmd.Execute()
	assert.NoError(t, err)

	result := output.String()

	// Verify output contains expected sections
	assert.Contains(t, result, "Skill Readiness Check")
	assert.Contains(t, result, "test-skill")
	assert.Contains(t, result, "Compliance Score:")
	assert.Contains(t, result, "Token Budget:")
	assert.Contains(t, result, "Evaluation Suite:")
	assert.Contains(t, result, "Overall Readiness")
	assert.Contains(t, result, "Next Steps")
}

func TestCheckCommandNoSkillMd(t *testing.T) {
	tmpDir := t.TempDir()

	cmd := newCheckCommand()
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetErr(&output)
	cmd.SetArgs([]string{tmpDir})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no SKILL.md found")
}

func TestCheckCommandWithEval(t *testing.T) {
	tmpDir := t.TempDir()

	// Create SKILL.md
	skillContent := `---
name: test-skill-with-eval
description: This is a test skill with evaluation suite for unit testing readiness check.
---

# Test Skill

Body content.
`
	skillPath := filepath.Join(tmpDir, "SKILL.md")
	err := os.WriteFile(skillPath, []byte(skillContent), 0644)
	require.NoError(t, err)

	// Create eval.yaml
	evalContent := `name: test-eval
tasks: []
`
	evalPath := filepath.Join(tmpDir, "eval.yaml")
	err = os.WriteFile(evalPath, []byte(evalContent), 0644)
	require.NoError(t, err)

	cmd := newCheckCommand()
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetErr(&output)
	cmd.SetArgs([]string{tmpDir})

	err = cmd.Execute()
	assert.NoError(t, err)

	result := output.String()
	assert.Contains(t, result, "Evaluation Suite: Found")
	assert.Contains(t, result, "eval.yaml detected")
}

func TestCheckCommandHighCompliance(t *testing.T) {
	tmpDir := t.TempDir()

	// Create SKILL.md with high compliance (includes triggers, anti-triggers, routing clarity)
	skillContent := `---
name: high-compliance-skill
description: "**WORKFLOW SKILL** - This is a comprehensive test skill that demonstrates high compliance with all requirements including sufficient length and proper formatting. USE FOR: testing, validation, compliance checks, readiness assessment, quality verification. DO NOT USE FOR: unrelated tasks, production use, deployment automation. INVOKES: internal validators, scoring engines. FOR SINGLE OPERATIONS: use direct commands for simple checks."
---

# High Compliance Skill

This skill has high compliance.
`
	skillPath := filepath.Join(tmpDir, "SKILL.md")
	err := os.WriteFile(skillPath, []byte(skillContent), 0644)
	require.NoError(t, err)

	cmd := newCheckCommand()
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetErr(&output)
	cmd.SetArgs([]string{tmpDir})

	err = cmd.Execute()
	assert.NoError(t, err)

	result := output.String()
	assert.Contains(t, result, "Compliance Score: High")
	assert.Contains(t, result, "Excellent!")
	// Should not suggest compliance improvements since it's already High
	assert.NotContains(t, strings.ToLower(result), "expand your description")
	assert.NotContains(t, strings.ToLower(result), "add a 'use for:' section")
}

func TestCheckCommandLowCompliance(t *testing.T) {
	tmpDir := t.TempDir()

	// Create SKILL.md with low compliance (short description, no triggers)
	skillContent := `---
name: low-compliance-skill
description: Short description.
---

# Low Compliance Skill

This skill has low compliance.
`
	skillPath := filepath.Join(tmpDir, "SKILL.md")
	err := os.WriteFile(skillPath, []byte(skillContent), 0644)
	require.NoError(t, err)

	cmd := newCheckCommand()
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetErr(&output)
	cmd.SetArgs([]string{tmpDir})

	err = cmd.Execute()
	assert.NoError(t, err)

	result := output.String()
	assert.Contains(t, result, "Compliance Score: Low")
	assert.Contains(t, result, "Next Steps")
	// Should suggest compliance improvements
	assert.Contains(t, result, "Expand your description")
	assert.Contains(t, result, "USE FOR:")
	assert.Contains(t, result, "DO NOT USE FOR:")
	assert.Contains(t, result, "waza dev")
}

func TestCheckReadiness(t *testing.T) {
	tmpDir := t.TempDir()

	skillContent := `---
name: readiness-test
description: This is a test skill for checking the readiness report generation with proper frontmatter and adequate description length.
---

# Readiness Test

Body content here.
`
	skillPath := filepath.Join(tmpDir, "SKILL.md")
	err := os.WriteFile(skillPath, []byte(skillContent), 0644)
	require.NoError(t, err)

	report, err := checkReadiness(tmpDir)
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, "readiness-test", report.skillName)
	assert.NotNil(t, report.complianceScore)
	assert.Greater(t, report.tokenCount, 0)
	assert.Equal(t, 500, report.tokenLimit) // Default limit
	assert.False(t, report.hasEval)
}

func TestGenerateNextSteps(t *testing.T) {
	// Test with high compliance - should have no steps
	t.Run("high compliance no issues", func(t *testing.T) {
		report := &readinessReport{
			complianceLevel: dev.AdherenceHigh,
			tokenCount:      400,
			tokenLimit:      500,
			tokenExceeded:   false,
			hasEval:         true,
			complianceScore: &dev.ScoreResult{
				Level:             dev.AdherenceHigh,
				DescriptionLen:    200,
				HasTriggers:       true,
				HasAntiTriggers:   true,
				HasRoutingClarity: true,
			},
		}
		steps := generateNextSteps(report)
		assert.Empty(t, steps)
	})

	// Test with low compliance - should have multiple steps
	t.Run("low compliance", func(t *testing.T) {
		report := &readinessReport{
			complianceLevel: dev.AdherenceLow,
			tokenCount:      400,
			tokenLimit:      500,
			tokenExceeded:   false,
			hasEval:         false,
			complianceScore: &dev.ScoreResult{
				Level:             dev.AdherenceLow,
				DescriptionLen:    100,
				HasTriggers:       false,
				HasAntiTriggers:   false,
				HasRoutingClarity: false,
			},
		}
		steps := generateNextSteps(report)
		assert.NotEmpty(t, steps)
		assert.Contains(t, strings.Join(steps, " "), "Expand your description")
		assert.Contains(t, strings.Join(steps, " "), "USE FOR:")
	})

	// Test with token issues
	t.Run("token exceeded", func(t *testing.T) {
		report := &readinessReport{
			complianceLevel: dev.AdherenceHigh,
			tokenCount:      600,
			tokenLimit:      500,
			tokenExceeded:   true,
			hasEval:         true,
			complianceScore: &dev.ScoreResult{
				Level:             dev.AdherenceHigh,
				DescriptionLen:    200,
				HasTriggers:       true,
				HasAntiTriggers:   true,
				HasRoutingClarity: true,
			},
		}
		steps := generateNextSteps(report)
		assert.NotEmpty(t, steps)
		assert.Contains(t, strings.Join(steps, " "), "Reduce SKILL.md by 100 tokens")
	})
}

func TestCheckCommandWithValidEvalSchema(t *testing.T) {
	tmpDir := t.TempDir()

	// Create SKILL.md
	skillContent := `---
name: schema-test
description: A test skill for schema validation.
---

# Schema Test Skill

Body content.
`
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "SKILL.md"), []byte(skillContent), 0644))

	// Create valid eval.yaml
	evalContent := `name: test-eval
skill: schema-test
version: "1.0"
config:
  trials_per_task: 1
  timeout_seconds: 60
  executor: mock
  model: gpt-4o
metrics:
  - name: accuracy
    weight: 1.0
    threshold: 0.8
tasks:
  - "tasks/*.yaml"
`
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "eval.yaml"), []byte(evalContent), 0644))

	// Create valid task
	tasksDir := filepath.Join(tmpDir, "tasks")
	require.NoError(t, os.MkdirAll(tasksDir, 0755))
	taskContent := `id: task-1
name: Basic task
inputs:
  prompt: "Test prompt"
`
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "basic.yaml"), []byte(taskContent), 0644))

	cmd := newCheckCommand()
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetErr(&output)
	cmd.SetArgs([]string{tmpDir})

	err := cmd.Execute()
	assert.NoError(t, err)

	result := output.String()
	assert.Contains(t, result, "eval.yaml schema valid")
	assert.Contains(t, result, "1 task file(s) validated")
	assert.NotContains(t, result, "Eval Schema:")
}

func TestCheckCommandWithInvalidEvalSchema(t *testing.T) {
	tmpDir := t.TempDir()

	// Create SKILL.md
	skillContent := `---
name: schema-test
description: A test skill for schema validation.
---

# Schema Test Skill

Body content.
`
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "SKILL.md"), []byte(skillContent), 0644))

	// Create invalid eval.yaml (bad executor)
	evalContent := `name: test-eval
skill: schema-test
version: "1.0"
config:
  trials_per_task: 1
  timeout_seconds: 60
  executor: invalid-engine
  model: gpt-4o
metrics:
  - name: accuracy
    weight: 1.0
    threshold: 0.8
tasks:
  - "tasks/*.yaml"
`
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "eval.yaml"), []byte(evalContent), 0644))

	cmd := newCheckCommand()
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetErr(&output)
	cmd.SetArgs([]string{tmpDir})

	err := cmd.Execute()
	assert.NoError(t, err)

	result := output.String()
	assert.Contains(t, result, "Eval Schema:")
	assert.Contains(t, result, "executor")
	assert.Contains(t, result, "Fix")
	assert.Contains(t, result, "schema error")
}

func TestCheckCommandNoEvalSkipsSchemaValidation(t *testing.T) {
	tmpDir := t.TempDir()

	// Create SKILL.md with no eval.yaml
	skillContent := `---
name: no-eval-test
description: A test skill without eval.
---

# No Eval Skill

Body content.
`
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "SKILL.md"), []byte(skillContent), 0644))

	cmd := newCheckCommand()
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetErr(&output)
	cmd.SetArgs([]string{tmpDir})

	err := cmd.Execute()
	assert.NoError(t, err)

	result := output.String()
	assert.NotContains(t, result, "Schema Validation")
	assert.NotContains(t, result, "Eval Schema")
	assert.NotContains(t, result, "Task Schema")
}

func TestCheckCommand_ScaffoldedSkillPassesSchemaValidation(t *testing.T) {
	// Build a workspace that mimics waza init + waza new output:
	//   dir/
	//     skills/my-skill/SKILL.md
	//     evals/my-skill/eval.yaml
	//     evals/my-skill/tasks/*.yaml
	//     evals/my-skill/fixtures/sample.py
	dir := t.TempDir()

	skillName := "my-skill"
	skillDir := filepath.Join(dir, "skills", skillName)
	evalsDir := filepath.Join(dir, "evals", skillName)
	tasksDir := filepath.Join(evalsDir, "tasks")
	fixturesDir := filepath.Join(evalsDir, "fixtures")

	require.NoError(t, os.MkdirAll(skillDir, 0o755))
	require.NoError(t, os.MkdirAll(tasksDir, 0o755))
	require.NoError(t, os.MkdirAll(fixturesDir, 0o755))

	// SKILL.md with high-compliance frontmatter (150+ char description, triggers, anti-triggers)
	skillContent := `---
name: my-skill
description: "A comprehensive skill that helps developers analyze and explain code patterns across multiple languages with clarity. USE FOR: code analysis, code explanation, pattern detection. DO NOT USE FOR: code execution, deployment, infrastructure management."
---

# My Skill

This skill analyzes code and provides explanations.
`
	require.NoError(t, os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(skillContent), 0o644))

	// eval.yaml from scaffold
	evalContent := scaffold.EvalYAML(skillName, "mock", "gpt-4o")
	require.NoError(t, os.WriteFile(filepath.Join(evalsDir, "eval.yaml"), []byte(evalContent), 0o644))

	// Task files from scaffold
	for name, content := range scaffold.TaskFiles(skillName) {
		require.NoError(t, os.WriteFile(filepath.Join(tasksDir, name), []byte(content), 0o644))
	}

	// Fixture file from scaffold
	require.NoError(t, os.WriteFile(filepath.Join(fixturesDir, "sample.py"), []byte(scaffold.Fixture()), 0o644))

	// cd into workspace root so workspace detection finds the eval
	t.Chdir(dir)

	// Directly call checkReadiness on the skill directory
	report, err := checkReadiness(skillDir)
	require.NoError(t, err)
	require.NotNil(t, report)
	require.True(t, report.hasEval, "scaffolded workspace should have eval detected")
	require.Empty(t, report.evalSchemaErrs, "scaffolded eval.yaml should have no schema errors: %v", report.evalSchemaErrs)
	require.Empty(t, report.taskSchemaErrs, "scaffolded task files should have no schema errors: %v", report.taskSchemaErrs)

	// Also run the check command and verify output
	cmd := newCheckCommand()
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetErr(&output)
	cmd.SetArgs([]string{skillName})

	err = cmd.Execute()
	require.NoError(t, err)

	result := output.String()
	assert.Contains(t, result, "eval.yaml schema valid")
	assert.Contains(t, result, "task file(s) validated")
	assert.NotContains(t, result, "Eval Schema:")
	assert.NotContains(t, result, "Task Schema:")
}

func TestCheckCommand_ScaffoldedEvalMatchesSchema(t *testing.T) {
	// Directly validate scaffold output against schema without full workspace
	dir := t.TempDir()
	tasksDir := filepath.Join(dir, "tasks")
	require.NoError(t, os.MkdirAll(tasksDir, 0o755))

	skillName := "schema-scaffold-test"

	// Write scaffold eval.yaml
	evalContent := scaffold.EvalYAML(skillName, "mock", "gpt-4o")
	evalPath := filepath.Join(dir, "eval.yaml")
	require.NoError(t, os.WriteFile(evalPath, []byte(evalContent), 0o644))

	// Write scaffold task files
	for name, content := range scaffold.TaskFiles(skillName) {
		require.NoError(t, os.WriteFile(filepath.Join(tasksDir, name), []byte(content), 0o644))
	}

	// Validate using the validation package directly
	evalErrs, taskErrs, err := validation.ValidateEvalFile(evalPath)
	require.NoError(t, err)
	require.Empty(t, evalErrs, "scaffolded eval.yaml should pass schema: %v", evalErrs)
	require.Empty(t, taskErrs, "scaffolded task files should pass schema: %v", taskErrs)
}
