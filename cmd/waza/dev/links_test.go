package dev

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/spboyer/waza/internal/skill"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// makeLinkSkill creates a temp skill directory with a SKILL.md and optional
// extra files. Returns a parsed *skill.Skill ready for LinkScorer.Score().
func makeLinkSkill(t *testing.T, skillMD string, files map[string]string) *skill.Skill {
	t.Helper()
	dir := t.TempDir()
	skillPath := filepath.Join(dir, "SKILL.md")
	require.NoError(t, os.WriteFile(skillPath, []byte(skillMD), 0644))
	for relPath, content := range files {
		abs := filepath.Join(dir, relPath)
		require.NoError(t, os.MkdirAll(filepath.Dir(abs), 0755))
		require.NoError(t, os.WriteFile(abs, []byte(content), 0644))
	}
	var sk skill.Skill
	require.NoError(t, sk.UnmarshalText([]byte(skillMD)))
	sk.Path = skillPath
	return &sk
}

// --- Local Link Validation ---

func TestLinkScorer_ValidLocalLink(t *testing.T) {
	sk := makeLinkSkill(t, `---
name: test-skill
description: A test skill
---
See [guide](references/guide.md) for details.
`, map[string]string{
		"references/guide.md": "# Guide\nSome content.",
	})

	r := (&LinkScorer{}).Score(sk)
	require.NotNil(t, r)
	assert.Empty(t, r.BrokenLinks)
	assert.Equal(t, 1, r.TotalLinks)
	assert.Equal(t, 1, r.ValidLinks)
}

func TestLinkScorer_BrokenLocalLink(t *testing.T) {
	sk := makeLinkSkill(t, `---
name: test-skill
description: A test skill
---
See [missing](references/missing.md) for details.
`, nil)

	r := (&LinkScorer{}).Score(sk)
	require.NotNil(t, r)
	require.Len(t, r.BrokenLinks, 1)
	assert.Contains(t, r.BrokenLinks[0].Target, "missing.md")
}

func TestLinkScorer_LinkPointsToDirectory(t *testing.T) {
	sk := makeLinkSkill(t, `---
name: test-skill
description: A test skill
---
See [refs](references) for details.
`, map[string]string{
		"references/guide.md": "# Guide",
	})

	r := (&LinkScorer{}).Score(sk)
	require.NotNil(t, r)
	require.Len(t, r.DirectoryLinks, 1)
	assert.Contains(t, r.DirectoryLinks[0].Target, "references")
}

func TestLinkScorer_LinkEscapesSkillDir(t *testing.T) {
	sk := makeLinkSkill(t, `---
name: test-skill
description: A test skill
---
See [other](../../other-skill/file.md) for details.
`, nil)

	r := (&LinkScorer{}).Score(sk)
	require.NotNil(t, r)
	require.Len(t, r.ScopeEscapes, 1)
	assert.Contains(t, r.ScopeEscapes[0].Target, "../../other-skill")
}

func TestLinkScorer_FragmentOnlyLink(t *testing.T) {
	sk := makeLinkSkill(t, `---
name: test-skill
description: A test skill
---
See [section](#my-section) for more info.
`, nil)

	r := (&LinkScorer{}).Score(sk)
	require.NotNil(t, r)
	assert.True(t, r.Passed())
	assert.Equal(t, 0, r.TotalLinks, "fragment-only links are not counted")
}

func TestLinkScorer_MailtoLink(t *testing.T) {
	sk := makeLinkSkill(t, `---
name: test-skill
description: A test skill
---
Contact [us](mailto:team@example.com).
`, nil)

	r := (&LinkScorer{}).Score(sk)
	require.NotNil(t, r)
	assert.True(t, r.Passed())
}

func TestLinkScorer_MdcLink(t *testing.T) {
	sk := makeLinkSkill(t, `---
name: test-skill
description: A test skill
---
See [component](mdc:components/MyComponent).
`, nil)

	r := (&LinkScorer{}).Score(sk)
	require.NotNil(t, r)
	assert.True(t, r.Passed())
}

func TestLinkScorer_LinkWithFragment(t *testing.T) {
	sk := makeLinkSkill(t, `---
name: test-skill
description: A test skill
---
See [section](references/guide.md#advanced) for advanced usage.
`, map[string]string{
		"references/guide.md": "# Guide\n## Advanced\nAdvanced content.",
	})

	r := (&LinkScorer{}).Score(sk)
	require.NotNil(t, r)
	assert.Empty(t, r.BrokenLinks, "link with fragment should resolve after stripping #")
	assert.Equal(t, 1, r.ValidLinks)
}

func TestLinkScorer_LinkInCodeBlock(t *testing.T) {
	sk := makeLinkSkill(t, "---\nname: test-skill\ndescription: A test skill\n---\n"+
		"```markdown\n"+
		"See [fake link](nonexistent.md) in code.\n"+
		"```\n", nil)

	r := (&LinkScorer{}).Score(sk)
	require.NotNil(t, r)
	assert.Empty(t, r.BrokenLinks, "link inside code block should be ignored")
}

func TestLinkScorer_ImageLink(t *testing.T) {
	sk := makeLinkSkill(t, `---
name: test-skill
description: A test skill
---
![screenshot](references/screenshot.png)
`, map[string]string{
		"references/screenshot.png": "fake-png-content",
	})

	r := (&LinkScorer{}).Score(sk)
	require.NotNil(t, r)
	assert.Empty(t, r.BrokenLinks, "valid image link should not produce issues")
}

func TestLinkScorer_BrokenImageLink(t *testing.T) {
	sk := makeLinkSkill(t, `---
name: test-skill
description: A test skill
---
![screenshot](references/missing.png)
`, nil)

	r := (&LinkScorer{}).Score(sk)
	require.NotNil(t, r)
	require.Len(t, r.BrokenLinks, 1)
}

