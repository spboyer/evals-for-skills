package suggest

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spboyer/waza/internal/execution"
	"github.com/spboyer/waza/internal/models"
	"github.com/spboyer/waza/internal/scaffold"
	"github.com/spboyer/waza/internal/skill"
	"gopkg.in/yaml.v3"
)

const defaultTimeoutSec = 120

// Options configures suggestion generation.
type Options struct {
	SkillPath  string
	TimeoutSec int
}

// GeneratedFile is a single generated artifact.
type GeneratedFile struct {
	Path    string `yaml:"path" json:"path"`
	Content string `yaml:"content" json:"content"`
}

// Suggestion is the structured output returned by the LLM.
type Suggestion struct {
	EvalYAML string          `yaml:"eval_yaml" json:"eval_yaml"`
	Tasks    []GeneratedFile `yaml:"tasks,omitempty" json:"tasks,omitempty"`
	Fixtures []GeneratedFile `yaml:"fixtures,omitempty" json:"fixtures,omitempty"`
}

// Generate runs the suggestion flow end-to-end.
func Generate(ctx context.Context, engine execution.AgentEngine, opts Options) (*Suggestion, error) {
	skillFile, err := resolveSkillFile(opts.SkillPath)
	if err != nil {
		return nil, err
	}

	skillContent, sk, err := loadSkill(skillFile)
	if err != nil {
		return nil, err
	}

	prompt := BuildPrompt(sk, skillContent)
	timeoutSec := opts.TimeoutSec
	if timeoutSec <= 0 {
		timeoutSec = defaultTimeoutSec
	}

	resp, err := engine.Execute(ctx, &execution.ExecutionRequest{
		TestID:     "waza-suggest",
		Message:    prompt,
		TimeoutSec: timeoutSec,
	})
	if err != nil {
		return nil, fmt.Errorf("getting suggestions: %w", err)
	}
	if resp == nil {
		return nil, errors.New("empty engine response")
	}

	suggestion, err := ParseResponse(resp.FinalOutput)
	if err != nil {
		return nil, fmt.Errorf("parsing suggest response: %w", err)
	}
	return suggestion, nil
}

// BuildPrompt builds the LLM prompt for eval suggestions.
func BuildPrompt(sk *skill.Skill, skillContent string) string {
	useFor, doNotUseFor := scaffold.ParseTriggerPhrases(sk.Frontmatter.Description)

	promptData := promptData{
		SkillName:      orDefault(sk.Frontmatter.Name, filepath.Base(filepath.Dir(sk.Path))),
		Description:    strings.TrimSpace(sk.Frontmatter.Description),
		Triggers:       phrasesToText(useFor),
		AntiTriggers:   phrasesToText(doNotUseFor),
		ContentSummary: summarizeBody(sk.Body),
		GraderTypes:    "- " + strings.Join(AvailableGraderTypes(), "\n- "),
		SkillContent:   skillContent,
	}
	return renderPrompt(promptData)
}

// AvailableGraderTypes returns supported grader kinds.
func AvailableGraderTypes() []string {
	return []string{
		string(models.GraderKindInlineScript),
		string(models.GraderKindPrompt),
		string(models.GraderKindRegex),
		string(models.GraderKindFile),
		string(models.GraderKindKeyword),
		string(models.GraderKindJSONSchema),
		string(models.GraderKindProgram),
		string(models.GraderKindBehavior),
		string(models.GraderKindActionSequence),
		string(models.GraderKindSkillInvocation),
		string(models.GraderKindDiff),
	}
}

// ParseResponse parses model YAML output into a Suggestion.
func ParseResponse(raw string) (*Suggestion, error) {
	normalized := extractYAML(raw)

	var s Suggestion
	if err := yaml.Unmarshal([]byte(normalized), &s); err == nil && strings.TrimSpace(s.EvalYAML) != "" {
		if err := validateEvalYAML(s.EvalYAML); err != nil {
			return nil, err
		}
		return &s, nil
	}

	if err := validateEvalYAML(normalized); err == nil {
		return &Suggestion{EvalYAML: normalized}, nil
	}

	return nil, errors.New("response is not valid suggestion YAML")
}

