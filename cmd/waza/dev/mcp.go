package dev

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spboyer/waza/internal/skill"
)

// builtinToolNames lists well-known Copilot CLI built-in tool names that
// skills should avoid shadowing.
var builtinToolNames = map[string]bool{
	"bash":     true,
	"edit":     true,
	"view":     true,
	"grep":     true,
	"glob":     true,
	"create":   true,
	"read":     true,
	"write":    true,
	"terminal": true,
	"search":   true,
	"fetch":    true,
	"run":      true,
}

// toolsTablePattern matches Markdown tables with at least a header row
// containing "tool" (case-insensitive) and a separator row.
var toolsTablePattern = regexp.MustCompile(`(?im)^\|.*tool.*\|\s*\n\|[\s\-:|]+\|`)

// McpResult holds the output from MCP integration checks.
type McpResult struct {
	HasInvokes           bool
	ToolsTablePresent    bool
	PrereqsDocumented    bool
	CliFallbackDescribed bool
	NameCollisions       []string
	Issues               []Issue
	SubScore             int // 0–4
}

// McpScorer evaluates MCP integration documentation quality for skills
// that declare INVOKES: in their description.
type McpScorer struct{}

// Score runs MCP integration checks. Returns nil if the skill does not
// use INVOKES: (i.e. not an MCP-dependent skill).
func (McpScorer) Score(sk *skill.Skill) *McpResult {
	if sk == nil {
		return nil
	}

	desc := sk.Frontmatter.Description
	body := sk.Body
	combined := desc + "\n" + body

	if !containsInvokes(desc) {
		return nil
	}

	r := &McpResult{HasInvokes: true}

	checkToolsTable(combined, r)
	checkPrerequisites(combined, r)
	checkCliFallback(combined, r)
	checkNameCollisions(desc, r)

	r.SubScore = computeMcpSubScore(r)
	return r
}

// containsInvokes returns true if the text contains "INVOKES:" (case-insensitive).
func containsInvokes(text string) bool {
	return strings.Contains(strings.ToLower(text), "invokes:")
}

// checkToolsTable looks for a Markdown table listing MCP tools.
func checkToolsTable(text string, r *McpResult) {
	if toolsTablePattern.MatchString(text) {
		r.ToolsTablePresent = true
		return
	}
	r.Issues = append(r.Issues, Issue{
		Rule:     "mcp-tools-table",
		Message:  "No MCP tools table found — add a Markdown table listing tools used",
		Severity: "warning",
	})
}

// prerequisitePatterns identifies documentation about MCP prerequisites.
var prerequisitePatterns = []string{
	"prerequisite",
	"requires mcp",
	"mcp server",
	"mcp servers needed",
	"required servers",
	"required mcp",
	"server setup",
	"install the",
}

// checkPrerequisites verifies that MCP prerequisites are documented.
func checkPrerequisites(text string, r *McpResult) {
	lower := strings.ToLower(text)
	for _, p := range prerequisitePatterns {
		if strings.Contains(lower, p) {
			r.PrereqsDocumented = true
			return
		}
	}
	r.Issues = append(r.Issues, Issue{
		Rule:     "mcp-prerequisites",
		Message:  "MCP prerequisites not documented — describe which MCP servers are needed",
		Severity: "warning",
	})
}

// fallbackPatterns identifies documentation about CLI fallback behavior.
var fallbackPatterns = []string{
	"fallback",
	"without mcp",
	"if mcp",
	"when mcp",
	"not available",
	"unavailable",
	"falls back",
	"graceful degradation",
	"cli alternative",
}

// checkCliFallback verifies that CLI fallback behavior is documented.
func checkCliFallback(text string, r *McpResult) {
	lower := strings.ToLower(text)
	for _, p := range fallbackPatterns {
		if strings.Contains(lower, p) {
			r.CliFallbackDescribed = true
			return
		}
	}
	r.Issues = append(r.Issues, Issue{
		Rule:     "mcp-cli-fallback",
		Message:  "No CLI fallback documented — describe behavior when MCP is unavailable",
		Severity: "warning",
	})
}

// checkNameCollisions extracts tool names from after "INVOKES:" and checks
// them against built-in tool names.
func checkNameCollisions(desc string, r *McpResult) {
	tools := extractInvokedTools(desc)
	for _, tool := range tools {
		if builtinToolNames[tool] {
			r.NameCollisions = append(r.NameCollisions, tool)
		}
	}
	if len(r.NameCollisions) > 0 {
		r.Issues = append(r.Issues, Issue{
			Rule:     "mcp-name-collision",
			Message:  fmt.Sprintf("Tool names conflict with built-ins: %s", strings.Join(r.NameCollisions, ", ")),
			Severity: "error",
		})
	}
}

// extractInvokedTools parses tool names from the INVOKES: section of a description.
// It expects comma-separated values after "INVOKES:" until a period, newline, or
// the next section marker.
func extractInvokedTools(desc string) []string {
	lower := strings.ToLower(desc)
	idx := strings.Index(lower, "invokes:")
	if idx < 0 {
		return nil
	}
	after := desc[idx+len("invokes:"):]

	// Stop at section markers or double newline.
	for _, stop := range []string{"FOR SINGLE OPERATIONS:", "DO NOT USE FOR:", "USE FOR:", "\n\n", "\n"} {
		if si := strings.Index(strings.ToUpper(after), strings.ToUpper(stop)); si >= 0 {
			after = after[:si]
		}
	}

	segments := strings.Split(after, ",")
	var tools []string
	for _, seg := range segments {
		candidate := strings.TrimSpace(seg)
		candidate = strings.TrimRight(candidate, ".")
		candidate = strings.Trim(candidate, "\"'`")
		// Take only the first word as the tool name.
		if parts := strings.Fields(candidate); len(parts) > 0 {
			tools = append(tools, strings.ToLower(parts[0]))
		}
	}
	return tools
}

// computeMcpSubScore returns 0-4 based on the four MCP checks.
func computeMcpSubScore(r *McpResult) int {
	score := 0
	if r.ToolsTablePresent {
		score++
	}
	if r.PrereqsDocumented {
		score++
	}
	if r.CliFallbackDescribed {
		score++
	}
	if len(r.NameCollisions) == 0 {
		score++
	}
	return score
}
