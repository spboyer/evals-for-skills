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
