### 2026-02-09: Sensei reference repo for E2 implementation
**By:** Squad (Coordinator)
**What:** The sensei engine (E2: #32-38) must adopt functionality from https://github.com/spboyer/sensei. Key patterns to port to Go:

**Scoring Algorithm (references/scoring.md):**
- Low: description < 150 chars OR no triggers
- Medium: description >= 150 chars AND has trigger keywords
- Medium-High: has "USE FOR:" AND "DO NOT USE FOR:"
- High: Medium-High + routing clarity (INVOKES/FOR SINGLE OPERATIONS)
- Rule-based checks: name validation, description length, trigger detection, anti-trigger detection, routing clarity
- MCP integration checks when INVOKES present (4-point sub-score)

**Ralph Loop (references/loop.md):**
1. READ — load SKILL.md + count tokens
2. SCORE — rule-based compliance check
3. CHECK — if >= Medium-High AND tests pass → done
4. SCAFFOLD — create tests from templates if missing
5. IMPROVE FRONTMATTER — add USE FOR/DO NOT USE FOR
6. IMPROVE TESTS — update shouldTrigger/shouldNotTrigger prompts
7. VERIFY — run tests (skip with --fast)
8. CHECK TOKENS — verify budget
9. SUMMARY — before/after comparison
10. PROMPT — Commit/Issue/Skip
11. REPEAT — until target score or max 5 iterations

**Token Management (scripts/src/tokens/):**
- count: count tokens in markdown files, per-file + total
- check: validate against .token-limits.json (soft/hard limits)
- suggest: optimization suggestions
- compare: diff against git refs
- --strict flag: exit 1 if limits exceeded
- Config: .token-limits.json with defaults and overrides

**Config defaults:**
- SKILL.md: 500 soft, 5000 hard token limit
- references/*.md: 2000 tokens each
- Max iterations: 5
- Target score: Medium-High

**Why:** User directive — spboyer/sensei is the reference implementation for waza's sensei engine
