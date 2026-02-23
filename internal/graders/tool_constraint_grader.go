package graders

import (
	"context"
	"fmt"
	"strings"

	"github.com/spboyer/waza/internal/models"
)

// toolConstraintGrader validates which tools an agent should/shouldn't use,
// plus turn and token limits.
type toolConstraintGrader struct {
	name        string
	expectTools []string
	rejectTools []string
	maxTurns    int
	maxTokens   int
}

// ToolConstraintGraderParams holds the mapstructure-decoded parameters.
type ToolConstraintGraderParams struct {
	ExpectTools []string `mapstructure:"expect_tools"`
	RejectTools []string `mapstructure:"reject_tools"`
	MaxTurns    int      `mapstructure:"max_turns"`
	MaxTokens   int      `mapstructure:"max_tokens"`
}

// NewToolConstraintGrader creates a toolConstraintGrader from decoded parameters.
func NewToolConstraintGrader(name string, params ToolConstraintGraderParams) (*toolConstraintGrader, error) {
	if len(params.ExpectTools) == 0 && len(params.RejectTools) == 0 &&
		params.MaxTurns == 0 && params.MaxTokens == 0 {
		return nil, fmt.Errorf("tool_constraint grader '%s' must have at least one constraint configured", name)
	}

	return &toolConstraintGrader{
		name:        name,
		expectTools: params.ExpectTools,
		rejectTools: params.RejectTools,
		maxTurns:    params.MaxTurns,
		maxTokens:   params.MaxTokens,
	}, nil
}

func (tc *toolConstraintGrader) Name() string            { return tc.name }
func (tc *toolConstraintGrader) Kind() models.GraderKind { return models.GraderKindToolConstraint }

func (tc *toolConstraintGrader) Grade(ctx context.Context, gradingContext *Context) (*models.GraderResults, error) {
	return measureTime(func() (*models.GraderResults, error) {
		session := gradingContext.Session
		if session == nil {
			return &models.GraderResults{
				Name:     tc.name,
				Type:     models.GraderKindToolConstraint,
				Score:    0.0,
				Passed:   false,
				Feedback: "No session digest available for tool constraint grading",
			}, nil
		}

		var failures []string

		failures = append(failures, tc.checkExpectTools(session)...)
		failures = append(failures, tc.checkRejectTools(session)...)
		failures = append(failures, tc.checkMaxTurns(session)...)
		failures = append(failures, tc.checkMaxTokens(session)...)

		totalChecks := tc.countTotalChecks()
		passedChecks := totalChecks - len(failures)

		score := 1.0
		if totalChecks > 0 {
			score = float64(passedChecks) / float64(totalChecks)
		}

		feedback := "All tool constraint checks passed"
		if len(failures) > 0 {
			feedback = strings.Join(failures, "; ")
		}

		return &models.GraderResults{
			Name:     tc.name,
			Type:     models.GraderKindToolConstraint,
			Score:    score,
			Passed:   len(failures) == 0,
			Feedback: feedback,
			Details: map[string]any{
				"expect_tools": tc.expectTools,
				"reject_tools": tc.rejectTools,
				"max_turns":    tc.maxTurns,
				"max_tokens":   tc.maxTokens,
				"failures":     failures,
				"tools_used":   session.ToolsUsed,
				"total_turns":  session.TotalTurns,
				"tokens_total": session.TokensTotal,
			},
		}, nil
	})
}

func (tc *toolConstraintGrader) checkExpectTools(session *models.SessionDigest) []string {
	if len(tc.expectTools) == 0 {
		return nil
	}

	toolSet := make(map[string]bool, len(session.ToolsUsed))
	for _, t := range session.ToolsUsed {
		toolSet[t] = true
	}

	var failures []string
	for _, expected := range tc.expectTools {
		if !toolSet[expected] {
			failures = append(failures, fmt.Sprintf("Expected tool not used: %s", expected))
		}
	}
	return failures
}

func (tc *toolConstraintGrader) checkRejectTools(session *models.SessionDigest) []string {
	if len(tc.rejectTools) == 0 {
		return nil
	}

	toolSet := make(map[string]bool, len(session.ToolsUsed))
	for _, t := range session.ToolsUsed {
		toolSet[t] = true
	}

	var failures []string
	for _, rejected := range tc.rejectTools {
		if toolSet[rejected] {
			failures = append(failures, fmt.Sprintf("Rejected tool was used: %s", rejected))
		}
	}
	return failures
}

func (tc *toolConstraintGrader) checkMaxTurns(session *models.SessionDigest) []string {
	if tc.maxTurns <= 0 {
		return nil
	}

	if session.TotalTurns > tc.maxTurns {
		return []string{fmt.Sprintf("Turn count %d exceeds max allowed %d", session.TotalTurns, tc.maxTurns)}
	}
	return nil
}

func (tc *toolConstraintGrader) checkMaxTokens(session *models.SessionDigest) []string {
	if tc.maxTokens <= 0 {
		return nil
	}

	if session.TokensTotal > tc.maxTokens {
		return []string{fmt.Sprintf("Token usage %d exceeds max allowed %d", session.TokensTotal, tc.maxTokens)}
	}
	return nil
}

func (tc *toolConstraintGrader) countTotalChecks() int {
	total := len(tc.expectTools) + len(tc.rejectTools)
	if tc.maxTurns > 0 {
		total++
	}
	if tc.maxTokens > 0 {
		total++
	}
	return total
}
