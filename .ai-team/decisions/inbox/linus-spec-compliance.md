# Decision: SpecScorer as separate scorer from HeuristicScorer

**By:** Linus (Backend Developer)
**Date:** 2026-02-20
**PR:** #322
**Issue:** #314

**What:** The agentskills.io spec compliance checks are implemented as a separate `SpecScorer` type rather than extending `HeuristicScorer`. Both run independently — `HeuristicScorer` handles heuristic quality scoring (triggers, anti-triggers, routing clarity) while `SpecScorer` handles formal spec validation (field presence, naming rules, length limits).

**Why:** The two scorers have different concerns: `HeuristicScorer` is about quality/adherence level (Low→High), while `SpecScorer` is about pass/fail conformance to the agentskills.io specification. Keeping them separate means each can evolve independently. The spec may change without affecting heuristic scoring, and vice versa.

**Impact:** Both `waza dev` and `waza check` now run both scorers. Any new spec checks should be added to `SpecScorer` in `cmd/waza/dev/spec.go`. The `SpecResult.Passed()` method only considers errors (not warnings) — warnings like missing license/version don't block readiness.
