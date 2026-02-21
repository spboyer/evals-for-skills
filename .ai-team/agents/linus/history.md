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

### Web API Architecture
- API types in `internal/webapi/types.go` are decoupled from internal models (no direct imports)
- `outcomeToDetail()` in `store.go` maps `models.EvaluationOutcome` → API response types
- JSON uses camelCase consistently across the API surface
- TranscriptEvent mapping uses direct field access (not marshal/unmarshal) due to MarshalJSON snake_case mismatch

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


## 📌 Team update (2026-02-20): Model policy overhaul

All code roles now use `claude-opus-4.6`. Docs/Scribe/diversity use `gemini-3-pro-preview`. Heavy code gen uses `gpt-5.2-codex`. Decided by Scott Boyer. See decisions.md for full details.

### #299 — Grader Weighting (PR pending)
- **Date:** 2026-02-20
- **Branch:** `squad/299-grader-weighting`
- **Files changed:** `internal/models/spec.go`, `internal/models/outcome.go`, `internal/orchestration/runner.go`, `internal/reporting/interpreter.go`, `internal/webapi/types.go`, `internal/webapi/store.go`, `internal/models/spec_test.go`, `internal/models/outcome_test.go`
- **What:** Added optional `weight` field to `GraderConfig` (default 1.0 via `EffectiveWeight()`). Added `ComputeWeightedRunScore()` to `RunResult`. Weighted composite score surfaces in `TestStats.AvgWeightedScore`, `OutcomeDigest.WeightedScore`, and the interpretation report. Web API `GraderResult` also carries weight. All existing eval.yaml files work unchanged — weight is optional and defaults to 1.0.
- **Key learning:** `ValidatorInline` (task-level graders) already had a `Weight` field before this change — only `GraderConfig` (spec-level) was missing it. The runner is the correct place to stamp weights onto `GraderResults` since graders themselves don't know their config weight.

### #314 — agentskills.io Spec Compliance Checks (PR #322)
- **Date:** 2026-02-20
- **Branch:** `squad/314-spec-compliance`
- **Files changed:** `cmd/waza/dev/spec.go` (new), `cmd/waza/dev/spec_test.go` (new), `cmd/waza/dev/display.go`, `cmd/waza/dev/display_test.go`, `cmd/waza/dev/score_test.go`, `cmd/waza/dev/loop_test.go`, `cmd/waza/cmd_check.go`
- **What:** Added `SpecScorer` with 8 formal agentskills.io spec checks (frontmatter, allowed-fields, name format, dir-match, description length, compatibility length, license recommendation, version recommendation). Integrated into both `waza dev` (inline in DisplayScore) and `waza check` (separate section, summary table column, readiness gate). 15 new test cases.
- **Key learning:** `makeSkill` test helper needed `FrontmatterRaw` field populated (was nil before) to properly test spec checks. Existing display/loop tests use exact string matching — any new output from `DisplayScore` requires updating all dependent test expected strings.

### #308 — Statistical Confidence Intervals (PR #323)
- **Date:** 2026-02-20
- **Branch:** `squad/308-statistical-ci`
- **Files changed:** `internal/statistics/bootstrap.go` (new), `internal/statistics/bootstrap_test.go` (new), `internal/models/outcome.go`, `internal/orchestration/runner.go`
- **What:** New `internal/statistics/` package with `BootstrapCI` (10k resamples, percentile method), `IsSignificant` (CI doesn't cross zero), and `NormalizedGain` (Hake 1998). Wired bootstrap CI into `computeTestStats` (per-task, when ≥2 runs) and `buildOutcome` (digest-level `StatisticalSummary`). Also populated previously-empty `TestStats` fields: `ScoreVariance`, `CI95Lo`, `CI95Hi`, `Flaky`. 11 test cases covering edge cases, determinism, and CI properties. Fully backward compatible via `omitempty`/pointer types.
- **Key learning:** `TestStats` already had `CI95Lo`/`CI95Hi`/`ScoreVariance`/`Flaky` fields defined but never populated — they were placeholders from initial model design. The `internal/metrics` package already had a normal-approximation `ConfidenceInterval95` function; the new bootstrap approach is more robust for small samples and non-normal distributions. Using `BootstrapCIWithSeed` for deterministic tests is essential — non-seeded bootstrap CIs are non-deterministic and will cause flaky tests.
