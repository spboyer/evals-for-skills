# Team Decisions

> Shared brain — all agents read this before starting work.

<!-- Scribe merges decisions from .ai-team/decisions/inbox/ into this file. -->
### 2026-02-09: Proper git workflow required
**By:** Shayne Boyer (via Copilot)
**What:** All issues must follow: feature branch → commit → push → PR with `Closes #N` → @copilot review → address feedback → merge. No direct commits to main.
**Why:** User directive

### 2026-02-09: Sensei reference repo for E2 implementation
**By:** Squad (Coordinator)
**What:** The sensei engine (E2: #32-38) must adopt functionality from https://github.com/spboyer/sensei. Key patterns to port to Go:

**Scoring Algorithm:** Low (description < 150 chars OR no triggers) → Medium (description >= 150 chars AND trigger keywords) → Medium-High (USE FOR + DO NOT USE FOR) → High (Medium-High + routing clarity). Rule-based checks: name validation, description length, trigger/anti-trigger detection, routing clarity, MCP integration.

**Ralph Loop:** READ → SCORE → CHECK → SCAFFOLD → IMPROVE FRONTMATTER → IMPROVE TESTS → VERIFY → CHECK TOKENS → SUMMARY → PROMPT → REPEAT (max 5 iterations, target: Medium-High).

**Token Management:** count/check/suggest/compare with .token-limits.json config. SKILL.md: 500 soft / 5000 hard. References: 2000 each. --strict exits 1 on limit breach.
**Why:** User directive — spboyer/sensei is the reference implementation

### 2026-02-09: Monitor human engineer comments
**By:** Shayne Boyer (via Copilot)
**What:** Periodically check for comments from Charles (@chlowell) and Richard (@richardpark-msft) on open issues. Follow up on responses.
**Why:** User directive

### 2026-02-09: Run command test patterns
**By:** Linus (Backend Dev)
**Related:** #24
**What:** Each test that calls `cmd.Execute()` must reset package-level flag vars (`contextDir`, `outputPath`, `verbose`) at the top of the test body to prevent state leakage. If new flags are added to `newRunCommand()`, the reset block must be updated accordingly.
**Why:** Package-level Cobra flag bindings persist across test cases in the same process.

### 2026-02-11: PR #103 — azd extension packaging approved
**By:** Rusty (Lead)
**Related:** E7 (#62), PR #103
**What:** Wallace Breza's azd extension PR approved and merged. Adds `extension.yaml`, build scripts, `registry.json`, `version.txt`. Purely additive — no changes to existing Go code or CI. Follow-up: update registry URLs to `spboyer/waza` for production; add trailing newlines (POSIX nit); consider CI for automated extension builds.
**Why:** E7 (AZD Extension) P2 roadmap item — establishes foundation without disruption.

### 2026-02-11: Grader nil-session error handling contract
**By:** Rusty (Lead)
**Related:** PR #110
**What:** All graders depending on `SessionDigest` must return a zero-score `GraderResults` with `nil` error when session is nil — not `(nil, error)`. Matches `behavior_grader.go` pattern (lines 56-62). Returning an error would abort the entire run instead of gracefully recording a zero score.
**Why:** Runner distinguishes grader errors (may abort) from zero-score results (contribute to scoring). Affects `action_sequence_grader.go` and future session-dependent graders.

### 2026-02-11: PR #111 approved — tokens compare command
**By:** Rusty (Lead)
**What:** Approved Charles Lowell's `waza tokens compare` implementation. New `internal/git` package lives under `cmd/waza/tokens/internal/git/` — correctly scoped as a tokens-internal dependency. Command follows established Cobra factory patterns (`newCompareCmd()`). Closes #51 (E4: Token Management).
**Why:** Clean architecture, comprehensive tests, CI green. One non-blocking nit: `RefExists()` is dead code. No security concerns.

### 2026-02-12: azd-publish skill location convention
**By:** Wallace Breza (via Copilot)
**What:** The azd-publish skill is a repo-level skill (lives at `.github/skills/azd-publish/`), NOT part of the project eval skills under `skills/`. Repo-level automation skills go under `.github/skills/`, project eval skills go under `skills/`.
**Why:** User request — captured for team memory

### 2026-02-12: azd extension uses non-standard tag pattern
**By:** Linus (Backend Dev)
**Related:** PR #113, E7
**What:** The azd extension release pipeline uses tags of the form `azd-ext-microsoft-azd-waza_VERSION` (e.g., `azd-ext-microsoft-azd-waza_0.2.0`), not `vVERSION`. Any tooling or documentation that references version tags for the azd extension must use this pattern. The SKILL.md comparison link examples have been updated accordingly.
**Why:** The `azd-publish` skill's Step 5 instructions referenced `vX.Y.Z` tags which don't match the actual tag convention, leading to broken comparison links in changelogs.

### 2026-02-12: PR #115 review feedback addressed
**By:** Linus (Backend Dev)
**Related:** PR #115, E7 (#62)
**What:** Rebased `feat/metadata-capability` onto latest `main` (797f72c), resolved `.golangci.yml` conflict (kept v2 format with `version: "2"` header), added doc comments on `metadataSchemaVersion` and `extensionID` constants per Rusty's review. All 4 metadata tests pass. Force-pushed to `wbreza/feat/metadata-capability`.
**Why:** PR had merge conflicts after main advanced; review requested clarifying comments on constants.

### 2026-02-11: Sensei dev command heuristics & test discipline
**By:** Rusty (Lead)
**Related:** E2 (#32-35), PR #117
**What:** `waza dev` heuristic scoring rules (from `spboyer/sensei` reference):
- **Low** — description < 150 chars OR no trigger phrases
- **Medium** — description ≥ 150 chars AND trigger keywords (USE FOR, USE THIS SKILL, TRIGGERS, etc.)
- **Medium-High** — Medium + anti-trigger phrases (DO NOT USE FOR, NOT FOR, Don't use, Instead use)
- **High** — Medium-High + routing clarity (INVOKES, FOR SINGLE OPERATIONS, **WORKFLOW SKILL**, **UTILITY SKILL**, **ANALYSIS SKILL**)

Test discipline: use table-driven subtests for pattern detection, validate against real fixture loading (code-explainer=Low, waza=High), test exact terminal output (box-drawing, emoji width awareness, rune counting), mock scorer interface for loop testing.
**Why:** Clear, testable compliance framework enables future skill compliance automation across the team's codebase. Reference implementation pattern reduces drift across similar tools.

### 2026-02-17: User directive — Code writing and PR review models (consolidated)
**By:** Shayne Boyer (via Copilot)
**What:** 
1. **Code generation:** All code-writing agents (Linus, Basher, any Go implementation work) MUST use claude-opus-4.6 model. No exceptions.
2. **PR review:** All PRs reviewed with BOTH Opus 4.6 AND Codex 5.3 (dual-model review for analytical diversity).
3. **Previous directives:** Review @copilot PRs (2026-02-13, pre-consolidated) now subsumed into comprehensive policy above.
**Why:** User request — premium quality assurance. Code generation (Opus 4.6 only), code review (dual-model: Opus 4.6 + Codex 5.3 for analytical diversity).

### 2026-02-14: User directive — Auto-assign unblocked work
**By:** Shayne Boyer (via Copilot)
**What:** Don't ask before assigning unblocked work to the squad or @copilot — just assign it and go.
**Why:** User request — reduces back-and-forth, keeps the pipeline moving

### 2026-02-14: User directive — Route doc updates to Saul after feature merges
**By:** Shayne Boyer (via Copilot)
**What:** After any feature PR merges that changes CLI commands, graders, eval YAML format, or examples, route a doc update task to Saul. Saul owns DEMO-GUIDE.md, GRADERS.md, TUTORIAL.md, examples/ READMEs, and the main README. Standing issue #148 tracks this.
**Why:** User wants documentation kept current automatically as features ship. Saul is the designated docs team member.

### 2026-02-14: User directive — Skills repo plugin bundle structure
**By:** Shayne Boyer (via Copilot)
**What:** The microsoft/skills repo is being reorganized from a flat `.github/skills/` layout (133 items) into plugin bundles (`.github/plugins/<bundle>/skills/<name>/`). Waza CI compatibility (#60) and any future skills integration must support both the current flat layout and the new nested plugin bundle layout. Key bundles: azure-skills (18 orchestration), azure-sdk-python (41), azure-sdk-dotnet (29), azure-sdk-typescript (24), azure-sdk-java (26), azure-sdk-rust (7), azure-core (6).
**Why:** User shared the distribution strategy gist (https://gist.github.com/spboyer/011190893f33d82d967180cdc5a2432d) — this is the planned future state and all CI work should be forward-compatible

### 2026-02-15: User directive — Don't take assigned work
**By:** Shayne Boyer (via Copilot)
**What:** Don't take anyone's work if it is assigned. Only pick up unassigned issues.
**Why:** User request

### 2026-02-15: Multi-model execution architecture
**By:** Linus (Backend Dev)
**Related:** #39 (E3), PR #152
**What:** When multiple `--model` flags are given, models are evaluated sequentially (not concurrently). Each model gets its own engine instance created fresh inside the loop. The `runSingleModel()` function encapsulates the full benchmark lifecycle for one model — config, engine, runner, execution, summary.
**Why:** Sequential execution is simpler, avoids resource contention between engines, and produces cleaner output (each model's progress prints in order). Parallel model execution can be added later as a separate flag if needed. The `modelResult` type and `printModelComparison()` function are ready for it — they operate on a collected slice regardless of execution order.

### 2026-02-15: Test failures in multi-model runs are non-fatal
**By:** Linus (Backend Dev)
**Related:** #39 (E3), PR #152
**What:** When running multiple models, a `TestFailureError` from one model doesn't abort the remaining models. The error is recorded and the last one is returned after all models complete. Infrastructure errors (load failure, unknown engine) still abort immediately.
**Why:** The whole point of comparison runs is to see how different models perform. Aborting on the first failure defeats the purpose. The user still gets a non-zero exit code if any model had failures.

### 2026-02-15: PR #152 review verdict
**By:** Rusty (Lead)
**What:** APPROVE WITH NITS. Two non-blocking issues:
1. **Comparison table formatting** — `%-10.1f%%` format string in `printModelComparison` produces `100.0     %` instead of `100.0%`. Fix: use `fmt.Sprintf("%.1f%%", passRate)` and print with `%-10s`.
2. **Engine shutdown** — `runSingleModel` creates an engine per model but never calls `engine.Shutdown()`. Pre-existing on main (not a regression), but should be addressed as a follow-up since multi-model runs now create N engines per invocation.

All acceptance criteria met: --model flag, multi-model loop, comparison table, per-model JSON output, backward compatibility preserved.
**Why:** Implementation is architecturally sound — clean extraction, correct error semantics (TestFailureError continues, infra errors abort), proper state isolation per model iteration. Tests are comprehensive (9+3 covering all paths). Build and all tests pass. The two nits are cosmetic/pre-existing and don't block merge.

### 2026-02-15: Engine shutdown must use context.Background() for defer cleanup
**By:** Linus
**What:** `defer engine.Shutdown(context.Background())` placed after engine creation in `runSingleModel()`. Uses `context.Background()` instead of the benchmark's `ctx` since shutdown is independent cleanup.
**Why:** Shutdown should not be cancelled if the benchmark context is cancelled — engines must always release resources. This also prevents a subtle ordering issue where `ctx` is declared later in the function.

### 2026-02-15: Engine Shutdown test strategy
**By:** Basher (Tester)
**Related:** #153 (engine.Shutdown() leak in runSingleModel)
**What:** Created two test files covering engine.Shutdown() lifecycle:
- `internal/execution/engine_shutdown_test.go` — unit tests for Shutdown contract on MockEngine, CopilotEngine, and a reusable SpyEngine test double
- `cmd/waza/cmd_run_shutdown_test.go` — integration tests verifying Shutdown runs in every runSingleModel exit path

**Key design decisions:**
1. **SpyEngine is exported** — so `cmd/waza` tests can import `execution.SpyEngine` if Linus adds an engine factory or injection hook. Currently the engine is hardcoded in a switch statement, so cmd-level tests use indirect verification.
2. **CopilotEngine workspace tests set internal state directly** — rather than requiring the full Copilot SDK, the test locks the mutex and sets `engine.workspace` to a temp dir, then verifies Shutdown clears it. This is a pragmatic tradeoff.
3. **Multi-model Shutdown test** — verifies each model iteration creates and shuts down its own engine independently. This is critical because the loop in `runCommandE` creates a new engine per model.

**Why:** Without these tests, Shutdown leaks are invisible — they don't cause test failures, they cause resource leaks in production (temp dirs, copilot client connections). The SpyEngine pattern makes future Shutdown contract violations immediately detectable.

### 2026-02-15: E3 Evaluation Framework Backlog Triage
**By:** Rusty (Lead)
**Related:** E3 (Epic), Issues #44, #106, #107, #138
**What:** Prioritized four unassigned E3 evaluation framework issues:
1. **#44 (P1) — LLM-powered improvement suggestions** — **READY NOW, assign to Linus**
   - No blockers; internal feature building on Charles's PR #117
   - Effort: 1-2 days (refactor + tests)
   - Architecture: Extract `internal/suggestions/` package, consolidate with `waza dev` logic
2. **#106, #107 (P1, tool_call & task rubrics)** — **Blocked on #104 (Prompt Grader)**
   - Parallel work after #104 merges
   - Recommend Livingston
   - Effort: 2-3 days per rubric set
   - Work: Azure ML `.prompty` schema translation to waza YAML
3. **#138 (P1, multi-model recommendation engine)** — **Blocked on #104 + #39 (now merged in PR #152)**
   - Capstone E3 feature (highest complexity)
   - Recommend Linus
   - Effort: 3-4 days (rubric design + aggregation + LLM judging)
   - Requires: result aggregation, statistical analysis, recommendation rubric design
4. **Critical path blocker:** #104 (Prompt Grader) unblocks 50% of E3 backlog. Recommend prioritizing merge in parallel track.

**Key decisions captured:**
- Suggestion engine must consolidate logic between `waza dev` (E2) and `waza run --suggestions` (E3)
- Rubric porting establishes reusable pattern for future evaluators (Azure ML schema mapping → YAML)
- #138 recommendation rubric needs design clarity: primary optimization target (cost vs. quality vs. balanced)?

**Why:** Unblocks sprint planning. Clear prioritization and dependency analysis reduces rework. #44 is ready immediately for squad momentum.

### 2026-02-17: User directive — deprecate Python release pipeline
**By:** Shayne Boyer (via Copilot)
**What:** Deprecate the Python build pipeline (release.yaml). The primary implementation is Go — releases should ship cross-platform Go binaries, not Python wheels.
**Why:** User request — devs are building from source because only Python artifacts are published. The Go CLI is the primary implementation.

### 2026-02-17: JSON-RPC server (#16) is independent from Web UI
**By:** Shayne Boyer (via Copilot)
**What:** #16 (JSON-RPC server) should be treated as a standalone workstream, separate from Web UI (#14). It enables bidirectional communication between waza and frontend clients (VS Code, JetBrains, etc.). The RPC layer must stay in sync with CLI/engine changes and have strong test coverage.
**Why:** User directive — JSON-RPC is an integration layer, not a Web UI dependency. Decoupling allows it to ship independently and serve multiple consumers.

### 2026-02-17: Recommendation engine normalization with 2 models (#138)
**By:** Linus (Backend Dev)
**Related:** #138, PR #165
**What:** The `internal/recommend/` package uses min-max normalization which inherently maps the best model to 10.0 and worst model to 0.0 for every metric. With exactly 2 models, winner always scores 10.0 and runner-up always scores 0.0 — making margin percentage meaningless (guarded to 0). With 3+ models, mid-range models produce meaningful intermediate scores.
**Why:** This is a known limitation of Phase 1's pure heuristic approach. Phase 2 (LLM-powered analysis, depends on #104) should consider alternative normalization (e.g., absolute scale, z-score) if 2-model comparisons are a common use case. For now, the recommendation itself (which model is best) is still correct — only the margin percentage is affected.

### 2026-02-17: Workspace resource setup is shared across engines
**By:** Linus (Backend Dev)
**Related:** #97, PR #159
**What:** `setupWorkspaceResources()` in `internal/execution/workspace.go` is the single source of truth for writing request resources into a workspace directory with path-traversal protection. Both CopilotEngine and MockEngine delegate to it. Any new engine that needs workspace isolation should use this function rather than rolling its own.
**Why:** Prevents divergence between engine implementations — if we change resource setup rules (e.g. permissions, symlink handling), it changes in one place.

### 2026-02-17: Go CLI release uses v* tags (not azd extension tag pattern)
**By:** Linus (Backend Dev)
**Related:** PR #155
**What:** The Go CLI release pipeline (`.github/workflows/go-release.yml`) triggers on `v*` tags (e.g., `v1.0.0`). This is intentionally different from the azd extension release which uses `azd-ext-microsoft-azd-waza_VERSION` tags. The two release pipelines are independent — pushing a `v*` tag releases Go CLI binaries, not the azd extension.
**Why:** Standard Go convention uses `v` prefixed semver tags. The azd extension has its own tag namespace to avoid collisions. Teams referencing release tags must use the correct pattern for the artifact they're targeting.

### 2026-02-17: Diff grader snapshot resolution requires context_dir
**By:** Linus (Backend Dev)
**Related:** #99, PR #158
**What:** The `diff` grader's `snapshot` paths resolve relative to a `context_dir` parameter. The runner must pass the fixtures/context directory as `context_dir` in the grader params so snapshot files can be found. Without it, snapshot paths are treated as-is (absolute or relative to cwd), which may break in sandboxed execution.
**Why:** Basher needs to know this when writing tests — snapshot fixture paths won't resolve without `context_dir`. Saul needs to document the `context_dir` param in GRADERS.md.

### 2026-02-17: Session review findings — PRs #154-#161
**By:** Rusty (Lead)
**Related:** PRs #154, #155, #158, #159, #160, #161
**What:** PR review session identified three issues:
1. **install.sh: `sha256sum` not available on macOS** (PR #155, line 81) — needs OS detection for sha256sum vs shasum. Severity: Medium.
2. **No diff_grader_test.go** (PR #158) — missing unit tests vs pattern of other graders. Severity: Low (logic sound, pattern inconsistency).
3. **No workspace_test.go** (PR #159) — path-traversal protection has no dedicated unit tests. Severity: Low (straightforward logic, security-adjacent).

Approved: PR #154 (shutdown), #160 (tool_call rubrics), #161 (task rubrics).
**Why:** Ensure quality before merging. macOS install blocker is actionable. Test coverage gaps documented for follow-up.

### 2026-02-17: Phase 1 heuristic recommendation engine (#138) — Design ready
**By:** Rusty (Lead)
**Related:** #138, #39, #104
**What:** Ship heuristic recommendation engine with `--recommend` flag (no LLM). Computes weighted avg of normalized metrics: aggregate score 40%, pass rate 30%, consistency 20%, speed 10%. Output: per-model recommendation + human reasoning + JSON metadata. Handles edge cases (single model, all fail, ties). New package: `internal/recommend/` (engine, normalizers). Data structures: `internal/models/recommendation.go`. Weights/normalization immutable Phase 1. Phase 2 (LLM-powered): deferred pending #104 completion.
**Why:** Unblocks multi-model recommendation without waiting for LLM grader. Two-model normalization limitation (winner 10.0, runner-up 0.0) noted but doesn't affect winner ID — only margin percentage (0%).

### 2026-02-17: Issue #65 scope analysis — azure.yaml integration MVP
**By:** Rusty (Lead)
**Related:** #65, E7 (AZD Extension)
**What:** Phase 1 MVP: Add `tools.waza` section to azure.yaml with runtime config (executor, model, timeout_seconds, trials_per_task, parallel, max_workers), skill_directories, tokens (enabled, soft/hard limits, mode). New package: `internal/azureconfig/` reads config. CLI applies as defaults; flags override. Phase 2 (lifecycle hooks: pre_deploy, validation) deferred — requires azd lifecycle event support. Scope: ~7.5d (schema, reader, CLI integration, validation, tests, docs).
**Why:** Centralize config, reduce flag noise, enable team standards. Backwards compatible (missing `tools.waza` silently ignored). Clear phase boundary prevents scope creep.

### 2026-02-17: Go release pipeline architecture — design spec
**By:** Rusty (Lead)
**Related:** E7 (#62)
**What:** Go CLI release (separate from azd extension). Trigger: `v*` tags or `workflow_dispatch`. Matrix: 6 platforms (linux/darwin/windows × amd64/arm64). Binaries: `waza-{os}-{arch}[.exe]`. Version injection: `-ldflags "-X main.version=$VERSION"`. Artifacts: 6 binaries + `waza-{VERSION}.sha256` + auto-generated release notes (from merged PRs). install.sh: platform detection + download + SHA256 verification. New files: `.github/workflows/go-release.yml` (250-300L), `install.sh` (50-70L). Existing: `release-python-legacy.yaml` renamed/disabled.
**Why:** Standard Go patterns. User-friendly install. Cross-platform tested. Verifiable checksums. Clear separation from Python legacy.
