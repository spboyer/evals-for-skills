package tokens

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCount_TableFormat(t *testing.T) {
	cmd := newCountCmd()
	cmd.SetArgs([]string{"testdata"})
	out := new(bytes.Buffer)
	cmd.SetOut(out)
	require.NoError(t, cmd.Execute())

	expected := `File                          Tokens     Chars   Lines
------------------------------------------------------
testdata/README.md                 7        27       2
testdata/SKILL.md                402      1608      84
testdata/references/one.md         9        35       2
testdata/references/two.md        10        40       2
------------------------------------------------------
Total                            428      1710      90

4 file(s) scanned
`
	require.Equal(t, expected, out.String())
}

func TestCount_JSONFormat(t *testing.T) {
	out := new(bytes.Buffer)
	cmd := newCountCmd()
	cmd.SetOut(out)
	cmd.SetArgs([]string{"--format", "json", "testdata"})
	require.NoError(t, cmd.Execute())

	output := out.String()

	var result countJSONOutput
	require.NoError(t, json.Unmarshal([]byte(output), &result), "invalid JSON output: %s", output)

	expectedFiles := []string{"testdata/README.md", "testdata/SKILL.md", "testdata/references/one.md", "testdata/references/two.md"}
	require.Equal(t, len(expectedFiles), result.TotalFiles)
	require.Equal(t, 428, result.TotalTokens) // 7 + 402 + 9 + 10

	for _, f := range expectedFiles {
		require.Contains(t, result.Files, f)
	}

	expected := map[string]countFileEntry{
		"testdata/README.md":         {Tokens: 7, Characters: 27, Lines: 2},
		"testdata/SKILL.md":          {Tokens: 402, Characters: 1608, Lines: 84},
		"testdata/references/one.md": {Tokens: 9, Characters: 35, Lines: 2},
		"testdata/references/two.md": {Tokens: 10, Characters: 40, Lines: 2},
	}
	for file, want := range expected {
		got := result.Files[file]
		require.Equal(t, want.Tokens, got.Tokens, "%s tokens", file)
		require.Equal(t, want.Characters, got.Characters, "%s characters", file)
		require.Equal(t, want.Lines, got.Lines, "%s lines", file)
	}

	require.NotContains(t, result.Files, "testdata/scripts/sample.py")
}

func TestCount_SortByTokens(t *testing.T) {
	out := new(bytes.Buffer)
	cmd := newCountCmd()
	cmd.SetOut(out)
	cmd.SetArgs([]string{"--sort", "tokens", "testdata"})
	require.NoError(t, cmd.Execute())

	lines := strings.Split(out.String(), "\n")
	var dataLines []string
	for _, line := range lines {
		if strings.Contains(line, ".md") {
			dataLines = append(dataLines, line)
		}
	}
	require.GreaterOrEqual(t, len(dataLines), 2)
	require.Contains(t, dataLines[0], "SKILL.md", "SKILL.md should be first when sorted by tokens")
}

func TestCount_SortByName(t *testing.T) {
	out := new(bytes.Buffer)
	cmd := newCountCmd()
	cmd.SetOut(out)
	cmd.SetArgs([]string{"--sort", "name", "testdata"})
	require.NoError(t, cmd.Execute())

	output := out.String()

	lines := strings.Split(output, "\n")
	var dataLines []string
	for _, line := range lines {
		if strings.Contains(line, ".md") {
			dataLines = append(dataLines, line)
		}
	}

	require.GreaterOrEqual(t, len(dataLines), 2)
	require.Contains(t, dataLines[0], "one.md", "first file alphabetically by name")
}

func TestCount_MinTokens(t *testing.T) {
	out := new(bytes.Buffer)
	cmd := newCountCmd()
	cmd.SetOut(out)
	cmd.SetArgs([]string{"--format", "json", "--min-tokens", "100", "testdata"})
	require.NoError(t, cmd.Execute())

	var result countJSONOutput
	require.NoError(t, json.Unmarshal(out.Bytes(), &result))

	for file, entry := range result.Files {
		require.GreaterOrEqual(t, entry.Tokens, 100, "%s should have >= 100 tokens", file)
	}
}

func TestCount_NoTotal(t *testing.T) {
	out := new(bytes.Buffer)
	cmd := newCountCmd()
	cmd.SetOut(out)
	cmd.SetArgs([]string{"--no-total", "testdata"})
	require.NoError(t, cmd.Execute())

	output := out.String()

	require.NotContains(t, output, "Total")
	require.NotContains(t, output, "file(s) scanned")
}

