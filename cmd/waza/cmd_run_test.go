package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// resetRunGlobals zeroes the package-level flag vars so prior tests don't leak.
func resetRunGlobals() {
	contextDir = ""
	outputPath = ""
	verbose = false
	taskFilters = nil
}

// helper creates a valid minimal eval spec YAML in a temp dir,
// including a matching task file, and returns the spec path.
func createTestSpec(t *testing.T, engine string) string {
	t.Helper()
	dir := t.TempDir()

	taskDir := filepath.Join(dir, "tasks")
	require.NoError(t, os.MkdirAll(taskDir, 0o755))

	task := `id: test-task-001
name: Test Task
inputs:
  prompt: "Explain this code"
`
	require.NoError(t, os.WriteFile(filepath.Join(taskDir, "task.yaml"), []byte(task), 0o644))

	spec := `name: test-eval
skill: test-skill
version: "1.0"
config:
  trials_per_task: 1
  timeout_seconds: 30
  parallel: false
  executor: ` + engine + `
  model: test-model
tasks:
  - "tasks/*.yaml"
`
	specPath := filepath.Join(dir, "eval.yaml")
	require.NoError(t, os.WriteFile(specPath, []byte(spec), 0o644))
	return specPath
}

// ---------------------------------------------------------------------------
// Argument validation
// ---------------------------------------------------------------------------

func TestRunCommand_RequiresExactlyOneArg(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"no args", []string{}},
		{"two args", []string{"a.yaml", "b.yaml"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newRunCommand()
			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			assert.Error(t, err, "expected error for args=%v", tt.args)
		})
	}
}

// ---------------------------------------------------------------------------
// Flag parsing
// ---------------------------------------------------------------------------

func TestRunCommand_FlagsParsed(t *testing.T) {
	tmpCtx := filepath.Join(t.TempDir(), "ctx")
	tmpOut := filepath.Join(t.TempDir(), "out.json")

	cmd := newRunCommand()
	cmd.SetArgs([]string{
		"--context-dir", tmpCtx,
		"--output", tmpOut,
		"--verbose",
		"spec.yaml",
	})

	// Don't execute — just parse flags to verify they bind.
	require.NoError(t, cmd.ParseFlags([]string{
		"--context-dir", tmpCtx,
		"--output", tmpOut,
		"--verbose",
	}))

	val, err := cmd.Flags().GetString("context-dir")
	require.NoError(t, err)
	assert.Equal(t, tmpCtx, val)

	val, err = cmd.Flags().GetString("output")
	require.NoError(t, err)
	assert.Equal(t, tmpOut, val)

	boolVal, err := cmd.Flags().GetBool("verbose")
	require.NoError(t, err)
	assert.True(t, boolVal)
}

func TestRunCommand_ShortFlags(t *testing.T) {
	tmpOut := filepath.Join(t.TempDir(), "out.json")

	cmd := newRunCommand()
	require.NoError(t, cmd.ParseFlags([]string{
		"-o", tmpOut,
		"-v",
	}))

	val, err := cmd.Flags().GetString("output")
	require.NoError(t, err)
	assert.Equal(t, tmpOut, val)

	boolVal, err := cmd.Flags().GetBool("verbose")
	require.NoError(t, err)
	assert.True(t, boolVal)
}

// ---------------------------------------------------------------------------
// Error handling
// ---------------------------------------------------------------------------

func TestRunCommand_MissingSpecFile(t *testing.T) {
	resetRunGlobals()

	cmd := newRunCommand()
	cmd.SetArgs([]string{"nonexistent.yaml"})
	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load spec")
}

func TestRunCommand_InvalidSpecFile(t *testing.T) {
	resetRunGlobals()

	dir := t.TempDir()
	badSpec := filepath.Join(dir, "bad.yaml")
	require.NoError(t, os.WriteFile(badSpec, []byte("foo: [bar"), 0o644))

	cmd := newRunCommand()
	cmd.SetArgs([]string{badSpec})
	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load spec")
}

