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
