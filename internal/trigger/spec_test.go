package trigger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTriggerTestSpec_ParseSpec(t *testing.T) {
	content := `
skill: code-explainer

should_trigger_prompts:
  - prompt: "Explain this code"
  - prompt: "What does this code do?"

should_not_trigger_prompts:
  - prompt: "Write me a sort function"
    reason: "Code writing, not explaining"
`
	spec, err := ParseSpec([]byte(content))
	require.NoError(t, err)
	require.Equal(t, "code-explainer", spec.Skill)
	require.Len(t, spec.ShouldTriggerPrompts, 2)
	require.Len(t, spec.ShouldNotTriggerPrompts, 1)
	require.Equal(t, "Explain this code", spec.ShouldTriggerPrompts[0].Prompt)
}

func TestTriggerTestSpec_UnmarshalText_MissingSkill(t *testing.T) {
	content := `
should_trigger_prompts:
  - prompt: "Explain this"
    reason: "test"
`
	_, err := ParseSpec([]byte(content))
	require.ErrorContains(t, err, "missing required 'skill' field")
}

func TestTriggerTestSpec_UnmarshalText_NoPrompts(t *testing.T) {
	content := `
skill: test-skill

should_trigger_prompts:
should_not_trigger_prompts:
`
	_, err := ParseSpec([]byte(content))
	require.ErrorContains(t, err, "at least one prompt")
}

func TestTriggerTestSpec_UnmarshalText_InvalidConfidence(t *testing.T) {
	content := `
skill: test-skill

should_trigger_prompts:
  - prompt: "Explain this"
    confidence: low
`
	_, err := ParseSpec([]byte(content))
	require.ErrorContains(t, err, "unrecognized confidence")
}
