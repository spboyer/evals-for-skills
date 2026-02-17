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

### 2026-02-15: GETTING-STARTED.md created (Issue #168)
- Created comprehensive getting started guide at docs/GETTING-STARTED.md
- Covers: Quick start (init → new → edit → run → check → push), project layouts (single-skill vs multi-skill), workspace detection, commands reference, evaluation patterns (per-skill and orchestration), graders, common workflows, troubleshooting
- Updated README.md Quick Start section with link to guide
- Added `waza new <name>` command to commands table (was missing)
- Added Getting Started to Documentation section in README
- Proactive doc work before CLI retrofit lands — ready for minor adjustments after implementation
- Branch: squad/168-getting-started-docs, Commit: e3605f5
