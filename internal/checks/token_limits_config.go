package checks

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	// FallbackLimit is the default token limit when no pattern matches.
	FallbackLimit    = 2000
	maxPatternLength = 500
)

// TokenLimitsConfig holds token limit configuration.
type TokenLimitsConfig struct {
	Description string         `json:"description,omitempty"`
	Defaults    map[string]int `json:"defaults"`
	Overrides   map[string]int `json:"overrides"`
}

// DefaultLimits is the built-in fallback when neither .waza.yaml nor .token-limits.json provides limits.
var DefaultLimits = TokenLimitsConfig{
	Defaults: map[string]int{
		"SKILL.md":           500,
		"references/**/*.md": 1000,
		"docs/**/*.md":       1500,
		"*.md":               2000,
	},
	Overrides: map[string]int{
		"README.md":       3000,
		"CONTRIBUTING.md": 2500,
	},
}

// LoadLimitsConfig unmarshals dir/.token-limits.json or returns [DefaultLimits].
// Note: callers should prefer .waza.yaml tokens.limits; this is the fallback path.
func LoadLimitsConfig(dir string) (TokenLimitsConfig, error) {
	path := filepath.Join(dir, ".token-limits.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return DefaultLimits, nil
		}
		return TokenLimitsConfig{}, fmt.Errorf("reading %q: %w", path, err)
	}
	cfg := TokenLimitsConfig{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("error parsing limits %q: %w", path, err)
	}
	if cfg.Defaults == nil {
		return cfg, errors.New(`missing or invalid "defaults" field in limits file: ` + path)
	}
	if cfg.Overrides == nil {
		cfg.Overrides = make(map[string]int)
	}
	return cfg, nil
}

// normalizePath converts backslashes to forward slashes.
func normalizePath(filePath string) string {
	return strings.ReplaceAll(filePath, `\`, "/")
}

// globToRegex converts a glob pattern to a compiled regex.
func globToRegex(pattern string) (*regexp.Regexp, error) {
	if len(pattern) > maxPatternLength {
		return nil, fmt.Errorf("pattern too long (max %d characters)", maxPatternLength)
	}

	re := regexp.QuoteMeta(pattern)
	re = strings.ReplaceAll(re, `\*\*`, "{{GLOBSTAR}}")
	re = strings.ReplaceAll(re, `\*`, `[^/]*`)
	re = strings.ReplaceAll(re, "{{GLOBSTAR}}", `.*?`)

	if strings.Contains(pattern, "/") {
		return regexp.Compile(`^/?` + re + `$`)
	}
	return regexp.Compile(`(?:^|/)` + re + `$`)
}

// matchesPattern checks whether filePath matches a glob pattern.
func matchesPattern(filePath, pattern string) bool {
	normalized := normalizePath(filePath)

	if !strings.Contains(pattern, "/") && !strings.Contains(pattern, "*") {
		return normalized == pattern || strings.HasSuffix(normalized, "/"+pattern)
	}

	re, err := globToRegex(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(normalized)
}

// patternSpecificity scores a glob pattern; higher means more specific.
func patternSpecificity(pattern string) int {
	score := 0

	if !strings.Contains(pattern, "*") {
		score += 10000
	}

	score += (strings.Count(pattern, "/")) * 100

	temp := strings.ReplaceAll(pattern, "**", "")
	score += strings.Count(temp, "*") * 10

	score -= strings.Count(pattern, "**") * 50
	score += len(pattern)

	return score
}

// LimitResult holds a resolved limit and the pattern that matched.
type LimitResult struct {
	Limit   int
	Pattern string
}

// GetLimitForFile determines the token limit for a file.
// An optional workspaceRelPrefix (e.g. "plugin/skills/azure-deploy") is
// prepended to filePath so workspace-root-relative patterns such as
// "plugin/skills/**/SKILL.md" can match alongside skill-relative patterns.
// Patterns are tried in specificity order; for each pattern both the
// skill-relative and workspace-relative paths are checked.
func GetLimitForFile(filePath string, cfg TokenLimitsConfig, workspaceRelPrefix ...string) LimitResult {
	normalized := normalizePath(filePath)

	// Compute workspace-relative path when a prefix is provided.
	prefixed := ""
	if len(workspaceRelPrefix) > 0 && workspaceRelPrefix[0] != "" {
		prefixed = normalizePath(workspaceRelPrefix[0] + "/" + filePath)
	}

	for overridePath, limit := range cfg.Overrides {
		if normalized == overridePath || strings.HasSuffix(normalized, "/"+overridePath) {
			return LimitResult{Limit: limit, Pattern: overridePath}
		}
		if prefixed != "" && (prefixed == overridePath || strings.HasSuffix(prefixed, "/"+overridePath)) {
			return LimitResult{Limit: limit, Pattern: overridePath}
		}
	}

	type entry struct {
		pattern string
		limit   int
	}
	entries := make([]entry, 0, len(cfg.Defaults))
	for p, l := range cfg.Defaults {
		entries = append(entries, entry{p, l})
	}
	sort.Slice(entries, func(i, j int) bool {
		return patternSpecificity(entries[i].pattern) > patternSpecificity(entries[j].pattern)
	})

	for _, e := range entries {
		if matchesPattern(normalized, e.pattern) {
			return LimitResult{Limit: e.limit, Pattern: e.pattern}
		}
		if prefixed != "" && matchesPattern(prefixed, e.pattern) {
			return LimitResult{Limit: e.limit, Pattern: e.pattern}
		}
	}

	return LimitResult{Limit: FallbackLimit, Pattern: "none"}
}
