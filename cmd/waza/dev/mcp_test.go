package dev

import (
	"strings"
	"testing"

	"github.com/spboyer/waza/internal/skill"
	"github.com/stretchr/testify/require"
)

func TestMcpScorer_NilSkill(t *testing.T) {
	r := (McpScorer{}).Score(nil)
	require.Nil(t, r, "nil skill should return nil result")
}

func TestMcpScorer_NoInvokes(t *testing.T) {
	sk := makeSkill("my-tool", "A simple tool. USE FOR: things. DO NOT USE FOR: other things.")
	r := (McpScorer{}).Score(sk)
	require.Nil(t, r, "skill without INVOKES: should return nil")
}

func TestMcpScorer_HasInvokes_MinimalDescription(t *testing.T) {
	sk := makeSkill("mcp-tool", "A tool that does things. INVOKES: some-tool, another-tool.")
	r := (McpScorer{}).Score(sk)
	require.NotNil(t, r)
	require.True(t, r.HasInvokes)
	// No tools table, no prereqs, no fallback → sub-score = 1 (only no-collision passes)
	require.False(t, r.ToolsTablePresent)
	require.False(t, r.PrereqsDocumented)
	require.False(t, r.CliFallbackDescribed)
	require.Empty(t, r.NameCollisions)
	require.Equal(t, 1, r.SubScore)
}

func TestMcpScorer_FullScore(t *testing.T) {
	desc := "A workflow skill. INVOKES: pdf-tools, format-engine."
	body := `## MCP Tools

| Tool | Purpose |
|------|---------|
| pdf-tools | Extract text from PDFs |
| format-engine | Format documents |

## Prerequisites

Requires MCP server pdf-tools running locally.

## Fallback

Without MCP, falls back to basic text extraction.
`
	sk := &skill.Skill{
		Frontmatter: skill.Frontmatter{
			Name:        "mcp-full",
			Description: desc,
		},
		Body: body,
	}
	r := (McpScorer{}).Score(sk)
	require.NotNil(t, r)
	require.True(t, r.HasInvokes)
	require.True(t, r.ToolsTablePresent)
	require.True(t, r.PrereqsDocumented)
	require.True(t, r.CliFallbackDescribed)
	require.Empty(t, r.NameCollisions)
	require.Equal(t, 4, r.SubScore)
	require.Empty(t, r.Issues)
}

func TestMcpScorer_NameCollision(t *testing.T) {
	sk := makeSkill("bad-tool", "Does stuff. INVOKES: bash, grep, custom-tool.")
	r := (McpScorer{}).Score(sk)
	require.NotNil(t, r)
	require.Contains(t, r.NameCollisions, "bash")
	require.Contains(t, r.NameCollisions, "grep")
	require.Equal(t, 0, r.SubScore) // collision means this check fails: max 3, but also no table/prereqs/fallback = 0

	// Verify the error issue is present
	var found bool
	for _, iss := range r.Issues {
		if iss.Rule == "mcp-name-collision" {
			found = true
			require.Equal(t, "error", iss.Severity)
		}
	}
	require.True(t, found, "expected mcp-name-collision issue")
}

func TestMcpScorer_ToolsTable(t *testing.T) {
	tests := []struct {
		name string
		text string
		want bool
	}{
		{
			"markdown table with Tool header",
			"| Tool | Description |\n|------|-------------|\n| my-tool | does stuff |",
			true,
		},
		{
			"no table",
			"INVOKES: my-tool. This uses some tools.",
			false,
		},
		{
			"table without tool column",
			"| Name | Description |\n|------|-------------|\n| foo | bar |",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, toolsTablePattern.MatchString(tt.text))
		})
	}
}

func TestMcpScorer_Prerequisites(t *testing.T) {
	tests := []struct {
		name string
		text string
		want bool
	}{
		{"explicit prerequisite", "## Prerequisites\nRequires MCP server.", true},
		{"requires mcp", "This requires mcp tools to run.", true},
		{"mcp server mention", "You need an MCP server running.", true},
		{"no prereqs", "This is a normal skill.", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &McpResult{}
			checkPrerequisites(tt.text, r)
			require.Equal(t, tt.want, r.PrereqsDocumented)
		})
	}
}

