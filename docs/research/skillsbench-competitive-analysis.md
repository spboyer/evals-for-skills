# SkillsBench Competitive Analysis

**Date:** February 2026  
**Audience:** Waza product team, leadership  
**Status:** Research complete, recommendations ready for prioritization

---

## Executive Summary

**SkillsBench** is an open-source benchmark framework that measures how effectively AI agents use skills—modular procedural knowledge—to perform specialized workflows. It has gained significant traction in the developer community (HackerNews discussion, active Discord) and represents the first systematic effort to answer the question: "Do skills actually help agents?"

Waza and SkillsBench are **complementary, not competitive**:

- **SkillsBench** answers the research question: "Do skills help agents *in general*?" (benchmarking, academic rigor)
- **Waza** answers the developer question: "Do *my* skills help *my* agents on *my* tasks?" (developer iteration, compliance, token management)

This document outlines the competitive landscape, feature gaps in both directions, and strategic recommendations to strengthen Waza's positioning.

---

## Product Overview: SkillsBench

### What It Is

SkillsBench is an open-source (Apache 2.0) benchmark that evaluates agent skill effectiveness at scale. Built on the Harbor framework, it measures task completion across multiple models and agents, with deterministic verification through pytest-based test suites.

**Repository:** https://github.com/benchflow-ai/skillsbench  
**Website:** https://www.skillsbench.ai/  
**Community:** HackerNews discussion (https://news.ycombinator.com/item?id=47040430), Discord syncs, weekly contributor meetings

### Architecture

#### Task Structure
Each SkillsBench task is a self-contained, deterministically verifiable workflow:

```
task-name/
├── instruction.md          # Human-authored task description
├── task.toml              # Metadata (category, difficulty, agents)
├── Dockerfile             # Container definition with env setup
├── solve.sh               # Oracle solution (ground truth execution)
├── tests/                 # pytest assertions for verification
├── environment/
│   └── skills/            # Modular procedural knowledge
│       ├── skill-1.md
│       └── skill-2.md
└── fixtures/              # Input data, reference files
```

#### Execution Model
1. **Containerized isolation:** Each task runs in a Docker container with frozen dependencies
2. **Agent-specific skill paths:** Skills copied to agent-specific directories (`.claude/skills`, `.codex/skills`, etc.)
3. **Deterministic verification:** pytest runs assertions on task outputs
4. **With/without comparison:** Built-in variant (`tasks-no-skills/`) for measuring A/B impact

#### Supported Agents
- Claude Code
- Codex (OpenAI)
- OpenCode
- Goose
- Factory (internal BenchFlow tool)

#### Task Diversity
50+ tasks spanning:
- **Science & Engineering:** 3D geometry, PID control, seismic analysis, crystallography
- **Cloud & Infrastructure:** Azure networking, security, system design
- **Data & ML:** Citation verification, data visualization, financial analysis
- **Specialized Domains:** PDF processing, video analysis, game mechanics, geospatial analysis

**Target Models:** Claude Opus 4.5, GPT-5.2, MiniMax M2.1, GLM-4.7

### Paper Findings (Key Numbers)

The SkillsBench paper ([PDF](https://www.skillsbench.ai/skillsbench.pdf)) tested 86 tasks across 11 domains under three conditions: no skills, curated (human-authored) skills, and self-generated skills. Key results:

| Finding | Data |
|---------|------|
| **Curated skills average improvement** | +16.2 percentage points pass rate |
| **Healthcare domain improvement** | +51.9pp (highest gain) |
| **Software engineering improvement** | +4.5pp (lowest gain) |
| **Tasks where skills hurt** | 16 of 84 tasks scored worse with skills |
| **Self-generated skills** | Zero average benefit, sometimes negative |
| **Optimal skill pack size** | 2–3 focused modules beat exhaustive docs |
| **Smaller model + skills** | ≈ larger model without skills |

**Critical implication for Waza:** Our primary audience writes skills for software engineering agents — the domain where SkillsBench measured the **least** skill impact (+4.5pp). This makes our compliance scoring and quality tooling even more important: in a domain with razor-thin margins, skill quality matters enormously. A bad skill in software engineering is more likely to hurt than help.

---

## Feature Comparison Matrix

| Dimension | SkillsBench | Waza | Gap? |
|-----------|-------------|------|------|
| **Primary Use Case** | Benchmark research: "Do skills help agents?" | Developer tool: "Create & evaluate my skills" | Different focus—complementary |
| **Target Audience** | Researchers, benchmark contributors | Skill authors (microsoft/skills ecosystem) | Different markets |
| **Execution Model** | Docker containers (isolated) | Temp workspaces (lightweight) | Trade-off: isolation vs. speed |
| **Verification** | pytest + oracle scripts | 8 grader types (regex, code, behavior, etc.) | Waza: richer grading variety |
| **Skill Format** | Markdown + task.toml metadata | SKILL.md (pure markdown) | Waza: simpler, lighter |
| **Task Scope** | Broad domains (science, finance, infra) | Coding-focused (Copilot agent tasks) | SkillsBench: broader examples |
| **Multi-Agent Support** | 5+ agents (Claude, Codex, etc.) | Copilot SDK only (extensible) | **Gap: Limited to Copilot** |
| **A/B Skill Impact** | ✅ Built-in (with/without skills) | ❌ Not supported | **Gap: No impact measurement** |
| **Compliance Scoring** | ❌ No | ✅ Yes (Sensei engine) | Waza advantage |
| **Token Management** | ❌ No | ✅ Yes (E4 roadmap) | Waza advantage |
| **Iterative Improvement** | ❌ No structured loop | ✅ Yes (`waza dev`) | Waza advantage |
| **Skill Generation** | ❌ Manual only | ✅ `waza new` scaffold | Waza advantage |
| **CI/CD Integration** | Planned | Planned (#156) | Parity |
| **Community & Governance** | Open-source (GitHub, Discord) | Internal + microsoft/skills | SkillsBench: larger OSS community |
| **Installation** | Requires Docker, uv, Harbor CLI | Single Go binary | Waza: lighter setup |

---

## Community Insights: What HN Learned Us

The HackerNews discussion (https://news.ycombinator.com/item?id=47040430) surfaced critical insights about skill effectiveness that validate Waza's direction:

### Key Finding 1: Skills as Reasoning Cache
**The most upvoted insight:** Skills aren't about teaching agents new facts—they're about **encoding hard-won reasoning into reusable, cheaper-to-invoke knowledge.**

> "Skills are a cache for LLM reasoning. Over time, they let you route simpler tasks to cheaper models because the hard knowledge is pre-encoded."

**Implication for Waza:** Our compliance scoring (Sensei engine) aligns perfectly with this. Skills should be high-quality, well-structured, documented reasoning—not just arbitrary instructions. `waza dev` helps authors iterate toward that standard.

### Key Finding 2: Post-Failure Skills > Pre-Task Skills
**Practitioner consensus:** Self-generated skills before a task provides zero benefit. The knowledge is already in the model's probability space.

> "After several failures, then a success, I had the agent create the skill. Next run it is successful first run."

**Implication for Waza:** This suggests a future roadmap feature: **post-evaluation skill generation**. After a `waza run` fails, Waza could automatically generate or refine skills from the failure transcript. We should design for this now.

### Key Finding 3: Skill Quality > Quantity
**Common warning:** Skills with factual errors or unclear instructions confuse agents more than no skill at all.

> "A bad skill teaching incorrect procedure is worse than the agent reasoning through it from first principles."

**Implication for Waza:** Sensei's compliance checks directly address this risk. High-quality skill metadata (clear triggers, anti-triggers, routing clarity) is the difference between helpful and harmful.

### Key Finding 4: Skills Should Contain Only Non-Training Data
**Clarifying principle:** Document only:
- (a) Learned-through-experience information
- (b) Context-specific information
- (c) Alignment information for future sessions

Not: Facts the model already knows from training.

**Implication for Waza:** Our skill templates and guidelines should emphasize this distinction. `waza dev` scoring could even detect "obvious training-data facts" and flag them.

### Key Finding 5: Skill Composition Matters
**SkillsBench design principle:** Many tasks intentionally require 2+ skills composed together, targeting SOTA performance <50%.

**Implication for Waza:** We should design examples and documentation to show skill composition patterns. Future work: compose-specific graders (e.g., "did the agent invoke both skill-1 AND skill-2?").

---

## Gap Analysis: What We're Missing (and Why It Matters)

### Gap 1: A/B Skill Impact Measurement ⭐ CRITICAL

**The Problem:**
Waza can tell you:
- ✅ "Your skill YAML is compliant" (Sensei)
- ✅ "Your eval passed on these tasks" (Runner)
- ❌ **"Your skill made things better"** (MISSING)

SkillsBench has explicit `tasks/` and `tasks-no-skills/` variants to measure impact. Waza has no way to answer: "Does removing this skill improve or degrade agent performance?"

**Why It Matters:**
This is the **money question** for every skill author and reviewer. HN discussion repeatedly returns to: "Show me the delta." Right now, Waza is a compliance tool and an evaluator—not a **measurement tool**. That's a fundamental gap.

**Gap Closure:**
- Add `--baseline` flag to `waza run` that strips `skill_directories` before execution
- Run tasks twice (with skills, without), produce comparison report
- Output: delta in pass rate, composite score, per-task impact
- Effort: ~150 lines of Go (small)
- **Recommendation: P0 — close this immediately**

**Impact:**
Closing this gap moves Waza from "compliance checker" to "effectiveness validator." It positions us credibly in SkillsBench's core value prop while keeping our dev-tool advantages.

---

### Gap 2: Multi-Agent Executor Support

**The Problem:**
Waza supports Copilot SDK only. SkillsBench validates across Claude Code, Codex, OpenCode, Goose, and Factory.

**Why It Matters:**
The microsoft/skills ecosystem targets Copilot, so single-agent is fine for core use. But increasingly, organizations want cross-agent validation: "Does my skill work on Claude Code too?" This validates skill robustness.

**Current State:**
The `AgentEngine` interface is well-abstracted (`internal/execution/engine.go`). Adding new agents is an extension point, not a core change.

**Gap Closure Path:**

1. **Phase 1 (Now): Decouple ExecutionResponse from Copilot SDK types**
   - Current: `ExecutionResponse` imports `copilot.SessionEvent` directly
   - Needed: Abstraction layer so response is agent-agnostic
   - Effort: ~200 lines
   - No new agents yet, just infrastructure for them

2. **Phase 2 (P1): Add Claude Code engine**
   - Implement `AgentEngine` interface using Claude CLI
   - Effort: ~300 lines
   - Proves the extension model works

3. **Phase 3+ (P2): Add Codex, others**
   - Reuse patterns from Claude Code engine
   - Each: ~300-400 lines

**Recommendation: P1 for Phase 1 (decoupling), P2 for Phase 2+ (new engines)**

---

### Gap 3: Containerized Isolation (Docker)

**The Problem:**
SkillsBench runs tasks in Docker containers. Waza uses temp workspaces (copy files, run, delete).

**Why It Matters (for them):**
Reproducible benchmarking requires frozen environments. Docker guarantees that.

**Why It Doesn't Matter (for us):**
Waza's audience is **skill authors iterating fast**. Docker adds 30-60 seconds per run (image builds, startup). Temp workspaces are actually *better* for iteration speed. This is an intentional differentiator, not a gap.

**Recommendation: DON'T close this gap.** Speed wins for developers. If someone needs containerized benchmarks, SkillsBench exists. We should lean into our speed advantage.

---

### Gap 4: Domain Breadth in Examples

**The Problem:**
Waza examples are coding-focused (code-explainer, grader-showcase). SkillsBench has 50+ tasks spanning science, finance, geospatial analysis, energy systems.

**Why It Matters (somewhat):**
Broader examples position Waza as general-purpose, not just for coding. Good for marketing/discovery.

**Why It's Not Critical:**
Our target audience (microsoft/skills contributors) is mostly coding-focused. The `waza new` scaffold works for any domain—we just lack example coverage.

**Gap Closure:**
Create 10-15 example evals spanning non-coding domains (data analysis, Azure operations, document processing). Each: ~1-2 hours to author.

**Recommendation: P2.** Nice for README/positioning, doesn't change capability. Defer until P0/P1 gaps closed.

---

### Gap 5: Skill Effectiveness Metric

**The Problem:**
No single metric answering: "By what percentage did this skill improve agent performance?"

**Why It Matters:**
Skill impact is what everyone cares about. A single number in CI (e.g., "+15% improvement on this skill") is powerful.

**Why It's Not a Separate Gap:**
This is a **derived output** of Gap 1 (A/B testing). Once we can run with/without skills, the math is trivial: `(pass_rate_with - pass_rate_without) / pass_rate_without * 100%`

**Gap Closure:**
Bundle with Gap 1 implementation. Once we have the two runs, compute delta as standard output.

**Recommendation: P0, but as part of Gap 1 closure.**

---

## Gaps Where Waza Wins

### Compliance Scoring (Sensei Engine)
SkillsBench has no equivalent to Sensei. Our heuristic scoring engine validates:
- Description clarity and length
- Trigger phrase clarity (USE FOR, DO NOT USE FOR)
- Routing specificity (INVOKES, WORKFLOW SKILL)
- Token budget compliance
- Anti-trigger detection

This directly addresses the HN insight: "Bad skills hurt more than no skill." Sensei is preventative quality control.

### Iterative Improvement Loop (`waza dev`)
SkillsBench requires manual skill creation and GitHub PR review. Waza offers guided refinement:
- Automated compliance scoring
- LLM-powered improvement suggestions (#44)
- Incremental skill validation

This is a developer experience win.

### Token Management
SkillsBench has no token budget awareness. Waza includes:
- Per-skill token counting
- Soft/hard limits (`.token-limits.json`)
- Budget reports in `waza tokens compare`
- Cost-aware grading

As skills scale, this becomes critical.

### Skill Scaffolding
`waza new` generates SKILL.md templates from a few prompts. SkillsBench requires manual authoring.

---

## Key Insights & Strategic Takeaways

### 1. Skills Are Reasoning Caches, Not Training Data
This validates our entire compliance focus. Skills should encode *hard-won knowledge from experience*, not facts the model already knows. Our Sensei engine's quality checks align perfectly with this principle.

### 2. Feedback-Driven Skill Generation Works; Feed-Forward Doesn't
Pre-task skill generation (before failure) provides zero benefit. Post-failure skill generation (after learning what went wrong) is highly valuable. This suggests a future Waza feature: automatic skill refinement from evaluation failures.

### 3. Quality > Quantity
A single, high-quality, well-scoped skill beats ten generic instructions. Our compliance scoring directly incentivizes this trade-off. Marketing point: "Waza helps you build skills that actually help."

### 4. Composability Testing Is Valuable
SkillsBench intentionally designs tasks requiring 2+ skills. Waza should design examples and documentation showing skill composition patterns. Future work: graders that verify multi-skill orchestration.

### 5. Lightweight Isolation Wins for Development
Docker is right for academic benchmarks. Temp workspaces are right for iterative development. We should double down on speed as a positioning differentiator.

---

## Recommendations: What Waza Should Do

### P0: Measure Skill Impact (#194)

**What:** Add A/B skill impact measurement to `waza run`.

**Scope:**
- `--baseline` flag strips skill directories before execution
- Runs tasks twice: with skills, without
- Outputs comparison report + delta metrics
- Integration with `waza compare` infrastructure

**Why:** This is the single most impactful feature we're missing. Closes the gap on SkillsBench's core value prop while keeping our dev-tool advantages. Positions Waza as both a *developer tool* and an *effectiveness validator*.

**Effort:** Small (~150 lines). Touches `cmd/waza/cmd_run.go` and `internal/orchestration/runner.go`.

**Timeline:** 2-3 days.

**Success Metric:** `waza run --baseline` produces side-by-side pass-rate delta.

---

### P1: Multi-Agent Architecture Decoupling (#195)

**What:** Decouple `ExecutionResponse` from Copilot SDK types.

**Scope:**
- Create agent-agnostic response model
- Current `copilot.SessionEvent` → abstraction layer
- No new agents yet, just infrastructure

**Why:** Unblocks future multi-agent support. Current tight coupling is technical debt regardless. Decoupling now prevents refactoring pain later.

**Effort:** Small-Medium (~200 lines). Scoped to `internal/execution/`.

**Timeline:** 2-3 days.

**Success Metric:** `ExecutionResponse` has no Copilot SDK imports. Tests pass.

**Follow-up (P2):** Implement Claude Code engine once decoupling is done.

---

### P2 (Skip): Docker Containerization

**Decision:** Don't pursue Docker isolation.

**Rationale:** Speed wins for developers. Temp workspaces are intentionally faster than containerization. This is a *differentiator*, not a gap. SkillsBench is the right tool for reproducible benchmarks. We're the right tool for fast iteration.

---

### P2 (Defer): Domain Example Expansion

**Decision:** Content creation, not priority.

**Scope:** 10-15 non-coding examples (data analysis, Azure ops, document processing).

**Effort:** Small-Medium (~15-20 hours). No code changes needed.

**Why:** Good for marketing. Doesn't unlock new capability. Defer until P0/P1 done.

---

## Positioning: How to Talk About Waza vs. SkillsBench

### The Frame
> **SkillsBench** is the academic benchmark. It answers: "Do skills help agents in general?" with rigor and breadth.
> 
> **Waza** is the developer tool. It answers: "Do my skills help my agent on my tasks?" with speed and clarity.

### Key Message Points

**For Skill Authors:**
- "Waza helps you build skills that actually work—with proof. Run your skills, measure their impact, and iterate toward higher effectiveness."

**For Skip Quality:**
- "Sensei compliance scoring ensures your skills help (not hurt) agents. High-quality skills are the difference between +15% and -5% agent performance."

**For the Microsoft/Skills Ecosystem:**
- "Waza is built for Copilot skill authors. It validates, measures, and improves your skills before they reach production."

**In Relation to SkillsBench:**
- "SkillsBench is the industry benchmark for skill effectiveness research. Waza is the developer environment for building and validating skills in practice. They're complementary—use SkillsBench to understand the research, use Waza to build better skills."

---

## References & Links

**SkillsBench:**
- Repository: https://github.com/benchflow-ai/skillsbench
- Website: https://www.skillsbench.ai/
- HackerNews Discussion: https://news.ycombinator.com/item?id=47040430

**Waza Issues:**
- #194 — A/B Impact Measurement (Gap 1)
- #195 — Multi-Agent Decoupling (Gap 2, Phase 1)
- #196 — Recommended Follow-ups (post-failure skill generation, composability testing)

**Related Waza Roadmap:**
- E3 (#39) — Evaluation Framework (multi-model comparison, LLM recommendations)
- E4 (#51) — Token Management (cost tracking, budget enforcement)
- E6 (#156) — CI/CD Integration

---

## Conclusion

SkillsBench raised important questions about skill effectiveness at scale. Waza is positioned to answer those questions for individual skill authors and teams. By closing the A/B impact measurement gap (P0) and decoupling from Copilot SDK (P1), Waza becomes both a *developer tool* and an *effectiveness validator*—complementary to SkillsBench, not competitive.

The two products will coexist healthily:
- **Waza** for fast, iterative skill development with compliance and impact measurement
- **SkillsBench** for academic benchmarking across agents and domains

Teams building skills will use both: SkillsBench to understand the research, Waza to build better skills in practice.
