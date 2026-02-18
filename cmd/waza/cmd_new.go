package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/spboyer/waza/internal/wizard"
)

func newNewCommand() *cobra.Command {
	var template string

	cmd := &cobra.Command{
		Use:   "new <skill-name>",
		Short: "Create a new skill with its eval suite",
		Long: `Create a new skill and its evaluation suite with a compliant directory structure.

Two modes of operation:

  Inside a project (skills/ directory detected):
    Creates skills/{name}/SKILL.md and evals/{name}/ with eval.yaml,
    task files, and fixtures.

  Standalone (no skills/ directory):
    Creates {name}/ with SKILL.md, evals/, .github/workflows/eval.yml,
    .gitignore, and README.md.

When running in a terminal (TTY), launches an interactive wizard for skill
metadata collection. In non-interactive environments (CI, pipes), uses defaults.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return newCommandE(cmd, args, template)
		},
	}

	cmd.Flags().StringVarP(&template, "template", "t", "", "Template pack to use (coming soon)")

	return cmd
}

func newCommandE(cmd *cobra.Command, args []string, templatePack string) error {
	skillName := args[0]

	if err := validateSkillName(skillName); err != nil {
		return err
	}

	if templatePack != "" {
		fmt.Fprintln(cmd.OutOrStdout(), "Note: template packs coming soon. Using default template.") //nolint:errcheck
	}

	// Determine mode based on skills/ directory presence
	projectRoot, inProject := findProjectRoot()

	var skillMDContent string
	// Check TTY from the command's input stream, not os.Stdin directly.
	inReader := cmd.InOrStdin()
	isTTY := false
	if f, ok := inReader.(*os.File); ok {
		isTTY = term.IsTerminal(int(f.Fd()))
	}
	if isTTY {
		spec, err := wizard.RunSkillWizard(cmd.InOrStdin(), cmd.OutOrStdout())
		if err != nil {
			return fmt.Errorf("wizard failed: %w", err)
		}
		// Validate wizard name against CLI arg
		if spec.Name != "" && spec.Name != skillName {
			return fmt.Errorf("wizard name %q does not match CLI argument %q", spec.Name, skillName)
		}
		spec.Name = skillName
		content, err := wizard.GenerateSkillMD(spec)
		if err != nil {
			return fmt.Errorf("failed to generate SKILL.md: %w", err)
		}
		skillMDContent = content
	} else {
		skillMDContent = defaultSkillMD(skillName)
	}

	if inProject {
		return scaffoldInProject(cmd, projectRoot, skillName, skillMDContent)
	}
	return scaffoldStandalone(cmd, skillName, skillMDContent)
}

// validateSkillName rejects names with path-traversal characters or empty names.
func validateSkillName(name string) error {
	if name == "" {
		return fmt.Errorf("skill name must not be empty")
	}
	cleaned := filepath.Clean(name)
	if cleaned == ".." || strings.Contains(cleaned, "/") || strings.Contains(cleaned, "\\") {
		return fmt.Errorf("skill name %q contains invalid path characters", name)
	}
	return nil
}

// findProjectRoot walks up from CWD looking for a skills/ directory.
// Returns the directory containing skills/ and true, or ("", false) if not found.
func findProjectRoot() (string, bool) {
	dir, err := os.Getwd()
	if err != nil {
		return "", false
	}

	for {
		candidate := filepath.Join(dir, "skills")
		info, err := os.Stat(candidate)
		if err == nil && info.IsDir() {
			return dir, true
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", false
}

// scaffoldInProject creates files within an existing project structure.
func scaffoldInProject(cmd *cobra.Command, projectRoot, skillName, skillMD string) error {
	skillDir := filepath.Join(projectRoot, "skills", skillName)
	evalDir := filepath.Join(projectRoot, "evals", skillName)
	tasksDir := filepath.Join(evalDir, "tasks")
	fixturesDir := filepath.Join(evalDir, "fixtures")

	// Create directories
	for _, d := range []string{skillDir, tasksDir, fixturesDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", d, err)
		}
	}

	files := []fileEntry{
		{filepath.Join(skillDir, "SKILL.md"), skillMD},
		{filepath.Join(evalDir, "eval.yaml"), defaultEvalYAML(skillName)},
		{filepath.Join(tasksDir, "basic-usage.yaml"), defaultBasicUsageTask()},
		{filepath.Join(tasksDir, "edge-case.yaml"), defaultEdgeCaseTask()},
		{filepath.Join(tasksDir, "should-not-trigger.yaml"), defaultShouldNotTriggerTask()},
		{filepath.Join(fixturesDir, "sample.py"), defaultFixture()},
	}

	return writeFiles(cmd, files)
}

// scaffoldStandalone creates a self-contained skill directory.
func scaffoldStandalone(cmd *cobra.Command, skillName, skillMD string) error {
	rootDir := skillName
	evalsDir := filepath.Join(rootDir, "evals")
	tasksDir := filepath.Join(evalsDir, "tasks")
	fixturesDir := filepath.Join(evalsDir, "fixtures")
	workflowDir := filepath.Join(rootDir, ".github", "workflows")

	// Create directories
	for _, d := range []string{tasksDir, fixturesDir, workflowDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", d, err)
		}
	}

	files := []fileEntry{
		{filepath.Join(rootDir, "SKILL.md"), skillMD},
		{filepath.Join(evalsDir, "eval.yaml"), defaultEvalYAML(skillName)},
		{filepath.Join(tasksDir, "basic-usage.yaml"), defaultBasicUsageTask()},
		{filepath.Join(tasksDir, "edge-case.yaml"), defaultEdgeCaseTask()},
		{filepath.Join(tasksDir, "should-not-trigger.yaml"), defaultShouldNotTriggerTask()},
		{filepath.Join(fixturesDir, "sample.py"), defaultFixture()},
		{filepath.Join(workflowDir, "eval.yml"), defaultCIWorkflow(skillName)},
		{filepath.Join(rootDir, ".gitignore"), defaultGitignore()},
		{filepath.Join(rootDir, "README.md"), defaultReadme(skillName)},
	}

	return writeFiles(cmd, files)
}

// fileEntry pairs a path with its content for batch writing.
type fileEntry struct {
	path    string
	content string
}

// writeFiles writes each file, skipping any that already exist.
func writeFiles(cmd *cobra.Command, files []fileEntry) error {
	fmt.Fprintln(cmd.OutOrStdout(), "Scaffolding skill:") //nolint:errcheck

	for _, f := range files {
		if _, err := os.Stat(f.path); err == nil {
			fmt.Fprintf(cmd.OutOrStdout(), "  skip %s (already exists)\n", f.path) //nolint:errcheck
			continue
		}

		if err := os.WriteFile(f.path, []byte(f.content), 0o644); err != nil {
			return fmt.Errorf("failed to write %s: %w", f.path, err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "  create %s\n", f.path) //nolint:errcheck
	}

	return nil
}

// titleCase converts a kebab-case name to Title Case.
func titleCase(s string) string {
	words := strings.Split(s, "-")
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

// --- Template content functions ---

func defaultSkillMD(name string) string {
	return fmt.Sprintf(`---
name: %s
type: utility
description: |
  USE FOR: %s tasks, ...
  DO NOT USE FOR: unrelated tasks, ...
---

# %s

## Overview

Describe what this skill does and when an agent should use it.

## Usage

Provide examples of prompts that should trigger this skill.

## References

List any reference documents or APIs this skill depends on.
`, name, name, titleCase(name))
}

func defaultEvalYAML(name string) string {
	engine, model := readProjectDefaults()
	return fmt.Sprintf(`name: %s-eval
description: Evaluation suite for %s.
skill: %s
version: "1.0"
config:
  trials_per_task: 1
  timeout_seconds: 300
  parallel: false
  executor: %s
  model: %s
graders:
  - type: code
    name: has_output
    config:
      assertions:
        - "len(output) > 0"
  - type: regex
    name: relevant_content
    config:
      pattern: "(?i)(explain|describe|analyze|implement)"
tasks:
  - "tasks/*.yaml"
`, name, name, name, engine, model)
}

// readProjectDefaults reads engine and model from .waza.yaml if it exists.
// Falls back to copilot-sdk and claude-sonnet-4.6.
func readProjectDefaults() (engine, model string) {
	engine = "copilot-sdk"
	model = "claude-sonnet-4.6"

	// Walk up from CWD looking for .waza.yaml
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	for i := 0; i < 10; i++ {
		configPath := filepath.Join(dir, ".waza.yaml")
		data, err := os.ReadFile(configPath)
		if err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "engine:") {
					if v := strings.TrimSpace(strings.TrimPrefix(line, "engine:")); v != "" {
						engine = v
					}
				}
				if strings.HasPrefix(line, "model:") {
					if v := strings.TrimSpace(strings.TrimPrefix(line, "model:")); v != "" {
						model = v
					}
				}
			}
			return
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return
}

func defaultBasicUsageTask() string {
	return `id: basic-usage-001
name: Basic Usage
description: |
  Test that the skill handles a typical request correctly.
tags:
  - basic
  - happy-path
inputs:
  prompt: "Help me with this task"
  files:
    - path: sample.py
expected:
  output_contains:
    - "function"
  outcomes:
    - type: task_completed
`
}

func defaultEdgeCaseTask() string {
	return `id: edge-case-001
name: Edge Case - Empty Input
description: |
  Test that the skill handles edge cases gracefully.
tags:
  - edge-case
inputs:
  prompt: ""
expected:
  outcomes:
    - type: task_completed
`
}

func defaultShouldNotTriggerTask() string {
	return `id: should-not-trigger-001
name: Should Not Trigger
description: |
  Test that the skill does NOT activate on unrelated prompts.
  This validates trigger specificity.
tags:
  - anti-trigger
  - negative-test
inputs:
  prompt: "What is the weather today?"
expected:
  output_not_contains:
    - "skill activated"
`
}

func defaultFixture() string {
	return `def hello(name):
    """Greet someone by name."""
    return f"Hello, {name}!"
`
}

func defaultCIWorkflow(name string) string {
	return fmt.Sprintf(`name: Eval %s

on:
  pull_request:
    branches: [main]

permissions:
  contents: read

jobs:
  run-eval:
    name: Run Evaluation
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install waza
        run: |
          curl -fsSL https://raw.githubusercontent.com/spboyer/waza/main/install.sh | bash

      - name: Run eval
        run: |
          waza run evals/eval.yaml --context-dir evals/fixtures -v
`, titleCase(name))
}

func defaultGitignore() string {
	return `results.json
.waza-cache/
coverage.txt
`
}

func defaultReadme(name string) string {
	return fmt.Sprintf(`# %s

A skill for agent evaluation with [waza](https://github.com/spboyer/waza).

## Quick Start

1. Edit `+"`SKILL.md`"+` with your skill's description and triggers.
2. Customize the task files in `+"`evals/tasks/`"+`.
3. Add real fixtures to `+"`evals/fixtures/`"+`.
4. Run the evaluation:

`+"```bash"+`
waza run evals/eval.yaml --context-dir evals/fixtures -v
`+"```"+`

## Structure

`+"```"+`
%s/
├── SKILL.md                  # Skill definition
├── evals/
│   ├── eval.yaml             # Eval configuration
│   ├── tasks/
│   │   ├── basic-usage.yaml
│   │   ├── edge-case.yaml
│   │   └── should-not-trigger.yaml
│   └── fixtures/
│       └── sample.py
├── .github/workflows/
│   └── eval.yml              # CI workflow
├── .gitignore
└── README.md
`+"```"+`

## Learn More

- [Waza Documentation](https://github.com/spboyer/waza)
`, titleCase(name), name)
}
