# Saul — History

## Project Context
- **Project:** waza — Go CLI for evaluating AI agent skills (Copilot Skills)
- **Tech:** Go, Copilot SDK, YAML eval specs, 7 grader types
- **User:** Shayne Boyer
- **Repo:** spboyer/waza

## Learnings

### 2026-02-14: Initial context
- docs/DEMO-GUIDE.md created with 7 demo scenarios (quick start, graders, tokens, sensei, CI/CD, multi-skill, cross-model)
- docs/GRADERS.md covers all 7 grader types: regex, file, code, behavior, action_sequence, skill_invocation, prompt (pending #104)
- docs/TUTORIAL.md has the getting-started flow
- DEMO-SCRIPT.md at repo root is a presenter narrative (complements DEMO-GUIDE.md)
- examples/code-explainer/ and examples/grader-showcase/ are the main demo examples
- Key recent features: skill_directories (#142), required_skills (#143), skill_invocation grader (#144), exit codes (#58), PR reporter (#59), skills CI compat (#60)
- microsoft/skills repo is reorganizing into plugin bundles — docs should reference both flat and nested layouts
- Tracking issue is #66 — always update checkboxes there when features land

📌 Team update (2026-02-15): After feature PRs merge (CLI, graders, YAML format, examples), you get routed doc update tasks. Update DEMO-GUIDE.md, GRADERS.md, TUTORIAL.md, examples/ READMEs, main README. — decided by Shayne Boyer
📌 Team update (2026-02-15): All developers use claude-opus-4.6. For code review, if developer isn't using Opus, reviewer uses it. — decided by Shayne Boyer
📌 Team update (2026-02-15): Don't take assigned work. Only pick up unassigned issues. — decided by Shayne Boyer
📌 Team update (2026-02-15): Multi-model execution is sequential (not parallel). Test failures non-fatal so all models complete. — decided by Linus
📌 Team update (2026-02-15): Microsoft/skills repo moving to plugin bundle structure. CI must support both flat and nested layouts. — decided by Shayne Boyer
📌 Team update (2026-02-15): Don't take assigned work — only pick up unassigned issues — decided by Shayne Boyer

📌 Team update (2026-02-17): Diff grader context_dir parameter must be documented in GRADERS.md — snapshot paths resolve relative to it. — decided by Linus

### 2026-02-18: SkillsBench Competitive Analysis
- Created `docs/research/` directory and published `skillsbench-competitive-analysis.md`
- Document covers: executive summary, product overview, feature comparison matrix, community insights (HN discussion), gap analysis (A/B impact P0, multi-agent P1, no Docker), positioning strategy
- Key findings: Skills as reasoning cache, post-failure generation > pre-task generation, quality > quantity, composability matters
- Strategic recommendations: P0 A/B impact (#194), P1 multi-agent decoupling (#195), defer Docker, P2 domain examples
- Positioning: Waza is developer tool ("does MY skill work?"), SkillsBench is research benchmark ("do skills help in general?") — complementary, not competitive
