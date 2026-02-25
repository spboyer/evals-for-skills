package dev

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spboyer/waza/internal/scoring"
	"github.com/stretchr/testify/require"
)

// the "want" strings in these tests are formatted with newlines and spacing to match the exact output of the
// display functions, including box-drawing characters and alignment. This ensures that the tests verify not
// just the content but also the formatting of the output. The literals may appear misaligned because they
// contain emojis whose apparent width depends on the font. A terminal monospaced font will probably render
// bounding boxes correctly.

func TestDisplayScore_LowScore(t *testing.T) {
	skill := makeSkill("my-skill", "Short description")
	result := (&scoring.HeuristicScorer{}).Score(skill)

	var buf bytes.Buffer
	DisplayScore(&buf, skill, result)

	want := `Skill: my-skill
Score: Low
Tokens: 14
Description: 17 chars
Triggers: 0
Anti-triggers: 0

Issues:
  ❌ Description is 17 chars (need 150+)

── Spec Compliance ──
  ✅ spec-frontmatter: Frontmatter structure valid with required fields
  ✅ spec-allowed-fields: All frontmatter fields are spec-allowed
  ✅ spec-name: Name follows spec naming rules
  ✅ spec-dir-match: Cannot validate (missing path or name)
  ✅ spec-description: Description is valid
  ✅ spec-compatibility: No compatibility field (optional)
  ⚠️ spec-license: No license field found
     📎 Best practice: include a license field (e.g., MIT, Apache-2.0)
  ⚠️ spec-version: No metadata.version field found
     📎 Best practice: include metadata.version for tracking and compatibility

── Advisory Checks ──
  ✅ module-count: Found 0 reference module(s)
  ✅ complexity: Complexity: compact (14 tokens, 0 modules)
  ✅ negative-delta-risk: No negative delta risk patterns detected
  ⚠️ procedural-content: Description lacks procedural language (no action verbs or procedure keywords found)
  ✅ over-specificity: No over-specificity patterns detected
`
	require.Equal(t, want, buf.String())
}

func TestDisplayScore_HighScore(t *testing.T) {
	skill := makeSkill("pdf-processor",
		`**WORKFLOW SKILL** - Process PDF files including text extraction.
USE FOR: "extract PDF text", "rotate PDF".
DO NOT USE FOR: creating PDFs (use document-creator).
INVOKES: pdf-tools MCP for extraction.
FOR SINGLE OPERATIONS: Use pdf-tools directly.`)
	result := (&scoring.HeuristicScorer{}).Score(skill)

	var buf bytes.Buffer
	DisplayScore(&buf, skill, result)

	want := `Skill: pdf-processor
Score: High
Tokens: 73
Description: 248 chars
Triggers: 2
Anti-triggers: 1

── Spec Compliance ──
  ✅ spec-frontmatter: Frontmatter structure valid with required fields
  ✅ spec-allowed-fields: All frontmatter fields are spec-allowed
  ✅ spec-name: Name follows spec naming rules
  ✅ spec-dir-match: Cannot validate (missing path or name)
  ✅ spec-description: Description is valid
  ✅ spec-compatibility: No compatibility field (optional)
  ⚠️ spec-license: No license field found
     📎 Best practice: include a license field (e.g., MIT, Apache-2.0)
  ⚠️ spec-version: No metadata.version field found
     📎 Best practice: include metadata.version for tracking and compatibility

MCP Integration: 1/4
  ⚠️ [mcp-tools-table] No MCP tools table found — add a Markdown table listing tools used
  ⚠️ [mcp-prerequisites] MCP prerequisites not documented — describe which MCP servers are needed
  ⚠️ [mcp-cli-fallback] No CLI fallback documented — describe behavior when MCP is unavailable

── Advisory Checks ──
  ✅ module-count: Found 0 reference module(s)
  ✅ complexity: Complexity: compact (73 tokens, 0 modules)
  ✅ negative-delta-risk: No negative delta risk patterns detected
  ✅ procedural-content: Description contains procedural language
  ✅ over-specificity: No over-specificity patterns detected
`
	require.Equal(t, want, buf.String())
}

func TestDisplayScore_ShowsTriggerCounts(t *testing.T) {
	skill := makeSkill("counter-test",
		`Process things with great care and attention to detail over many steps.
USE FOR: "process data", "transform files", "validate input".
DO NOT USE FOR: deleting files (use file-manager).`)
	result := (&scoring.HeuristicScorer{}).Score(skill)

	var buf bytes.Buffer
	DisplayScore(&buf, skill, result)

	want := `Skill: counter-test
Score: Medium-High
Tokens: 57
Description: 184 chars
Triggers: 3
Anti-triggers: 1

── Spec Compliance ──
  ✅ spec-frontmatter: Frontmatter structure valid with required fields
  ✅ spec-allowed-fields: All frontmatter fields are spec-allowed
  ✅ spec-name: Name follows spec naming rules
  ✅ spec-dir-match: Cannot validate (missing path or name)
  ✅ spec-description: Description is valid
  ✅ spec-compatibility: No compatibility field (optional)
  ⚠️ spec-license: No license field found
     📎 Best practice: include a license field (e.g., MIT, Apache-2.0)
  ⚠️ spec-version: No metadata.version field found
     📎 Best practice: include metadata.version for tracking and compatibility

── Advisory Checks ──
  ✅ module-count: Found 0 reference module(s)
  ✅ complexity: Complexity: compact (57 tokens, 0 modules)
  ✅ negative-delta-risk: No negative delta risk patterns detected
  ✅ procedural-content: Description contains procedural language
  ✅ over-specificity: No over-specificity patterns detected
`
	require.Equal(t, want, buf.String())
}

