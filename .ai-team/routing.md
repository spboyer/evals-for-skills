# Work Routing

How to decide who handles what.

## Routing Table

| Work Type | Route To | Examples |
|-----------|----------|----------|
| Architecture, scope, decisions | Rusty | Design new commands, review approaches, dependency decisions |
| Go implementation, CLI commands, graders, engines | Linus | Build cmd files, internal packages, business logic |
| Tests, quality, edge cases | Basher | Write *_test.go files, verify CI, coverage |
| Documentation, README, SKILL.md, PR descriptions | Livingston | Update docs, write SKILL.md, PR body text |
| Code review | Rusty | Review PRs, check quality, suggest improvements |
| Graders, Copilot SDK | Richard Park 👤 | #28 all 8 graders, #29 SDK executor, #22 eval testing, #23 Cobra |
| Sensei engine, compliance, tokens | Charles Lowell 👤 | #32-38 sensei, #47-51 tokens, #33 scoring |
| Session logging | Scribe | Automatic — never needs routing |

## Rules

1. **Eager by default** — spawn all agents who could usefully start work, including anticipatory downstream work.
2. **Scribe always runs** after substantial work, always as `mode: "background"`. Never blocks.
3. **Quick facts → coordinator answers directly.** Don't spawn an agent for "what port does the server run on?"
4. **When two agents could handle it**, pick the one whose domain is the primary concern.
5. **"Team, ..." → fan-out.** Spawn all relevant agents in parallel as `mode: "background"`.
6. **Anticipate downstream work.** If a feature is being built, spawn Basher to write test cases from requirements simultaneously.
7. **Human routing → pause.** When work routes to Richard or Charles, present it to Shayne and wait for their input.
