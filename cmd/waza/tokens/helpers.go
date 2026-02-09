package tokens

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

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
