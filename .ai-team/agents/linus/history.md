# Linus — History

## Learnings

### 2026-02-18: Retry/attempts mechanism (#184)
**What:** Added `MaxAttempts` field to `Config` struct and an inner retry loop inside `runTestUncached()`. Each trial (run) now retries up to `max_attempts` times on grader failure before recording the final result. The `RunResult.Attempts` field tracks how many attempts were consumed.

**Key design decisions:**
- Default `MaxAttempts` is treated as 1 when omitted or zero — preserves backward compatibility with no behavioral change.
- Retry only on `StatusFailed` (grader failures). `StatusPassed` exits immediately (success), `StatusError` exits immediately (infrastructure errors shouldn't be retried — they indicate engine/grader setup problems, not flaky agent responses).
- The retry loop lives inside the existing trial loop in `runTestUncached()`, not in `executeRun()`. This keeps `executeRun()` as a pure single-execution function and makes the retry boundary explicit.
- Retry logging (`[RETRY]` prefix) only emits in verbose mode to keep non-verbose output clean.
- The `Attempts` field on `RunResult` records the attempt that produced the final result (1 = first try succeeded, N = took N attempts).
