# Rusty — Decision History

## Learnings

### 2026-02-17: Multi-part session summary
- **Unscoped issue triage (#2, #10, #14, #16, #21)**: Decomposed Web UI cluster. JSON-RPC (#16) is independent foundation. All 5 issues P2. Recommended closing #2 (Python-era duplicate of #14).
- **#65 scoping (azure.yaml integration)**: Phase 1 config-only (~7.5d). New `tools.waza` section in azure.yaml with runtime config, skill directories, token limits. New `internal/azureconfig/` package. Fully backwards compatible. Phase 2 defers lifecycle hooks pending azd support.
- **E3 evaluation framework backlog triage**: #44 (improvement suggestions) ready now for Linus. #106 (tool_call rubrics), #107 (task rubrics) blocked on #104 (Prompt Grader). #138 (multi-model recommendation) is capstone feature.
- **Phase 1 recommendation engine design (#138)**: Heuristic-only (no LLM). Weighted average: 40% aggregate + 30% pass rate + 20% consistency + 10% speed. New `--recommend` flag. Ships independently of #104. All scores normalized 0–10.
- **Code review session (PRs #154–#161)**: Shutdown fix approved. install.sh macOS sha256sum issue flagged (fixed in #163). diff_grader and workspace missing unit tests noted. Rubric ports (#160, #161) approved. All 6 PRs approved with noted gaps for follow-up.
- **Documentation review (PRs #162–#163)**: CHANGELOG inaccuracies identified (wrong PR numbers, missing skill_invocation entry). README/GRADERS accurate. install.sh platform-detection fix is clean and ships well.

### 2026-02-15: Multi-model execution and shutdown testing
- Sequential model execution (not concurrent). Each model gets fresh engine. `runSingleModel()` encapsulates full lifecycle. Test failures non-fatal (continue to next model); infrastructure errors abort immediately.
- SpyEngine test double pattern exported for use in cmd-level tests. Engine shutdown lifecycle covers all exit paths. Workspace resource setup is shared (`setupWorkspaceResources()` in `workspace.go`).

### 2026-02-15: E3 backlog prioritization and dependency analysis
- #44 (suggestions) has no blockers — assign to Linus. #106, #107 (rubrics) blocked on #104 (Prompt Grader). #138 (recommendation) unblocks after #39 merges. Consolidation point: `waza dev` and `waza run --suggestions` must share suggestion logic.

### 2026-02-14: User directives consolidated
- Code generation: claude-opus-4.6 only (no exceptions). PR review: dual-model (Opus 4.6 + Codex 5.3 for analytical diversity). Auto-assign unblocked work without asking. Route doc updates to Saul after feature merges. Don't take assigned work — only unassigned issues.

### 2026-02-12: Config patterns and extension conventions
- Functional options for config (e.g., `WithModel()`, `WithTimeout()`). Interface-based design (AgentEngine, Validator). Registry pattern for validators. azd extension uses non-standard tag pattern: `azd-ext-microsoft-azd-waza_VERSION` (e.g., `azd-ext-microsoft-azd-waza_0.2.0`).

### 2026-02-11: Grader contract and error semantics
- Graders depending on SessionDigest must return zero-score result (nil error) when session is nil, not `(nil, error)`. Prevents grader errors from aborting entire run. Zero-score results contribute to scoring; errors abort.

### 2026-02-09: Git workflow and reference patterns
- All issues follow: feature branch → commit → push → PR with `Closes #N` → @copilot review → address feedback → merge. No direct commits to main. Monitor human engineer comments from Charles (@chlowell) and Richard (@richardpark-msft).

### 2026-02-09: Sensei reference patterns
- Compliance scoring: Low (desc < 150 chars OR no triggers) → Medium → Medium-High → High. Ralph loop: READ → SCORE → CHECK → SCAFFOLD → IMPROVE → TESTS → VERIFY → TOKENS → SUMMARY (max 5 iterations). Token management: count/check/suggest/compare with `.token-limits.json`.
