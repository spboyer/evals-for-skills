package dev

import (
	"strings"
	"testing"

	"github.com/spboyer/waza/internal/skill"
	"github.com/stretchr/testify/require"
)

func makeSpecSkill(name, description string, rawFields map[string]any, path string) *skill.Skill {
	raw := make(map[string]any)
	for k, v := range rawFields {
		raw[k] = v
	}
	raw["name"] = name
	raw["description"] = description
	return &skill.Skill{
		Frontmatter:    skill.Frontmatter{Name: name, Description: description},
		FrontmatterRaw: raw,
		Path:           path,
	}
}

func TestSpecScorer_NilSkill(t *testing.T) {
	r := (SpecScorer{}).Score(nil)
	require.False(t, r.Passed())
	require.Equal(t, 1, r.Total)
}

func TestSpecFrontmatter(t *testing.T) {
	tests := []struct {
		name    string
		sk      *skill.Skill
		wantErr bool
	}{
		{
			"valid",
			makeSpecSkill("my-skill", "A description", nil, ""),
			false,
		},
		{
			"missing frontmatter",
			&skill.Skill{},
			true,
		},
		{
			"missing name",
			makeSpecSkill("", "A description", nil, ""),
			true,
		},
		{
			"missing description",
			makeSpecSkill("my-skill", "", nil, ""),
			true,
		},
		{
			"whitespace description",
			makeSpecSkill("my-skill", "   \n  ", nil, ""),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &SpecResult{}
			specFrontmatter(tt.sk, r)
			if tt.wantErr {
				require.NotEmpty(t, r.Issues)
				require.Equal(t, "spec-frontmatter", r.Issues[0].Rule)
			} else {
				require.Empty(t, r.Issues)
				require.Equal(t, 1, r.Pass)
			}
		})
	}
}

func TestSpecAllowedFields(t *testing.T) {
	tests := []struct {
		name   string
		raw    map[string]any
		wantOK bool
	}{
		{
			"all allowed",
			map[string]any{"name": "x", "description": "y", "license": "MIT"},
			true,
		},
		{
			"unknown field",
			map[string]any{"name": "x", "description": "y", "author": "someone"},
			false,
		},
		{
			"nil raw",
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sk := &skill.Skill{FrontmatterRaw: tt.raw}
			r := &SpecResult{}
			specAllowedFields(sk, r)
			if tt.wantOK {
				hasError := false
				for _, iss := range r.Issues {
					if iss.Severity == "error" {
						hasError = true
					}
				}
				require.False(t, hasError)
			} else {
				require.NotEmpty(t, r.Issues)
				require.Equal(t, "spec-allowed-fields", r.Issues[0].Rule)
			}
		})
	}
}

func TestSpecName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid", "my-skill", false},
		{"valid single", "skill", false},
		{"valid with digits", "tool-v2", false},
		{"leading hyphen", "-my-skill", true},
		{"trailing hyphen", "my-skill-", true},
		{"consecutive hyphens", "my--skill", true},
		{"uppercase", "My-Skill", true},
		{"underscore", "my_skill", true},
		{"empty", "", false}, // empty handled by spec-frontmatter
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sk := makeSpecSkill(tt.input, "desc", nil, "")
			r := &SpecResult{}
			specName(sk, r)
			if tt.wantErr {
				require.NotEmpty(t, r.Issues, "expected error for name %q", tt.input)
				require.Equal(t, "spec-name", r.Issues[0].Rule)
			} else {
				require.Empty(t, r.Issues, "unexpected issues for name %q: %v", tt.input, r.Issues)
			}
		})
	}
}

func TestSpecDirMatch(t *testing.T) {
	tests := []struct {
		name    string
		skName  string
		path    string
		wantErr bool
	}{
		{"match", "my-skill", "/skills/my-skill/SKILL.md", false},
		{"mismatch", "my-skill", "/skills/other-skill/SKILL.md", true},
		{"no path", "my-skill", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sk := makeSpecSkill(tt.skName, "desc", nil, tt.path)
			r := &SpecResult{}
			specDirMatch(sk, r)
			if tt.wantErr {
				require.NotEmpty(t, r.Issues)
				require.Equal(t, "spec-dir-match", r.Issues[0].Rule)
			} else {
				require.Empty(t, r.Issues)
			}
		})
	}
}

func TestSpecDescription(t *testing.T) {
	tests := []struct {
		name    string
		desc    string
		wantErr bool
	}{
		{"valid", "A good description", false},
		{"empty", "", false}, // handled by spec-frontmatter
		{"too long", strings.Repeat("a", 1025), true},
		{"exactly 1024", strings.Repeat("a", 1024), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sk := makeSpecSkill("test", tt.desc, nil, "")
			r := &SpecResult{}
			specDescription(sk, r)
			if tt.wantErr {
				require.NotEmpty(t, r.Issues)
				require.Equal(t, "spec-description", r.Issues[0].Rule)
			} else {
				require.Empty(t, r.Issues)
			}
		})
	}
}

