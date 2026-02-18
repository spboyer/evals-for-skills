### 2026-02-18: Baseline A/B comparison execution model
**By:** Linus
**What:** The `--baseline` flag runs two sequential passes (not parallel): Pass 1 with skills, Pass 2 without. Each pass gets a fresh engine instance and workspace. The config mutation pattern (save → clear → restore) is used to disable skills for Pass 2. Exit code semantics change when baseline is enabled: 0 = improvement, 1 = regression/neutral.
**Why:** Sequential execution simplifies state management, avoids resource contention, and produces cleaner output. The config mutation pattern is proven and used elsewhere (task filters). Exit code branching allows CI to gate on skill quality — pipelines can require skills to improve or maintain pass rates.
