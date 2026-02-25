package tokens

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spboyer/waza/internal/checks"
	"github.com/spboyer/waza/internal/projectconfig"
)

// resolveLimitsConfig loads token limits config with proper precedence:
// 1. .waza.yaml tokens.limits section
// 2. .token-limits.json file
// 3. built-in DefaultLimits
func resolveLimitsConfig(rootDir string) (checks.TokenLimitsConfig, error) {
	cfg, err := projectconfig.Load(rootDir)
	if err != nil {
		return checks.TokenLimitsConfig{}, fmt.Errorf("loading project config: %w", err)
	}
	if cfg.Tokens.Limits != nil && cfg.Tokens.Limits.Defaults != nil {
		return checks.TokenLimitsConfig{
			Defaults:  cfg.Tokens.Limits.Defaults,
			Overrides: cfg.Tokens.Limits.Overrides,
		}, nil
	}
	return checks.LoadLimitsConfig(rootDir)
}

// FileResult holds token count results for a single file.
type FileResult struct {
	Path       string
	Tokens     int
	Characters int
	Lines      int
}

// nowISO returns the current time in ISO 8601 format.
func nowISO() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}

// countLines returns the number of lines in s. An empty string has 0 lines.
// A trailing newline does not count as an additional line (matches wc -l behavior
// for files that end with a newline).
func countLines(s string) int {
	if s == "" {
		return 0
	}
	n := strings.Count(s, "\n")
	if !strings.HasSuffix(s, "\n") {
		n++
	}
	return n
}

var excludedDirs = map[string]bool{
	"node_modules": true,
	".git":         true,
	"dist":         true,
	"coverage":     true,
}

// findMarkdownFiles takes a user-provided path (file or directory) and
// returns a list of markdown file paths. If paths is empty, scans rootDir.
func findMarkdownFiles(paths []string, rootDir string) ([]string, error) {
	if len(paths) == 0 {
		paths = []string{rootDir}
	}

	var result []string
	for _, p := range paths {
		if !filepath.IsAbs(p) {
			p = filepath.Join(rootDir, p)
		}

		info, err := os.Stat(p)
		if err != nil {
			return nil, fmt.Errorf("stat %q: %w", p, err)
		}

		if !info.IsDir() {
			result = append(result, p)
			continue
		}

		err = filepath.WalkDir(p, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() && excludedDirs[d.Name()] {
				return filepath.SkipDir
			}
			if !d.IsDir() {
				switch strings.ToLower(filepath.Ext(d.Name())) {
				case ".md", ".mdx":
					result = append(result, path)
				}
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("walking %q: %w", p, err)
		}
	}

	return result, nil
}
