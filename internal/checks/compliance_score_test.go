package checks

import (
	"testing"

	"github.com/spboyer/waza/internal/scoring"
	"github.com/spboyer/waza/internal/skill"
	"github.com/stretchr/testify/require"
)

var _ ComplianceChecker = (*ComplianceScoreChecker)(nil)

func TestComplianceScoreChecker_HighCompliance(t *testing.T) {
	sk := skill.Skill{
		Frontmatter: skill.Frontmatter{
			Name:        "test-skill",
			Description: `**WORKFLOW SKILL** - A comprehensive test skill for validating compliance checks. USE FOR: testing, validation, compliance checks, readiness assessment. DO NOT USE FOR: unrelated tasks, production use, deployment automation. INVOKES: internal validators, scoring engines. FOR SINGLE OPERATIONS: use direct commands.`,
		},
	}

	checker := &ComplianceScoreChecker{}
	result, err := checker.Check(sk)
	require.NoError(t, err)
	require.Equal(t, "compliance-score", result.Name)
	require.True(t, result.Passed)

	data, ok := result.Data.(*ComplianceScoreData)
	require.True(t, ok)
	require.Equal(t, scoring.AdherenceHigh, data.Level)
}

func TestComplianceScoreChecker_LowCompliance(t *testing.T) {
	sk := skill.Skill{
		Frontmatter: skill.Frontmatter{
			Name:        "test-skill",
			Description: "Short description.",
		},
	}

	checker := &ComplianceScoreChecker{}
	result, err := checker.Check(sk)
	require.NoError(t, err)
	require.False(t, result.Passed)

	data, ok := result.Data.(*ComplianceScoreData)
	require.True(t, ok)
	require.Equal(t, scoring.AdherenceLow, data.Level)
}

func TestComplianceScoreChecker_CustomScorer(t *testing.T) {
	sk := skill.Skill{
		Frontmatter: skill.Frontmatter{
			Name:        "test-skill",
			Description: "Any description.",
		},
	}

	scorer := &scoring.HeuristicScorer{}
	checker := &ComplianceScoreChecker{Scorer: scorer}
	result, err := checker.Check(sk)
	require.NoError(t, err)
	require.NotNil(t, result.Data)
}
