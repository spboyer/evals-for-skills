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
- **2026-02-12:** Wrote 9 tests (+ 3 subtests) in `cmd/waza/cmd_run_test.go` for multi-model support (#39). Tests cover: flag parsing (single, multiple, edge cases), single-model override of YAML spec, multi-model execution with per-model output files, backward compatibility (no flag preserves YAML model), model name in JSON output, identity override (--model matching spec), and comparison table stdout capture. All tests pass against Linus's existing implementation.
- Multi-model runs save output as `<base>_<model>.json` per model, NOT to the original `--output` path. Tests must check per-model file paths.
- `resetRunGlobals()` must include `modelOverrides = nil` — Cobra `StringArrayVar` persists across test cases in the same process.
- For capturing stdout from `fmt.Printf` (not Cobra output), use `os.Pipe()` redirect since `cmd.SetOut(io.Discard)` only affects Cobra's own writer.
📌 Team update (2026-02-15): All developers must use claude-opus-4.6 for code. For code review, if developer isn't using Opus, reviewer uses it. — decided by Shayne Boyer
📌 Team update (2026-02-15): Don't take assigned work. Only pick up unassigned issues. — decided by Shayne Boyer
📌 Team update (2026-02-15): Multi-model execution is sequential (not parallel). Test failures non-fatal so all models complete. — decided by Linus
📌 Team update (2026-02-15): Microsoft/skills repo moving to plugin bundle structure. CI must support both flat and nested layouts. — decided by Shayne Boyer
