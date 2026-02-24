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

func TestLoadConfig_WazaYAML(t *testing.T) {
	dir := t.TempDir()
	yaml := `tokens:
  limits:
    defaults:
      "*.md": 800
    overrides:
      "README.md": 1600
`
	err := os.WriteFile(filepath.Join(dir, ".waza.yaml"), []byte(yaml), 0644)
	require.NoError(t, err)

	cfg, err := LoadConfig(dir)
	require.NoError(t, err)
	require.Equal(t, 800, cfg.Defaults["*.md"])
	require.Equal(t, 1600, cfg.Overrides["README.md"])
}

func TestLoadConfig_JSONFallback(t *testing.T) {
	dir := t.TempDir()
	jsonData := `{"defaults":{"*.md":900},"overrides":{"A.md":1800}}`
	err := os.WriteFile(filepath.Join(dir, ".token-limits.json"), []byte(jsonData), 0644)
	require.NoError(t, err)
	// No .waza.yaml → should fall back to JSON file

	cfg, err := LoadConfig(dir)
	require.NoError(t, err)
	require.Equal(t, 900, cfg.Defaults["*.md"])
	require.Equal(t, 1800, cfg.Overrides["A.md"])
}

func TestLoadConfig_WazaYAMLWins(t *testing.T) {
	dir := t.TempDir()
	yaml := `tokens:
  limits:
    defaults:
      "*.md": 700
    overrides:
      "B.md": 1400
`
	err := os.WriteFile(filepath.Join(dir, ".waza.yaml"), []byte(yaml), 0644)
	require.NoError(t, err)

	jsonData := `{"defaults":{"*.md":999},"overrides":{"B.md":9999}}`
	err = os.WriteFile(filepath.Join(dir, ".token-limits.json"), []byte(jsonData), 0644)
	require.NoError(t, err)

	// Capture stderr for the note
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	cfg, err := LoadConfig(dir)
	require.NoError(t, err)

	w.Close()
	var buf [512]byte
	n, _ := r.Read(buf[:])
	os.Stderr = oldStderr

	// .waza.yaml values should win
	require.Equal(t, 700, cfg.Defaults["*.md"])
	require.Equal(t, 1400, cfg.Overrides["B.md"])

	// Note should be logged
	require.Contains(t, string(buf[:n]), "token limits loaded from .waza.yaml")
}

func TestLoadConfig_NeitherFile(t *testing.T) {
	dir := t.TempDir()
	// No .waza.yaml, no .token-limits.json → built-in defaults

	cfg, err := LoadConfig(dir)
	require.NoError(t, err)
	require.Equal(t, defaultLimits, cfg)
}
