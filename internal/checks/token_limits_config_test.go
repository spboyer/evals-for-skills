package checks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGlobToRegex(t *testing.T) {
	tests := []struct {
		pattern string
		match   []string
		noMatch []string
	}{
		{
			pattern: "*.md",
			match:   []string{"README.md", "foo/bar.md", "/a/b/c.md"},
			noMatch: []string{"README.txt", "md", "README.md.bak"},
		},
		{
			pattern: "**/*.md",
			match:   []string{"docs/foo.md", "a/b/c.md", "/x/y.md"},
			noMatch: []string{"README.txt"},
		},
		{
			pattern: "references/**/*.md",
			match:   []string{"references/sub/two.md"},
			noMatch: []string{"refs/one.md", "references_extra/one.md", "references/one.md", "x/references/deep/f.md"},
		},
		{
			pattern: "docs/*.md",
			match:   []string{"docs/guide.md"},
			noMatch: []string{"docs/sub/guide.md", "mydocs/guide.md", "/root/docs/guide.md"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.pattern, func(t *testing.T) {
			re, err := globToRegex(tc.pattern)
			require.NoError(t, err)
			for _, m := range tc.match {
				require.True(t, re.MatchString(m), "%q should match %q", tc.pattern, m)
			}
			for _, m := range tc.noMatch {
				require.False(t, re.MatchString(m), "%q should not match %q", tc.pattern, m)
			}
		})
	}
}

func TestGlobToRegex_PatternTooLong(t *testing.T) {
	long := strings.Repeat("a", maxPatternLength+1)
	_, err := globToRegex(long)
	require.ErrorContains(t, err, "pattern too long")
}

func TestMatchesPattern(t *testing.T) {
	tests := []struct {
		filePath string
		pattern  string
		want     bool
	}{
		{"SKILL.md", "SKILL.md", true},
		{"sub/SKILL.md", "SKILL.md", true},
		{"README.md", "SKILL.md", false},

		{"foo.md", "*.md", true},
		{"sub/foo.md", "*.md", true},
		{"foo.txt", "*.md", false},

		{"references/sub/two.md", "references/**/*.md", true},
		{"references/one.md", "references/**/*.md", false},
		{"other/one.md", "references/**/*.md", false},

		{"docs/guide.md", "docs/*.md", true},
		{"docs/sub/guide.md", "docs/*.md", false},

		{`docs\guide.md`, "docs/*.md", true},
	}

	for _, tc := range tests {
		name := tc.filePath + " ~ " + tc.pattern
		t.Run(name, func(t *testing.T) {
			got := matchesPattern(tc.filePath, tc.pattern)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestPatternSpecificity(t *testing.T) {
	require.Greater(t, patternSpecificity("SKILL.md"), patternSpecificity("*.md"))
	require.Greater(t, patternSpecificity("docs/*.md"), patternSpecificity("*.md"))
	require.Greater(t, patternSpecificity("a/b/*.md"), patternSpecificity("a/*.md"))
}

func TestLoadLimitsConfig_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, ".token-limits.json"), []byte(`{not json}`), 0644)
	require.NoError(t, err)

	_, err = LoadLimitsConfig(dir)
	require.ErrorContains(t, err, "error parsing limits")
}

func TestLoadLimitsConfig_MissingDefaults(t *testing.T) {
	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, ".token-limits.json"), []byte(`{"overrides":{"a.md":1}}`), 0644)
	require.NoError(t, err)

	_, err = LoadLimitsConfig(dir)
	require.ErrorContains(t, err, `missing or invalid "defaults"`)
}

func TestLoadLimitsConfig_NoFile(t *testing.T) {
	dir := t.TempDir()

	cfg, err := LoadLimitsConfig(dir)
	require.NoError(t, err)
	require.Equal(t, DefaultLimits, cfg)
}

// TestGetLimitForFile_WorkspaceRelPrefix verifies that workspace-root-relative
// patterns (e.g. "plugin/skills/**/SKILL.md") match when a workspace prefix is
// supplied, even though the file path itself is skill-directory-relative.
func TestGetLimitForFile_WorkspaceRelPrefix(t *testing.T) {
	cfg := TokenLimitsConfig{
		Defaults: map[string]int{
			"plugin/skills/**/SKILL.md":        1000,
			"plugin/skills/**/references/*.md": 800,
			"*.md":                             2000,
		},
		Overrides: map[string]int{
			"plugin/skills/special/OVERRIDE.md": 42,
		},
	}

	// Without prefix: SKILL.md only matches *.md (2000).
	lr := GetLimitForFile("SKILL.md", cfg)
	require.Equal(t, 2000, lr.Limit)
	require.Equal(t, "*.md", lr.Pattern)

	// With prefix: SKILL.md → plugin/skills/azure-deploy/SKILL.md → matches 1000.
	lr = GetLimitForFile("SKILL.md", cfg, "plugin/skills/azure-deploy")
	require.Equal(t, 1000, lr.Limit)
	require.Equal(t, "plugin/skills/**/SKILL.md", lr.Pattern)

	// Nested reference file with prefix.
	lr = GetLimitForFile("references/doc.md", cfg, "plugin/skills/azure-deploy")
	require.Equal(t, 800, lr.Limit)
	require.Equal(t, "plugin/skills/**/references/*.md", lr.Pattern)

	// Override matched via prefixed path.
	lr = GetLimitForFile("OVERRIDE.md", cfg, "plugin/skills/special")
	require.Equal(t, 42, lr.Limit)

	// Empty prefix behaves like no prefix.
	lr = GetLimitForFile("SKILL.md", cfg, "")
	require.Equal(t, 2000, lr.Limit)

	// Skill-relative match still wins when pattern matches without prefix.
	cfgLocal := TokenLimitsConfig{
		Defaults: map[string]int{
			"SKILL.md":                  500,
			"plugin/skills/**/SKILL.md": 1000,
		},
		Overrides: map[string]int{},
	}
	lr = GetLimitForFile("SKILL.md", cfgLocal, "plugin/skills/azure-deploy")
	require.Equal(t, 500, lr.Limit, "skill-relative match should take precedence")
}

// Note: .waza.yaml integration tests for token limits live in
// cmd/waza/tokens/internal/limits_test.go, not here.
