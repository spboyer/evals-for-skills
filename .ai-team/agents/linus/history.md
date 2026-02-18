# History — Linus

## Project Context
- **Project:** waza — CLI tool for evaluating Agent Skills
- **Stack:** Go (primary), React 19 + Tailwind CSS v4 (web UI)
- **User:** Shayne Boyer (spboyer)
- **Repo:** spboyer/waza
- **Universe:** The Usual Suspects

## Key Learnings

### Go Architecture
- **Model directive:** Coding in Claude Opus 4.6 (user requirement)
- **Code structure:** Functional options pattern for configuration
- **Interfaces:** AgentEngine, Validator, Grader (extensible design)
- **Testing:** Unit tests in internal packages, integration tests for CLI

### Waza-specific
- Fixture isolation: temp workspace created per task, original fixtures never modified
- TestCase, BenchmarkSpec, EvaluationOutcome models
- ValidatorRegistry pattern for pluggable graders
- CLI flags: -v (verbose), -o (output), --context-dir (fixtures)

### Integration
- Copilot SDK integration (via AgentEngine interface)
- Web UI gets results from CLI JSON output
- Makefile for build/test/lint automation

### Dev Package Patterns (cmd/waza/dev/)
- Package-level var functions as test hooks: `newDevEngine`, `startDevSpinner`, `promptConfirm`
- Test helpers follow `withDevTest*` naming: `withDevTestEngine`, `withDevTestSpinner`, `withDevTestConfirm`
- huh library used for interactive prompts (confirm, select, input) — see cmd_init.go, prompt.go
- TTY detection via `term.IsTerminal` to skip prompts in non-interactive environments
- PR #225 migrated PromptConfirm to huh.NewConfirm (closes #221)

## Learnings

### PR #226 — Generated CI workflow now uses azd extension (#222)
- `initCIWorkflow()` in `cmd/waza/cmd_init.go` holds the embedded YAML string for the generated GitHub Actions workflow
- Replaced `actions/setup-go` + `go install` with `Azure/setup-azd@v2` + azd extension install pattern
- Test in `cmd_init_test.go` (`TestInitCommand_CIWorkflowContent`) validates generated workflow content — must be updated alongside template changes
- azd extension install requires 3 commands: enable alpha extensions, add source registry, install extension
- `TestInitCommand_ReadmeContent` also asserts `waza run` but that's from the README template, not the CI workflow — left unchanged

### PR — Make waza new interactive by default, migrate wizard to huh (#220)
- Removed `--interactive` flag from `waza new` and `waza init`; TTY detection (`term.IsTerminal`) now drives wizard activation
- Migrated `internal/wizard/wizard.go` from `bufio.Scanner` to `huh.NewForm` (Input + Select fields)
- bubbletea (huh's backend) gobbles all input from `strings.Reader` at once — wizard tests must use `io.Pipe` with delayed writes for multi-field forms
- huh gracefully handles EOF/incomplete input by using default values (no error returned) — tests should assert default values, not errors
- `initCommandE` no longer calls `newCommandE` to create skills; instead calls `scaffoldInProject` directly to avoid triggering the wizard
