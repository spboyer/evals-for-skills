### 2026-02-18: Retry/attempts mechanism — design decisions
**By:** Linus (Backend Dev)
**Related:** #184
**What:** `max_attempts` config field (YAML: `max_attempts`, default: 1) adds an inner retry loop inside each trial. On grader failure (`StatusFailed`), the runner re-executes the task up to `max_attempts` times before recording the final result. Infrastructure errors (`StatusError`) are never retried. `RunResult.Attempts` tracks the attempt count in JSON output.
**Why:** Flaky agent responses shouldn't fail an entire trial on first attempt. Retrying grader failures within a trial gives non-deterministic agents a fair shot while keeping infrastructure errors fatal. The retry loop is inside `runTestUncached()` (not `executeRun()`) to maintain single-responsibility — `executeRun()` stays a pure single-execution function.