func TestDisplayIssues_ShowsAllIssues(t *testing.T) {
	issues := []scoring.Issue{
		{Rule: "description-length", Message: "Description is 30 chars (need 150+)", Severity: "error"},
		{Rule: "name-format", Message: "Name must be lowercase", Severity: "error"},
		{Rule: "token-soft-limit", Message: "Over soft limit", Severity: "warning"},
	}

	var buf bytes.Buffer
	DisplayIssues(&buf, issues)

	want := `Issues:
  ❌ Description is 30 chars (need 150+)
  ❌ Name must be lowercase
  ⚠️ Over soft limit
`
	require.Equal(t, want, buf.String())
}

func TestDisplayIssues_ErrorIcon(t *testing.T) {
	issues := []scoring.Issue{
		{Rule: "test", Message: "An error issue", Severity: "error"},
	}

	var buf bytes.Buffer
	DisplayIssues(&buf, issues)

	want := `Issues:
  ❌ An error issue
`
	require.Equal(t, want, buf.String())
}

func TestDisplayIssues_WarningIcon(t *testing.T) {
	issues := []scoring.Issue{
		{Rule: "test", Message: "A warning issue", Severity: "warning"},
	}

	var buf bytes.Buffer
	DisplayIssues(&buf, issues)

	want := `Issues:
  ⚠️ A warning issue
`
	require.Equal(t, want, buf.String())
}

func TestDisplaySummary_BoxFormat(t *testing.T) {
	before := &scoring.ScoreResult{
		Level:            scoring.AdherenceLow,
		TriggerCount:     0,
		AntiTriggerCount: 0,
	}
	after := &scoring.ScoreResult{
		Level:            scoring.AdherenceMediumHigh,
		TriggerCount:     5,
		AntiTriggerCount: 3,
	}

	var buf bytes.Buffer
	DisplaySummary(&buf, "my-skill", before, after, 142, 385)

	want := `╔══════════════════════════════════════════════════════════════════╗
║  SENSEI SUMMARY: my-skill                                        ║
╠══════════════════════════════════════════════════════════════════╣
║  BEFORE                          AFTER                           ║
║  ──────                          ─────                           ║
║  Score: Low                      Score: Medium-High              ║
║  Tokens: 142                     Tokens: 385                     ║
║  Triggers: 0                     Triggers: 5                     ║
║  Anti-triggers: 0                Anti-triggers: 3                ║
║                                                                  ║
║  TOKEN STATUS: ✅ Under budget (385 < 500)                        ║
╚══════════════════════════════════════════════════════════════════╝
`
	require.Equal(t, want, buf.String())
}

func TestDisplaySummary_ContainsBoxCharacters(t *testing.T) {
	before := &scoring.ScoreResult{Level: scoring.AdherenceLow}
	after := &scoring.ScoreResult{Level: scoring.AdherenceMedium}

	var buf bytes.Buffer
	DisplaySummary(&buf, "box-test", before, after, 100, 200)

	want := `╔══════════════════════════════════════════════════════════════════╗
║  SENSEI SUMMARY: box-test                                        ║
╠══════════════════════════════════════════════════════════════════╣
║  BEFORE                          AFTER                           ║
║  ──────                          ─────                           ║
║  Score: Low                      Score: Medium                   ║
║  Tokens: 100                     Tokens: 200                     ║
║  Triggers: 0                     Triggers: 0                     ║
║  Anti-triggers: 0                Anti-triggers: 0                ║
║                                                                  ║
║  TOKEN STATUS: ✅ Under budget (200 < 500)                        ║
╚══════════════════════════════════════════════════════════════════╝
`
	require.Equal(t, want, buf.String())
}

