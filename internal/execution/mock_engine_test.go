package execution

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockEngine_Initialize(t *testing.T) {
	engine := NewMockEngine("test-model")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := engine.Initialize(ctx)
	require.NoError(t, err)
}

func TestMockEngine_Execute_WritesResources(t *testing.T) {
	engine := NewMockEngine("test-model")

	resp, err := engine.Execute(context.Background(), &ExecutionRequest{
		Message: "hello",
		Resources: []ResourceFile{{
			Path:    "fixtures/input.txt",
			Content: "test-content",
		}},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Contains(t, resp.FinalOutput, "Mock response for: hello")
	assert.Contains(t, resp.FinalOutput, "Analyzed 1 file(s)")

	content, err := os.ReadFile(filepath.Join(resp.WorkspaceDir, "fixtures", "input.txt"))
	require.NoError(t, err)
	assert.Equal(t, "test-content", string(content))

	require.NoError(t, engine.Shutdown(context.Background()))
}

func TestMockEngine_Execute_ReplacesWorkspace(t *testing.T) {
	engine := NewMockEngine("test-model")

	resp1, err := engine.Execute(context.Background(), &ExecutionRequest{Message: "one"})
	require.NoError(t, err)
	firstWorkspace := resp1.WorkspaceDir

	resp2, err := engine.Execute(context.Background(), &ExecutionRequest{Message: "two"})
	require.NoError(t, err)
	secondWorkspace := resp2.WorkspaceDir

	assert.NotEqual(t, firstWorkspace, secondWorkspace)

	// Both workspaces exist until Shutdown (safe for concurrent use)
	_, err1 := os.Stat(firstWorkspace)
	_, err2 := os.Stat(secondWorkspace)
	assert.NoError(t, err1, "first workspace should still exist before Shutdown")
	assert.NoError(t, err2, "second workspace should still exist before Shutdown")

	require.NoError(t, engine.Shutdown(context.Background()))

	// Both workspaces removed after Shutdown
	_, statErr := os.Stat(firstWorkspace)
	assert.True(t, os.IsNotExist(statErr), "first workspace should be removed after Shutdown")
	_, statErr2 := os.Stat(secondWorkspace)
	assert.True(t, os.IsNotExist(statErr2), "second workspace should be removed after Shutdown")
}

func TestMockEngine_Execute_SetupResourcesError(t *testing.T) {
	engine := NewMockEngine("test-model")

	resp, err := engine.Execute(context.Background(), &ExecutionRequest{
		Message: "hello",
		Resources: []ResourceFile{{
			Path:    "/absolute/path.txt",
			Content: "x",
		}},
	})
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to setup mock workspace resources")

	require.NoError(t, engine.Shutdown(context.Background()))
}
