# Project Context

- **Owner:** Shayne Boyer (spboyer@live.com)
- **Project:** Waza — Go CLI for evaluating AI agent skills (scaffolding, compliance scoring, cross-model testing)
- **Project:** Waza — Go CLI for evaluating AI agent skills
- **Stack:** Go, Cobra CLI, Copilot SDK, YAML specs
- **Created:** 2026-02-09

## Learnings

<!-- Append new learnings below. Each entry is something lasting about the project. -->
<!-- Append new learnings below. -->

- **2026-02-09:** Wrote `internal/graders/behavior_grader_test.go` with 30 table-driven tests for the behavior grader (#98). Tests cover: max_tool_calls (pass/fail/skip-at-zero), max_tokens (pass/fail), required_tools (pass/fail/empty), forbidden_tools (pass/fail/empty), max_duration_ms (pass/fail), combined rules (all pass, partial fail, multiple fail, independent checking), edge cases (nil session, zero values, nil tools_used, result details, duration recording). Tests follow the exact patterns from regex_grader_test.go and file_grader_test.go. File is gofmt-clean. Compilation depends on Linus's `BehaviorGrader` implementation landing.
- `graders.Context` already has `Session *models.SessionDigest` field (lines 61-63 of grader.go) — no modification needed.
- The test uses `require` from testify throughout, matching project convention. All tests use `context.Background()` and `graders.Context{}` struct literals.
- errcheck with check-blank:true is enforced by golangci-lint — always capture and assert on returned errors.
