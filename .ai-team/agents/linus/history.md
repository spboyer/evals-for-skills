# Linus — History

## Learnings

### 2026-02-16: Diff grader implementation (PR #158, Issue #99)
- Implemented `diff` grader following file grader's exact pattern: `NewDiffGrader` constructor with args struct, `Grade()` wrapped in `measureTime`, `WorkspaceDir` from grading context, `validatePathInWorkspace` reuse, proportional scoring from `countTotalChecks`.
- Registered in `Create()` factory with mapstructure decoding — same switch-case pattern as all other graders.
- Added `GraderKindDiff` constant to `models/outcome.go`.
- Contains-line fragments use `+`/`-` prefix convention borrowed from unified diff format. No prefix defaults to must-appear.
- Snapshot paths resolve relative to `context_dir` param, which the runner should set to the fixtures directory.

### 2026-02-17: Shared workspace resource setup (PR #159, Issue #97)
- Extracted `setupWorkspaceResources()` from CopilotEngine into `internal/execution/workspace.go` as a package-level function. Both CopilotEngine and MockEngine now delegate to it for resource writing + path-traversal protection.
- MockEngine.Execute() now creates a `waza-mock-*` temp dir, writes resources, and sets `ExecutionResponse.WorkspaceDir`. MockEngine.Shutdown() cleans it up.
- This enables FileGrader (and any future workspace-dependent graders) to work in mock/test scenarios.

### 2026-02-17: Heuristic recommendation engine (PR #165, Issue #138)
- Created `internal/recommend/` package with `Engine` struct and `Recommend()` method. Uses min-max normalization (0–10) with weighted scoring: 40% aggregate score, 30% pass rate, 20% consistency (inverse stddev), 10% speed (inverse duration).
- Added `Recommendation`, `RecommendationWeights`, `ModelScore` structs in `internal/models/recommendation.go`.
- Integrated `--recommend` flag in `cmd_run.go` — fires after `printModelComparison()`, attaches recommendation to outcome `Metadata` map for JSON output.
- Key edge case: with only 2 models, min-max normalization pushes winner to 10.0 and runner-up to 0.0 on every axis, making margin percentage always 0 (division by zero guard). 3+ models produce meaningful margins.
- Added `recommendFlag` to `resetRunGlobals()` per team convention on Cobra flag state isolation.

📌 Team update (2026-02-17): Go CLI release uses v* tags; azd extension uses azd-ext-microsoft-azd-waza_VERSION tags (separate pipelines). — decided by Linus
📌 Team update (2026-02-17): Workspace resource setup is shared via setupWorkspaceResources() in internal/execution/workspace.go — both CopilotEngine and MockEngine delegate to it. — decided by Linus
📌 Team update (2026-02-17): Diff grader snapshot paths resolve relative to context_dir parameter. — decided by Linus
