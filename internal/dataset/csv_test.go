package dataset

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeCSV(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(p, []byte(content), 0o644))
	return p
}

func TestLoadCSV(t *testing.T) {
	tests := []struct {
		name     string
		csv      string
		wantRows int
		wantCols int
		wantErr  string
	}{
		{
			name:     "happy path 3 rows 3 columns",
			csv:      "name,prompt,model\nauth-fix,Fix the auth bug,claude-sonnet-4.5\ndb-migration,Refactor DB logic,gpt-5.2-codex\nui-polish,Polish the UI,gpt-4o\n",
			wantRows: 3,
			wantCols: 3,
		},
		{
			name:     "single row",
			csv:      "id,prompt\nonly-one,Do something\n",
			wantRows: 1,
			wantCols: 2,
		},
		{
			name:     "empty CSV headers only",
			csv:      "name,prompt,model\n",
			wantRows: 0,
			wantCols: 0,
		},
		{
			name:    "mismatched column count",
			csv:     "name,prompt\nok,fine\nbad\n",
			wantErr: "wrong number of fields",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := writeCSV(t, dir, "test.csv", tt.csv)

			rows, err := LoadCSV(path)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Len(t, rows, tt.wantRows)
			if tt.wantRows > 0 {
				assert.Len(t, rows[0], tt.wantCols)
			}
		})
	}
}

func TestLoadCSV_HappyPathValues(t *testing.T) {
	dir := t.TempDir()
	path := writeCSV(t, dir, "data.csv", "name,prompt,model\nauth-fix,Fix the auth bug,claude-sonnet-4.5\ndb-migration,Refactor DB logic,gpt-5.2-codex\n")

	rows, err := LoadCSV(path)
	require.NoError(t, err)
	require.Len(t, rows, 2)

	assert.Equal(t, "auth-fix", rows[0]["name"])
	assert.Equal(t, "Fix the auth bug", rows[0]["prompt"])
	assert.Equal(t, "claude-sonnet-4.5", rows[0]["model"])

	assert.Equal(t, "db-migration", rows[1]["name"])
	assert.Equal(t, "Refactor DB logic", rows[1]["prompt"])
	assert.Equal(t, "gpt-5.2-codex", rows[1]["model"])
}

func TestLoadCSV_MissingFile(t *testing.T) {
	_, err := LoadCSV("/nonexistent/path/data.csv")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "csv: open")
}

func TestLoadCSVRange(t *testing.T) {
	tests := []struct {
		name     string
		csv      string
		start    int
		end      int
		wantRows int
		wantErr  string
	}{
		{
			name:     "range 2-3 of 5",
			csv:      "name,prompt\na,p1\nb,p2\nc,p3\nd,p4\ne,p5\n",
			start:    2,
			end:      3,
			wantRows: 2,
		},
		{
			name:     "range 1-1 single row",
			csv:      "name,prompt\na,p1\nb,p2\n",
			start:    1,
			end:      1,
			wantRows: 1,
		},
		{
			name:     "range beyond available rows clamps",
			csv:      "name,prompt\na,p1\nb,p2\n",
			start:    1,
			end:      100,
			wantRows: 2,
		},
		{
			name:     "start beyond available returns empty",
			csv:      "name,prompt\na,p1\n",
			start:    5,
			end:      10,
			wantRows: 0,
		},
		{
			name:    "invalid range start < 1",
			csv:     "name,prompt\na,p1\n",
			start:   0,
			end:     1,
			wantErr: "range start must be >= 1",
		},
		{
			name:    "invalid range end < start",
			csv:     "name,prompt\na,p1\n",
			start:   3,
			end:     1,
			wantErr: "range end (1) must be >= start (3)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := writeCSV(t, dir, "test.csv", tt.csv)

			rows, err := LoadCSVRange(path, tt.start, tt.end)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Len(t, rows, tt.wantRows)
		})
	}
}

func TestLoadCSVRange_Values(t *testing.T) {
	dir := t.TempDir()
	path := writeCSV(t, dir, "data.csv", "name,prompt\na,p1\nb,p2\nc,p3\nd,p4\ne,p5\n")

	rows, err := LoadCSVRange(path, 2, 3)
	require.NoError(t, err)
	require.Len(t, rows, 2)

	assert.Equal(t, "b", rows[0]["name"])
	assert.Equal(t, "p2", rows[0]["prompt"])
	assert.Equal(t, "c", rows[1]["name"])
	assert.Equal(t, "p3", rows[1]["prompt"])
}