func TestMcpScorer_CliFallback(t *testing.T) {
	tests := []struct {
		name string
		text string
		want bool
	}{
		{"fallback mentioned", "Falls back to CLI when MCP unavailable.", true},
		{"without mcp", "Without MCP, uses basic processing.", true},
		{"graceful degradation", "Supports graceful degradation.", true},
		{"no fallback info", "INVOKES: my-tool for processing.", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &McpResult{}
			checkCliFallback(tt.text, r)
			require.Equal(t, tt.want, r.CliFallbackDescribed)
		})
	}
}

func TestExtractInvokedTools(t *testing.T) {
	tests := []struct {
		name string
		desc string
		want []string
	}{
		{
			"simple list",
			"INVOKES: pdf-tools, format-engine.",
			[]string{"pdf-tools", "format-engine"},
		},
		{
			"with descriptive text",
			"INVOKES: pdf-tools for extraction, format-engine for output.",
			[]string{"pdf-tools", "format-engine"},
		},
		{
			"stops at section boundary",
			"INVOKES: my-tool. FOR SINGLE OPERATIONS: use directly.",
			[]string{"my-tool"},
		},
		{
			"no invokes",
			"Just a regular description.",
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractInvokedTools(tt.desc)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestContainsInvokes(t *testing.T) {
	require.True(t, containsInvokes("INVOKES: tools"))
	require.True(t, containsInvokes("invokes: tools"))
	require.False(t, containsInvokes("no invocation here"))
}

func TestComputeMcpSubScore(t *testing.T) {
	tests := []struct {
		name string
		r    McpResult
		want int
	}{
		{"all pass", McpResult{ToolsTablePresent: true, PrereqsDocumented: true, CliFallbackDescribed: true}, 4},
		{"no collisions only", McpResult{}, 1},
		{"table + no collision", McpResult{ToolsTablePresent: true}, 2},
		{"has collision", McpResult{NameCollisions: []string{"bash"}}, 0},
		{"all but collision", McpResult{ToolsTablePresent: true, PrereqsDocumented: true, CliFallbackDescribed: true, NameCollisions: []string{"edit"}}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, computeMcpSubScore(&tt.r))
		})
	}
}

func TestMcpScorer_WazaSkill(t *testing.T) {
	sk, err := readSkillFile("../../../skills/waza/SKILL.md")
	require.NoError(t, err)
	r := (McpScorer{}).Score(sk)
	// Waza skill has INVOKES: in its description
	require.NotNil(t, r, "waza skill has INVOKES: — should produce result")
	require.True(t, r.HasInvokes)
	require.Empty(t, r.NameCollisions, "waza skill should not have name collisions")
}

func TestMcpScorer_IssueCount(t *testing.T) {
	// Skill with INVOKES but missing all documentation
	sk := makeSkill("bare-tool", "INVOKES: custom-tool.")
	r := (McpScorer{}).Score(sk)
	require.NotNil(t, r)
	// Should have warnings for: tools-table, prerequisites, cli-fallback
	warnCount := 0
	for _, iss := range r.Issues {
		if iss.Severity == "warning" {
			warnCount++
		}
	}
	require.Equal(t, 3, warnCount, "expected 3 warning issues for missing MCP docs")
}

func TestMcpScorer_PartialDocs(t *testing.T) {
	desc := strings.Repeat("MCP tool for processing. ", 8) +
		"INVOKES: processor-tool. Requires MCP server processor-tool."
	sk := makeSkill("partial-tool", desc)
	r := (McpScorer{}).Score(sk)
	require.NotNil(t, r)
	require.True(t, r.PrereqsDocumented)
	require.False(t, r.ToolsTablePresent)
	require.False(t, r.CliFallbackDescribed)
	require.Equal(t, 2, r.SubScore) // prereqs + no collision
}
