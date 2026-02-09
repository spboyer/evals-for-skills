package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/spboyer/waza/internal/models"
)

func TestInitCommand_DefaultDir(t *testing.T) {
	dir := t.TempDir()

	// Change to temp dir so default "." works in an isolated location
	origDir, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() {
		err := os.Chdir(origDir)
		if err != nil {
			t.Logf("failed to restore working directory: %v", err)
		}
	})

	var buf bytes.Buffer
	cmd := newInitCommand()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{})
	require.NoError(t, cmd.Execute())

	// Verify files exist
	assert.FileExists(t, filepath.Join(dir, "eval.yaml"))
	assert.FileExists(t, filepath.Join(dir, "tasks", "example-task.yaml"))
	assert.FileExists(t, filepath.Join(dir, "fixtures", "example.py"))

	// Verify output mentions files
	output := buf.String()
	assert.Contains(t, output, "eval.yaml")
	assert.Contains(t, output, "example-task.yaml")
	assert.Contains(t, output, "example.py")
}

func TestInitCommand_NamedDir(t *testing.T) {
	parent := t.TempDir()
	target := filepath.Join(parent, "my-eval")

	var buf bytes.Buffer
	cmd := newInitCommand()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{target})
	require.NoError(t, cmd.Execute())

	assert.FileExists(t, filepath.Join(target, "eval.yaml"))
	assert.FileExists(t, filepath.Join(target, "tasks", "example-task.yaml"))
	assert.FileExists(t, filepath.Join(target, "fixtures", "example.py"))
}

func TestInitCommand_ValidSpec(t *testing.T) {
	dir := t.TempDir()

	cmd := newInitCommand()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetArgs([]string{dir})
	require.NoError(t, cmd.Execute())

	// Read and parse the generated eval.yaml
	data, err := os.ReadFile(filepath.Join(dir, "eval.yaml"))
	require.NoError(t, err)

	var spec models.BenchmarkSpec
	require.NoError(t, yaml.Unmarshal(data, &spec))

	assert.Equal(t, "my-skill-eval", spec.Name)
	assert.Equal(t, "my-skill", spec.SkillName)
	assert.Equal(t, "1.0", spec.Version)
	assert.Equal(t, "mock", spec.Config.EngineType)
	assert.Equal(t, 1, spec.Config.RunsPerTest)
	assert.Equal(t, 300, spec.Config.TimeoutSec)
	assert.Equal(t, []string{"tasks/*.yaml"}, spec.Tasks)
}

func TestInitCommand_ExistingDir(t *testing.T) {
	dir := t.TempDir()

	// Run init twice — second should succeed (overwrite)
	cmd := newInitCommand()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetArgs([]string{dir})
	require.NoError(t, cmd.Execute())

	cmd2 := newInitCommand()
	cmd2.SetOut(&bytes.Buffer{})
	cmd2.SetArgs([]string{dir})
	require.NoError(t, cmd2.Execute())
}

func TestInitCommand_TooManyArgs(t *testing.T) {
	cmd := newInitCommand()
	cmd.SetArgs([]string{"a", "b"})
	err := cmd.Execute()
	assert.Error(t, err)
}

func TestRootCommand_HasInitSubcommand(t *testing.T) {
	root := newRootCommand()
	found := false
	for _, c := range root.Commands() {
		if c.Name() == "init" {
			found = true
			break
		}
	}
	assert.True(t, found, "root command should have 'init' subcommand")
}
