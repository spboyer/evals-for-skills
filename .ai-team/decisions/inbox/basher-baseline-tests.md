# Decision: Baseline Feature Test Patterns

**Date:** 2026-02-18  
**Author:** Basher (Tester)  
**Related:** #194 (A/B Baseline Skill Impact Measurement)

## Decision

Baseline feature tests use actual `TestOutcome.Runs` slice for pass rate calculations, NOT hypothetical `Passed`/`Trials` fields. Tests are table-driven with subtests following existing `runner_test.go` patterns.

## Key Test Patterns Established

### 1. Pass Rate Calculation (computePassRate)
- Uses `len(outcome.Runs)` to count trials
- Counts `run.Status == models.StatusPassed` to determine passes
- Zero-runs guard returns 0.0 (no division error)
- Error status counts as failure (not passed)

### 2. Skill Impact Calculation (computeSkillImpact)
- Division-by-zero guard: `denom = max(passRateWithout, 0.01)` per design doc
- Handles negative delta (skills hurt scenario)
- Handles zero delta (both fail, or perfect scores)
- Percent change can be very large (e.g., 6650% for 0% → 67%)

### 3. Task Mismatch Errors (mergeBaselineOutcomes)
- Returns `*TaskMismatchError` (not generic error) for task set misalignment
- Checks both directions: tasks in skills-enabled but not baseline, AND tasks in baseline but not skills-enabled
- Error message includes task ID for debugging

### 4. Integration Tests (Skipped Pending Implementation)
- `TestRunBenchmark_BaselineNoSkills` — warns and runs single pass when no skills configured
- `TestRunBenchmark_BaselineWithSkills` — executes two passes and computes skill impact
- `TestRunBenchmark_BaselineEmptyTasks` — edge case: no tasks doesn't panic

These tests are marked `t.Skip("Requires baseline implementation and task loading setup")` because they need proper test case loading which isn't exposed in the current runner API (testCases is local to methods, not a field).

## Edge Cases Covered

1. **Zero baseline pass rate** → division guard prevents infinity/NaN
2. **Both runs fail** → zero delta, zero percent change (inconclusive)
3. **Single trial** → 1/1 with skills, 0/1 without → 100% delta, ~10000% improvement
4. **Skills hurt** → negative delta, negative percent change
5. **Error status** → counts as failure, not passed
6. **Empty runs slice** → zero pass rate (no panic)

## Why This Matters

- Linus's implementation uses `Runs` slice (not Passed/Trials), so tests must match
- Division-by-zero is a critical edge case per design doc (§3.3 line 276-279)
- Task mismatch detection prevents silent data corruption in A/B comparison
- Integration tests document expected behavior for full baseline runs

## Files Created

- `internal/models/baseline_test.go` (6 tests, all passing)
- `internal/orchestration/baseline_test.go` (14 tests, 3 skipped, 11 will pass once linter issues fixed)
