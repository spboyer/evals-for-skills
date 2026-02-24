package dev

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/spboyer/waza/internal/scoring"
	"github.com/spboyer/waza/internal/skill"
)

// allowedFrontmatterFields lists the fields permitted by the agentskills.io spec.
var allowedFrontmatterFields = map[string]bool{
	"name":          true,
	"description":   true,
	"license":       true,
	"allowed-tools": true,
	"metadata":      true,
	"compatibility": true,
}

// namePattern validates the spec name format: alphanumeric + hyphens,
// no leading/trailing hyphens, no consecutive --.
var namePattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// SpecScorer validates a SKILL.md against the agentskills.io specification.
type SpecScorer struct{}

// SpecResult holds the output from spec compliance checks.
type SpecResult struct {
	Issues []scoring.Issue
	Pass   int // number of checks that passed
	Total  int // total checks run
}

// Passed returns true when every check passed without errors.
func (r *SpecResult) Passed() bool {
	for _, iss := range r.Issues {
		if iss.Severity == "error" {
			return false
		}
	}
	return true
}

// Score runs all agentskills.io spec checks on the skill.
func (SpecScorer) Score(sk *skill.Skill) *SpecResult {
	r := &SpecResult{}

	if sk == nil {
		r.Total = 1
		r.Issues = append(r.Issues, scoring.Issue{
			Rule:     "spec-frontmatter",
			Message:  "Skill is nil",
			Severity: "error",
		})
		return r
	}

	specFrontmatter(sk, r)
	specAllowedFields(sk, r)
	specName(sk, r)
	specDirMatch(sk, r)
	specDescription(sk, r)
	specCompatibility(sk, r)
	specLicense(sk, r)
	specVersion(sk, r)

	return r
}

// specFrontmatter checks that YAML frontmatter exists and has name + description.
func specFrontmatter(sk *skill.Skill, r *SpecResult) {
	r.Total++
	if sk.FrontmatterRaw == nil {
		r.Issues = append(r.Issues, scoring.Issue{
			Rule:     "spec-frontmatter",
			Message:  "YAML frontmatter is missing",
			Severity: "error",
		})
		return
	}
	missing := []string{}
	if sk.Frontmatter.Name == "" {
		missing = append(missing, "name")
	}
	if strings.TrimSpace(sk.Frontmatter.Description) == "" {
		missing = append(missing, "description")
	}
	if len(missing) > 0 {
		r.Issues = append(r.Issues, scoring.Issue{
			Rule:     "spec-frontmatter",
			Message:  fmt.Sprintf("Required frontmatter fields missing: %s", strings.Join(missing, ", ")),
			Severity: "error",
		})
		return
	}
	r.Pass++
}

// specAllowedFields checks that only spec-allowed fields are present.
func specAllowedFields(sk *skill.Skill, r *SpecResult) {
	r.Total++
	if sk.FrontmatterRaw == nil {
		return // already flagged by spec-frontmatter
	}
	var unknown []string
	for key := range sk.FrontmatterRaw {
		if !allowedFrontmatterFields[key] {
			unknown = append(unknown, key)
		}
	}
	if len(unknown) > 0 {
		r.Issues = append(r.Issues, scoring.Issue{
			Rule:     "spec-allowed-fields",
			Message:  fmt.Sprintf("Unknown frontmatter fields: %s", strings.Join(unknown, ", ")),
			Severity: "warning",
		})
		return
	}
	r.Pass++
}

// specName validates the name format per agentskills.io rules.
func specName(sk *skill.Skill, r *SpecResult) {
	r.Total++
	name := sk.Frontmatter.Name
	if name == "" {
		// Already caught by spec-frontmatter; don't double-report.
		return
	}
	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
		r.Issues = append(r.Issues, scoring.Issue{
			Rule:     "spec-name",
			Message:  "Name must not start or end with a hyphen",
			Severity: "error",
		})
		return
	}
	if strings.Contains(name, "--") {
		r.Issues = append(r.Issues, scoring.Issue{
			Rule:     "spec-name",
			Message:  "Name must not contain consecutive hyphens (--)",
			Severity: "error",
		})
		return
	}
	if !namePattern.MatchString(name) {
		r.Issues = append(r.Issues, scoring.Issue{
			Rule:     "spec-name",
			Message:  "Name must contain only lowercase alphanumeric characters and hyphens",
			Severity: "error",
		})
		return
	}
	r.Pass++
}

