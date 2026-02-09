# Team Decisions

> Shared brain — all agents read this before starting work.

<!-- Scribe merges decisions from .ai-team/decisions/inbox/ into this file. -->
### 2026-02-09: Proper git workflow required
**By:** Shayne Boyer (via Copilot)
**What:** All issues must follow: feature branch → commit → push → PR with `Closes #N` → @copilot review → address feedback → merge. No direct commits to main.
**Why:** User directive

### 2026-02-09: Sensei reference repo
**By:** Squad (Coordinator)
**What:** E2 sensei engine must adopt patterns from https://github.com/spboyer/sensei. Scoring: Low/Medium/Medium-High/High. Ralph loop: READ→SCORE→CHECK→SCAFFOLD→IMPROVE→VERIFY→TOKENS→SUMMARY→PROMPT→REPEAT. Token management: count/check/suggest/compare with .token-limits.json config.
**Why:** User directive — spboyer/sensei is the reference implementation

### 2026-02-09: Monitor human engineer comments
**By:** Shayne Boyer (via Copilot)
**What:** Periodically check for comments from Charles and Richard on open issues. Follow up on responses.
**Why:** User directive
