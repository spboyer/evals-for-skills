# Project Context

- **Owner:** Shayne Boyer (spboyer@live.com)
- **Project:** Waza — Go CLI for evaluating AI agent skills (scaffolding, compliance scoring, cross-model testing)
- **Project:** Waza — Go CLI for evaluating AI agent skills
- **Stack:** Go, Cobra CLI, Copilot SDK, YAML specs
- **Created:** 2026-02-09

## Learnings

<!-- Append new learnings below. Each entry is something lasting about the project. -->
📌 Team update (2026-02-11): Run command tests must reset package-level flag vars (contextDir, outputPath, verbose) at top of each test body to prevent state leakage. — decided by Linus
<!-- Append new learnings below. -->

### 2026-02-11: PR #111 Review — tokens compare command
- **Author:** Charles Lowell (chlowell), branch `tokens-compare`
- **Verdict:** Approved. Clean implementation of `waza tokens compare` (E4, closes #51).
- **Architecture:** New `internal/git` package under `cmd/waza/tokens/internal/git/` — well-bounded, not importable outside tokens tree.
- **Quality:** Comprehensive tests with real git repos in temp dirs, table-driven subtests. Reuses existing `EstimatingCounter`, `NormalizePath`, `nowISO`.
- **Nit:** `RefExists()` is dead code (defined but never called). Non-blocking.
- **CI:** Green — both build/test and lint passed.

### 2026-02-11: PR #112 Review — --show-unchanged applied to JSON output
- **Author:** Charles Lowell (chlowell), branch `unchanged-json`
- **Verdict:** Approved. Tight follow-up to PR #111 (+12/-13, single file).
- **Change:** Lifts `--show-unchanged` filtering from `compareTable` up to `runCompare`, so it applies to both table and JSON output. Summary computed before filtering so totals remain correct.
- **Quality:** Uses `slices.DeleteFunc` (Go 1.21+ stdlib) — replaces manual filter loop. `compareTable` simplified by removing `showUnchanged` parameter.
- **Tests:** Existing tests cover both paths. No new tests needed — the filtering is now shared code exercised by table-output tests.
- **CI:** Green — both build/test and lint passed.

### 2026-02-11: PR #113 Review — azd extension release pipeline
- **Author:** Wallace Breza (wbreza), branch `feat/azd-ext-release-pipeline`
- **Verdict:** Changes requested. Two blocking issues, three suggestions.
- **Blocking:** (1) Version downgrade 0.1.0→0.0.2 — semver should only move forward, needs justification or fix. (2) Registry checksum/tag mismatch — URLs reference `_0.1.0` tag but version is being set to 0.0.2.
- **Suggestions:** Move validation scripts out of repo root (`scripts/`); clarify `.github/skills/` vs `skills/` convention for workflow-automation skills; add trailing newline to `version.txt`.
- **Good:** Pipeline structure (workflow_dispatch → build → pack → release → publish → auto-merge registry PR), permissions minimized, all 6 platform targets, `GH_TOKEN` from `secrets.GITHUB_TOKEN`, SKILL.md well-structured with user prompts at decision points, both bash and PowerShell validation scripts.
- **Alignment:** Directly advances E7 (AZD Extension). Completes release automation story started in PR #103.
- **CI:** No checks reported on the branch (new workflow only, no Go code changes).

### 2026-02-11: PR #114 Review — tokens suggest command
- **Author:** Charles Lowell (chlowell), branch `tokens-suggest`
- **Verdict:** Changes requested. Three lint issues blocking CI.
- **Blocking:** (1) `errcheck` — `engine.Shutdown` return value unchecked in `suggest.go`. (2) `errcheck` — `filepath.Rel` return value unchecked in copilot goroutine. (3) `misspell` — `analyses`/`Analyses` flagged as misspelling of `analyzes`/`Analyzes` (6 occurrences across suggest.go and suggest_test.go).
- **Architecture:** Two-mode design (heuristic + copilot) with `newChatEngine` function variable for test injection. Semaphore-bounded concurrency (`maxCopilotWorkers=4`). Prompt embedded via `//go:embed`. Refactored `countFile` → `countTokens` as pure function shared across count/check/suggest. Moved `countLines` from `compare.go` to `helpers.go`.
- **Quality:** 17 test functions, comprehensive fixture set under `testdata/suggest/`, mock engine integration, JSON/text output, edge cases. Heuristic checks align with sensei reference (emojis, code blocks, tables, duplicates, horizontal rules, limit violations).
- **Size:** +1137/-34 — substantial but well-scoped.
- **CI:** Build/test green. Lint failing (3 categories above).
- **Lesson:** golangci-lint's misspell checker treats "analyses" (valid English noun) as a misspelling of "analyzes". Watch for this in future PRs — either rename variables or suppress with nolint directive.

### 2026-02-11: PR #115 Review — azd extension metadata capability
- **Author:** Wallace Breza (wbreza), branch `feat/metadata-capability`
- **Verdict:** Changes requested. Two blocking CI failures, three non-blocking suggestions.
- **Blocking:** (1) `gofmt` — both `cmd_metadata.go` and `cmd_metadata_test.go` have formatting issues. (2) `go 1.25` version bump in `go.mod` breaks golangci-lint v1.64.8 (built with Go 1.24, refuses Go 1.25 targets). Either pin to a Go 1.24-compatible azd module version or upgrade golangci-lint in CI.
- **Architecture:** Hidden `metadata` Cobra command calls `azdext.GenerateExtensionMetadata()` — pure introspection, no side effects, writes JSON to stdout. Uses canonical azd types, no custom converters. Wired via `cmd.AddCommand(newMetadataCommand(cmd))` in root.go. `extension.yaml` adds `metadata` to capabilities list.
- **Quality:** 4 tests covering JSON validity/schema, expected commands, flag population, and hidden status. Clean separation — single 32-line file for the command.
- **Concern:** The `azd` module pulls ~60 transitive dependencies (OpenTelemetry, gRPC, protobuf, Azure SDK). Significant weight increase for a previously lightweight CLI. Acceptable for canonical integration, but should migrate to standalone `azdext` module if one is published.
- **Alignment:** Directly advances E7 (AZD Extension). Completes metadata discovery story alongside PR #113 (release pipeline).
- **CI:** Both build/test and lint failing (gofmt + golangci-lint version mismatch).
- **Lesson:** When adding dependencies that require a Go version bump, check that CI toolchain (especially golangci-lint) supports the new version. Coordinate go.mod and CI workflow changes in the same PR.

📌 Team update (2026-02-12): azd-publish skill location convention — repo-level skills go under `.github/skills/`, project eval skills go under `skills/`. — decided by Wallace Breza
📌 Team update (2026-02-12): azd extension uses tag pattern `azd-ext-microsoft-azd-waza_VERSION`, not `vVERSION`. — decided by Linus
📌 Team update (2026-02-12): PR #115 review feedback addressed — Linus rebased, resolved conflicts, added doc comments per review. — decided by Linus
### 2026-02-11: PR #117 Review — waza dev command (E2: Sensei Engine)
- **Author:** Charles Lowell (chlowell), branch `waza-dev`
- **Verdict:** APPROVED. Clean, well-architected implementation of the Sensei development loop. Closes #32, #33, #35.
- **Epic:** E2 (Sensei Engine) — Iterative skill improvement with heuristic scoring
- **Architecture:** Four-package structure:
  - `cmd/waza/dev/` — CLI (root.go, loop.go, score.go, display.go, prompt.go) with clear separation: orchestration, heuristics, formatting, user input
  - `internal/skill/` — New SKILL.md parser with `TextMarshaler`/`TextUnmarshaler` for YAML round-trip, preserves unknown fields
  - `internal/tokens/` — Extracted token estimation logic (surgical refactor: import path changes only, no logic modifications)
  - Tests: 61 functions across 6 test files (display, score, loop, prompt, skill, tokens)
- **Heuristic Scoring:** Correctly implements Sensei reference rules: Low → Medium (desc 150+ chars + triggers) → Medium-High (+ anti-triggers) → High (+ routing clarity). Pattern detection validates against real fixtures (code-explainer=Low, waza=High).
- **Ralph Loop:** Iterates through description expansion → triggers → anti-triggers → routing clarity. Correctly skips inapplicable steps; declining suggestion advances to next step (not terminating). Enforces soft (500) and hard (5000) token limits.
- **Code Quality:** Idiomatic Go — interface-based scorer for test injection, functional error wrapping (%w), clean naming, well-commented. Prompt state handling (shared bufio.Scanner) prevents stdin clobbering across multiple user prompts. Display formatting includes box-drawing characters, emoji width awareness, rune-aware truncation.
- **Test Coverage:** Table-driven pattern matching tests, edge cases (nil skill, over-length descriptions, token budgets), real fixture loading, comprehensive testdata (high, valid, minimal, no-frontmatter).
- **Alignment:** ✅ E2 epic, ✅ issues #32-35, ✅ sensei reference patterns, ✅ Ralph loop discipline
- **CI:** Both build/test and lint passing (no errcheck, gofmt, or misspell violations)
- **Recommendation:** Merge immediately. Implementation demonstrates deep understanding of Sensei architecture and Go conventions. No rework needed. Post-merge: consider future --strict flag for hard-limit exit code (E4 scope).

### 2026-02-11: PR #117 Deep Review (second opinion, opus-4.6)
- **Verdict:** Confirmed approval. First review was accurate on architecture, scoring correctness, and test quality.
- **New findings (all non-blocking):**
  1. **TriggerCount vs HasTriggers mismatch** — `HasTriggers` matches 4 patterns but `TriggerCount` only counts after "USE FOR:". Display can show "Triggers: 0" while scoring Medium. UX confusion, not a scoring bug.
  2. **parseFrontmatter closing delimiter fragile** — `strings.Index(rest, "\n---")` could split prematurely on multiline YAML scalars containing `---`. Safe during round-trip (yaml.Marshal escapes), edge case for hand-crafted files.
  3. **writeSkillFile not atomic** — `os.WriteFile` without temp+rename. Ctrl+C during write could truncate. Low risk.
  4. **No context.Context** — Will need retrofit when Copilot-based suggestions (#36) land.
  5. **suggestTriggers semantic duplicates** — Name + heading overlap produces redundant phrases. Cosmetic.
  6. **boxLine emoji width** — Rune count ≠ terminal column width for ✅/❌. Acknowledged in code comments.
- **Suggestions:** count all trigger patterns (not just "USE FOR:"), atomic writes, context.Context plumbing, --dry-run flag, tests for non-existent path and --target low.
- **Copilot reviewer alignment:** 3 of 4 inline comments were valid (trigger count mismatch, parser fragility, comment/code mismatch). The 4th was a typo fix.
- **Lesson:** Deep reviews catch UX inconsistencies and future-proofing gaps that fast-model reviews miss, but the fast model correctly identified all structural and correctness aspects.

📌 Team update (2026-02-15): Review @copilot PRs with claude-opus-4.6 before merging — quality gate for doc PRs. — decided by Shayne Boyer
📌 Team update (2026-02-15): Auto-assign unblocked work to squad/@copilot. Don't ask, just assign and go. — decided by Shayne Boyer
📌 Team update (2026-02-15): After feature PRs merge (CLI, graders, YAML format, examples), route doc updates to Saul. Issue #148 tracks this. — decided by Shayne Boyer
📌 Team update (2026-02-15): All developers use claude-opus-4.6. For code review, if developer isn't using Opus, reviewer uses it. — decided by Shayne Boyer
📌 Team update (2026-02-15): Don't take assigned work. Only pick up unassigned issues. — decided by Shayne Boyer
📌 Team update (2026-02-15): Multi-model execution is sequential (not parallel). Test failures non-fatal so all models complete. — decided by Linus
📌 Team update (2026-02-15): Microsoft/skills repo moving to plugin bundle structure. CI must support both flat and nested layouts. — decided by Shayne Boyer

### 2026-02-15: PR #152 Review — --model flag for multi-model evaluation (#39)
- **Author:** Linus (Backend Dev), branch `squad/39-multi-model-support`
- **Verdict:** APPROVE WITH NITS. Clean implementation, two non-blocking issues.
- **Epic:** E3 (Evaluation Framework) — Closes #39 [E3-01] Support multiple model execution.
- **Architecture:** `runSingleModel()` extraction is well-structured — each model gets its own engine instance, runner, and context. Spec mutation (`spec.Config.ModelID = modelID`) in the loop is safe because `runSingleModel` reads it only at creation time. Error handling distinguishes `TestFailureError` (continue in multi-model) from infrastructure errors (abort immediately). `sanitizeModelName()` handles `/`, `\`, `:`, space. Per-model JSON output uses `base_model.ext` naming.
- **Nits found:** (1) Comparison table `%-10.1f%%` format puts a gap between value and `%` sign (cosmetic). (2) `runSingleModel` creates engine but never calls `engine.Shutdown()` — pre-existing issue (also absent on main), not introduced by this PR.
- **Test coverage:** 9 tests + 3 subtests covering flag parsing, override, multi-model, backward compat, comparison table output, edge cases. `resetRunGlobals()` correctly includes `modelOverrides = nil`.
- **CI:** Build clean, go vet clean, all tests pass (including all pre-existing tests unchanged).
- **Lesson:** When extracting a loop body into a helper, always verify resource lifecycle (init/defer-shutdown) is preserved per iteration. In this case, engine.Shutdown() was already missing upstream, so the extraction didn't regress.

### 2026-02-15: E3 Evaluation Framework Backlog Triage
- **Task:** Assess four unassigned E3 issues and recommend prioritization
- **Issues triaged:** #44 (suggestions), #106 (tool_call rubrics), #107 (task rubrics), #138 (multi-model recommendations)
- **Key findings:**
  1. **#44 ready NOW** — Charles already extracted suggestion engine in PR #117. Main work: consolidate into `internal/suggestions/` package shared by `waza dev` and `waza run --suggestions`. No blockers. Assign to Linus.
  2. **#106 & #107 blocked by #104** — Both port Azure ML `.prompty` templates to YAML rubrics. Identical structure, work in parallel after #104 merges. Assign to Livingston (docs-integration role fits).
  3. **#138 blocked by #39 + #104** — Capstone E3 feature (multi-model recommendations). #39 merged (PR #152). Waiting on #104 prompt grader stability. Needs design clarity: What are optimization dimensions? (cost/quality/latency/consistency?)
- **Critical blocker:** #104 (prompt grader) unblocks 3 of 4 issues downstream. Recommend prioritizing in parallel track.
- **Architecture decisions captured:**
  - Suggestion engine must be a shared package, not duplicated in `waza dev` and `waza run`.
  - Azure ML porting establishes reusable rubric YAML pattern — capture for future evaluators.
  - Recommendation engine requires explicit rubric design (dimensions) before prompt engineering.
- **Deliverable:** Triage analysis written to `.ai-team/decisions/inbox/rusty-eval-backlog-triage.md`
📌 Team update (2026-02-15): Don't take assigned work — only pick up unassigned issues — decided by Shayne Boyer
📌 Team update (2026-02-15): Engine shutdown fix (#153) complete — merged 21 tests from Basher covering all exit paths. Critical path blocker #104 (Prompt Grader) unblocks 50% of E3 backlog. — Rusty (Lead)

### 2026-02-16: Go Release Pipeline Architecture (Design Spec)
- **Task:** Design the Go binary release pipeline to replace Python release workflow
- **Actions taken:**
  1. Deprecated `.github/workflows/release.yaml` → renamed to `release-python-legacy.yaml` with `if: false` and deprecation header
  2. Designed comprehensive Go release pipeline architecture
  3. Documented spec at `.ai-team/decisions/inbox/rusty-go-release-pipeline.md`
- **Key architecture decisions:**
  - Trigger: `v*` Git tags + workflow_dispatch with manual version
  - Matrix: 6 platforms (linux/darwin/windows × amd64/arm64)
  - Binary naming: `waza-{os}-{arch}[.exe]` (standalone CLI, not "microsoft-azd-waza" extension)
  - Version injection: `-ldflags "-X main.version=$VERSION"`
  - Artifacts: 6 binaries + SHA256 checksums + install.sh + release notes
  - Release notes: Auto-generated from merged PRs since last tag
  - Install script: Platform detection, download, verify checksum, add to PATH
- **Why this design:**
  - Standard Go release pattern (matrix build)
  - User-friendly (single `curl | bash` install)
  - Orthogonal to azd extension workflow (separate versioning concern)
  - Cross-platform tested before release
  - Verifiable (SHA256 checksums)
- **Outcome:** Spec ready for Linus to implement. Acceptance criteria and next steps documented. Orthogonality to azd-ext-release.yml maintained — they share `version.txt` but operate independently.
- **Implementation timing:** After design review, Linus will implement `go-release.yml` and `install.sh`.