func TestRunCommand_InvalidEngineType(t *testing.T) {
	resetRunGlobals()

	dir := t.TempDir()
	taskDir := filepath.Join(dir, "tasks")
	require.NoError(t, os.MkdirAll(taskDir, 0o755))
	require.NoError(t, os.WriteFile(
		filepath.Join(taskDir, "t.yaml"),
		[]byte("id: t1\nname: t\ninputs:\n  prompt: hi\n"),
		0o644,
	))

	spec := `name: test
skill: s
config:
  trials_per_task: 1
  timeout_seconds: 10
  executor: nonexistent-engine
  model: m
tasks:
  - "tasks/*.yaml"
`
	specPath := filepath.Join(dir, "eval.yaml")
	require.NoError(t, os.WriteFile(specPath, []byte(spec), 0o644))

	cmd := newRunCommand()
	cmd.SetArgs([]string{specPath})
	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown engine type")
}

// ---------------------------------------------------------------------------
// Integration with mock engine — full run
// ---------------------------------------------------------------------------

func TestRunCommand_MockEngineRun(t *testing.T) {
	resetRunGlobals()

	specPath := createTestSpec(t, "mock")

	cmd := newRunCommand()
	cmd.SetArgs([]string{specPath})

	// Suppress stdout/stderr noise during test
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)

	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestRunCommand_MockEngineVerbose(t *testing.T) {
	resetRunGlobals()

	specPath := createTestSpec(t, "mock")

	cmd := newRunCommand()
	cmd.SetArgs([]string{specPath, "--verbose"})

	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestRunCommand_OutputJSON(t *testing.T) {
	resetRunGlobals()

	specPath := createTestSpec(t, "mock")
	outFile := filepath.Join(t.TempDir(), "results.json")

	cmd := newRunCommand()
	cmd.SetArgs([]string{specPath, "--output", outFile})

	err := cmd.Execute()
	require.NoError(t, err)

	// Verify JSON output was written and is valid
	data, err := os.ReadFile(outFile)
	require.NoError(t, err)
	assert.Greater(t, len(data), 0)

	var result map[string]any
	require.NoError(t, json.Unmarshal(data, &result))
	assert.Equal(t, "test-eval", result["eval_name"])
	assert.Equal(t, "test-skill", result["skill"])
}

func TestRunCommand_ContextDirFlag(t *testing.T) {
	resetRunGlobals()

	specPath := createTestSpec(t, "mock")

	// Pass an explicit --context-dir (the spec already has fixtures alongside it,
	// but the flag should override without error).
	cmd := newRunCommand()
	cmd.SetArgs([]string{specPath, "--context-dir", t.TempDir()})

	// This will succeed because the mock engine doesn't need real fixture files.
	err := cmd.Execute()
	assert.NoError(t, err)
}

// ---------------------------------------------------------------------------
// Root command wiring
// ---------------------------------------------------------------------------

func TestRootCommand_HasRunSubcommand(t *testing.T) {
	root := newRootCommand()
	found := false
	for _, c := range root.Commands() {
		if c.Name() == "run" {
			found = true
			break
		}
	}
	assert.True(t, found, "root command should have 'run' subcommand")
}

// ---------------------------------------------------------------------------
// Task filtering via --task flag
// ---------------------------------------------------------------------------

func TestRunCommand_TaskFlagParsed(t *testing.T) {
	cmd := newRunCommand()
	require.NoError(t, cmd.ParseFlags([]string{
		"--task", "Create*",
		"--task", "tc-002",
	}))

	vals, err := cmd.Flags().GetStringArray("task")
	require.NoError(t, err)
	assert.Equal(t, []string{"Create*", "tc-002"}, vals)
}

func TestRunCommand_TaskFilterRunsMock(t *testing.T) {
	resetRunGlobals()

	specPath := createTestSpec(t, "mock")

	cmd := newRunCommand()
	cmd.SetArgs([]string{specPath, "--task", "Test Task"})
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)

	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestRunCommand_TaskFilterNoMatch(t *testing.T) {
	resetRunGlobals()

	specPath := createTestSpec(t, "mock")

	cmd := newRunCommand()
	cmd.SetArgs([]string{specPath, "--task", "nonexistent-task"})
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no test cases found")
}
