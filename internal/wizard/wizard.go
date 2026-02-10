package wizard

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"text/template"
)

// SkillType represents the category of a skill.
type SkillType string

const (
	SkillTypeWorkflow SkillType = "workflow"
	SkillTypeUtility  SkillType = "utility"
	SkillTypeAnalysis SkillType = "analysis"
)

// SkillSpec holds all fields collected during the interactive wizard.
type SkillSpec struct {
	Name         string
	Description  string
	Triggers     []string
	AntiTriggers []string
	Type         SkillType
}

var validSkillTypes = map[string]SkillType{
	"workflow": SkillTypeWorkflow,
	"utility":  SkillTypeUtility,
	"analysis": SkillTypeAnalysis,
}

const skillMDTemplate = `---
name: {{ .Name }}
type: {{ .Type }}
description: >
  {{ .Description }}
---

# {{ .Name }}

{{ .Description }}

## Usage

**USE FOR:**
{{- range .Triggers }}
- {{ . }}
{{- end }}

**DO NOT USE FOR:**
{{- range .AntiTriggers }}
- {{ . }}
{{- end }}
`

// RunSkillWizard runs an interactive prompt session to collect skill metadata.
// It reads from the provided reader and writes prompts to the provided writer.
func RunSkillWizard(in io.Reader, out io.Writer) (*SkillSpec, error) {
	scanner := bufio.NewScanner(in)
	spec := &SkillSpec{}

	name, err := prompt(scanner, out, "Skill name: ")
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, fmt.Errorf("skill name is required")
	}
	spec.Name = name

	desc, err := prompt(scanner, out, "Description: ")
	if err != nil {
		return nil, err
	}
	if desc == "" {
		return nil, fmt.Errorf("description is required")
	}
	spec.Description = desc

	triggersRaw, err := prompt(scanner, out, "Trigger phrases (comma-separated): ")
	if err != nil {
		return nil, err
	}
	spec.Triggers = splitAndTrim(triggersRaw)

	antiTriggersRaw, err := prompt(scanner, out, "Anti-trigger phrases (comma-separated): ")
	if err != nil {
		return nil, err
	}
	spec.AntiTriggers = splitAndTrim(antiTriggersRaw)

	skillType, err := prompt(scanner, out, "Skill type (workflow/utility/analysis): ")
	if err != nil {
		return nil, err
	}
	st, ok := validSkillTypes[strings.ToLower(strings.TrimSpace(skillType))]
	if !ok {
		return nil, fmt.Errorf("invalid skill type %q: must be workflow, utility, or analysis", skillType)
	}
	spec.Type = st

	return spec, nil
}

// GenerateSkillMD renders a SKILL.md from the given spec.
func GenerateSkillMD(spec *SkillSpec) (string, error) {
	tmpl, err := template.New("skillmd").Parse(skillMDTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, spec); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}
	return buf.String(), nil
}

func prompt(scanner *bufio.Scanner, out io.Writer, question string) (string, error) {
	if _, err := fmt.Fprint(out, question); err != nil {
		return "", err
	}
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		}
		return "", fmt.Errorf("unexpected end of input")
	}
	return strings.TrimSpace(scanner.Text()), nil
}

func splitAndTrim(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var result []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
