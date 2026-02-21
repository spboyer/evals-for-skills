package suggest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spboyer/waza/internal/skill"
	"github.com/stretchr/testify/require"
)

func TestBuildPromptIncludesSkillMetadata(t *testing.T) {
	raw := `---
name: prompt-skill
description: "Useful skill. USE FOR: summarize, explain. DO NOT USE FOR: coding, deployment."
---

# Prompt Skill

## Overview
This skill summarizes docs.
`

	var sk skill.Skill
	require.NoError(t, sk.UnmarshalText([]byte(raw)))
	sk.Path = filepath.Join(t.TempDir(), "SKILL.md")

	prompt := BuildPrompt(&sk, raw)
	require.Contains(t, prompt, "Name: prompt-skill")
	require.Contains(t, prompt, "Triggers (USE FOR): summarize, explain")
	require.Contains(t, prompt, "Anti-triggers (DO NOT USE FOR): coding, deployment")
	require.Contains(t, prompt, "Available grader types")
	require.Contains(t, prompt, "waza eval YAML schema summary")
	require.Contains(t, prompt, "Skill content (SKILL.md)")
}

func TestParseResponseStructuredYAML(t *testing.T) {
	resp := "```yaml\neval_yaml: |\n  name: generated-eval\n  description: generated\n  skill: sample\n  version: \"1.0\"\n  config:\n    trials_per_task: 1\n    timeout_seconds: 120\n    parallel: false\n    executor: mock\n    model: test\n  graders:\n    - type: code\n      name: has_output\n      config:\n        assertions:\n          - \\\"len(output) > 0\\\"\n  metrics:\n    - name: completion\n      weight: 1.0\n      threshold: 0.8\n  tasks:\n    - \"tasks/*.yaml\"\ntasks:\n  - path: tasks/basic.yaml\n    content: |\n      id: basic-001\n      name: Basic\n      inputs:\n        prompt: \"hello\"\nfixtures:\n  - path: fixtures/sample.txt\n    content: |\n      sample\n```"

	s, err := ParseResponse(resp)
	require.NoError(t, err)
	require.Equal(t, 1, len(s.Tasks))
	require.Equal(t, "tasks/basic.yaml", s.Tasks[0].Path)
	require.Equal(t, 1, len(s.Fixtures))
}

func TestParseResponseInvalid(t *testing.T) {
	_, err := ParseResponse("not valid yaml")
	require.Error(t, err)
}

func TestWriteToDirWritesFiles(t *testing.T) {
	s := &Suggestion{
		EvalYAML: `name: generated-eval
description: generated
skill: sample
version: "1.0"
config:
  trials_per_task: 1
  timeout_seconds: 120
  parallel: false
  executor: mock
  model: test
graders:
  - type: code
    name: has_output
    config:
      assertions:
        - "len(output) > 0"
metrics:
  - name: completion
    weight: 1.0
    threshold: 0.8
tasks:
  - "tasks/*.yaml"`,
		Tasks: []GeneratedFile{
			{Path: "tasks/basic.yaml", Content: "id: basic-001\nname: Basic\ninputs:\n  prompt: \"hello\""},
		},
		Fixtures: []GeneratedFile{
			{Path: "fixtures/sample.txt", Content: "sample"},
		},
	}

	outDir := t.TempDir()
	written, err := s.WriteToDir(outDir)
	require.NoError(t, err)
	require.Len(t, written, 3)

	evalData, err := os.ReadFile(filepath.Join(outDir, "eval.yaml"))
	require.NoError(t, err)
	require.Contains(t, string(evalData), "name: generated-eval")

	taskData, err := os.ReadFile(filepath.Join(outDir, "tasks", "basic.yaml"))
	require.NoError(t, err)
	require.True(t, strings.Contains(string(taskData), "id: basic-001"))
}
