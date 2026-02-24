// Package projectconfig provides the ProjectConfig struct and loader for
// .waza.yaml project-level configuration files.
package projectconfig

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// PathsConfig holds directory paths for skills, evals, and results.
type PathsConfig struct {
	Skills  string `yaml:"skills,omitempty"`
	Evals   string `yaml:"evals,omitempty"`
	Results string `yaml:"results,omitempty"`
}

// DefaultsConfig holds default execution parameters.
type DefaultsConfig struct {
	Engine     string `yaml:"engine,omitempty"`
	Model      string `yaml:"model,omitempty"`
	JudgeModel string `yaml:"judge_model,omitempty"`
	Timeout    int    `yaml:"timeout,omitempty"`
	Parallel   *bool  `yaml:"parallel,omitempty"`
	Workers    int    `yaml:"workers,omitempty"`
	Verbose    *bool  `yaml:"verbose,omitempty"`
	SessionLog *bool  `yaml:"session_log,omitempty"`
}

// CacheConfig holds cache settings.
type CacheConfig struct {
	Enabled *bool  `yaml:"enabled,omitempty"`
	Dir     string `yaml:"dir,omitempty"`
}

// ServerConfig holds dashboard server settings.
type ServerConfig struct {
	Port       int    `yaml:"port,omitempty"`
	ResultsDir string `yaml:"results_dir,omitempty"`
}

// DevConfig holds waza dev command settings.
type DevConfig struct {
	Model         string `yaml:"model,omitempty"`
	Target        string `yaml:"target,omitempty"`
	MaxIterations int    `yaml:"max_iterations,omitempty"`
}

// TokenLimitsConfig holds per-model token limit maps.
type TokenLimitsConfig struct {
	Defaults  map[string]int `yaml:"defaults,omitempty"`
	Overrides map[string]int `yaml:"overrides,omitempty"`
}

// TokensConfig holds token budget settings.
type TokensConfig struct {
	WarningThreshold int               `yaml:"warning_threshold,omitempty"`
	FallbackLimit    int               `yaml:"fallback_limit,omitempty"`
	Limits           *TokenLimitsConfig `yaml:"limits,omitempty"`
}

// GradersConfig holds grader execution settings.
type GradersConfig struct {
	ProgramTimeout int `yaml:"program_timeout,omitempty"`
}

// ProjectConfig is the top-level configuration loaded from .waza.yaml.
type ProjectConfig struct {
	Paths    PathsConfig    `yaml:"paths,omitempty"`
	Defaults DefaultsConfig `yaml:"defaults,omitempty"`
	Cache    CacheConfig    `yaml:"cache,omitempty"`
	Server   ServerConfig   `yaml:"server,omitempty"`
	Dev      DevConfig      `yaml:"dev,omitempty"`
	Tokens   TokensConfig   `yaml:"tokens,omitempty"`
	Graders  GradersConfig  `yaml:"graders,omitempty"`
}

// New returns a ProjectConfig with all hard-coded defaults populated.
func New() *ProjectConfig {
	return &ProjectConfig{
		Paths: PathsConfig{
			Skills:  "skills/",
			Evals:   "evals/",
			Results: "results/",
		},
		Defaults: DefaultsConfig{
			Engine:     "copilot-sdk",
			Model:      "claude-sonnet-4.6",
			JudgeModel: "",
			Timeout:    300,
			Parallel:   boolPtr(false),
			Workers:    4,
			Verbose:    boolPtr(false),
			SessionLog: boolPtr(false),
		},
		Cache: CacheConfig{
			Enabled: boolPtr(false),
			Dir:     ".waza-cache",
		},
		Server: ServerConfig{
			Port:       3000,
			ResultsDir: ".",
		},
		Dev: DevConfig{
			Model:         "claude-sonnet-4-20250514",
			Target:        "medium-high",
			MaxIterations: 5,
		},
		Tokens: TokensConfig{
			WarningThreshold: 2500,
			FallbackLimit:    2000,
		},
		Graders: GradersConfig{
			ProgramTimeout: 30,
		},
	}
}

