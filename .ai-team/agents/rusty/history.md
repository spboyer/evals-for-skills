# Rusty — Decision History

## Learnings

### 2026-02-17: Unscoped Issue Triage — #2, #10, #14, #16, #21

**What I Did:**
- Triaged 5 unscoped issues (no priority, no epic) against current project state (v0.4.0-alpha.1, Go CLI).
- Identified Web UI cluster (#2, #14, #16) as the primary overlap. #2 is a Python-era duplicate of #14.
- Commented on all 5 issues with structured triage analysis. Applied priority:p2 + epic labels + release:backlog to 4 issues. Recommended closing #2.

**Key Architecture Decisions:**

1. **Web UI cluster decomposition:** #16 (JSON-RPC) is the foundation layer, #14 (Web UI) wraps it via HTTP. These must remain separate issues — JSON-RPC delivers IDE integration value independently of any Web UI.

2. **Epic gap identified:** No existing epic covers Web UI / IDE integration. Used `epic:go-cli` as closest fit since `waza serve` and `waza jsonrpc` are CLI commands. Noted that `epic:web-ui` (E8) may be needed if this work becomes a priority.

3. **All 5 issues are P2:** None of these block core eval workflows. CLI-first remains the correct priority. Web UI, IDE integration, evaluation-first workflow, and session logging are all stretch goals relative to E1-E4 completion.

4. **Evaluation-first workflow (#10):** Option B (`--no-skill` flag) is the right approach — minimal CLI surface, no new commands, fits existing `waza run` pattern.

5. **Session logging (#21):** Scoped to data format only (Copilot CLI-compatible JSON). A standalone viewer would be scope creep — Copilot CLI tools already handle visualization.

**Output:** Decision doc at `.ai-team/decisions/inbox/rusty-unscoped-triage.md`. Comments posted on all 5 issues. Labels applied to #10, #14, #16, #21.

### 2026-02-17: Issue #65 Scoping — azure.yaml Integration (E7-04)

**What I Did:**
- Investigated current state of azure.yaml in azd projects + waza CLI
- Examined existing `BenchmarkSpec` and `BenchmarkConfig` to understand config patterns
- Reviewed extension.yaml for azd packaging context
- Analyzed all existing waza commands to map which settings apply where
- Designed Phase 1 (config-only) vs Phase 2 (lifecycle hooks) boundary

**Key Scoping Decisions:**

1. **MVP = Config-Only (Phase 1)**
   - New `tools.waza` section in azure.yaml for runtime config + skill directories + token limits
   - No automatic execution; no lifecycle hooks in Phase 1
   - Rationale: Unblocks users immediately, minimal disruption to `azd up`/`azd deploy`, clear phase boundary

2. **What Goes in `tools.waza`:**
   - Runtime config: `executor`, `model`, `timeout_seconds`, `trials_per_task`, `parallel`, `max_workers`
   - Paths: `skill_directories` (array of directories to scan)
   - Token enforcement: `tokens.enabled`, `skill_soft_limit`, `skill_hard_limit`, `mode`
   - NOT in Phase 1: hooks for pre_deploy, validation, etc.

3. **CLI Integration Pattern:**
   - New `internal/azureconfig/` package loads YAML
   - Settings applied as defaults; CLI flags override (standard precedence)
   - Backwards compatible: missing `tools.waza` → silent ignore
   - `waza run`, `waza check`, `waza tokens` all inherit config

4. **Phase 2 Deferred (Lifecycle Hooks):**
   - `hooks.pre_deploy` — run evals before `azd deploy`
   - `hooks.validation` — run quick checks during `azd up`
   - Requires azd extension lifecycle event support (not available yet in Phase 1)
   - Prevents scope explosion; gives time for proper design

**Work Breakdown (Phase 1):**
- W1: Schema design (0.5d)
- W2: Config reader `internal/azureconfig/` (2d)
- W3: CLI integration (1.5d)
- W4: Validation + defaults (1d)
- W5: Unit + integration tests (1.5d)
- W6: Documentation (1d)
- **Total: ~7.5 days**

**Risks Identified:**
- Scope creep on hooks — clearly communicate Phase 2 deferral
- Token mode ambiguity — Phase 1 logs warnings only; Phase 2 ties to `waza tokens` command
- Backwards compatibility — handled via optional section read

**Acceptance Criteria:**
- [ ] Schema in azure.yaml documented
- [ ] Reader package with unit tests
- [ ] CLI integration with override precedence tests
- [ ] Example azure.yaml file
- [ ] README updated
- [ ] All existing tests pass

**Output:** Detailed scoping spec written to `.ai-team/decisions/inbox/rusty-65-azure-yaml-scope.md` — ready for Linus to implement.

**Why This Matters:**
- E7-04 is vague ("waza integrates with azure.yaml"). Scoping transforms it into concrete Phase 1 + Phase 2 work
- Establishes config pattern reusable for future azd extensions
- Prevents over-engineering; hooks can be added later when azd support is clearer

### 2026-02-17: Phase 1 Recommendation Engine Design (#138)

**What I Did:**
- Read existing team decisions to understand architecture patterns
- Analyzed current multi-model run flow (PR #152) and `waza compare` command
- Examined `EvaluationOutcome` and `OutcomeDigest` structures
- Designed Phase 1 (heuristic-only) recommendation engine
- Wrote full design spec for Linus to implement from

**Key Design Decisions:**

1. **Scope: Heuristic-Only (No LLM)**
   - Phase 1 ships NOW without #104 (Prompt Grader)
   - Uses weighted average: 40% aggregate score + 30% pass rate + 20% consistency + 10% speed
   - All components normalized to 0–10 scale
   - Rationale: Unblocks multi-model comparison use case; LLM analysis deferred to Phase 2

2. **Architecture: New `internal/recommend/` Package**
   - Separates recommendation logic from CLI layer
   - `RecommendationEngine` orchestrates heuristic scoring
   - Normalizers handle per-component 0–10 scaling
   - `Recommendation` struct holds structured output for JSON + terminal

3. **Integration: Non-Invasive CLI Change**
   - Single new flag: `--recommend` (boolean, optional)
   - Computes after comparison table in `runCommandE()`
   - Attaches recommendation to outcome `metadata` field for `--output`
   - Backwards compatible: defaults false, affects no existing behavior

4. **Output: Dual Format**
   - **Terminal:** Human-readable summary with reasoning, component scores table, margin of victory
   - **JSON:** Full `Recommendation` struct in `outcome.Metadata["recommendation"]`
   - Design mirrors `waza compare` command patterns for consistency

5. **Edge Case Handling**
   - Single model + `--recommend` → silently skips (no error)
   - All models fail → no recommendation computed (nil)
   - Tied scores → first-in-order wins (stable sort)
   - Partial failures → uses available data, handles nil gracefully

**Why This Design:**

- **Minimal disruption:** `--recommend` is optional, existing single-model workflows unaffected
- **Clear separation:** Recommendation logic in dedicated package, not mixed into cmd layer
- **Reusable components:** Normalizers and scoring can be extended/customized for Phase 2 LLM weights
- **Testability:** Each component independently testable; integration tests verify flag + output
- **Transparency:** Weights and per-component scores exposed in JSON for audit/debugging
- **Extensibility:** Weights can become configurable in Phase 2; foundation set for custom rubrics

**What Linus Gets:**

Full spec with:
- Struct definitions for `Recommendation`, `RecommendationWeights`, `ModelScore`
- Algorithm pseudocode for `scoreModels()`, `normalizeMetrics()`, `findWinner()`
- Integration points in `cmd_run.go` (where to call engine, how to attach metadata)
- Terminal output format with ASCII box drawing
- JSON schema for recommendation metadata
- Edge case handling rules
- Test strategy (unit + integration)

**Why This Matters:**

#138 is a **capstone E3 feature** — highest complexity in evaluation framework. Phase 1 unblocks:
- Users can run `waza run eval.yaml --model gpt-4o,claude-sonnet-4 --recommend` today
- Establishes recommendation engine pattern for future phases
- No dependency on #104 — parallel work path
- Prepares foundation for Phase 2 LLM-powered analysis without rework

### 2026-02-17: Post-Merge Review — PRs #162, #163

**What I Reviewed:**
- PR #162: docs: update documentation for v0.4.0-alpha.1 features
- PR #163: fix: install.sh checksum verification on macOS

**Findings:**

1. **PR #163 (install.sh macOS fix)** — ✅ Approved. The `sha256sum → shasum -a 256` fallback at lines 81-87 is correct. Logic: try `sha256sum` first (Linux), fall back to `shasum -a 256` (macOS), warn if neither exists. Both branches use `--status` for silent verification. This directly fixes the issue I flagged in my PR #155 review.

2. **PR #162 (docs update)** — ✅ Approved with 2 issues noted:
   - **CHANGELOG PR numbers are wrong:** Line 16 says diff grader is `#156` and line 21 says MockEngine workspace is `#157`. Actual PRs are `#158` (diff grader) and `#159` (MockEngine workspace). PR `#157` was the extension registry chore commit. These features also merged *after* the v0.4.0-alpha.1 tag — they should be listed under `[Unreleased]`, not under `[0.4.0-alpha.1]`.
   - **CHANGELOG missing `skill_invocation` grader:** The `skill_invocation` grader (#146) shipped in v0.4.0-alpha.1 (per git log) but is not mentioned in the CHANGELOG. It IS correctly documented in README.md (line 438) and GRADERS.md (lines 495-577).

3. **README.md** — ✅ Accurate. All 8 grader types listed in the table match registered graders in `internal/graders/grader.go`. Commands (`check`, `cache clear`, `dev`, `init`, `generate`, `compare`, `tokens count/suggest`) all verified in source. Version references consistent with `version.txt` (0.4.0-alpha.1).

4. **GRADERS.md** — ✅ Comprehensive. `diff` grader documentation (lines 203-293) accurately describes the implementation in `diff_grader.go`. `skill_invocation` grader documentation (lines 495-577) correctly documents modes, scoring, and data source. One observation: the "Creating Custom Graders" section at the bottom (lines 953-976) still shows Python code — this is legacy content from the Python era. Non-blocking but should be updated to Go eventually.

5. **DEMO-GUIDE.md** — ✅ Good. 8 demos covering all major features. Demo 8 (Azure ML Rubrics) correctly notes the `prompt` grader dependency. Commands use `./waza-bin` consistently.

6. **examples/rubrics/README.md** — ✅ Accurate. Correctly notes prompt grader is blocked on #104. All 8 rubric YAML files verified present on disk. Schema documentation matches actual YAML structure.

**Summary:**
- install.sh fix: Ship it. Clean platform detection, correct fallback chain.
- Docs: Two CHANGELOG inaccuracies need fixing (wrong PR numbers, missing skill_invocation entry). Everything else is accurate and well-structured.

### 2026-02-17: Session Code Review — PRs #154-#161

**What I Reviewed:**
- PR #154: engine.Shutdown() leak fix in runSingleModel
- PR #155: Go cross-platform release workflow + install.sh
- PR #158: diff grader for workspace changes (#99)
- PR #159: MockEngine workspace directory support (#97)
- PR #160: port Azure ML tool_call evaluation rubrics (#106)
- PR #161: port Azure ML task evaluation rubrics (#107)

**Findings:**

1. **PR #154 (Shutdown fix)** — ✅ Approved. `defer engine.Shutdown(ctx)` at line 210-213 of `cmd_run.go` is correct. Covers the error return paths that previously leaked. Thorough test coverage in `engine_shutdown_test.go` with SpyEngine, idempotency, and cancelled-context tests.

2. **PR #155 (Release workflow + install.sh)** — ⚠️ Two concerns:
   - `install.sh` line 81 uses `sha256sum` which doesn't exist on macOS. Should use `shasum -a 256` or detect platform. This will break the install experience for Darwin users.
   - `go-release.yml` line 44 specifies `go-version: '1.25'` — should verify this is correct; current stable is likely different. The `-ldflags "-X main.version=${VERSION}"` correctly links to `root.go:12`. Legacy Python workflow properly disabled with `if: false`.

3. **PR #158 (Diff grader)** — ✅ Approved with one gap noted. Clean architecture: path traversal protection via `validatePathInWorkspace`, proper partial scoring, graceful handling of missing workspace. Score math in `countTotalChecks`/`buildResult` is correct. Registration in `grader.go` follows existing patterns. **Gap:** No `diff_grader_test.go` file exists — this is the only grader without dedicated unit tests.

4. **PR #159 (MockEngine workspace)** — ✅ Approved. Shared `setupWorkspaceResources` in `workspace.go` with path-traversal protection mirrors CopilotEngine behavior. Clean separation. Previous workspace cleanup in `Execute()` prevents accumulation. **Gap:** No dedicated `workspace_test.go` for `setupWorkspaceResources` — the path-traversal protection logic should have unit tests.

5. **PR #160 (tool_call rubrics)** — ✅ Approved. Well-structured YAML: `tool_call_accuracy` (ordinal 1-5), `tool_selection` (binary), `tool_input_accuracy` (binary), `tool_output_utilization` (binary). Consistent schema, proper chain_of_thought sections. These are reference rubrics for future `prompt` grader (#104) — no runtime integration yet.

6. **PR #161 (task rubrics)** — ✅ Approved. `task_adherence` (binary_flag), `task_completion` (binary), `intent_resolution` (ordinal 1-5), `response_completeness` (ordinal 1-5). Good Azure ML attribution. Consistent with tool_call rubric structure. README.md present.

**Key Decisions:**
- install.sh sha256sum issue needs a fix before anyone ships a release targeting macOS users
- diff_grader and workspace_test gaps should be tracked as follow-up work
- Rubrics are documentation-only until #104 (prompt grader) lands — no integration risk

