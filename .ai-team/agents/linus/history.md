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
- **Multi-skill results:** `skillRunResult` struct captures per-model `EvaluationOutcome` data
- **Function signature:** `runCommandForSpec()` returns `([]modelResult, error)` to enable outcome capture

### Integration
- Copilot SDK integration (via AgentEngine interface)
- Web UI gets results from CLI JSON output
- Makefile for build/test/lint automation

### Web API Architecture
- API types in `internal/webapi/types.go` are decoupled from internal models (no direct imports)
- `outcomeToDetail()` in `store.go` maps `models.EvaluationOutcome` → API response types
- JSON uses camelCase consistently across the API surface
- TranscriptEvent mapping uses direct field access (not marshal/unmarshal) due to MarshalJSON snake_case mismatch

### Key Files
- `cmd/waza/cmd_run.go` — run command implementation, multi-skill orchestration, skillRunResult
- `internal/models/outcome.go` — EvaluationOutcome data model
- `internal/orchestration/runner.go` — TestRunner, benchmark execution

## Completed Work

### #237 — Expose transcript & session digest in web API (PR #242)
- **Date:** 2026-02-19
- **Branch:** `squad/237-api-transcript`
- **Files changed:** `internal/webapi/types.go`, `internal/webapi/store.go`, `internal/webapi/handlers_test.go`, `web/src/api/client.ts`
- **What:** Added `TranscriptEventResponse`, `SessionDigestResponse` API types; wired them into `TaskResult`; mapped from `RunResult` in `outcomeToDetail()`; added TS interfaces; added test

### #239 — Trajectory Diffing (PR #244)
- **Date:** 2026-02-19
- **Branch:** `squad/239-trajectory-diffing`
- **Files changed:** `web/src/components/TrajectoryDiff.tsx` (new), `web/src/components/TaskTrajectoryCompare.tsx` (new), `web/src/components/CompareView.tsx` (modified)
- **What:** Added trajectory diffing to CompareView — aligns ToolExecutionStart events by tool name+index, renders matched/changed/insertion/deletion with color coding and expandable JSON diffs. No backend changes needed.

### #271 — Per-skill output files for --output (PR TBD)
- **Date:** 2026-02-20
- **Branch:** `squad/271-per-skill-output`
- **Files changed:** `cmd/waza/cmd_run.go`, `cmd/waza/cmd_run_test.go`
- **What:** Implemented per-skill output files for multi-skill runs when using `--output`. Created `buildOutputPath()` to handle naming patterns: `{base}_{skill}_{model}.json` for multi-skill+multi-model, `{base}_{skill}.json` for multi-skill only, `{base}_{model}.json` for multi-model only. Renamed `sanitizeModelName` → `sanitizePathSegment` for clarity. Suppressed output during multi-skill loop, wrote per-skill files after. Added tests for all path patterns.
- **Key learning:** Multi-skill output pattern follows the existing multi-model template (line 323). Suppress during loop by clearing `outputPath`, restore after. Single-skill behavior unchanged.

### #272 — Fix skillRunResult to capture EvaluationOutcome (PR TBD)
- **Date:** 2026-02-20
- **Branch:** `squad/272-fix-skill-run-result`
- **Files changed:** `cmd/waza/cmd_run.go`, `cmd/waza/cmd_run_test.go`
- **What:** Enhanced `skillRunResult` to include `outcomes []modelResult` field. Modified `runCommandForSpec` to return `([]modelResult, error)` so multi-skill runs can capture per-model outcome data. Updated `printSkillRunSummary` to show pass rates and aggregate scores from captured outcomes. Added test `TestSkillRunResult_CapturesOutcomes`.

### #273 — Combined summary.json for multi-skill runs (PR TBD)
- **Date:** 2026-02-20
- **Branch:** `squad/273-combined-summary`
- **Files changed:** `internal/models/summary.go` (new), `cmd/waza/cmd_run.go`, `cmd/waza/cmd_run_test.go`
- **What:** Implemented combined summary.json output for multi-skill runs. Created `MultiSkillSummary` struct with skill-level and overall aggregated metrics (pass rate, aggregate score). Added `--no-summary` flag to skip summary generation. When `--output` is set, writes `{base}_summary{ext}` file after multi-skill run completes. Added `buildMultiSkillSummary()` helper that aggregates pass rates and scores across all models for each skill. Added comprehensive tests covering single-skill, multi-skill, and multi-model scenarios.

### #274 — Cross-product multi-skill × multi-model output naming (PR TBD)
- **Date:** 2026-02-20
- **Branch:** `squad/274-cross-product-naming`
- **Files changed:** `cmd/waza/cmd_run.go`, `cmd/waza/cmd_run_test.go`
- **What:** Verified and enhanced cross-product naming for multi-skill × multi-model runs. Fixed `runCommandForSpec` (line 356) to use `buildOutputPath` instead of hardcoded format for consistency. Added comprehensive integration test `TestRunCommand_CrossProductMultiSkillMultiModel` that simulates 2 skills × 2 models = 4 output files. Added edge case tests for partial failures and special character sanitization. The `buildOutputPath` function from #271 already handled the cross-product case correctly with pattern `{base}_{skill}_{model}.json`.
- **Key learning:** Multi-skill loop suppresses `outputPath` to prevent double-writing, then restores it to write per-skill files. Single-model single-skill writes to exact path. Multi-model uses `buildOutputPath` with `multiSkill=false` for single-skill context. Cross-product works because multi-skill loop passes both flags correctly.

### #270 — Add --output-dir flag for structured directory output (PR TBD)
- **Date:** 2026-02-20
- **Branch:** `squad/270-output-dir-flag`
- **Files changed:** `cmd/waza/cmd_run.go`, `cmd/waza/cmd_run_test.go`
- **What:** Added `--output-dir` flag for structured directory output as an alternative to `--output`. Multi-skill runs create subdirectories per skill: `{output-dir}/{skill-name}/{model-name}.json`. Single-skill runs write directly: `{output-dir}/{model-name}.json`. Implemented mutual exclusion validation with `--output` flag (validated early in `runCommandE`). Added `writeOutputDir()` helper that creates directories, sanitizes paths using existing `sanitizePathSegment()`, and writes per-model JSON files. Added comprehensive tests for mutual exclusion, single-skill mode, multi-skill subdirectories, and path sanitization.
- **Key learning:** The `--output-dir` flag provides a cleaner alternative to the `--output` filename-based approach for multi-skill/multi-model runs. Directory structure is more intuitive and avoids filename clutter. Single-skill detection (based on `len(results) > 1`) determines whether to create skill subdirectories or write directly to the output directory. Reused existing `sanitizePathSegment()` and `saveOutcome()` helpers for consistency.


## 📌 Team update (2026-02-20): Model policy overhaul

All code roles now use `claude-opus-4.6`. Docs/Scribe/diversity use `gemini-3-pro-preview`. Heavy code gen uses `gpt-5.2-codex`. Decided by Scott Boyer. See decisions.md for full details.
