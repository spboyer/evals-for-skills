### 2026-02-18: Lifecycle hooks implementation (#185)
**By:** Linus (Backend Dev)
**Related:** #185
**What:** Lifecycle hooks execute shell commands at four points: before_run, after_run, before_task, after_task. The `internal/hooks` package owns execution logic. `HooksConfig` is embedded in `BenchmarkSpec` via yaml tag `hooks`. Runner orchestration calls hooks with these semantics: before_run failure aborts the run; before_task failure skips that task (marks failed); after_run always fires (defer); after_task failures are warnings only. No hooks = no-op, fully backward compatible.
**Why:** Enables eval authors to run setup/teardown scripts (e.g. starting services, cleaning state, collecting metrics) without modifying the waza codebase. Error semantics follow the principle that pre-hooks gate execution while post-hooks are observational.
