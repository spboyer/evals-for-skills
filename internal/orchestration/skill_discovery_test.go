package orchestration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscoverSkills(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Create skill directories with SKILL.md files
	skill1Dir := filepath.Join(tmpDir, "skill1")
	skill2Dir := filepath.Join(tmpDir, "skill2")
	skill3Dir := filepath.Join(tmpDir, "skill3")
	require.NoError(t, os.MkdirAll(skill1Dir, 0755))
	require.NoError(t, os.MkdirAll(skill2Dir, 0755))
	require.NoError(t, os.MkdirAll(skill3Dir, 0755))

	// Write SKILL.md for skill1
	skill1Content := `---
name: test-skill-1
description: Test skill 1
---

# Test Skill 1
This is a test skill.
`
	require.NoError(t, os.WriteFile(filepath.Join(skill1Dir, "SKILL.md"), []byte(skill1Content), 0644))

	// Write SKILL.md for skill2
	skill2Content := `---
name: test-skill-2
description: Test skill 2
---

# Test Skill 2
Another test skill.
`
	require.NoError(t, os.WriteFile(filepath.Join(skill2Dir, "SKILL.md"), []byte(skill2Content), 0644))

	// skill3 has no SKILL.md

	t.Run("discovers all skills", func(t *testing.T) {
		discovered, err := discoverSkills([]string{skill1Dir, skill2Dir, skill3Dir})
		require.NoError(t, err)
		assert.Len(t, discovered, 2)
		assert.Contains(t, discovered, "test-skill-1")
		assert.Contains(t, discovered, "test-skill-2")
		assert.Equal(t, filepath.Join(skill1Dir, "SKILL.md"), discovered["test-skill-1"])
		assert.Equal(t, filepath.Join(skill2Dir, "SKILL.md"), discovered["test-skill-2"])
	})

	t.Run("handles non-existent directory", func(t *testing.T) {
		nonExistentDir := filepath.Join(tmpDir, "does-not-exist")
		discovered, err := discoverSkills([]string{skill1Dir, nonExistentDir})
		require.NoError(t, err)
		assert.Len(t, discovered, 1)
		assert.Contains(t, discovered, "test-skill-1")
	})

	t.Run("handles empty directory list", func(t *testing.T) {
		discovered, err := discoverSkills([]string{})
		require.NoError(t, err)
		assert.Len(t, discovered, 0)
	})
}

func TestParseSkillName(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("valid SKILL.md", func(t *testing.T) {
		skillPath := filepath.Join(tmpDir, "valid-skill.md")
		content := `---
name: my-awesome-skill
description: Does awesome things
---

# My Awesome Skill
Content here.
`
		require.NoError(t, os.WriteFile(skillPath, []byte(content), 0644))

		name, err := parseSkillName(skillPath)
		require.NoError(t, err)
		assert.Equal(t, "my-awesome-skill", name)
	})

	t.Run("SKILL.md with spaces in name", func(t *testing.T) {
		skillPath := filepath.Join(tmpDir, "spaces-skill.md")
		content := `---
name: "  my-skill  "
description: Has spaces
---

Content.
`
		require.NoError(t, os.WriteFile(skillPath, []byte(content), 0644))

		name, err := parseSkillName(skillPath)
		require.NoError(t, err)
		assert.Equal(t, "my-skill", name)
	})

	t.Run("non-existent file", func(t *testing.T) {
		_, err := parseSkillName(filepath.Join(tmpDir, "does-not-exist.md"))
		require.Error(t, err)
	})

	t.Run("invalid frontmatter", func(t *testing.T) {
		skillPath := filepath.Join(tmpDir, "invalid.md")
		content := `---
name: [invalid yaml
---
`
		require.NoError(t, os.WriteFile(skillPath, []byte(content), 0644))

		_, err := parseSkillName(skillPath)
		require.Error(t, err)
	})
}

