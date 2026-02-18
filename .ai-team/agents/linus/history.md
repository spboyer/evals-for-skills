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
