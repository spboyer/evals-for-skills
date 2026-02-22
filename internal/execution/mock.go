package execution

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	copilot "github.com/github/copilot-sdk/go"
	"github.com/spboyer/waza/internal/models"
)

// MockEngine is a simple mock implementation for testing
type MockEngine struct {
	modelID    string
	mu         sync.Mutex
	workspaces []string
}

// NewMockEngine creates a new mock engine
func NewMockEngine(modelID string) *MockEngine {
	return &MockEngine{
		modelID: modelID,
	}
}

func (m *MockEngine) Initialize(ctx context.Context) error {
	return nil
}

func (m *MockEngine) Execute(ctx context.Context, req *ExecutionRequest) (*ExecutionResponse, error) {
	start := time.Now()

	// Create a temp workspace so graders that inspect files (e.g. FileGrader) have
	// a directory to work with, mirroring CopilotEngine behavior.
	tmpDir, err := os.MkdirTemp("", "waza-mock-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create mock workspace: %w", err)
	}

	m.mu.Lock()
	m.workspaces = append(m.workspaces, tmpDir)
	m.mu.Unlock()

	// Write request resources into the workspace
	if err := setupWorkspaceResources(tmpDir, req.Resources); err != nil {
		return nil, fmt.Errorf("failed to setup mock workspace resources: %w", err)
	}

	// Simple mock response
	output := fmt.Sprintf("Mock response for: %s", req.Message)

	// Add some context if files are present
	if len(req.Resources) > 0 {
		output += fmt.Sprintf("\nAnalyzed %d file(s)", len(req.Resources))
	}

	resp := &ExecutionResponse{
		FinalOutput:  output,
		Events:       []copilot.SessionEvent{},
		ModelID:      m.modelID,
		DurationMs:   time.Since(start).Milliseconds(),
		ToolCalls:    []models.ToolCall{},
		Success:      true,
		WorkspaceDir: tmpDir,
	}

	return resp, nil
}

func (m *MockEngine) Shutdown(ctx context.Context) error {
	m.mu.Lock()
	dirs := m.workspaces
	m.workspaces = nil
	m.mu.Unlock()

	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			return fmt.Errorf("failed to remove mock workspace %s: %w", dir, err)
		}
	}
	return nil
}
