package dev

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunDev_SkillNameArg(t *testing.T) {
	dir := t.TempDir()
	skillsDir := filepath.Join(dir, "skills")

	for _, name := range []string{"one", "two"} {
		skillDir := filepath.Join(skillsDir, name)
		require.NoError(t, os.MkdirAll(skillDir, 0o755))
		content := "---\nname: " + name + "\ndescription: \"desc\"\n---\n# Body\n"
		require.NoError(t, os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644))
	}
	t.Chdir(dir)

	out := new(bytes.Buffer)
	cmd := NewCommand()
	cmd.SetOut(out)
	cmd.SetErr(new(bytes.Buffer))
	cmd.SetArgs([]string{"two"})
	err := cmd.Execute()
	// Should resolve skill name "two" and find SKILL.md
	// The command will run the dev loop (or fail for other reasons), but should not
	// fail with "skill not found"
	if err != nil {
		assert.NotContains(t, err.Error(), "not found in workspace")
	}
}

func TestRunDev_RequiresArg(t *testing.T) {
	cmd := NewCommand()
	cmd.SetOut(new(bytes.Buffer))
	cmd.SetErr(new(bytes.Buffer))
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	require.Error(t, err)
}
