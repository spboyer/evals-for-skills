package tokens

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheck_SkillFlag(t *testing.T) {
	dir := t.TempDir()
	skillsDir := filepath.Join(dir, "skills")

	// Create two skills with different content
	for _, name := range []string{"skill-x", "skill-y"} {
		skillDir := filepath.Join(skillsDir, name)
		require.NoError(t, os.MkdirAll(skillDir, 0o755))
		content := "---\nname: " + name + "\ndescription: \"desc for " + name + "\"\n---\n# Body of " + name + "\n"
		require.NoError(t, os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644))
	}

	t.Chdir(dir)

	out := new(bytes.Buffer)
	cmd := newCheckCmd()
	cmd.SetOut(out)
	cmd.SetArgs([]string{"skill-x"})
	require.NoError(t, cmd.Execute())

	result := out.String()
	assert.Contains(t, result, "SKILL.md")
	// Should only see one file (the skill-x SKILL.md)
	assert.Contains(t, result, "1/1 files within limits")
}

func TestCheck_SkillFlagNotFound(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	cmd := newCheckCmd()
	cmd.SetOut(new(bytes.Buffer))
	cmd.SetArgs([]string{"nonexistent"})
	err := cmd.Execute()
	require.Error(t, err)
}
