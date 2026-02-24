// Package checks provides the ComplianceChecker interface and implementations
// for validating skill readiness across multiple dimensions.
package checks

import "github.com/spboyer/waza/internal/skill"

// CheckResult holds the outcome of a single compliance check.
type CheckResult struct {
	// Name is a stable check identifier used in output and downstream processing.
	Name string
	// Passed indicates whether the check met its acceptance criteria.
	Passed bool
	// Summary is a human-readable one-line result intended for concise display.
	Summary string
	// Details provides optional supporting lines for diagnostics or remediation.
	Details []string
	// Data carries an optional checker-specific payload for structured consumers.
	Data any
}

// ComplianceChecker runs a single compliance check.
type ComplianceChecker interface {
	Name() string
	Check(skill.Skill) (*CheckResult, error)
}
