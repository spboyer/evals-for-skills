package graders

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spboyer/waza/internal/models"
)

type ScriptGraderArgs struct {
	Path    string `mapstructure:"path"`
	SpecDir string `mapstructure:"-"`
}

// ScriptGrader launches an external script (.py or .js) and reads a
// {score, passed, message, details} JSON result from its stdout.
type ScriptGrader struct {
	name       string
	scriptPath string // relative path from eval spec directory
	scriptBin  string // "python" or "node"
}

// scriptOutput is the JSON contract for external grader scripts.
type scriptOutput struct {
	Score   float64        `json:"score"`
	Passed  bool           `json:"passed"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

func NewScriptGrader(name string, args ScriptGraderArgs) (*ScriptGrader, error) {
	ext := filepath.Ext(args.Path)

	var bin string
	switch ext {
	case ".py":
		bin = resolvePythonBin()
	case ".js":
		bin = "node"
	default:
		return nil, fmt.Errorf("unsupported script extension '%s' for script grader (supported: .py, .js)", ext)
	}

	scriptPath := filepath.Join(args.SpecDir, args.Path)

	return &ScriptGrader{
		name:       name,
		scriptPath: scriptPath,
		scriptBin:  bin,
	}, nil
}

func (sg *ScriptGrader) Name() string            { return sg.name }
func (sg *ScriptGrader) Kind() models.GraderKind { return models.GraderKindScript }

func (sg *ScriptGrader) Grade(ctx context.Context, gradingContext *Context) (*models.GraderResults, error) {
	return measureTime(func() (*models.GraderResults, error) {
		// TODO: fix
		stdinJSONText, err := getStdinTextForScript(gradingContext, []string{})
		if err != nil {
			return nil, fmt.Errorf("script grader '%s': failed to build stdin: %w", sg.name, err)
		}

		cmd := exec.CommandContext(ctx, sg.scriptBin, sg.scriptPath)
		cmd.Stdin = bytes.NewReader(stdinJSONText)
		cmd.Stderr = os.Stderr

		outputBytes, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("script grader '%s': script execution failed (%s): %w", sg.name, string(outputBytes), err)
		}

		var result scriptOutput
		if err := json.Unmarshal(outputBytes, &result); err != nil {
			return nil, fmt.Errorf("script grader '%s': failed to parse script output (%s): %w", sg.name, string(outputBytes), err)
		}

		return &models.GraderResults{
			Name:     sg.name,
			Type:     models.GraderKindScript,
			Score:    result.Score,
			Passed:   result.Passed,
			Feedback: result.Message,
			Details:  result.Details,
		}, nil
	})
}
