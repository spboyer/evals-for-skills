package wizard

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"

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

// pipeInput feeds lines to an io.Pipe with small delays so bubbletea
// processes each field before the next line arrives.
func pipeInput(t *testing.T, lines ...string) io.Reader {
	t.Helper()
	r, w := io.Pipe()
	go func() {
		defer w.Close() //nolint:errcheck
		for _, line := range lines {
			w.Write([]byte(line + "\n")) //nolint:errcheck
			time.Sleep(50 * time.Millisecond)
		}
	}()
	return r
}

func TestRunSkillWizard_ValidInput(t *testing.T) {
	// huh accessible mode: Input reads a line, Select reads option number (1-indexed)
	in := pipeInput(t, "my-skill", "A great skill", "trigger1, trigger2", "anti1, anti2", "1")
	out := &bytes.Buffer{}

	spec, err := RunSkillWizard(in, out)
	require.NoError(t, err)

	assert.Equal(t, "my-skill", spec.Name)
	assert.Equal(t, "A great skill", spec.Description)
	assert.Equal(t, []string{"trigger1", "trigger2"}, spec.Triggers)
	assert.Equal(t, []string{"anti1", "anti2"}, spec.AntiTriggers)
	assert.Equal(t, SkillTypeWorkflow, spec.Type)
}

func TestRunSkillWizard_EmptyInput(t *testing.T) {
	// huh gracefully handles EOF by using default (empty) values
	in := strings.NewReader("")
	out := &bytes.Buffer{}

	spec, err := RunSkillWizard(in, out)
	require.NoError(t, err)
	assert.Empty(t, spec.Name)
	assert.Equal(t, SkillType("workflow"), spec.Type) // default to first option
}

func TestRunSkillWizard_IncompleteInput(t *testing.T) {
	// Only name provided; remaining fields use defaults
	in := pipeInput(t, "my-skill")
	out := &bytes.Buffer{}

	spec, err := RunSkillWizard(in, out)
	require.NoError(t, err)
	assert.Equal(t, "my-skill", spec.Name)
	assert.Empty(t, spec.Description)
}

func TestRunSkillWizard_SelectAnalysis(t *testing.T) {
	// 3 = analysis (third option)
	in := pipeInput(t, "test-skill", "A test skill", "", "", "3")
	out := &bytes.Buffer{}

	spec, err := RunSkillWizard(in, out)
	require.NoError(t, err)
	assert.Equal(t, SkillTypeAnalysis, spec.Type)
}
