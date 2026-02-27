# History — Livingston

## Project Context
- **Project:** waza — CLI tool for evaluating Agent Skills
- **Stack:** Go (primary), React 19 + Tailwind CSS v4 (web UI)
- **User:** Shayne Boyer (spboyer)
- **Repo:** spboyer/waza
- **Universe:** The Usual Suspects

## Key Learnings

### Documentation Structure
- **Main files:** README.md, docs/, waza-go/README.md
- **Key sections:** Usage, examples, CLI flags, agent integration
- **API docs:** BenchmarkSpec, TestCase, EvaluationOutcome, Validator interface
- **Update requirement:** Must stay in sync with code changes

### Waza Concepts
- Evaluation specs (YAML format)
- Task definitions with fixtures
- Validator registry (extensible grading)
- Agent execution (Go engine, fixture isolation)
- Results and scoring

### CI/CD
- Workflows defined in .github/workflows/
- Branch protection enforces docs stay current
- Changelog tracking for releases

### GUIDE.md Patterns (Issue #253)
- **Structure:** Overview → Installation → Quick Start → Command Reference → Advanced → Dashboard → Troubleshooting
- **Key principle:** All examples use Go CLI only (no Python, no venv, no legacy references)
- **Installation methods:** Binary (recommended), from source, azd extension
- **Quick Start:** 5-step workflow — init → new → define → run → serve
- **Command reference:** Detailed flags and examples for init, new, run, check, serve
- **Advanced sections:** Caching, filtering, parallel execution, multi-model comparison, CI/CD, session logging
- **Dashboard:** Pages are home/dashboard, run details, compare, trends, live view (from web/src/App.tsx routing)
- **Troubleshooting:** Port conflicts, missing results, fixture paths, validation issues
- **File paths:** docs/GUIDE.md is the canonical user guide; links to GETTING-STARTED.md for step-by-step; references examples/ for runnable code

### CLI Command Implementation Details
- **init:** Creates skills/, evals/, .github/workflows/eval.yml, defaults to claude-sonnet-4.6
- **new:** Two modes (project vs standalone), interactive wizard for TTY, non-interactive for CI/CD
- **run:** Accepts eval.yaml OR skill-name OR auto-detect; supports filtering (--task, --tags), parallel (--workers), multi-model (--model), caching (--cache, --cache-dir), output formats (--format), session logging
- **check:** Validates compliance (Low/Medium/Medium-High/High), token count, eval presence; supports auto-detect
- **serve:** HTTP dashboard (default port 3000), can also run JSON-RPC TCP (--tcp :9000) or stdio
- **Exit codes:** 0 = success, 1 = test failed, 2 = configuration/runtime error

### Web Dashboard Routing
- Pages in App.tsx: home (Dashboard), run (RunDetail), compare (CompareView), trends (TrendsPage), live (LiveView)
- Features: live updates, search, filtering by status/tags/date, export, dark mode

📌 Team update (2026-02-19): Documentation maintenance gates established (Saul reviews PRs for doc impact) — decided by Saul (#256)


## 📌 Team update (2026-02-20): Model policy overhaul

All code roles now use `claude-opus-4.6`. Docs/Scribe/diversity use `gemini-3-pro-preview`. Heavy code gen uses `gpt-5.2-codex`. Decided by Scott Boyer. See decisions.md for full details.

## Learnings — Azure Storage Feature Documentation (2026-02-26)

### What was documented
- **Feature:** Azure Blob Storage integration for eval results (auto-upload, `waza results list`, `waza results compare`)
- **Scope:** 
  - README.md — Added Cloud Storage section with configuration example and workflow
  - site/src/content/docs/reference/cli.mdx — Added `waza results list` and `waza results compare` command reference
  - site/src/content/docs/guides/azure-storage.mdx — Created comprehensive guide with prerequisites, setup, troubleshooting
  - site/src/content/docs/reference/schema.mdx — Added `storage:` section to .waza.yaml schema

### Key documentation patterns observed
1. **Configuration-first approach** — Always show actual YAML examples readers can copy-paste
2. **Audience context** — Never assume reader knows Azure Storage or terminology (e.g., "blob container" explained briefly)
3. **Troubleshooting as essential** — Every guide ends with ❌ error patterns + fixes
4. **Prerequisites first** — Setup guides list required tools/accounts before any steps
5. **Mermaid for diagrams** — No ASCII art; Astro/Starlight renders Mermaid inline
6. **Frontmatter convention** — All guides use `---title/description---` header format
7. **Cross-linking** — Guides reference CLI docs and schema, schema references guides
8. **Code examples are actionable** — Every bash/yaml example can run exactly as shown

### Site build process
- Astro/Starlight auto-discovers guides in `site/src/content/docs/guides/*.mdx`
- Build succeeds if all markdown frontmatter is valid YAML
- Search index rebuilt automatically (Pagefind)
- Test with `npm run build` — takes ~250ms
- Verify new guide appears in static route list during build

### Future: When adding new cloud storage providers
- Add new `provider:` option to storage section in schema.mdx
- Create new guide in guides/ (e.g., s3-storage.mdx) with identical structure: Why / Prerequisites / Configuration / Getting Started / Using Results / Troubleshooting / Advanced
- Update README.md "Cloud Storage" section with provider matrix
- CLI commands (`waza results list/compare`) remain provider-agnostic — auth is handled by provider-specific adapters
