package scaffold

import (
	"fmt"
	"regexp"
	"strings"
)

// TriggerPhrase holds a parsed trigger or anti-trigger phrase.
type TriggerPhrase struct {
	Prompt string
	Reason string
}

// ParseTriggerPhrases extracts USE FOR and DO NOT USE FOR phrases from a
// skill description field. Phrases may be quoted or comma-separated.
func ParseTriggerPhrases(description string) (useFor []TriggerPhrase, doNotUseFor []TriggerPhrase) {
	// Normalize newlines and collapse to single line for regex matching.
	text := strings.ReplaceAll(description, "\r\n", " ")
	text = strings.ReplaceAll(text, "\n", " ")

	useFor = extractSection(text, `USE FOR:`)
	doNotUseFor = extractSection(text, `DO NOT USE FOR:`)
	return useFor, doNotUseFor
}

// extractSection finds a labeled section (e.g. "USE FOR:") and parses
// comma-separated phrases until the next label or end of string.
func extractSection(text, label string) []TriggerPhrase {
	upper := strings.ToUpper(text)
	labelUpper := strings.ToUpper(label)

	idx := -1
	searchFrom := 0
	for {
		i := strings.Index(upper[searchFrom:], labelUpper)
		if i < 0 {
			break
		}
		pos := searchFrom + i
		// For "USE FOR:", reject matches inside "DO NOT USE FOR:".
		if labelUpper == "USE FOR:" && pos >= 4 {
			preceding := upper[max(0, pos-11):pos]
			if strings.Contains(preceding, "NOT ") {
				searchFrom = pos + len(labelUpper)
				continue
			}
		}
		idx = pos
		break
	}
	if idx < 0 {
		return nil
	}

	// Take everything after the label.
	after := text[idx+len(label):]

	// Terminate at the next known label or period-space boundary.
	terminators := []string{"DO NOT USE FOR:", "USE FOR:", "INVOKES:", "FOR SINGLE OPERATIONS:"}
	for _, t := range terminators {
		if t == label {
			continue
		}
		if ti := strings.Index(strings.ToUpper(after), strings.ToUpper(t)); ti >= 0 {
			after = after[:ti]
		}
	}

	after = strings.TrimSpace(after)
	// Remove trailing period.
	after = strings.TrimRight(after, ".")

	return splitPhrases(after)
}

// quotedPhraseRE matches "quoted text" (double or curly quotes).
var quotedPhraseRE = regexp.MustCompile(`[""\x{201C}\x{201D}]([^""\x{201C}\x{201D}]+)[""\x{201C}\x{201D}]`)

// splitPhrases splits a comma-separated list of phrases. If the text
// contains quoted strings, only the quoted portions are extracted.
// Parenthetical notes like "(use X)" are stripped from anti-trigger reasons.
func splitPhrases(text string) []TriggerPhrase {
	// Try quoted extraction first.
	matches := quotedPhraseRE.FindAllStringSubmatch(text, -1)
	if len(matches) > 0 {
		var phrases []TriggerPhrase
		for _, m := range matches {
			p := strings.TrimSpace(m[1])
			if p != "" {
				phrases = append(phrases, TriggerPhrase{
					Prompt: p,
					Reason: fmt.Sprintf("Matches frontmatter trigger phrase: %q", p),
				})
			}
		}
		return phrases
	}

	// Fall back to comma-separated unquoted phrases.
	parts := strings.Split(text, ",")
	var phrases []TriggerPhrase
	for _, part := range parts {
		p := strings.TrimSpace(part)
		p = stripParenthetical(p)
		p = strings.TrimSpace(p)
		if p != "" {
			phrases = append(phrases, TriggerPhrase{
				Prompt: p,
				Reason: fmt.Sprintf("Matches frontmatter trigger phrase: %q", p),
			})
		}
	}
	return phrases
}

// parentheticalRE matches (parenthetical notes) at the end of a phrase.
var parentheticalRE = regexp.MustCompile(`\s*\([^)]*\)\s*$`)

// stripParenthetical removes trailing parenthetical notes like "(use X)".
func stripParenthetical(s string) string {
	return parentheticalRE.ReplaceAllString(s, "")
}

// TriggerTestsYAML generates a trigger_tests.yaml file from parsed phrases.
func TriggerTestsYAML(skillName string, useFor, doNotUseFor []TriggerPhrase) string {
	var b strings.Builder

	fmt.Fprintf(&b, "# Trigger accuracy tests for %s skill\n", skillName)
	fmt.Fprintf(&b, "# Auto-generated from SKILL.md frontmatter\n\n")
	fmt.Fprintf(&b, "skill: %s\n", skillName)

	if len(useFor) > 0 {
		b.WriteString("\n# Prompts that SHOULD trigger this skill\n")
		b.WriteString("should_trigger_prompts:\n")
		for _, p := range useFor {
			fmt.Fprintf(&b, "  - prompt: %q\n", p.Prompt)
			fmt.Fprintf(&b, "    reason: %q\n", p.Reason)
			b.WriteString("    confidence: high\n\n")
		}
	}

	if len(doNotUseFor) > 0 {
		b.WriteString("# Prompts that should NOT trigger this skill\n")
		b.WriteString("should_not_trigger_prompts:\n")
		for _, p := range doNotUseFor {
			fmt.Fprintf(&b, "  - prompt: %q\n", p.Prompt)
			fmt.Fprintf(&b, "    reason: %q\n", p.Reason)
			b.WriteString("    confidence: high\n\n")
		}
	}

	return b.String()
}
