package scaffold

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTriggerPhrases_QuotedPhrases(t *testing.T) {
	desc := `Process PDF files including text extraction, rotation, and merging.
USE FOR: "extract PDF text", "rotate PDF", "merge PDFs", "split PDF pages".
DO NOT USE FOR: creating PDFs from scratch (use document-creator), editing PDF forms (use form-filler).`

	useFor, doNotUseFor := ParseTriggerPhrases(desc)

	require.Len(t, useFor, 4)
	assert.Equal(t, "extract PDF text", useFor[0].Prompt)
	assert.Equal(t, "rotate PDF", useFor[1].Prompt)
	assert.Equal(t, "merge PDFs", useFor[2].Prompt)
	assert.Equal(t, "split PDF pages", useFor[3].Prompt)

	require.Len(t, doNotUseFor, 2)
	assert.Equal(t, "creating PDFs from scratch", doNotUseFor[0].Prompt)
	assert.Equal(t, "editing PDF forms", doNotUseFor[1].Prompt)
}

func TestParseTriggerPhrases_UnquotedPhrases(t *testing.T) {
	desc := `USE FOR: run waza, evaluate skill, test agent.
DO NOT USE FOR: improving skill frontmatter (use waza dev), creating new skills (use skill-creator).`

	useFor, doNotUseFor := ParseTriggerPhrases(desc)

	require.Len(t, useFor, 3)
	assert.Equal(t, "run waza", useFor[0].Prompt)
	assert.Equal(t, "evaluate skill", useFor[1].Prompt)
	assert.Equal(t, "test agent", useFor[2].Prompt)

	require.Len(t, doNotUseFor, 2)
	assert.Equal(t, "improving skill frontmatter", doNotUseFor[0].Prompt)
	assert.Equal(t, "creating new skills", doNotUseFor[1].Prompt)
}

func TestParseTriggerPhrases_NoTriggers(t *testing.T) {
	desc := "A simple tool that does things."

	useFor, doNotUseFor := ParseTriggerPhrases(desc)

	assert.Empty(t, useFor)
	assert.Empty(t, doNotUseFor)
}

func TestParseTriggerPhrases_OnlyUseFor(t *testing.T) {
	desc := `USE FOR: "task a", "task b".`

	useFor, doNotUseFor := ParseTriggerPhrases(desc)

	require.Len(t, useFor, 2)
	assert.Equal(t, "task a", useFor[0].Prompt)
	assert.Equal(t, "task b", useFor[1].Prompt)
	assert.Empty(t, doNotUseFor)
}

func TestParseTriggerPhrases_OnlyDoNotUseFor(t *testing.T) {
	desc := `DO NOT USE FOR: unrelated tasks, dangerous operations.`

	useFor, doNotUseFor := ParseTriggerPhrases(desc)

	assert.Empty(t, useFor)
	require.Len(t, doNotUseFor, 2)
	assert.Equal(t, "unrelated tasks", doNotUseFor[0].Prompt)
	assert.Equal(t, "dangerous operations", doNotUseFor[1].Prompt)
}

func TestParseTriggerPhrases_MultilineDescription(t *testing.T) {
	desc := `**WORKFLOW SKILL** - Process PDF files.
USE FOR: "extract PDF text", "rotate PDF",
"merge PDFs".
DO NOT USE FOR: creating PDFs from scratch (use document-creator).`

	useFor, doNotUseFor := ParseTriggerPhrases(desc)

	require.Len(t, useFor, 3)
	assert.Equal(t, "extract PDF text", useFor[0].Prompt)
	assert.Equal(t, "rotate PDF", useFor[1].Prompt)
	assert.Equal(t, "merge PDFs", useFor[2].Prompt)

	require.Len(t, doNotUseFor, 1)
	assert.Equal(t, "creating PDFs from scratch", doNotUseFor[0].Prompt)
}

func TestParseTriggerPhrases_WithInvokesSection(t *testing.T) {
	desc := `USE FOR: "task a", "task b".
DO NOT USE FOR: bad tasks.
INVOKES: some-tool MCP.
FOR SINGLE OPERATIONS: do it directly.`

	useFor, doNotUseFor := ParseTriggerPhrases(desc)

	require.Len(t, useFor, 2)
	require.Len(t, doNotUseFor, 1)
	assert.Equal(t, "bad tasks", doNotUseFor[0].Prompt)
}

