package graders

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/spboyer/waza/internal/models"
	"github.com/stretchr/testify/require"
)

func TestScriptGrader_Python(t *testing.T) {
	skipIfNoPython(t)

	dir := t.TempDir()

	err := os.WriteFile(filepath.Join(dir, "grader.py"), []byte(`
import json, sys
ctx = json.load(sys.stdin)
score = 1.0 if "hello" in ctx.get("output", "") else 0.0
print(json.dumps({"score": score, "passed": score >= 0.5, "message": "checked", "details": {"len": len(ctx.get("output", ""))}}))
`), 0600)
	require.NoError(t, err)

	grader, err := NewScriptGrader("py_test", ScriptGraderArgs{
		SpecDir: dir,
		Path:    "grader.py",
	})
	require.NoError(t, err)
	require.Equal(t, "py_test", grader.Name())
	require.Equal(t, models.GraderKindScript, grader.Kind())

	t.Run("pass", func(t *testing.T) {
		result, err := grader.Grade(context.Background(), &Context{
			Output: "hello world",
		})
		require.NoError(t, err)
		require.Greater(t, result.DurationMs, int64(0))
		result.DurationMs = 0

		require.Equal(t, &models.GraderResults{
			Name:     "py_test",
			Type:     models.GraderKindScript,
			Score:    1.0,
			Passed:   true,
			Feedback: "checked",
			Details:  map[string]any{"len": float64(len("hello world"))},
		}, result)
	})

	t.Run("fail", func(t *testing.T) {
		result, err := grader.Grade(context.Background(), &Context{
			Output: "goodbye",
		})
		require.NoError(t, err)
		require.Equal(t, 0.0, result.Score)
		require.False(t, result.Passed)
	})
}

func TestScriptGrader_Javascript(t *testing.T) {
	skipIfNoJavascript(t)

	dir := t.TempDir()

	graderJS := filepath.Join(dir, "grader.js")
	err := os.WriteFile(graderJS, []byte(`
import {readFileSync} from "node:fs";
const text = readFileSync(process.stdin.fd, "utf-8");
const ctx = JSON.parse(text);
const score = ctx.output.includes("hello") ? 1.0 : 0.0;

console.log(JSON.stringify({score, passed: score >= 0.5, message: "js check", details: {}}));
`), 0600)
	require.NoError(t, err)

	grader, err := NewScriptGrader("js_test", ScriptGraderArgs{
		SpecDir: dir,
		Path:    "grader.js",
	})
	require.NoError(t, err)

	result, err := grader.Grade(context.Background(), &Context{
		Output: "hello from js",
	})
	require.NoError(t, err)
	require.Equal(t, 1.0, result.Score)
	require.True(t, result.Passed)
}

func TestScriptGrader_ScriptNotFound(t *testing.T) {
	grader, err := NewScriptGrader("missing", ScriptGraderArgs{
		Path: "nonexistent.py",
	})
	require.NoError(t, err)

	_, err = grader.Grade(context.Background(), &Context{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "script not found")
}

func TestScriptGrader_UnsupportedExtension(t *testing.T) {
	_, err := NewScriptGrader("bad", ScriptGraderArgs{
		Path: "grader.rb",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported script extension")
}

func TestScriptGrader_InvalidJSON(t *testing.T) {
	skipIfNoPython(t)

	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, "bad.py"), []byte(`print("not json")`), 0600)
	require.NoError(t, err)

	grader, err := NewScriptGrader("bad_json", ScriptGraderArgs{
		Path:    "bad.py",
		SpecDir: dir,
	})
	require.NoError(t, err)

	_, err = grader.Grade(context.Background(), &Context{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to parse script output")
}

func TestScriptGrader_NonZeroExit(t *testing.T) {
	skipIfNoPython(t)

	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, "crash.py"), []byte(`import sys; sys.exit(1)`), 0600)
	require.NoError(t, err)

	grader, err := NewScriptGrader("crash", ScriptGraderArgs{
		Path:    "crash.py",
		SpecDir: dir,
	})
	require.NoError(t, err)

	_, err = grader.Grade(context.Background(), &Context{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "script execution failed")
}

func TestScriptGrader_WithTranscript(t *testing.T) {
	skipIfNoPython(t)

	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, "check_ctx.py"), []byte(`
import json, sys
ctx = json.load(sys.stdin)
checks = []
if ctx.get("duration_ms") == 101:
    checks.append("duration ok")
if isinstance(ctx.get("tool_calls"), list):
    checks.append("tool_calls ok")
if isinstance(ctx.get("transcript"), list):
    checks.append("transcript ok")
score = float(len(checks)) / 3.0
print(json.dumps({"score": score, "passed": score == 1.0, "message": "; ".join(checks), "details": {}}))
`), 0600)
	require.NoError(t, err)

	sessionEvents := loadSampleEvents(t)
	transcript := convertToTranscriptEvents(sessionEvents)

	grader, err := NewScriptGrader("ctx_check", ScriptGraderArgs{
		Path:    "check_ctx.py",
		SpecDir: dir,
	})
	require.NoError(t, err)

	result, err := grader.Grade(context.Background(), &Context{
		Output:     "test",
		Transcript: transcript,
		DurationMS: 101,
	})
	require.NoError(t, err)
	require.Equal(t, 1.0, result.Score)
	require.True(t, result.Passed)
}

func TestScriptGrader_Create(t *testing.T) {
	grader, err := Create(CreateArgs{
		GraderKind: models.GraderKindScript,
		Identifier: "test_script",
		Config: map[string]any{
			"script": "graders/my_script.py",
		},
	})
	require.NoError(t, err)
	require.Equal(t, "test_script", grader.Name())
	require.Equal(t, models.GraderKindScript, grader.Kind())
}

func TestScriptGrader_CreateMissingScript(t *testing.T) {
	_, err := Create(CreateArgs{
		GraderKind: models.GraderKindScript,
		Identifier: "bad",
		Config:     map[string]any{},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "requires a 'script' path")
}

func TestScriptGrader_ProgramKindRemoved(t *testing.T) {
	_, err := Create(CreateArgs{
		GraderKind: "program",
		Identifier: "test",
		Config:     map[string]any{},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "is not a valid grader type")
}
