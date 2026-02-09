# Decision: Run command test patterns

**Date:** 2026-02-09
**Author:** Linus (Backend Dev)
**Issue:** #24

## Context

The `waza run` command uses package-level variables for Cobra flag binding (`contextDir`, `outputPath`, `verbose`). These persist across test cases in the same process.

## Decision

Each test that calls `cmd.Execute()` must explicitly reset the package-level flag vars at the top of the test body:

```go
contextDir = ""
outputPath = ""
verbose = false
```

This prevents state leakage between tests and keeps them order-independent.

## Alternatives Considered

- Refactoring to pass flags through a struct: cleaner, but would be a larger change to the existing working implementation. Deferred to a future refactor.

## Consequences

- All future `cmd_run_test.go` tests must follow this pattern.
- If new flags are added to `newRunCommand()`, the reset block in tests must be updated accordingly.
