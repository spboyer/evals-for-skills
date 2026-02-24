package checks

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spboyer/waza/internal/skill"
	"github.com/spboyer/waza/internal/workspace"
)

// EvalSuiteChecker checks for the presence of an eval.yaml file.
type EvalSuiteChecker struct {
	WorkDir string
}

// EvalSuiteData holds the structured output of an eval suite check.
type EvalSuiteData struct {
	Found bool
}

var _ ComplianceChecker = (*EvalSuiteChecker)(nil)

func (c *EvalSuiteChecker) Name() string { return "eval-suite" }

func (c *EvalSuiteChecker) Check(sk skill.Skill) (*CheckResult, error) {
	found := false
	skillDir := filepath.Dir(sk.Path)

	if wd := c.WorkDir; wd != "" {
		if ctx, ctxErr := workspace.DetectContext(wd); ctxErr == nil && ctx.Type != workspace.ContextNone {
			if evalPath, findErr := workspace.FindEval(ctx, sk.Frontmatter.Name); findErr == nil && evalPath != "" {
				found = true
			}
		}
	}

	// Fallback: co-located eval.yaml
	if !found {
		evalPath := filepath.Join(skillDir, "eval.yaml")
		if _, err := os.Stat(evalPath); err == nil {
			found = true
		}
	}

	summary := "Evaluation Suite: Not Found"
	if found {
		summary = "Evaluation Suite: Found"
	}

	return &CheckResult{
		Name:    c.Name(),
		Passed:  true, // eval is a recommendation, not a requirement
		Summary: summary,
		Data:    &EvalSuiteData{Found: found},
	}, nil
}

// Eval is a convenience wrapper that returns the typed data directly.
func (c *EvalSuiteChecker) Eval(sk skill.Skill) (*EvalSuiteData, error) {
	result, err := c.Check(sk)
	if err != nil {
		return nil, err
	}
	d, ok := result.Data.(*EvalSuiteData)
	if !ok {
		err = fmt.Errorf("unexpected data type: %T", result.Data)
	}
	return d, err
}
