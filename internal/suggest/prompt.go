package suggest

import (
	"fmt"
	"strings"
)

const evalYAMLSchemaSummary = `Top-level eval.yaml fields:
- name (string)
- description (string)
- skill (string)
- version (string)
- config:
  - trials_per_task (int >= 1)
  - timeout_seconds (int >= 1)
  - parallel (bool)
  - executor (mock|copilot-sdk)
  - model (string)
- graders[]:
  - type (code|prompt|regex|file|keyword|json_schema|program|behavior|action_sequence|skill_invocation|diff)
  - name (string)
  - config (map)
- metrics[]:
  - name (string)
  - weight (float)
  - threshold (float)
- tasks[] (glob patterns, usually "tasks/*.yaml")`

const exampleEvalYAML = `name: example-skill-eval
description: Evaluation suite for example-skill
skill: example-skill
version: "1.0"
config:
  trials_per_task: 1
  timeout_seconds: 300
  parallel: false
  executor: copilot-sdk
  model: claude-opus-4.6
graders:
  - type: code
    name: has_output
    config:
      assertions:
        - "len(output) > 0"
  - type: regex
    name: no_errors
    config:
      must_not_match:
        - "(?i)error|exception"
metrics:
  - name: task_completion
    weight: 1.0
    threshold: 0.8
tasks:
  - "tasks/*.yaml"`

type promptData struct {
	SkillName      string
	Description    string
	Triggers       string
	AntiTriggers   string
	ContentSummary string
	GraderTypes    string
	SkillContent   string
}

func renderPrompt(data promptData) string {
	var b strings.Builder
	b.WriteString("You are generating a waza evaluation suite for a skill.\\n")
	b.WriteString("Return ONLY YAML in this exact schema:\\n\\n")
	b.WriteString("eval_yaml: |\\n")
	b.WriteString("  <full eval.yaml content>\\n")
	b.WriteString("tasks:\\n")
	b.WriteString("  - path: tasks/<task-file>.yaml\\n")
	b.WriteString("    content: |\\n")
	b.WriteString("      <task yaml>\\n")
	b.WriteString("fixtures:\\n")
	b.WriteString("  - path: fixtures/<fixture-file>\\n")
	b.WriteString("    content: |\\n")
	b.WriteString("      <fixture content>\\n\\n")
	b.WriteString("Requirements:\\n")
	b.WriteString("- Ensure eval_yaml is valid waza BenchmarkSpec YAML.\\n")
	b.WriteString("- Include at least 3 diverse tasks and at least 1 negative/anti-trigger task.\\n")
	b.WriteString("- Use grader types from the allowed list only.\\n")
	b.WriteString("- Keep task IDs deterministic and kebab-case.\\n")
	b.WriteString("- Make fixtures minimal and realistic for the tasks.\\n\\n")
	b.WriteString("Skill metadata:\\n")
	b.WriteString(fmt.Sprintf("- Name: %s\\n", data.SkillName))
	b.WriteString(fmt.Sprintf("- Description: %s\\n", data.Description))
	b.WriteString(fmt.Sprintf("- Triggers (USE FOR): %s\\n", data.Triggers))
	b.WriteString(fmt.Sprintf("- Anti-triggers (DO NOT USE FOR): %s\\n", data.AntiTriggers))
	b.WriteString(fmt.Sprintf("- Content summary: %s\\n\\n", data.ContentSummary))
	b.WriteString("Available grader types:\\n")
	b.WriteString(data.GraderTypes)
	b.WriteString("\\n\\n")
	b.WriteString("waza eval YAML schema summary:\\n")
	b.WriteString(evalYAMLSchemaSummary)
	b.WriteString("\\n\\n")
	b.WriteString("Example eval.yaml:\\n")
	b.WriteString(exampleEvalYAML)
	b.WriteString("\\n\\n")
	b.WriteString("Skill content (SKILL.md):\\n")
	b.WriteString(data.SkillContent)
	b.WriteString("\\n")
	return b.String()
}