func TestSpecCompatibility(t *testing.T) {
	tests := []struct {
		name    string
		raw     map[string]any
		wantErr bool
	}{
		{"not present", map[string]any{"name": "x"}, false},
		{"valid", map[string]any{"name": "x", "compatibility": "Works with v1"}, false},
		{"too long", map[string]any{"name": "x", "compatibility": strings.Repeat("a", 501)}, true},
		{"exactly 500", map[string]any{"name": "x", "compatibility": strings.Repeat("a", 500)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sk := &skill.Skill{FrontmatterRaw: tt.raw}
			r := &SpecResult{}
			specCompatibility(sk, r)
			if tt.wantErr {
				require.NotEmpty(t, r.Issues)
				require.Equal(t, "spec-compatibility", r.Issues[0].Rule)
			} else {
				require.Empty(t, r.Issues)
			}
		})
	}
}

func TestSpecLicense(t *testing.T) {
	tests := []struct {
		name     string
		raw      map[string]any
		wantWarn bool
	}{
		{"has license", map[string]any{"license": "MIT"}, false},
		{"no license", map[string]any{"name": "x"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sk := &skill.Skill{FrontmatterRaw: tt.raw}
			r := &SpecResult{}
			specLicense(sk, r)
			if tt.wantWarn {
				require.NotEmpty(t, r.Issues)
				require.Equal(t, "spec-license", r.Issues[0].Rule)
				require.Equal(t, "warning", r.Issues[0].Severity)
			} else {
				require.Empty(t, r.Issues)
			}
		})
	}
}

func TestSpecVersion(t *testing.T) {
	tests := []struct {
		name     string
		raw      map[string]any
		wantWarn bool
	}{
		{
			"has metadata.version",
			map[string]any{"metadata": map[string]any{"version": "1.0.0"}},
			false,
		},
		{
			"no metadata",
			map[string]any{"name": "x"},
			true,
		},
		{
			"metadata without version",
			map[string]any{"metadata": map[string]any{"author": "me"}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sk := &skill.Skill{FrontmatterRaw: tt.raw}
			r := &SpecResult{}
			specVersion(sk, r)
			if tt.wantWarn {
				require.NotEmpty(t, r.Issues)
				require.Equal(t, "spec-version", r.Issues[0].Rule)
				require.Equal(t, "warning", r.Issues[0].Severity)
			} else {
				require.Empty(t, r.Issues)
			}
		})
	}
}

func TestSpecScorer_FullCompliant(t *testing.T) {
	sk := makeSpecSkill("my-skill", "A valid description for the skill", map[string]any{
		"license":  "MIT",
		"metadata": map[string]any{"version": "1.0.0"},
	}, "/skills/my-skill/SKILL.md")

	r := (SpecScorer{}).Score(sk)
	require.True(t, r.Passed(), "fully compliant skill should pass; issues: %v", r.Issues)
	require.Equal(t, r.Pass, r.Total)
}

func TestSpecScorer_MultipleIssues(t *testing.T) {
	sk := makeSpecSkill("-bad--name-", strings.Repeat("x", 1100), map[string]any{
		"author":        "someone",
		"compatibility": strings.Repeat("c", 600),
	}, "/skills/wrong-dir/SKILL.md")

	r := (SpecScorer{}).Score(sk)
	require.False(t, r.Passed())

	rules := map[string]bool{}
	for _, iss := range r.Issues {
		rules[iss.Rule] = true
	}
	require.True(t, rules["spec-name"], "should flag spec-name")
	require.True(t, rules["spec-dir-match"], "should flag spec-dir-match")
	require.True(t, rules["spec-description"], "should flag spec-description")
	require.True(t, rules["spec-compatibility"], "should flag spec-compatibility")
	require.True(t, rules["spec-allowed-fields"], "should flag spec-allowed-fields")
}

func TestSpecScorer_RealTestdataHigh(t *testing.T) {
	sk, err := readSkillFile("testdata/high/SKILL.md")
	require.NoError(t, err)
	r := (SpecScorer{}).Score(sk)
	// High testdata has dir-match error (testdata dir "high" != skill name "pdf-processor")
	for _, iss := range r.Issues {
		if iss.Rule == "spec-dir-match" {
			continue // expected: testdata dir doesn't match skill name
		}
		require.NotEqual(t, "error", iss.Severity,
			"high testdata should have no spec errors (except dir-match), got: %s — %s", iss.Rule, iss.Message)
	}
}
