package execution

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	copilot "github.com/github/copilot-sdk/go"
	"github.com/spboyer/waza/internal/utils"
)

// CopilotEngine integrates with GitHub Copilot SDK
type CopilotEngine struct {
	defaultModelID string

	client copilotClient

	workspacesMu sync.Mutex
	workspaces   []string // workspaces to clean up at Shutdown
}

// CopilotEngineBuilder builds a CopilotEngine with options
type CopilotEngineBuilder struct {
	engine *CopilotEngine
}

type CopilotEngineBuilderOptions struct {
	NewCopilotClient func(clientOptions *copilot.ClientOptions) copilotClient
}

// NewCopilotEngineBuilder creates a builder for CopilotEngine
//   - defaultModelID - used if no model ID is specified in session creation. Can be blank, which means the copilot
//     CLI will choose its own fallback model.
func NewCopilotEngineBuilder(defaultModelID string, options *CopilotEngineBuilderOptions) *CopilotEngineBuilder {
	var client copilotClient

	copilotOptions := &copilot.ClientOptions{
		// workspace is set at the session level, instead of at the client.
		LogLevel: "error",
	}

	if options == nil || options.NewCopilotClient == nil {
		client = newCopilotClient(copilotOptions)
	} else {
		client = options.NewCopilotClient(copilotOptions)
	}

	builder := &CopilotEngineBuilder{
		engine: &CopilotEngine{
			defaultModelID: defaultModelID,
		},
	}

	builder.engine.client = client
	return builder
}

func (b *CopilotEngineBuilder) Build() *CopilotEngine {
	return b.engine
}

// Initialize sets up the Copilot client
func (e *CopilotEngine) Initialize(ctx context.Context) error {
	// NOTE: the copilot client auto-starts on first call, so no need to explicitly init
	// or start it here.
	return nil
}

// Execute runs a test with Copilot SDK
func (e *CopilotEngine) Execute(ctx context.Context, req *ExecutionRequest) (*ExecutionResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("nil req was passed to CopilotEngine.Execute")
	}

	modelID, sourceDir, err := e.extractReqParams(req)

	if err != nil {
		return nil, err
	}

	start := time.Now()

	workspaceDir, err := e.setupWorkspace(req.Resources)

	if err != nil {
		return nil, err
	}

	// Build skill directories list: start with CWD, then add any from request
	skillDirs := e.getSkillDirs(sourceDir, req)

	ctx, cancel := context.WithTimeout(ctx, req.Timeout)
	defer cancel()

	var session copilotSession

	if req.SessionID == "" {
		// Create session with updated API
		session, err = e.client.CreateSession(ctx, &copilot.SessionConfig{
			Model:            modelID,
			SkillDirectories: skillDirs,
			WorkingDirectory: workspaceDir,
		})

		if err != nil {
			return nil, fmt.Errorf("failed to create session: %w", err)
		}
	} else {
		session, err = e.client.ResumeSessionWithOptions(ctx, req.SessionID, &copilot.ResumeSessionConfig{
			Model: modelID,
			// these are the directory for the skill itself.
			SkillDirectories: skillDirs,
			WorkingDirectory: workspaceDir,
		})

		if err != nil {
			return nil, fmt.Errorf("failed to resume session (%s): %w", req.SessionID, err)
		}
	}

	eventsCollector := NewSessionEventsCollector()

	// Event handler with updated API
	unsubscribe := session.On(eventsCollector.On)
	defer unsubscribe()

	unsubscribe = session.On(utils.SessionToSlog)
	defer unsubscribe()

	// Send prompt with updated API
	_, err = session.SendAndWait(ctx, copilot.MessageOptions{
		Prompt: req.Message,
	})

	var errMsg string

	if err != nil {
		// errors that are returned inline, as part of the conversation, also come back
		// in the returned error. Rather than having one of those fun functions that returns
		// both an error and a result, I'll just put the error message in the ExecutionResponse.
		errMsg = err.Error()
	}

	duration := time.Since(start)

	// Build response
	resp := &ExecutionResponse{
		FinalOutput:      joinStrings(eventsCollector.OutputParts()),
		Events:           eventsCollector.SessionEvents(),
		ModelID:          modelID,
		SkillInvocations: eventsCollector.SkillInvocations,
		DurationMs:       duration.Milliseconds(),
		ToolCalls:        eventsCollector.ToolCalls(),
		ErrorMsg:         errMsg,
		Success:          err == nil,
		WorkspaceDir:     workspaceDir,
		SessionID:        session.SessionID(),
	}

	return resp, nil
}

// Shutdown cleans up resources
func (e *CopilotEngine) Shutdown(ctx context.Context) error {
	if err := e.client.Stop(); err != nil {
		// Log but continue cleanup
		slog.Info("failed to stop client", "error", err)
	}

	// remove the workspace folders - should be safe now that all the copilot sessions are shut down
	// and the tests are complete.
	workspaces := func() []string {
		e.workspacesMu.Lock()
		defer e.workspacesMu.Unlock()
		workspaces := e.workspaces
		e.workspaces = nil
		return workspaces
	}()

	for _, ws := range workspaces {
		if ws != "" {
			if err := os.RemoveAll(ws); err != nil {
				// errors here probably indicate some issue with our code continuing to lock files
				// even after tests have completed...
				slog.Warn("failed to cleanup stale workspace", "path", ws, "error", err)
			}
		}
	}

	return nil
}

func (e *CopilotEngine) extractReqParams(req *ExecutionRequest) (modelID string, sourceDir string, err error) {
	modelID = e.defaultModelID

	if req.ModelID != "" {
		modelID = req.ModelID // override the default model for the engine
	}

	sourceDir = req.SourceDir

	if req.SourceDir == "" {
		cwd, err := os.Getwd()

		if err != nil {
			return "", "", fmt.Errorf("failed to get current directory: %w", err)
		}

		sourceDir = cwd
	}

	if req.Timeout <= 0 {
		return "", "", fmt.Errorf("positive Timeout is required")
	}

	return modelID, sourceDir, nil
}

func (*CopilotEngine) getSkillDirs(cwd string, req *ExecutionRequest) []string {
	skillDirs := []string{cwd}

	seen := map[string]bool{
		cwd: true,
	}

	// Add skill directories from request, avoiding duplicates
	for _, path := range req.SkillPaths {
		if !seen[path] {
			seen[path] = true
			skillDirs = append(skillDirs, path)
		} else {
			slog.Warn("Skill directory included more than once in request", "path", path)
		}
	}

	// Log skill directories in verbose mode
	for _, dir := range skillDirs {
		slog.Debug("Adding skill directory", "path", dir)
	}

	return skillDirs
}

func (e *CopilotEngine) setupWorkspace(resources []ResourceFile) (string, error) {
	workspaceDir, err := os.MkdirTemp("", "waza-*")

	if err != nil {
		return "", fmt.Errorf("failed to create temp workspace: %w", err)
	}

	e.workspacesMu.Lock()
	e.workspaces = append(e.workspaces, workspaceDir)
	e.workspacesMu.Unlock()

	// Write resource files to workspace
	if err := setupWorkspaceResources(workspaceDir, resources); err != nil {
		return "", fmt.Errorf("failed to setup resources at workspace %s: %w", workspaceDir, err)
	}

	return workspaceDir, nil
}

func joinStrings(parts []string) string {
	var builder strings.Builder
	for _, p := range parts {
		builder.WriteString(p)
	}
	return builder.String()
}