func TestCount_SpecificPath(t *testing.T) {
	out := new(bytes.Buffer)
	cmd := newCountCmd()
	cmd.SetOut(out)
	cmd.SetArgs([]string{"--format", "json", filepath.Join("testdata", "SKILL.md")})
	require.NoError(t, cmd.Execute())

	var result countJSONOutput
	require.NoError(t, json.Unmarshal(out.Bytes(), &result))

	require.Equal(t, 1, result.TotalFiles)
	require.Equal(t, 402, result.TotalTokens)

	require.Contains(t, result.Files, "testdata/SKILL.md")
	entry := result.Files["testdata/SKILL.md"]
	require.Equal(t, 402, entry.Tokens)
	require.Equal(t, 1608, entry.Characters)
	require.Equal(t, 84, entry.Lines)
}

func TestCount_DirectoryPath(t *testing.T) {
	out := new(bytes.Buffer)
	cmd := newCountCmd()
	cmd.SetOut(out)
	cmd.SetArgs([]string{"--format", "json", filepath.Join("testdata", "references")})
	require.NoError(t, cmd.Execute())

	var result countJSONOutput
	require.NoError(t, json.Unmarshal(out.Bytes(), &result))

	require.Equal(t, 2, result.TotalFiles)
	require.Equal(t, 19, result.TotalTokens) // 9 + 10
}

func TestCount_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()

	out := new(bytes.Buffer)
	cmd := newCountCmd()
	cmd.SetOut(out)
	cmd.SetArgs([]string{dir})
	require.NoError(t, cmd.Execute())

	output := out.String()

	require.Contains(t, output, "No markdown files found")
}

func TestCount_NonMarkdownFilesExcluded(t *testing.T) {
	out := new(bytes.Buffer)
	cmd := newCountCmd()
	cmd.SetOut(out)
	cmd.SetArgs([]string{filepath.Join("testdata", "scripts")})
	require.NoError(t, cmd.Execute())

	output := out.String()

	require.Contains(t, output, "No markdown files found")
}

func TestCount_AbsoluteDirectoryPath(t *testing.T) {
	absDir, err := filepath.Abs(filepath.Join("testdata", "references"))
	require.NoError(t, err)

	out := new(bytes.Buffer)
	cmd := newCountCmd()
	cmd.SetOut(out)
	cmd.SetArgs([]string{"--format", "json", absDir})
	require.NoError(t, cmd.Execute())

	var result countJSONOutput
	require.NoError(t, json.Unmarshal(out.Bytes(), &result))

	require.Equal(t, 2, result.TotalFiles)
	require.Equal(t, 19, result.TotalTokens)
}

func TestCount_AbsoluteFilePath(t *testing.T) {
	absFile, err := filepath.Abs(filepath.Join("testdata", "SKILL.md"))
	require.NoError(t, err)

	out := new(bytes.Buffer)
	cmd := newCountCmd()
	cmd.SetOut(out)
	cmd.SetArgs([]string{"--format", "json", absFile})
	require.NoError(t, cmd.Execute())

	var result countJSONOutput
	require.NoError(t, json.Unmarshal(out.Bytes(), &result))

	require.Equal(t, 1, result.TotalFiles)
	require.Contains(t, result.Files, "testdata/SKILL.md")

	t.Run("multiple files", func(t *testing.T) {
		a, err := filepath.Abs(filepath.Join("testdata", "SKILL.md"))
		require.NoError(t, err)
		b, err := filepath.Abs(filepath.Join("testdata", "references", "one.md"))
		require.NoError(t, err)
		out := new(bytes.Buffer)
		cmd := newCountCmd()
		cmd.SetOut(out)
		cmd.SetArgs([]string{"--format", "json", a, b})
		require.NoError(t, cmd.Execute())

		var result countJSONOutput
		require.NoError(t, json.Unmarshal(out.Bytes(), &result))

		require.Equal(t, 2, result.TotalFiles)
		require.Contains(t, result.Files, "testdata/SKILL.md")
		require.Contains(t, result.Files, "testdata/references/one.md")
	})
}

func TestCount_SortWithJSONErrors(t *testing.T) {
	cmd := newCountCmd()
	cmd.SetOut(new(bytes.Buffer))
	cmd.SetArgs([]string{"--format", "json", "--sort", "tokens", "testdata"})
	require.ErrorContains(t, cmd.Execute(), "--sort is only supported with table output")
}

func TestCount_NoTotalWithJSONErrors(t *testing.T) {
	cmd := newCountCmd()
	cmd.SetOut(new(bytes.Buffer))
	cmd.SetArgs([]string{"--format", "json", "--no-total", "testdata"})
	require.ErrorContains(t, cmd.Execute(), "--no-total is only supported with table output")
}

func TestCount_NonexistentPath(t *testing.T) {
	out := new(bytes.Buffer)
	errOut := new(bytes.Buffer)
	cmd := newCountCmd()
	cmd.SetOut(out)
	cmd.SetErr(errOut)
	cmd.SetArgs([]string{"no-such-dir"})
	require.ErrorContains(t, cmd.Execute(), "no-such-dir")
}
