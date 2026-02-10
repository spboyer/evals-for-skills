package wizard

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateSkillMD_BasicSpec(t *testing.T) {
	spec := &SkillSpec{
		Name:         "code-reviewer",
		Description:  "Reviews code for quality and best practices.",
		Triggers:     []string{"review code", "check quality"},
		AntiTriggers: []string{"deploy code", "run tests"},
		Type:         SkillTypeWorkflow,
	}

	result, err := GenerateSkillMD(spec)
	require.NoError(t, err)

	assert.Contains(t, result, "name: code-reviewer")
	assert.Contains(t, result, "type: workflow")
	assert.Contains(t, result, "Reviews code for quality and best practices.")
	assert.Contains(t, result, "# code-reviewer")
	assert.Contains(t, result, "**USE FOR:**")
	assert.Contains(t, result, "- review code")
	assert.Contains(t, result, "- check quality")
	assert.Contains(t, result, "**DO NOT USE FOR:**")
	assert.Contains(t, result, "- deploy code")
	assert.Contains(t, result, "- run tests")
}

func TestGenerateSkillMD_AllTypes(t *testing.T) {
	tests := []struct {
		skillType SkillType
		expected  string
	}{
		{SkillTypeWorkflow, "type: workflow"},
		{SkillTypeUtility, "type: utility"},
		{SkillTypeAnalysis, "type: analysis"},
	}

	for _, tt := range tests {
		t.Run(string(tt.skillType), func(t *testing.T) {
			spec := &SkillSpec{
				Name:        "test-skill",
				Description: "A test skill.",
				Type:        tt.skillType,
			}
			result, err := GenerateSkillMD(spec)
			require.NoError(t, err)
			assert.Contains(t, result, tt.expected)
		})
	}
}

func TestGenerateSkillMD_EmptyTriggers(t *testing.T) {
	spec := &SkillSpec{
		Name:        "minimal-skill",
		Description: "A minimal skill with no triggers.",
		Type:        SkillTypeUtility,
	}

	result, err := GenerateSkillMD(spec)
	require.NoError(t, err)

	assert.Contains(t, result, "name: minimal-skill")
	assert.Contains(t, result, "**USE FOR:**")
	assert.Contains(t, result, "**DO NOT USE FOR:**")
}

func TestSplitAndTrim(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"empty", "", nil},
		{"single", "hello", []string{"hello"}},
		{"multiple", "a, b, c", []string{"a", "b", "c"}},
		{"with blanks", "a,, b, ,c", []string{"a", "b", "c"}},
		{"whitespace only", "  ,  ,  ", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitAndTrim(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRunSkillWizard_ValidInput(t *testing.T) {
	input := "my-skill\nA great skill\ntrigger1, trigger2\nanti1, anti2\nworkflow\n"
	in := strings.NewReader(input)
	out := &bytes.Buffer{}

	spec, err := RunSkillWizard(in, out)
	require.NoError(t, err)

	assert.Equal(t, "my-skill", spec.Name)
	assert.Equal(t, "A great skill", spec.Description)
	assert.Equal(t, []string{"trigger1", "trigger2"}, spec.Triggers)
	assert.Equal(t, []string{"anti1", "anti2"}, spec.AntiTriggers)
	assert.Equal(t, SkillTypeWorkflow, spec.Type)
}

func TestRunSkillWizard_EmptyName(t *testing.T) {
	input := "\n"
	in := strings.NewReader(input)
	out := &bytes.Buffer{}

	_, err := RunSkillWizard(in, out)
	assert.EqualError(t, err, "skill name is required")
}

func TestRunSkillWizard_EmptyDescription(t *testing.T) {
	input := "my-skill\n\n"
	in := strings.NewReader(input)
	out := &bytes.Buffer{}

	_, err := RunSkillWizard(in, out)
	assert.EqualError(t, err, "description is required")
}

func TestRunSkillWizard_InvalidType(t *testing.T) {
	input := "my-skill\nDesc\nt1\na1\nbadtype\n"
	in := strings.NewReader(input)
	out := &bytes.Buffer{}

	_, err := RunSkillWizard(in, out)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid skill type")
}

func TestRunSkillWizard_UnexpectedEOF(t *testing.T) {
	input := "my-skill\n"
	in := strings.NewReader(input)
	out := &bytes.Buffer{}

	_, err := RunSkillWizard(in, out)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected end of input")
}
