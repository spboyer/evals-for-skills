package orchestration

import (
	"testing"

	"github.com/spboyer/waza/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sampleCases() []*models.TestCase {
	return []*models.TestCase{
		{TestID: "tc-001", DisplayName: "Create a REST API"},
		{TestID: "tc-002", DisplayName: "Fix login bug"},
		{TestID: "tc-003", DisplayName: "Create a CLI tool"},
		{TestID: "tc-004", DisplayName: "Optimize SQL query"},
	}
}

func TestFilterTestCases_NoPatterns(t *testing.T) {
	cases := sampleCases()
	result, err := FilterTestCases(cases, nil)
	require.NoError(t, err)
	assert.Len(t, result, 4, "empty patterns should return all cases")
}

func TestFilterTestCases_ExactName(t *testing.T) {
	cases := sampleCases()
	result, err := FilterTestCases(cases, []string{"Fix login bug"})
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "tc-002", result[0].TestID)
}

func TestFilterTestCases_ExactID(t *testing.T) {
	cases := sampleCases()
	result, err := FilterTestCases(cases, []string{"tc-003"})
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "Create a CLI tool", result[0].DisplayName)
}

func TestFilterTestCases_GlobPattern(t *testing.T) {
	cases := sampleCases()
	result, err := FilterTestCases(cases, []string{"Create*"})
	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "tc-001", result[0].TestID)
	assert.Equal(t, "tc-003", result[1].TestID)
}

func TestFilterTestCases_MultiplePatterns(t *testing.T) {
	cases := sampleCases()
	result, err := FilterTestCases(cases, []string{"tc-001", "Optimize*"})
	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "tc-001", result[0].TestID)
	assert.Equal(t, "tc-004", result[1].TestID)
}

func TestFilterTestCases_NoMatch(t *testing.T) {
	cases := sampleCases()
	result, err := FilterTestCases(cases, []string{"nonexistent"})
	require.NoError(t, err)
	assert.Len(t, result, 0)
}

func TestFilterTestCases_InvalidPattern(t *testing.T) {
	cases := sampleCases()
	_, err := FilterTestCases(cases, []string{"["})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid task filter pattern")
}

func TestFilterTestCases_IDGlob(t *testing.T) {
	cases := sampleCases()
	result, err := FilterTestCases(cases, []string{"tc-00?"})
	require.NoError(t, err)
	assert.Len(t, result, 4, "? should match single character in IDs")
}