// Load finds .waza.yaml by walking up from startDir (max 10 levels),
// unmarshals it, and fills in missing fields with defaults.
// If no config file is found, returns defaults with a nil error.
func Load(startDir string) (*ProjectConfig, error) {
	cfg := New()

	data, err := findConfigFile(startDir)
	if err != nil {
		return cfg, nil // no file found → return defaults
	}

	var fileCfg ProjectConfig
	if err := yaml.Unmarshal(data, &fileCfg); err != nil {
		return nil, err
	}

	// Merge file values onto defaults.
	mergeConfig(cfg, &fileCfg)
	return cfg, nil
}

// findConfigFile walks up from dir looking for .waza.yaml (max 10 levels).
func findConfigFile(dir string) ([]byte, error) {
	for i := 0; i < 10; i++ {
		p := filepath.Join(dir, ".waza.yaml")
		data, err := os.ReadFile(p)
		if err == nil {
			return data, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break // reached filesystem root
		}
		dir = parent
	}
	return nil, os.ErrNotExist
}

// mergeConfig overlays non-zero values from src onto dst.
func mergeConfig(dst, src *ProjectConfig) {
	// Paths
	if src.Paths.Skills != "" {
		dst.Paths.Skills = src.Paths.Skills
	}
	if src.Paths.Evals != "" {
		dst.Paths.Evals = src.Paths.Evals
	}
	if src.Paths.Results != "" {
		dst.Paths.Results = src.Paths.Results
	}

	// Defaults
	if src.Defaults.Engine != "" {
		dst.Defaults.Engine = src.Defaults.Engine
	}
	if src.Defaults.Model != "" {
		dst.Defaults.Model = src.Defaults.Model
	}
	if src.Defaults.JudgeModel != "" {
		dst.Defaults.JudgeModel = src.Defaults.JudgeModel
	}
	if src.Defaults.Timeout != 0 {
		dst.Defaults.Timeout = src.Defaults.Timeout
	}
	if src.Defaults.Parallel != nil {
		dst.Defaults.Parallel = src.Defaults.Parallel
	}
	if src.Defaults.Workers != 0 {
		dst.Defaults.Workers = src.Defaults.Workers
	}
	if src.Defaults.Verbose != nil {
		dst.Defaults.Verbose = src.Defaults.Verbose
	}
	if src.Defaults.SessionLog != nil {
		dst.Defaults.SessionLog = src.Defaults.SessionLog
	}

	// Cache
	if src.Cache.Enabled != nil {
		dst.Cache.Enabled = src.Cache.Enabled
	}
	if src.Cache.Dir != "" {
		dst.Cache.Dir = src.Cache.Dir
	}

	// Server
	if src.Server.Port != 0 {
		dst.Server.Port = src.Server.Port
	}
	if src.Server.ResultsDir != "" {
		dst.Server.ResultsDir = src.Server.ResultsDir
	}

	// Dev
	if src.Dev.Model != "" {
		dst.Dev.Model = src.Dev.Model
	}
	if src.Dev.Target != "" {
		dst.Dev.Target = src.Dev.Target
	}
	if src.Dev.MaxIterations != 0 {
		dst.Dev.MaxIterations = src.Dev.MaxIterations
	}

	// Tokens
	if src.Tokens.WarningThreshold != 0 {
		dst.Tokens.WarningThreshold = src.Tokens.WarningThreshold
	}
	if src.Tokens.FallbackLimit != 0 {
		dst.Tokens.FallbackLimit = src.Tokens.FallbackLimit
	}
	if src.Tokens.Limits != nil {
		dst.Tokens.Limits = src.Tokens.Limits
	}

	// Graders
	if src.Graders.ProgramTimeout != 0 {
		dst.Graders.ProgramTimeout = src.Graders.ProgramTimeout
	}
}

func boolPtr(b bool) *bool {
	return &b
}
