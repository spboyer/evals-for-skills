package execution

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupWorkspaceResources_WritesFiles(t *testing.T) {
	workspace := t.TempDir()

	resources := []ResourceFile{
		{Path: "root.txt", Content: "root"},
		{Path: "nested/child.txt", Content: "child"},
		{Path: "", Content: "ignored"},
	}

	err := setupWorkspaceResources(workspace, resources)
	require.NoError(t, err)

	rootContent, err := os.ReadFile(filepath.Join(workspace, "root.txt"))
	require.NoError(t, err)
	assert.Equal(t, "root", string(rootContent))

	childContent, err := os.ReadFile(filepath.Join(workspace, "nested", "child.txt"))
	require.NoError(t, err)
	assert.Equal(t, "child", string(childContent))
}

func TestSetupWorkspaceResources_RejectsAbsolutePath(t *testing.T) {
	err := setupWorkspaceResources(t.TempDir(), []ResourceFile{{Path: "/etc/passwd", Content: "x"}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must be relative")
}

func TestSetupWorkspaceResources_RejectsPathTraversal(t *testing.T) {
	err := setupWorkspaceResources(t.TempDir(), []ResourceFile{{Path: "../outside.txt", Content: "x"}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "escapes workspace")
}

func TestSetupWorkspaceResources_EmptyWorkspace(t *testing.T) {
	err := setupWorkspaceResources("", []ResourceFile{{Path: "file.txt", Content: "x"}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "escapes workspace")
}