// --- Orphaned File Detection ---

func TestLinkScorer_AllReferencesLinked(t *testing.T) {
	sk := makeLinkSkill(t, `---
name: test-skill
description: A test skill
---
See [guide](references/guide.md) and [faq](references/faq.md).
`, map[string]string{
		"references/guide.md": "# Guide",
		"references/faq.md":   "# FAQ",
	})

	r := (&LinkScorer{}).Score(sk)
	require.NotNil(t, r)
	assert.Empty(t, r.OrphanedFiles)
}

func TestLinkScorer_OrphanedFile(t *testing.T) {
	sk := makeLinkSkill(t, `---
name: test-skill
description: A test skill
---
See [guide](references/guide.md).
`, map[string]string{
		"references/guide.md":    "# Guide",
		"references/orphaned.md": "# Nobody links here",
	})

	r := (&LinkScorer{}).Score(sk)
	require.NotNil(t, r)
	require.Len(t, r.OrphanedFiles, 1)
	assert.True(t, strings.Contains(r.OrphanedFiles[0], "orphaned.md"))
}

func TestLinkScorer_TransitiveLink(t *testing.T) {
	sk := makeLinkSkill(t, `---
name: test-skill
description: A test skill
---
See [a](references/a.md).
`, map[string]string{
		"references/a.md": "# A\nSee [b](b.md) for more.",
		"references/b.md": "# B\nDeep content.",
	})

	r := (&LinkScorer{}).Score(sk)
	require.NotNil(t, r)
	assert.Empty(t, r.OrphanedFiles, "transitively linked file should not be orphaned")
}

func TestLinkScorer_NoReferencesDir(t *testing.T) {
	sk := makeLinkSkill(t, `---
name: test-skill
description: A test skill
---
A simple skill with no references directory.
`, nil)

	r := (&LinkScorer{}).Score(sk)
	require.NotNil(t, r)
	assert.Empty(t, r.OrphanedFiles)
}

// --- External URL Validation ---

func TestLinkScorer_ValidExternalURL(t *testing.T) {
	skipSSRFCheck = true
	defer func() { skipSSRFCheck = false }()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	sk := makeLinkSkill(t, fmt.Sprintf(`---
name: test-skill
description: A test skill
---
See [docs](%s/page) for more.
`, srv.URL), nil)

	r := (&LinkScorer{}).Score(sk)
	require.NotNil(t, r)
	assert.Empty(t, r.DeadURLs)
}

func TestLinkScorer_URL404(t *testing.T) {
	skipSSRFCheck = true
	defer func() { skipSSRFCheck = false }()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	sk := makeLinkSkill(t, fmt.Sprintf(`---
name: test-skill
description: A test skill
---
See [broken](%s/missing) for more.
`, srv.URL), nil)

	r := (&LinkScorer{}).Score(sk)
	require.NotNil(t, r)
	require.Len(t, r.DeadURLs, 1)
	assert.Contains(t, r.DeadURLs[0].Target, srv.URL)
}

func TestLinkScorer_DuplicateURLsDeduped(t *testing.T) {
	skipSSRFCheck = true
	defer func() { skipSSRFCheck = false }()

	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	sk := makeLinkSkill(t, fmt.Sprintf(`---
name: test-skill
description: A test skill
---
See [one](%s/page) and [two](%s/page) — same URL.
`, srv.URL, srv.URL), nil)

	r := (&LinkScorer{}).Score(sk)
	require.NotNil(t, r)
	assert.Empty(t, r.DeadURLs)
	// HEAD + possible GET fallback = at most 2 requests for 1 unique URL.
	assert.LessOrEqual(t, callCount, 2)
}

func TestLinkScorer_NilSkill(t *testing.T) {
	r := (&LinkScorer{}).Score(nil)
	assert.Nil(t, r)
}

func TestLinkScorer_EmptyPath(t *testing.T) {
	sk := &skill.Skill{}
	r := (&LinkScorer{}).Score(sk)
	assert.Nil(t, r)
}

func TestLinkScorer_EmptySkillMD(t *testing.T) {
	sk := makeLinkSkill(t, `---
name: test-skill
description: A test skill
---
`, nil)

	r := (&LinkScorer{}).Score(sk)
	require.NotNil(t, r)
	assert.True(t, r.Passed())
	assert.Equal(t, 0, r.TotalLinks)
}

func TestLinkScorer_WindowsPathNormalization(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows-specific test")
	}
	sk := makeLinkSkill(t, `---
name: test-skill
description: A test skill
---
See [guide](references/guide.md) for details.
`, map[string]string{
		"references/guide.md": "# Guide",
	})

	r := (&LinkScorer{}).Score(sk)
	require.NotNil(t, r)
	assert.Empty(t, r.BrokenLinks)
}

func TestLinkScorer_Passed(t *testing.T) {
	r := &LinkResult{}
	assert.True(t, r.Passed())

	r.BrokenLinks = []LinkIssue{{Source: "a", Target: "b", Reason: "missing"}}
	assert.False(t, r.Passed())
}

func TestLinkScorer_SummaryIssues(t *testing.T) {
	r := &LinkResult{
		BrokenLinks:   []LinkIssue{{Source: "a", Target: "b", Reason: "missing"}},
		OrphanedFiles: []string{"references/orphan.md"},
	}
	buildLinkSummaryIssues(r)
	require.Len(t, r.Issues, 2)
	assert.Equal(t, "error", r.Issues[0].Severity)
	assert.Equal(t, "warning", r.Issues[1].Severity)
}
