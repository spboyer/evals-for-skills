# Session: 2026-02-26 Docs Audit & Fixes

**Requested by:** Shayne Boyer

## Summary

Saul (Documentation Lead) and Livingston (Documentation Specialist) performed comprehensive documentation freshness audit.

**Findings:** 8 total
- **3 blocking:** tool_constraint missing from graders guide, AZD extension URL points to v0.8.0 (should be v0.9.0), token commands documented but not implemented
- **5 stale:** eval-yaml.mdx missing 7 config fields, no v0.9.0 feature examples, graders guide incomplete, workspace-aware tokens not documented, camelCase keys in schema example

## Action Taken

**Livingston (PR #454):**
- Fixed `site/src/content/docs/reference/cli.mdx` — added `--format json` flag to `waza check` command
- Fixed `site/src/content/docs/reference/schema.mdx` — updated `.waza.yaml` example from snake_case to camelCase keys (timeout_seconds → programTimeout, skills_dir → skillsDir, evals_dir → evalsDir, fixtures_dir → fixturesDir)

**Linus (PR #455):**
- Fixed `site/src/content/docs/guides/graders.mdx` — added tool_constraint grader section with config options and use cases
- Fixed `site/src/content/docs/guides/eval-yaml.mdx` — added 6 missing config fields (group_by, inputs, tasks_from, max_attempts, hooks, weight) with examples
- Fixed `site/src/content/docs/reference/releases.mdx` — updated AZD extension download URL from v0.8.0 to v0.9.0

Both PRs created from separate worktrees in parallel.

---

**Session Date:** 2026-02-26  
**Decision Files:** Merged from `.ai-team/decisions/inbox/saul-docs-freshness-audit.md` and `livingston-site-audit.md`