// WriteToDir writes suggested files to outputDir and returns written paths.
func (s *Suggestion) WriteToDir(outputDir string) ([]string, error) {
	if err := validateEvalYAML(s.EvalYAML); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, fmt.Errorf("creating output directory: %w", err)
	}

	var written []string
	evalPath := filepath.Join(outputDir, "eval.yaml")
	if err := os.WriteFile(evalPath, []byte(strings.TrimSpace(s.EvalYAML)+"\n"), 0o644); err != nil {
		return nil, fmt.Errorf("writing eval.yaml: %w", err)
	}
	written = append(written, evalPath)

	for i, task := range s.Tasks {
		path, err := normalizeGeneratedPath(task.Path, fmt.Sprintf("tasks/task-%02d.yaml", i+1))
		if err != nil {
			return nil, err
		}
		target := filepath.Join(outputDir, path)
		if err := writeGeneratedFile(target, task.Content); err != nil {
			return nil, err
		}
		written = append(written, target)
	}

	for i, fixture := range s.Fixtures {
		path, err := normalizeGeneratedPath(fixture.Path, fmt.Sprintf("fixtures/fixture-%02d.txt", i+1))
		if err != nil {
			return nil, err
		}
		target := filepath.Join(outputDir, path)
		if err := writeGeneratedFile(target, fixture.Content); err != nil {
			return nil, err
		}
		written = append(written, target)
	}

	return written, nil
}

func loadSkill(skillFile string) (string, *skill.Skill, error) {
	data, err := os.ReadFile(skillFile)
	if err != nil {
		return "", nil, fmt.Errorf("reading SKILL.md: %w", err)
	}
	var sk skill.Skill
	if err := sk.UnmarshalText(data); err != nil {
		return "", nil, fmt.Errorf("parsing SKILL.md: %w", err)
	}
	sk.Path = skillFile
	return string(data), &sk, nil
}

func resolveSkillFile(input string) (string, error) {
	if strings.TrimSpace(input) == "" {
		return "", errors.New("skill path is required")
	}
	resolved := input
	if !filepath.IsAbs(resolved) {
		wd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("getting working directory: %w", err)
		}
		resolved = filepath.Join(wd, resolved)
	}

	info, err := os.Stat(resolved)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("skill path does not exist: %s", input)
		}
		return "", fmt.Errorf("checking skill path: %w", err)
	}

	if info.IsDir() {
		resolved = filepath.Join(resolved, "SKILL.md")
	}

	if filepath.Base(resolved) != "SKILL.md" {
		return "", fmt.Errorf("expected SKILL.md or skill directory, got %s", input)
	}
	if _, err := os.Stat(resolved); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("no SKILL.md found in %s", input)
		}
		return "", fmt.Errorf("checking SKILL.md: %w", err)
	}
	return resolved, nil
}

func extractYAML(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}

	start := strings.Index(trimmed, "```")
	if start < 0 {
		return trimmed
	}

	rest := trimmed[start+3:]
	if nl := strings.Index(rest, "\n"); nl >= 0 {
		rest = rest[nl+1:]
	}
	if end := strings.Index(rest, "```"); end >= 0 {
		return strings.TrimSpace(rest[:end])
	}

	return trimmed
}

func validateEvalYAML(raw string) error {
	var spec models.BenchmarkSpec
	if err := yaml.Unmarshal([]byte(raw), &spec); err != nil {
		return fmt.Errorf("invalid eval_yaml: %w", err)
	}
	if err := spec.Validate(); err != nil {
		return fmt.Errorf("invalid eval_yaml: %w", err)
	}
	return nil
}

func phrasesToText(phrases []scaffold.TriggerPhrase) string {
	if len(phrases) == 0 {
		return "none"
	}
	items := make([]string, 0, len(phrases))
	for _, p := range phrases {
		if strings.TrimSpace(p.Prompt) != "" {
			items = append(items, p.Prompt)
		}
	}
	if len(items) == 0 {
		return "none"
	}
	return strings.Join(items, ", ")
}

func summarizeBody(body string) string {
	lines := strings.Split(body, "\n")
	var highlights []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "#") {
			highlights = append(highlights, trimmed)
			continue
		}
		if len(highlights) < 8 {
			highlights = append(highlights, trimmed)
		}
		if len(highlights) >= 8 {
			break
		}
	}
	if len(highlights) == 0 {
		return "No body content"
	}
	return strings.Join(highlights, " | ")
}

func normalizeGeneratedPath(path, fallback string) (string, error) {
	clean := strings.TrimSpace(path)
	if clean == "" {
		clean = fallback
	}
	clean = filepath.Clean(clean)
	if filepath.IsAbs(clean) || strings.HasPrefix(clean, "..") {
		return "", fmt.Errorf("invalid generated path: %s", path)
	}
	return clean, nil
}

func writeGeneratedFile(path string, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating directory for %s: %w", path, err)
	}
	if err := os.WriteFile(path, []byte(strings.TrimSpace(content)+"\n"), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}

func orDefault(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
