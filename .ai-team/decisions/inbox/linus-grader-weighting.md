# Decision: Grader Weighting Design

**By:** Linus (Backend Dev)
**Date:** 2026-02-20
**Issue:** #299

## What

Added optional `weight` field to grader configs for weighted composite scoring. Key design choices:

1. **Weight lives on config, not on the grader interface.** Graders don't know their own weight — the runner stamps it onto `GraderResults` after grading. This keeps grader implementations simple and weight-unaware.

2. **Default weight is 1.0** (via `EffectiveWeight()`). Zero and negative values are treated as 1.0. This means all existing eval.yaml files produce identical results — no migration needed.

3. **Weighted score is additive, not a replacement.** `AggregateScore` (unweighted) is preserved. `WeightedScore` is a new parallel field. The interpretation report only shows weighted score when it differs from unweighted.

4. **Weight flows through the full pipeline:** `GraderConfig.Weight` → `GraderResults.Weight` → `RunResult.ComputeWeightedRunScore()` → `TestStats.AvgWeightedScore` → `OutcomeDigest.WeightedScore`. Web API also carries weight per grader result.

## Why

Weighted scoring lets users express that some graders matter more than others (e.g., correctness 3× more important than style). Without breaking existing pass/fail semantics.

## Impact

- `internal/models/` — new fields on `GraderConfig`, `GraderResults`, `TestStats`, `OutcomeDigest`
- `internal/orchestration/runner.go` — weight stamping in `runGraders`, weighted stats in `computeTestStats`/`buildOutcome`
- `internal/reporting/` — conditional weighted score display
- `internal/webapi/` — weight exposed in API responses
- JSON schema unchanged (eval.yaml schema is separate from waza-config.schema.json)
