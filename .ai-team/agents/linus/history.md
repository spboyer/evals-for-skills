# Linus — History

## Learnings

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