func TestDisplaySummary_TokenStatus_UnderBudget(t *testing.T) {
	before := &scoring.ScoreResult{Level: scoring.AdherenceLow}
	after := &scoring.ScoreResult{Level: scoring.AdherenceMediumHigh}

	var buf bytes.Buffer
	DisplaySummary(&buf, "token-test", before, after, 100, 385)

	want := `╔══════════════════════════════════════════════════════════════════╗
║  SENSEI SUMMARY: token-test                                      ║
╠══════════════════════════════════════════════════════════════════╣
║  BEFORE                          AFTER                           ║
║  ──────                          ─────                           ║
║  Score: Low                      Score: Medium-High              ║
║  Tokens: 100                     Tokens: 385                     ║
║  Triggers: 0                     Triggers: 0                     ║
║  Anti-triggers: 0                Anti-triggers: 0                ║
║                                                                  ║
║  TOKEN STATUS: ✅ Under budget (385 < 500)                        ║
╚══════════════════════════════════════════════════════════════════╝
`
	require.Equal(t, want, buf.String())
}

func TestDisplaySummary_TokenStatus_OverSoftLimit(t *testing.T) {
	before := &scoring.ScoreResult{Level: scoring.AdherenceLow}
	after := &scoring.ScoreResult{Level: scoring.AdherenceMediumHigh}

	var buf bytes.Buffer
	DisplaySummary(&buf, "token-test", before, after, 100, 600)

	want := `╔══════════════════════════════════════════════════════════════════╗
║  SENSEI SUMMARY: token-test                                      ║
╠══════════════════════════════════════════════════════════════════╣
║  BEFORE                          AFTER                           ║
║  ──────                          ─────                           ║
║  Score: Low                      Score: Medium-High              ║
║  Tokens: 100                     Tokens: 600                     ║
║  Triggers: 0                     Triggers: 0                     ║
║  Anti-triggers: 0                Anti-triggers: 0                ║
║                                                                  ║
║  TOKEN STATUS: ⚠️ Over soft limit (600 > 500)                    ║
╚══════════════════════════════════════════════════════════════════╝
`
	require.Equal(t, want, buf.String())
}

func TestDisplaySummary_TokenStatus_OverHardLimit(t *testing.T) {
	before := &scoring.ScoreResult{Level: scoring.AdherenceLow}
	after := &scoring.ScoreResult{Level: scoring.AdherenceMediumHigh}

	var buf bytes.Buffer
	DisplaySummary(&buf, "token-test", before, after, 100, 6000)

	want := `╔══════════════════════════════════════════════════════════════════╗
║  SENSEI SUMMARY: token-test                                      ║
╠══════════════════════════════════════════════════════════════════╣
║  BEFORE                          AFTER                           ║
║  ──────                          ─────                           ║
║  Score: Low                      Score: Medium-High              ║
║  Tokens: 100                     Tokens: 6000                    ║
║  Triggers: 0                     Triggers: 0                     ║
║  Anti-triggers: 0                Anti-triggers: 0                ║
║                                                                  ║
║  TOKEN STATUS: ❌ Over hard limit (6000 > 5000)                   ║
╚══════════════════════════════════════════════════════════════════╝
`
	require.Equal(t, want, buf.String())
}

func TestDisplaySummary_ShowsBEFOREandAFTER(t *testing.T) {
	before := &scoring.ScoreResult{Level: scoring.AdherenceLow, TriggerCount: 0, AntiTriggerCount: 0}
	after := &scoring.ScoreResult{Level: scoring.AdherenceHigh, TriggerCount: 7, AntiTriggerCount: 2}

	var buf bytes.Buffer
	DisplaySummary(&buf, "summary-test", before, after, 50, 400)

	want := `╔══════════════════════════════════════════════════════════════════╗
║  SENSEI SUMMARY: summary-test                                    ║
╠══════════════════════════════════════════════════════════════════╣
║  BEFORE                          AFTER                           ║
║  ──────                          ─────                           ║
║  Score: Low                      Score: High                     ║
║  Tokens: 50                      Tokens: 400                     ║
║  Triggers: 0                     Triggers: 7                     ║
║  Anti-triggers: 0                Anti-triggers: 2                ║
║                                                                  ║
║  TOKEN STATUS: ✅ Under budget (400 < 500)                        ║
╚══════════════════════════════════════════════════════════════════╝
`
	require.Equal(t, want, buf.String())
}

func TestDisplayIterationHeader(t *testing.T) {
	var buf bytes.Buffer
	DisplayIterationHeader(&buf, 2, 5)

	want := `
── Iteration 2/5 ──────────────────────────────────────────

`
	require.Equal(t, want, buf.String())
}

func TestDisplayTargetReached(t *testing.T) {
	var buf bytes.Buffer
	DisplayTargetReached(&buf, scoring.AdherenceHigh)

	want := `
✅ Target adherence level High reached!
`
	require.Equal(t, want, buf.String())
}

func TestDisplayMaxIterations(t *testing.T) {
	var buf bytes.Buffer
	DisplayMaxIterations(&buf, scoring.AdherenceMedium)

	want := `
⏱️  Max iterations reached. Current level: Medium
`
	require.Equal(t, want, buf.String())
}

func TestBoxLine_TruncatesLongText(t *testing.T) {
	line := boxLine(strings.Repeat("a", boxWidth))
	require.Equal(t, boxWidth+2, len([]rune(line)))
	require.Contains(t, line, "...")
}