func TestParseTriggerPhrases_ReasonField(t *testing.T) {
	desc := `USE FOR: "run eval".`

	useFor, _ := ParseTriggerPhrases(desc)

	require.Len(t, useFor, 1)
	assert.Contains(t, useFor[0].Reason, "run eval")
	assert.Contains(t, useFor[0].Reason, "frontmatter trigger phrase")
}

func TestParseTriggerPhrases_EmptyDescription(t *testing.T) {
	useFor, doNotUseFor := ParseTriggerPhrases("")

	assert.Empty(t, useFor)
	assert.Empty(t, doNotUseFor)
}

func TestStripParenthetical(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"creating PDFs (use doc-creator)", "creating PDFs"},
		{"plain text", "plain text"},
		{"nested (one (two))", "nested (one (two))"},
		{"", ""},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			assert.Equal(t, tc.want, stripParenthetical(tc.input))
		})
	}
}

func TestTriggerTestsYAML(t *testing.T) {
	useFor := []TriggerPhrase{
		{Prompt: "extract PDF text", Reason: `Matches frontmatter trigger phrase: "extract PDF text"`},
		{Prompt: "rotate PDF", Reason: `Matches frontmatter trigger phrase: "rotate PDF"`},
	}
	doNotUseFor := []TriggerPhrase{
		{Prompt: "creating PDFs from scratch", Reason: `Matches frontmatter trigger phrase: "creating PDFs from scratch"`},
	}

	yaml := TriggerTestsYAML("pdf-processor", useFor, doNotUseFor)

	assert.Contains(t, yaml, "skill: pdf-processor")
	assert.Contains(t, yaml, "should_trigger_prompts:")
	assert.Contains(t, yaml, `prompt: "extract PDF text"`)
	assert.Contains(t, yaml, `prompt: "rotate PDF"`)
	assert.Contains(t, yaml, "should_not_trigger_prompts:")
	assert.Contains(t, yaml, `prompt: "creating PDFs from scratch"`)
	assert.Contains(t, yaml, "confidence: high")
	assert.Contains(t, yaml, "Auto-generated from SKILL.md frontmatter")
}

func TestTriggerTestsYAML_EmptyPhrases(t *testing.T) {
	yaml := TriggerTestsYAML("my-skill", nil, nil)

	assert.Contains(t, yaml, "skill: my-skill")
	assert.NotContains(t, yaml, "should_trigger_prompts:")
	assert.NotContains(t, yaml, "should_not_trigger_prompts:")
}

func TestTriggerTestsYAML_OnlyUseFor(t *testing.T) {
	useFor := []TriggerPhrase{
		{Prompt: "run eval", Reason: "test"},
	}

	yaml := TriggerTestsYAML("test-skill", useFor, nil)

	assert.Contains(t, yaml, "should_trigger_prompts:")
	assert.NotContains(t, yaml, "should_not_trigger_prompts:")
}

func TestTriggerTestsYAML_OnlyDoNotUseFor(t *testing.T) {
	doNotUseFor := []TriggerPhrase{
		{Prompt: "bad task", Reason: "test"},
	}

	yaml := TriggerTestsYAML("test-skill", nil, doNotUseFor)

	assert.NotContains(t, yaml, "should_trigger_prompts:")
	assert.Contains(t, yaml, "should_not_trigger_prompts:")
}

func TestParseTriggerPhrases_RealWorldWazaSkill(t *testing.T) {
	// Real-world description from skills/waza/SKILL.md
	desc := `**WORKFLOW SKILL** - Evaluate AI agent skills using structured benchmarks. USE FOR: run waza, waza help, run eval, run benchmark, evaluate skill, test agent. DO NOT USE FOR: improving skill frontmatter (use waza dev), creating new skills from scratch (use skill-creator), token counting or budget checks (use waza tokens). INVOKES: Copilot SDK executor.`

	useFor, doNotUseFor := ParseTriggerPhrases(desc)

	require.Len(t, useFor, 6)
	assert.Equal(t, "run waza", useFor[0].Prompt)
	assert.Equal(t, "waza help", useFor[1].Prompt)
	assert.Equal(t, "test agent", useFor[5].Prompt)

	require.Len(t, doNotUseFor, 3)
	assert.Equal(t, "improving skill frontmatter", doNotUseFor[0].Prompt)
	assert.Equal(t, "creating new skills from scratch", doNotUseFor[1].Prompt)
	assert.Equal(t, "token counting or budget checks", doNotUseFor[2].Prompt)
}
