package wizard

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/charmbracelet/huh"
	"golang.org/x/term"
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

// RunSkillWizard runs an interactive huh form to collect skill metadata.
// It reads from the provided reader and writes prompts to the provided writer.
func RunSkillWizard(in io.Reader, out io.Writer) (*SkillSpec, error) {
	var (
		name            string
		description     string
		triggersRaw     string
		antiTriggersRaw string
		skillType       string
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Skill name").
				Description("A kebab-case name for your skill").
				Placeholder("my-skill").
				Value(&name).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("skill name is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Description").
				Description("What does this skill do?").
				Placeholder("Describe your skill").
				Value(&description).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("description is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Trigger phrases").
				Description("Comma-separated phrases that should activate this skill").
				Placeholder("deploy app, create resource").
				Value(&triggersRaw),
			huh.NewInput().
				Title("Anti-trigger phrases").
				Description("Comma-separated phrases that should NOT activate this skill").
				Placeholder("unrelated task, wrong domain").
				Value(&antiTriggersRaw),
			huh.NewSelect[string]().
				Title("Skill type").
				Options(
					huh.NewOption("workflow", "workflow"),
					huh.NewOption("utility", "utility"),
					huh.NewOption("analysis", "analysis"),
				).
				Value(&skillType),
		),
	).
		WithInput(in).
		WithOutput(out)

	// Use accessible mode for non-TTY input (e.g., tests, piped input).
	if f, ok := in.(*os.File); !ok || !term.IsTerminal(int(f.Fd())) {
		form = form.WithAccessible(true)
	}

	if err := form.Run(); err != nil {
		return nil, fmt.Errorf("wizard failed: %w", err)
	}

	return &SkillSpec{
		Name:         strings.TrimSpace(name),
		Description:  strings.TrimSpace(description),
		Triggers:     splitAndTrim(triggersRaw),
		AntiTriggers: splitAndTrim(antiTriggersRaw),
		Type:         SkillType(skillType),
	}, nil
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
