# Linus — History

## Learnings

### 2026-02-18: Result groupBy categorization (#188)
**What:** Added `GroupBy` field to `Config`, `Group` field to `TestOutcome`, `GroupStats` struct, and grouped CLI output. The runner resolves groups via `resolveGroup()` — currently only `"model"` is supported (uses `spec.Config.ModelID`). Group stats are computed in `computeGroupStats()` using insertion-order-preserving accumulation and attached to `OutcomeDigest.Groups`. CLI prints a "RESULTS BY GROUP" section when groups are present. No GroupBy = unchanged flat output (backward compatible).

**Key design decisions:**
- `GroupStats` lives in `models/outcome.go` alongside `OutcomeDigest` — keeps the data model cohesive.
- `Groups []GroupStats` added directly to `OutcomeDigest` (not a separate top-level field) so it serializes naturally under `"summary"` in JSON output.
- `resolveGroup()` is a method on `TestRunner` with a switch on `GroupBy` value — extensible for future CSV column grouping (#187) without touching the stats computation.
- `computeGroupStats()` is a package-level function (not a method) since it only needs the outcomes slice — makes it independently testable.
- Group insertion order is preserved via a separate `order` slice to ensure deterministic output regardless of map iteration order.

### 2026-02-18: Template variable support (#186)
**What:** Created `internal/template/` package with `Context` struct and `Render()` function. Uses Go `text/template` with `missingkey=error` for strict variable resolution. Fast path skips template parsing when input contains no `{{` delimiters.

**Key design decisions:**
- `Context` holds system variables (`JobID`, `TaskName`, `Iteration`, `Attempt`, `Timestamp`) as typed fields and user-defined variables in a `Vars map[string]string`. This separates known system state from arbitrary user inputs.
- `Option("missingkey=error")` ensures unresolved variables produce clear errors rather than silent empty strings. This catches typos in eval YAML early.
- Fast path: if the input string contains no `{{`, return it unchanged with zero allocation. Most eval YAML fields won't use templates, so this avoids unnecessary parsing overhead.
- Error wrapping uses `fmt.Errorf("template: %w", err)` with `parse:` or `render:` sub-prefix to distinguish parse-time vs execution-time failures.
- Package is standalone — no integration into runner.go yet. CSV support (#187) will wire `Render` into task expansion.

### 2026-02-18: Lifecycle hooks for eval runs (#185)
**What:** Added lifecycle hooks (before_run, after_run, before_task, after_task) to the evaluation runner. Created `internal/hooks/hooks.go` with `Runner.Execute()` that runs shell commands via `os/exec` with configurable exit code validation, working directory, and error-on-fail semantics.

**Key design decisions:**
- `HooksConfig` types live in `internal/hooks/` and are referenced from `BenchmarkSpec.Hooks` via import — keeps hook logic self-contained and testable.
- `hooks.Runner` is stored on `TestRunner` struct so both `runSequential` and `runConcurrent` can access it without passing through every method signature.
- `after_run` hooks always fire via `defer` — even if `before_run` or task execution errors out. This ensures cleanup hooks (e.g. teardown scripts) run reliably.
- `before_run` failure aborts the entire benchmark with a wrapped error. `before_task` failure marks that single task as `StatusFailed` and skips execution but continues to the next task.
- `after_task` failures are logged as warnings and never abort — after-hooks are observational (e.g. metrics collection) and shouldn't block the run.
- Empty `ExitCodes` defaults to `[0]` — standard UNIX convention. Non-matching exit codes with `ErrorOnFail: false` log a warning and continue.
- No hooks configured = zero overhead — all hook call sites guard with `len(spec.Hooks.X) > 0` checks.

### 2026-02-18: Retry/attempts mechanism (#184)
**What:** Added `MaxAttempts` field to `Config` struct and an inner retry loop inside `runTestUncached()`. Each trial (run) now retries up to `max_attempts` times on grader failure before recording the final result. The `RunResult.Attempts` field tracks how many attempts were consumed.

**Key design decisions:**
- Default `MaxAttempts` is treated as 1 when omitted or zero — preserves backward compatibility with no behavioral change.
- Retry only on `StatusFailed` (grader failures). `StatusPassed` exits immediately (success), `StatusError` exits immediately (infrastructure errors shouldn't be retried — they indicate engine/grader setup problems, not flaky agent responses).
- The retry loop lives inside the existing trial loop in `runTestUncached()`, not in `executeRun()`. This keeps `executeRun()` as a pure single-execution function and makes the retry boundary explicit.
- Retry logging (`[RETRY]` prefix) only emits in verbose mode to keep non-verbose output clean.
- The `Attempts` field on `RunResult` records the attempt that produced the final result (1 = first try succeeded, N = took N attempts).

### 2026-02-18: CSV dataset support (#187)
**What:** Created `internal/dataset/` package with `LoadCSV()` and `LoadCSVRange()` functions. Added `TasksFrom` and `Range` fields to `BenchmarkSpec`. Modified `loadTestCases()` in runner.go to branch: if `spec.TasksFrom` is set, load CSV rows and generate in-memory `TestCase` objects with template-resolved prompts; otherwise fall through to existing file-based task loading.

**Key design decisions:**
- `loadTestCases()` delegates to `loadTestCasesFromCSV()` or `loadTestCasesFromFiles()` — clean separation keeps the original file-loading path completely untouched for backward compatibility.
- CSV row data merges into `template.Context.Vars` with CSV values overriding `spec.Inputs` on conflict. This lets global inputs provide defaults while per-row CSV columns specialize.
- `TestID` resolution: prefers "id" column, then "name" column, then "row-N". `DisplayName` prefers "name" column, then "row-N". This gives CSV authors control over test identity without requiring specific columns.
- `Range` is `[2]int` with 1-based inclusive bounds. Zero values mean "no range filtering" — the zero value of `[2]int` is `[0, 0]`, which naturally means "load all rows".
- Go's `csv.Reader.ReadAll()` enforces consistent field count natively, so our manual column count check is a belt-and-suspenders safety net that only fires if `FieldsPerRecord` is set to -1.
- Model template resolution happens per-row via `resolveModelForRow()` — enables `model: "{{.Vars.model}}"` in eval YAML where each CSV row specifies its own model.
