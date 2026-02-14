package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spboyer/waza/cmd/waza/dev"
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
