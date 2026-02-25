package checks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spboyer/waza/internal/skill"
	"github.com/stretchr/testify/require"
)

func TestModuleCountChecker(t *testing.T) {
	tmp := t.TempDir()
	refsDir := filepath.Join(tmp, "references")
	require.NoError(t, os.MkdirAll(refsDir, 0o755))

	tests := []struct {
		name    string
		mdCount int
		status  CheckStatus
		passed  bool
	}{
		{"zero modules", 0, StatusOK, true},
		{"one module", 1, StatusOK, true},
		{"two modules optimal", 2, StatusOptimal, true},
		{"three modules optimal", 3, StatusOptimal, true},
		{"four modules warning", 4, StatusWarning, false},
		{"five modules warning", 5, StatusWarning, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entries, _ := os.ReadDir(refsDir)
			for _, e := range entries {
				_ = os.Remove(filepath.Join(refsDir, e.Name()))
			}

			for i := range tt.mdCount {
				f := filepath.Join(refsDir, "ref"+string(rune('A'+i))+".md")
				require.NoError(t, os.WriteFile(f, []byte("# ref"), 0o644))
			}

			sk := skill.Skill{Path: filepath.Join(tmp, "SKILL.md")}
			checker := &ModuleCountChecker{}
			result, err := checker.Check(sk)
			require.NoError(t, err)
			require.Equal(t, tt.passed, result.Passed)
			data, ok := result.Data.(*ModuleCountData)
			require.True(t, ok)
			require.Equal(t, tt.status, data.Status)
			require.Equal(t, tt.mdCount, data.Count)
		})
	}
}

func TestComplexityChecker(t *testing.T) {
	tests := []struct {
		name           string
		tokens         int
		mdCount        int
		classification string
		status         CheckStatus
		passed         bool
	}{
		{"compact", 100, 0, "compact", StatusOK, true},
		{"detailed", 300, 2, "detailed", StatusOptimal, true},
		{"comprehensive by tokens", 600, 1, "comprehensive", StatusWarning, false},
		{"comprehensive by modules", 300, 4, "comprehensive", StatusWarning, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp := t.TempDir()
			refsDir := filepath.Join(tmp, "references")
			require.NoError(t, os.MkdirAll(refsDir, 0o755))
			for i := range tt.mdCount {
				f := filepath.Join(refsDir, "ref"+string(rune('A'+i))+".md")
				require.NoError(t, os.WriteFile(f, []byte("# ref"), 0o644))
			}

			sk := skill.Skill{
				Tokens: tt.tokens,
				Path:   filepath.Join(tmp, "SKILL.md"),
			}
			checker := &ComplexityChecker{}
			result, err := checker.Check(sk)
			require.NoError(t, err)
			require.Equal(t, tt.passed, result.Passed)
			data, ok := result.Data.(*ComplexityData)
			require.True(t, ok)
			require.Equal(t, tt.classification, data.Classification)
			require.Equal(t, tt.status, data.Status)
		})
	}
}

func TestNegativeDeltaRiskChecker(t *testing.T) {
	tests := []struct {
		name   string
		raw    string
		passed bool
		risks  int
	}{
		{
			name:   "clean content",
			raw:    "This is a normal skill description with no issues.",
			passed: true,
			risks:  0,
		},
		{
			name:   "conflicting paths",
			raw:    "Do X first. But alternatively you could do Y.",
			passed: false,
			risks:  1,
		},
		{
			name:   "duplicate step 1 blocks",
			raw:    "Step 1: Do this\nStep 2: Then this\nStep 1: Or start here",
			passed: false,
			risks:  1,
		},
		{
			name:   "excessive constraints",
			raw:    "You must not do A. Never do B. Always do C. It is forbidden to do D. Prohibited from E. You must not do F.",
			passed: false,
			risks:  1,
		},
		{
			name:   "multiple risks",
			raw:    "But alternatively try X.\nStep 1: first\nStep 1: second\nYou must not, never, always, forbidden, prohibited, must not do it.",
			passed: false,
			risks:  3,
		},
	}

	checker := &NegativeDeltaRiskChecker{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sk := skill.Skill{RawContent: tt.raw}
			result, err := checker.Check(sk)
			require.NoError(t, err)
			require.Equal(t, tt.passed, result.Passed)
			data, ok := result.Data.(*NegativeDeltaRiskData)
			require.True(t, ok)
			require.Len(t, data.Risks, tt.risks)
		})
	}
}

func TestProceduralContentChecker(t *testing.T) {
	tests := []struct {
		name   string
		desc   string
		passed bool
	}{
		{
			name:   "has action verb",
			desc:   "This skill extracts data from PDF files",
			passed: true,
		},
		{
			name:   "has procedure keyword",
			desc:   "A workflow for handling requests step by step",
			passed: true,
		},
		{
			name:   "no procedural language",
			desc:   "A general purpose tool for data",
			passed: false,
		},
		{
			name:   "empty description",
			desc:   "",
			passed: false,
		},
	}

	checker := &ProceduralContentChecker{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sk := skill.Skill{
				Frontmatter: skill.Frontmatter{Description: tt.desc},
			}
			result, err := checker.Check(sk)
			require.NoError(t, err)
			require.Equal(t, tt.passed, result.Passed, "desc: %q", tt.desc)
		})
	}
}

func TestOverSpecificityChecker(t *testing.T) {
	tests := []struct {
		name       string
		raw        string
		passed     bool
		categories int
	}{
		{
			name:       "clean content",
			raw:        "This skill processes markdown files and generates output.",
			passed:     true,
			categories: 0,
		},
		{
			name:       "unix path",
			raw:        "Files are stored in /usr/local/bin for access.",
			passed:     false,
			categories: 1,
		},
		{
			name:       "windows path",
			raw:        `Install to C:\Program Files\MyApp`,
			passed:     false,
			categories: 1,
		},
		{
			name:       "IP address",
			raw:        "Connect to 192.168.1.1 for the database.",
			passed:     false,
			categories: 1,
		},
		{
			name:       "hardcoded URL",
			raw:        "Download from https://example.com/releases/latest",
			passed:     false,
			categories: 1,
		},
		{
			name:       "doc URL allowed",
			raw:        "See https://github.com/owner/repo for details.",
			passed:     true,
			categories: 0,
		},
		{
			name:       "port number",
			raw:        "The server runs on :8080 by default.",
			passed:     false,
			categories: 1,
		},
		{
			name:       "multiple categories",
			raw:        "Use /home/user/app on 192.168.1.1:3000",
			passed:     false,
			categories: 3, // unix path, IP, port
		},
	}

	checker := &OverSpecificityChecker{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sk := skill.Skill{RawContent: tt.raw}
			result, err := checker.Check(sk)
			require.NoError(t, err)
			require.Equal(t, tt.passed, result.Passed)
			data, ok := result.Data.(*OverSpecificityData)
			require.True(t, ok)
			require.Len(t, data.Categories, tt.categories)
		})
	}
}
