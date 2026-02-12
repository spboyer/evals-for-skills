package dev

import (
	"fmt"
	"io"
	"strings"

	"github.com/spboyer/waza/internal/skill"
)

const boxWidth = 66

// DisplayIterationHeader shows iteration progress.
func DisplayIterationHeader(w io.Writer, iteration, maxIterations int) {
	fmt.Fprintf(w, "\n── Iteration %d/%d ──────────────────────────────────────────\n\n", iteration, maxIterations)
}

// DisplayScore shows the current score with issues.
func DisplayScore(w io.Writer, sk *skill.Skill, score *ScoreResult) {
	name := sk.Frontmatter.Name
	fmt.Fprintf(w, "Skill: %s\n", name)
	fmt.Fprintf(w, "Score: %s\n", score.Level)
	fmt.Fprintf(w, "Tokens: %d\n", sk.Tokens)
	fmt.Fprintf(w, "Description: %d chars\n", score.DescriptionLen)
	fmt.Fprintf(w, "Triggers: %d\n", score.TriggerCount)
	fmt.Fprintf(w, "Anti-triggers: %d\n", score.AntiTriggerCount)

	if len(score.Issues) > 0 {
		fmt.Fprintln(w)
		DisplayIssues(w, score.Issues)
	}
}

// DisplayIssues lists all issues found.
func DisplayIssues(w io.Writer, issues []Issue) {
	fmt.Fprintln(w, "Issues:")
	for _, iss := range issues {
		icon := "⚠️"
		if iss.Severity == "error" {
			icon = "❌"
		}
		fmt.Fprintf(w, "  %s %s\n", icon, iss.Message)
	}
}

// DisplaySummary shows before/after comparison box.
func DisplaySummary(w io.Writer, skillName string, before, after *ScoreResult, beforeTokens, afterTokens int) {
	top := "╔" + strings.Repeat("═", boxWidth) + "╗"
	mid := "╠" + strings.Repeat("═", boxWidth) + "╣"
	bot := "╚" + strings.Repeat("═", boxWidth) + "╝"

	fmt.Fprintln(w, top)
	fmt.Fprintln(w, boxLine(fmt.Sprintf("SENSEI SUMMARY: %s", skillName)))
	fmt.Fprintln(w, mid)
	fmt.Fprintln(w, boxLine("BEFORE                          AFTER"))
	fmt.Fprintln(w, boxLine("──────                          ─────"))
	fmt.Fprintln(w, boxLine(fmt.Sprintf("Score: %-24s Score: %s", before.Level, after.Level)))
	fmt.Fprintln(w, boxLine(fmt.Sprintf("Tokens: %-23d Tokens: %d", beforeTokens, afterTokens)))
	fmt.Fprintln(w, boxLine(fmt.Sprintf("Triggers: %-21d Triggers: %d", before.TriggerCount, after.TriggerCount)))
	fmt.Fprintln(w, boxLine(fmt.Sprintf("Anti-triggers: %-16d Anti-triggers: %d", before.AntiTriggerCount, after.AntiTriggerCount)))
	fmt.Fprintln(w, boxLine(""))

	tokenStatus := fmt.Sprintf("TOKEN STATUS: ✅ Under budget (%d < %d)", afterTokens, tokenSoftLimit)
	if afterTokens > tokenSoftLimit {
		tokenStatus = fmt.Sprintf("TOKEN STATUS: ⚠️ Over soft limit (%d > %d)", afterTokens, tokenSoftLimit)
	}
	if afterTokens > tokenHardLimit {
		tokenStatus = fmt.Sprintf("TOKEN STATUS: ❌ Over hard limit (%d > %d)", afterTokens, tokenHardLimit)
	}
	fmt.Fprintln(w, boxLine(tokenStatus))
	fmt.Fprintln(w, bot)
}

// boxLine pads text inside box borders (║ ... ║).
func boxLine(text string) string {
	maxText := boxWidth - 2
	text = truncateText(text, maxText)
	padding := boxWidth - 2 - len([]rune(text))
	if padding < 0 {
		padding = 0
	}
	return "║  " + text + strings.Repeat(" ", padding) + "║"
}

func truncateText(text string, max int) string {
	if max <= 0 {
		return ""
	}
	runes := []rune(text)
	if len(runes) <= max {
		return text
	}
	if max <= 3 {
		return string(runes[:max])
	}
	return string(runes[:max-3]) + "..."
}

// DisplayImprovement shows a suggested improvement.
func DisplayImprovement(w io.Writer, section, suggestion string) {
	fmt.Fprintf(w, "\n📝 Suggested improvement (%s):\n", section)
	fmt.Fprintln(w, "────────────────────────────────────────")
	fmt.Fprintln(w, suggestion)
	fmt.Fprintln(w, "────────────────────────────────────────")
}

// DisplayTargetReached shows success message.
func DisplayTargetReached(w io.Writer, level AdherenceLevel) {
	fmt.Fprintf(w, "\n✅ Target adherence level %s reached!\n", level)
}

// DisplayMaxIterations shows timeout message.
func DisplayMaxIterations(w io.Writer, currentLevel AdherenceLevel) {
	fmt.Fprintf(w, "\n⏱️  Max iterations reached. Current level: %s\n", currentLevel)
}
