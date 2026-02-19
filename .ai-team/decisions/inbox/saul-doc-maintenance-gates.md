# Decision: Documentation Maintenance Routing (Issue #256)

**By:** Saul (Documentation Lead)

**Date:** 2026-02-19

**Status:** Implemented

## What

Established Saul (Documentation Lead) as the documentation quality gate. Added two new PR review rules to ensure docs stay current as features change:

1. **Doc-review gate** (Rule 9): Saul reviews PRs touching CLI code (`cmd/waza/`, `internal/`, `web/src/`) for documentation impact
2. **Doc-consistency gate** (Rule 10): Saul reviews PRs touching documentation files for style consistency and accuracy

Added Documentation Impact Matrix mapping code paths to required doc updates, showing which doc files must be checked when specific code changes.

## Why

**Problem:** Documentation was reactive rather than proactive. Code changes happened without corresponding doc updates. Screenshots became stale. Examples diverged from actual behavior. No clear responsibility for doc freshness.

**Solution:** Make documentation review a first-class routing rule, like code review. Saul owns ongoing doc-freshness verification across all PRs. The Impact Matrix provides clear guidance on what needs checking for each code path.

## Scope

- **routing.md:** Added Rules 9–10 and Documentation Impact Matrix
- **charter.md:** Added doc-freshness reviews to "What I Own" and PR monitoring to "How I Work"
- **AGENTS.md:** Added Documentation Maintenance section with tables for "When to Update Docs" and screenshot regeneration steps
- **history.md:** Recorded doc-freshness reviews as a key learning

## Impact

- All code PRs (`cmd/waza/`, `internal/`, `web/src/`) now automatically routed to Saul for doc-impact review
- All doc PRs (`docs/`, `README.md`, `DEMO-SCRIPT.md`) routed to Saul for consistency check
- Clear accountability: Saul owns the matrix and updates it as new paths are discovered
- Screenshot maintenance can be automated via Playwright tests referenced in AGENTS.md

## Next Steps

- Communicate the new gates to the team (especially Linus, Rusty, Livingston)
- Start monitoring PRs against the Impact Matrix
- Iterate on the matrix as patterns emerge