// specDirMatch checks that the directory name matches the skill name field.
func specDirMatch(sk *skill.Skill, r *SpecResult) {
	r.Total++
	if sk.Path == "" || sk.Frontmatter.Name == "" {
		return // can't check without path or name
	}
	dir := filepath.Base(filepath.Dir(sk.Path))
	if dir != sk.Frontmatter.Name {
		r.Issues = append(r.Issues, scoring.Issue{
			Rule:     "spec-dir-match",
			Message:  fmt.Sprintf("Directory %q does not match skill name %q", dir, sk.Frontmatter.Name),
			Severity: "error",
		})
		return
	}
	r.Pass++
}

// specDescription validates description is non-empty and ≤1024 characters.
func specDescription(sk *skill.Skill, r *SpecResult) {
	r.Total++
	desc := strings.TrimSpace(sk.Frontmatter.Description)
	if desc == "" {
		// Already caught by spec-frontmatter; skip.
		return
	}
	length := utf8.RuneCountInString(desc)
	if length > 1024 {
		r.Issues = append(r.Issues, scoring.Issue{
			Rule:     "spec-description",
			Message:  fmt.Sprintf("Description is %d characters (max 1024)", length),
			Severity: "error",
		})
		return
	}
	r.Pass++
}

// specCompatibility checks the compatibility field length if present.
func specCompatibility(sk *skill.Skill, r *SpecResult) {
	r.Total++
	if sk.FrontmatterRaw == nil {
		r.Pass++
		return
	}
	val, ok := sk.FrontmatterRaw["compatibility"]
	if !ok {
		r.Pass++ // field is optional
		return
	}
	s, isStr := val.(string)
	if !isStr {
		r.Pass++ // non-string compat values are valid (could be a map)
		return
	}
	if utf8.RuneCountInString(s) > 500 {
		r.Issues = append(r.Issues, scoring.Issue{
			Rule:     "spec-compatibility",
			Message:  fmt.Sprintf("Compatibility field is %d characters (max 500)", utf8.RuneCountInString(s)),
			Severity: "error",
		})
		return
	}
	r.Pass++
}

// specLicense recommends adding a license field.
func specLicense(sk *skill.Skill, r *SpecResult) {
	r.Total++
	if sk.FrontmatterRaw == nil {
		return
	}
	if _, ok := sk.FrontmatterRaw["license"]; !ok {
		r.Issues = append(r.Issues, scoring.Issue{
			Rule:     "spec-license",
			Message:  "Consider adding a 'license' field (e.g., MIT, Apache-2.0)",
			Severity: "warning",
		})
		return
	}
	r.Pass++
}

// specVersion recommends adding metadata.version.
func specVersion(sk *skill.Skill, r *SpecResult) {
	r.Total++
	if sk.FrontmatterRaw == nil {
		return
	}
	meta, ok := sk.FrontmatterRaw["metadata"]
	if !ok {
		r.Issues = append(r.Issues, scoring.Issue{
			Rule:     "spec-version",
			Message:  "Consider adding 'metadata.version' for versioning",
			Severity: "warning",
		})
		return
	}
	metaMap, isMap := meta.(map[string]any)
	if !isMap {
		r.Issues = append(r.Issues, scoring.Issue{
			Rule:     "spec-version",
			Message:  "Consider adding 'metadata.version' for versioning",
			Severity: "warning",
		})
		return
	}
	if _, hasVersion := metaMap["version"]; !hasVersion {
		r.Issues = append(r.Issues, scoring.Issue{
			Rule:     "spec-version",
			Message:  "Consider adding 'metadata.version' for versioning",
			Severity: "warning",
		})
		return
	}
	r.Pass++
}