func TestValidateRequiredSkills(t *testing.T) {
	t.Run("all required skills found", func(t *testing.T) {
		required := []string{"skill-a", "skill-b"}
		discovered := map[string]string{
			"skill-a": "/path/to/skill-a/SKILL.md",
			"skill-b": "/path/to/skill-b/SKILL.md",
			"skill-c": "/path/to/skill-c/SKILL.md",
		}
		searchedDirs := []string{"/path/to/skill-a", "/path/to/skill-b", "/path/to/skill-c"}

		err := validateRequiredSkills(required, discovered, searchedDirs)
		assert.NoError(t, err)
	})

	t.Run("some required skills missing", func(t *testing.T) {
		required := []string{"skill-a", "skill-b", "skill-c"}
		discovered := map[string]string{
			"skill-a": "/path/to/skill-a/SKILL.md",
		}
		searchedDirs := []string{"/path/to/skill-a", "/path/to/skill-b"}

		err := validateRequiredSkills(required, discovered, searchedDirs)
		require.Error(t, err)
		errMsg := err.Error()
		assert.Contains(t, errMsg, "required skills not found")
		assert.Contains(t, errMsg, "skill-b")
		assert.Contains(t, errMsg, "skill-c")
		assert.Contains(t, errMsg, "Searched directories")
		assert.Contains(t, errMsg, "/path/to/skill-a")
		assert.Contains(t, errMsg, "/path/to/skill-b")
		assert.Contains(t, errMsg, "Found skills")
		assert.Contains(t, errMsg, "skill-a")
	})

	t.Run("no skills discovered but required", func(t *testing.T) {
		required := []string{"skill-a", "skill-b"}
		discovered := map[string]string{}
		searchedDirs := []string{"/path/to/empty"}

		err := validateRequiredSkills(required, discovered, searchedDirs)
		require.Error(t, err)
		errMsg := err.Error()
		assert.Contains(t, errMsg, "required skills not found")
		assert.Contains(t, errMsg, "skill-a")
		assert.Contains(t, errMsg, "skill-b")
		assert.Contains(t, errMsg, "No skills were found")
	})

	t.Run("empty required skills list", func(t *testing.T) {
		required := []string{}
		discovered := map[string]string{
			"skill-a": "/path/to/skill-a/SKILL.md",
		}
		searchedDirs := []string{"/path/to/skill-a"}

		err := validateRequiredSkills(required, discovered, searchedDirs)
		assert.NoError(t, err)
	})

	t.Run("nil required skills list", func(t *testing.T) {
		var required []string
		discovered := map[string]string{
			"skill-a": "/path/to/skill-a/SKILL.md",
		}
		searchedDirs := []string{"/path/to/skill-a"}

		err := validateRequiredSkills(required, discovered, searchedDirs)
		assert.NoError(t, err)
	})

	t.Run("error message formatting", func(t *testing.T) {
		required := []string{"missing-skill"}
		discovered := map[string]string{
			"found-skill": "/path/SKILL.md",
		}
		searchedDirs := []string{"/dir1", "/dir2"}

		err := validateRequiredSkills(required, discovered, searchedDirs)
		require.Error(t, err)

		// Check that error message is well-formatted
		errMsg := err.Error()
		lines := strings.Split(errMsg, "\n")

		// Should have multiple lines with clear sections
		assert.Greater(t, len(lines), 5, "Error message should be multi-line")

		// Check structure
		foundRequiredSection := false
		foundSearchedSection := false
		foundDiscoveredSection := false

		for _, line := range lines {
			if strings.Contains(line, "required skills not found") {
				foundRequiredSection = true
			}
			if strings.Contains(line, "Searched directories") {
				foundSearchedSection = true
			}
			if strings.Contains(line, "Found skills") {
				foundDiscoveredSection = true
			}
		}

		assert.True(t, foundRequiredSection, "Should have required skills section")
		assert.True(t, foundSearchedSection, "Should have searched directories section")
		assert.True(t, foundDiscoveredSection, "Should have found skills section")
	})
}
