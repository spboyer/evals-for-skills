package tokens

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindMarkdownFiles(t *testing.T) {
	dir := t.TempDir()

	require.NoError(t, os.WriteFile(filepath.Join(dir, "README.md"), []byte("# Hi"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "doc.mdx"), []byte("# Doc"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "script.ts"), []byte("//"), 0644))

	sub := filepath.Join(dir, "sub")
	require.NoError(t, os.MkdirAll(sub, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(sub, "nested.md"), []byte("# Nested"), 0644))

	nm := filepath.Join(dir, "node_modules")
	require.NoError(t, os.MkdirAll(nm, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(nm, "excluded.md"), []byte("# No"), 0644))

	files, err := findMarkdownFiles(nil, dir)
	require.NoError(t, err)
	sort.Strings(files)

	require.Len(t, files, 3)

	for _, f := range files {
		require.NotEqual(t, "excluded.md", filepath.Base(f), "node_modules/excluded.md should be excluded")
	}
}

func TestFindMarkdownFilesEmpty(t *testing.T) {
	dir := t.TempDir()
	files, err := findMarkdownFiles(nil, dir)
	require.NoError(t, err)
	require.Empty(t, files)
}

func TestFindMarkdownFilesNonexistent(t *testing.T) {
	_, err := findMarkdownFiles(nil, "/nonexistent/path")
	require.Error(t, err)
}
