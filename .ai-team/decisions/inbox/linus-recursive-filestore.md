# Decision: FileStore Schema Validation Pattern

**Date:** 2025-01-26  
**Author:** Linus (Backend Developer)  
**Status:** Implemented in PR #280

## Context

The FileStore loads JSON files from a results directory, but it was accepting any JSON that could be unmarshaled into an `EvaluationOutcome` struct—even files like `summary.json` that aren't actually evaluation outcomes. Go's `json.Unmarshal` silently ignores unknown fields and produces zero-value structs, leading to phantom runs appearing in the dashboard.

## Decision

**Validate unmarshaled JSON to ensure it matches the expected schema.**

For `EvaluationOutcome`, a valid outcome must have at least one of:
- `BenchName != ""` (the `eval_name` field in JSON)
- `Digest.TotalTests > 0` (from `summary.total_tests` in JSON)

Files that don't meet this criteria are silently skipped (not loaded into the store).

## Implementation

```go
var outcome models.EvaluationOutcome
if err := json.Unmarshal(data, &outcome); err != nil {
    return nil
}

// Validate that this is a real EvaluationOutcome
if outcome.BenchName == "" && outcome.Digest.TotalTests == 0 {
    return nil
}
```

## Alternatives Considered

1. **Require a version field** — Rejected: would break existing result files
2. **Check for specific required fields** — Rejected: too brittle, evaluation outcomes can vary
3. **Use JSON schema validation** — Rejected: overkill for this simple check
4. **Log warnings for invalid files** — Rejected: silent skipping is fine for file scanning

## Consequences

✅ **Good:**
- Filters out `summary.json` and other non-outcome JSON files automatically
- No breaking changes to existing valid outcome files
- Simple and maintainable validation logic

⚠️ **Watch for:**
- If we ever want to support "minimal" outcomes with no bench name and zero tests, this validation would need adjustment
- Tests must always have either a name or at least one test case

## Related Changes

This decision was implemented alongside recursive directory scanning (PR #280), which replaced `os.ReadDir()` with `filepath.WalkDir()` to support nested result directories.
