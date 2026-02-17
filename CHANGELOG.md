# Changelog

All notable changes to waza will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **`prompt` grader (LLM-as-judge)** — LLM-based evaluation with rubrics, tool-based grading, and session management modes (#177, closes #104)
  - Two modes: `clean` (fresh context) and `continue_session` (resumes test session)
  - Tool-based grading: `set_waza_grade_pass` and `set_waza_grade_fail` tools for LLM graders
  - Separate judge model configuration: run evaluation with a different model than the executor
  - Pre-built rubric templates adapted from Azure ML evaluators
- **`trigger_tests.yaml` auto-discovery** — measure prompt trigger accuracy for skills (#166, closes #36)
  - New `internal/trigger/` package for trigger testing
  - Automatically discovered alongside `eval.yaml`
  - Confidence weighting: `high` (weight 1.0) and `medium` (weight 0.5) for borderline cases
  - `trigger_accuracy` metric with configurable cutoff threshold
  - Metrics: accuracy, precision, recall, F1, error count
- **`diff` grader** — new grader type for workspace file comparison with snapshot matching and contains-line fragment checks (#158)
- **Azure ML evaluation rubrics** — 8 pre-built rubric YAMLs in `examples/rubrics/` adapted from Azure ML evaluators (#160, #161):
  - Tool call rubrics: `tool_call_accuracy`, `tool_selection`, `tool_input_accuracy`, `tool_output_utilization`
  - Task evaluation rubrics: `task_completion`, `task_adherence`, `intent_resolution`, `response_completeness`
- **MockEngine WorkspaceDir support** — test infrastructure for graders that need workspace access (#159)

### Fixed

- **install.sh macOS checksum** — added `shasum -a 256` fallback for macOS (which lacks `sha256sum`) (#163)

## [0.4.0-alpha.1] - 2026-02-17

### Added

- **Go cross-platform release pipeline** — `go-release.yml` workflow builds binaries for linux/darwin/windows on amd64 and arm64 (#155)
- **`install.sh` installer** — one-line binary install with checksum verification: `curl -fsSL https://raw.githubusercontent.com/spboyer/waza/main/install.sh | bash`
- **`skill_invocation` grader** — validates orchestration workflows by checking which skills were invoked (#146)
- **`required_skills` preflight validation** — verifies skill dependencies before evaluation (#147)
- **Multi-model `--model` flag** — run evaluations across multiple models in a single command (#39)
- **`waza check` command** — skill submission readiness checks (#151)
- **Evaluation result caching** — incremental testing with cache invalidation (#150)
- **GitHub PR comment reporter** — post eval results as PR comments (#140)
- **Skills CI integration** — GitHub Actions workflow for microsoft/skills (#141)

### Fixed

- **Engine shutdown leak** — `runSingleModel()` now calls `engine.Shutdown(context.Background())` via defer after engine creation (#153, #154)

### Changed

- **Python release deprecated** — the Python release workflow is no longer maintained; Go binaries are the official distribution
- **First Go binary release** — v0.4.0-alpha.1 is the first release distributed as pre-built binaries

## [0.3.0] - 2026-02-13

### Added

- Grader showcase examples demonstrating all grader types (#134)
- Reusable GitHub Actions workflow for waza evaluations (#132)
- Documentation for prompt and action_sequence grader types (#133)
- Documentation for `waza dev` command and compliance scoring (#131)
- Auto-loading of skills for testing (#129)
- Debug logging support (`--debug` flag) (#130)

### Fixed

- Always output test run errors to help debug failures (#128)
- Include cwd as a skill folder when running waza (workspace fix)

### Changed

- Exit codes for CI/CD integration: 0=success, 1=test failure, 2=config error (#135)
- Reordered azd-publish skill workflow steps (#127)
- Auto-merge bot registry PRs in release workflow

## [0.2.1] - 2026-02-12

### Added

- `waza dev` command for interactive skill development and testing (#117)
- Prerelease input to azd publish workflow
- CHANGELOG.md as release notes source for azd extension releases
- `waza generate --skill <name>` - Filter to specific skill when using `--repo` or `--scan`

### Fixed

- Fixed azd extensions documentation link
- Corrected `azd ext source add` command syntax
- Branch release PR from origin/main to avoid workflow permission error (#121)

### Changed

- Removed path filters from Go CI to unblock non-code PRs
- Removed auto-merge from azd publish PR workflow
- Added azd extension installation instructions to README

## [0.2.0] - 2026-02-02

### Added

- **Skill Discovery** (#3)
  - `waza generate --repo <org/repo>` - Scan GitHub repos for SKILL.md files
  - `waza generate --scan` - Scan local directory for skills
  - `waza generate --all` - Generate evals for all discovered skills (CI-friendly)
  - Interactive skill selection with checkboxes when not using `--all`

- **GitHub Issue Creation** (#3)
  - Post-run prompt to create GitHub issues with eval results
  - Options: create for failed tasks only, all tasks, or none
  - Issues include results table, failed task details, and suggestions
  - `--no-issues` flag to skip prompts (CI-friendly)

- **New Modules**
  - `waza/scanner.py` - Skill discovery from GitHub repos and local directories
  - `waza/issues.py` - GitHub issue creation and formatting

### Changed

- Improved documentation with new feature guides
- Added skill discovery section to DEMO-SCRIPT.md
- Updated TUTORIAL.md with discovery and issue creation steps

## [0.1.0] - 2026-02-02

### Changed

- **Renamed project from `skill-eval` to `waza`** (技 - Japanese for "technique/skill")
  - New CLI command: `waza` (previously `skill-eval`)
  - New package name: `waza` (previously `skill_eval`)
  - Repository renamed to `waza`
- Bumped version to 0.1.0 to mark the rename milestone

### Migration

If you were using `skill-eval`, update your scripts:

```bash
# Old
skill-eval run ./eval.yaml
pip install skill-eval

# New
waza run ./eval.yaml
pip install waza
```

## [0.0.2] - 2026-02-01

### Added

- `--suggestions-file` option to save improvement suggestions to markdown file
- Improved progress display with step-by-step status (tool counts, activity indicators)
- Copilot SDK usage guide in AGENTS.md

### Fixed

- Fixed Copilot SDK import (`from copilot import CopilotClient` not `copilot_sdk`)
- Fixed Windows glob pattern in release workflow
- Fixed linting issues across codebase (import sorting, exception chaining, etc.)
- Clarified fixture isolation between tasks (each task gets fresh temp workspace)

## [0.0.1] - 2026-02-01

### Added

- **CLI Commands**
  - `waza run` - Run evaluation suites against skills
  - `waza generate` - Auto-generate evals from SKILL.md files
  - `waza init` - Initialize new eval suites interactively
  - `waza report` - Generate reports from results

- **Eval Generation**
  - Pattern-based generation from SKILL.md files
  - LLM-assisted generation with `--assist` flag for better tasks/fixtures
  - Support for multiple models (Claude, GPT-4, etc.)

- **Executors**
  - Mock executor for testing without LLM calls
  - Copilot SDK executor for real integration testing

- **Graders**
  - Code graders with Python assertions
  - Regex graders for pattern matching
  - LLM graders for semantic evaluation

- **Features**
  - Real-time progress display with conversation streaming (`-v`)
  - Transcript logging (`--log`)
  - Project context support (`--context-dir`)
  - LLM-powered improvement suggestions (`--suggestions`)

- **Documentation**
  - Comprehensive README with examples
  - Tutorial guide
  - Grader reference
  - Demo script for walkthroughs

### Fixed

- Grader eval context now includes `str`, `int`, `bool`, etc.
- Transcript normalization for proper tool call detection
- YAML escaping for regex patterns with backslashes
- Progress bar now shows 100% on completion

[Unreleased]: https://github.com/spboyer/waza/compare/v0.4.0-alpha.1...HEAD
[0.4.0-alpha.1]: https://github.com/spboyer/waza/compare/azd-ext-microsoft-azd-waza_0.3.0...v0.4.0-alpha.1
[0.3.0]: https://github.com/spboyer/waza/compare/azd-ext-microsoft-azd-waza_0.2.1...azd-ext-microsoft-azd-waza_0.3.0
[0.2.1]: https://github.com/spboyer/waza/compare/azd-ext-microsoft-azd-waza_0.2.0...azd-ext-microsoft-azd-waza_0.2.1
[0.2.0]: https://github.com/spboyer/waza/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/spboyer/waza/compare/v0.0.2...v0.1.0
[0.0.2]: https://github.com/spboyer/waza/compare/v0.0.1...v0.0.2
[0.0.1]: https://github.com/spboyer/waza/releases/tag/v0.0.1
