package checks

import (
	"fmt"

	"github.com/spboyer/waza/internal/scoring"
	"github.com/spboyer/waza/internal/skill"
)

// TokenBudgetChecker validates that SKILL.md is within its token budget.
type TokenBudgetChecker struct {
	// Limit for SKILL.md tokens; 0 means use scoring.TokenSoftLimit
	Limit int
}

// TokenBudgetData holds the structured output of a token budget check.
type TokenBudgetData struct {
	TokenCount int
	TokenLimit int
	Exceeded   bool
}

var _ ComplianceChecker = (*TokenBudgetChecker)(nil)

func (c *TokenBudgetChecker) Name() string { return "token-budget" }

func (c *TokenBudgetChecker) Check(sk skill.Skill) (*CheckResult, error) {
	limit := c.Limit
	if limit == 0 {
		limit = scoring.TokenSoftLimit
	}

	count := sk.Tokens
	exceeded := count > limit

	summary := fmt.Sprintf("Token Budget: %d / %d tokens", count, limit)
	if exceeded {
		summary += fmt.Sprintf(" (exceeded by %d)", count-limit)
	}

	return &CheckResult{
		Name:    c.Name(),
		Passed:  !exceeded,
		Summary: summary,
		Data:    &TokenBudgetData{TokenCount: count, TokenLimit: limit, Exceeded: exceeded},
	}, nil
}

// Budget is a convenience wrapper that returns the typed data directly.
func (c *TokenBudgetChecker) Budget(sk skill.Skill) (*TokenBudgetData, error) {
	result, err := c.Check(sk)
	if err != nil {
		return nil, err
	}
	d, ok := result.Data.(*TokenBudgetData)
	if !ok {
		return nil, fmt.Errorf("unexpected data type: %T", result.Data)
	}
	return d, nil
}
