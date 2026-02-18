# Positioning Strategy: Waza vs. SkillsBench

**Date:** February 2026  
**Audience:** Waza team, skill authors, product leads  
**Purpose:** Ready-to-use language for positioning Waza in relation to SkillsBench across different contexts

---

## Core Positioning Statement

Waza and SkillsBench answer different questions for different audiences.

**SkillsBench** is the industry-standard research benchmark that answers: "Do skills help AI agents *in general*, and by how much?" It runs rigorous, containerized evaluations across multiple agents and domains to build scientific evidence about skill effectiveness.

**Waza** is the developer tool that answers: "Do *my* skills help *my* agent on *my* tasks, right now?" It prioritizes iteration speed, compliance quality, and practical impact measurement so skill authors can build better skills before they ship to production.

They are **complementary, not competitive**. Teams building skills use both: SkillsBench to understand the research, Waza to validate and improve in practice.

---

## Elevator Pitch Variants

### For Developers: "I use SkillsBench, why do I need Waza?"

> SkillsBench is the research tool—it measures skill effectiveness at scale across agents. Waza is the dev tool—it helps you build skills that actually work, *right now*, without waiting for Docker builds or full benchmarks. Run a task in seconds, see if your skill helped, iterate. When you're ready, publish to SkillsBench for broader validation. Two different tools for two different jobs.

---

### For Managers: "How does this compare to SkillsBench?"

> SkillsBench is a 50+ task academic benchmark that validates skill effectiveness across agents and domains—published research. Waza is our internal dev environment for building skills faster, catching quality problems early, and measuring real impact before production. We use SkillsBench results to inform our skill design principles; we use Waza to enforce those principles in practice. Think of SkillsBench as the reference library and Waza as the workshop.

---

### For Contributors: "I'm writing skills for microsoft/skills — should I use Waza or SkillsBench?"

> **Short answer: Waza for development, SkillsBench for validation.**  
> Use Waza while you're building—it's fast, it scores skill quality, and it tells you whether your skill actually helps your agent. Use SkillsBench if you want to test your skill against multiple agents or publish research on skill effectiveness. Most skills ship with Waza evals because they're practical and tight; some get published to SkillsBench for the community.

---

## Key Differentiators

When asked "What makes Waza different?", lead with these:

### 1. Developer Iteration Speed

**The differentiator:** Waza runs evaluations in seconds; SkillsBench takes minutes to hours.

- **Waza:** Temp workspaces, no Docker, single binary → run tasks in seconds
- **SkillsBench:** Containerized isolation, frozen environments → reproducible research

**Why it matters:** Skill authors iterate on dozens of eval runs per day. Waza is built for that workflow. SkillsBench is built for publishing results.

**Language to use:**
> "Waza keeps you in flow. No Docker build times, no container overhead—just define a task, run it, see the result, iterate. That's how modern tool development works."

---

### 2. Compliance Scoring (Sensei Engine)

**The differentiator:** Waza has automated quality checking; SkillsBench doesn't.

Waza's Sensei engine validates:
- Description clarity and length
- Trigger phrase clarity ("USE FOR", "DO NOT USE FOR")
- Routing specificity ("INVOKES", "WORKFLOW SKILL")
- Token budget compliance
- Anti-trigger detection

**Why it matters:** SkillsBench's own research found that bad skills hurt more than they help. Low-quality triggers or routing confuse agents more than no skill at all. Waza catches those problems before they ship.

**Language to use:**
> "Sensei compliance scoring is preventative quality control. In software engineering (the domain where skills are narrowest margin), skill quality is everything. We catch routing bugs, vague triggers, and missing anti-triggers before they hit production."

**Support with research:**
> "SkillsBench measured that software engineering skills only improve agent performance by +4.5 percentage points—the lowest of any domain. That thin margin means *every* skill must be high-quality. Sensei is built for that constraint."

---

### 3. Token Management

**The differentiator:** Waza tracks and enforces token budgets; SkillsBench doesn't.

Waza provides:
- Per-skill token counting
- Soft/hard budget enforcement (`.token-limits.json`)
- Budget comparison reports (`waza tokens compare`)
- Cost-aware grading

**Why it matters:** As agents scale and skills proliferate, token budgets explode. Waza helps authors stay lean and avoid the "skill tax" that degrades performance.

**Language to use:**
> "Token management isn't optional at scale. Waza lets you set budgets, track against them, and optimize ruthlessly. A single-page skill that costs 300 tokens is better than a 10-page skill with the same information."

---

### 4. Copilot SDK Native (First-Class Copilot Integration)

**The differentiator:** Waza is built specifically for Copilot; SkillsBench is agent-agnostic.

- **Waza:** Copilot SDK native, Skill.md format, SKILL.md compliance, skill_directories routing
- **SkillsBench:** Claude Code, Codex, OpenCode, Goose, Factory—no Copilot

**Why it matters:** The microsoft/skills ecosystem is Copilot-focused. Waza speaks that language natively.

**Language to use:**
> "We're built for microsoft/skills, not trying to be everything. Skill authors get first-class Copilot tooling: SKILL.md validation, skill composition scoring, Copilot-native execution. If you're writing for Copilot, you're home."

---

## What We DON'T Say (Anti-Patterns)

### ❌ "Waza is better than SkillsBench"

Never. They answer different questions for different audiences. Better → alienates the research community and makes us look insecure.

**Instead say:**
> "Waza is the tool for *development*; SkillsBench is the tool for *research*. They're complementary."

---

### ❌ "We're a benchmark"

We're not. We have evals, but we don't claim to be a benchmark—SkillsBench owns that term.

**Instead say:**
> "Waza is a developer tool for validating skills. SkillsBench is the research benchmark for skill effectiveness."

---

### ❌ "Docker isolation is outdated / unnecessary"

Containerization serves a real purpose for academic benchmarks. Dismissing it makes us look ignorant of our own constraints.

**Instead say:**
> "For development iteration, Docker adds overhead that slows feedback loops. For reproducible benchmarks, Docker is the right choice. We optimize for different needs."

---

### ❌ "We replaced/deprecated/superseded SkillsBench"

Absolutely not. SkillsBench is ongoing research. We're a different tool in a different context.

**Instead say:**
> "SkillsBench is the ongoing research benchmark. Waza is the developer environment. Both are active, both matter."

---

## Recommended README Language

**Suggested addition to main README.md** (under "Quick Start" or "What is Waza?"):

---

### **Waza vs. SkillsBench: Complementary Tools**

[SkillsBench](https://www.skillsbench.ai/) is the industry-standard benchmark for measuring skill effectiveness across agents and domains—published research with 50+ tasks and peer-reviewed rigor.

**Waza** is the developer environment for building skills *before* they reach SkillsBench. It answers: "Does this skill actually help my agent on my task?" with speed and clarity.

| Aspect | Waza | SkillsBench |
|--------|------|-------------|
| **Primary Use Case** | Build & validate skills locally | Benchmark research & publish results |
| **Execution Speed** | Seconds (no Docker) | Minutes to hours (containerized) |
| **Compliance Scoring** | ✅ Sensei engine | ❌ No |
| **Token Management** | ✅ Budget tracking | ❌ No |
| **Multi-Agent Support** | Copilot SDK (extensible) | 5+ agents (Claude, Codex, etc.) |
| **Target Audience** | Skill authors | Researchers, academia |
| **Publication** | Internal dev environment | Academic papers, public benchmarks |

**Recommended workflow:**
1. Use Waza to develop skills locally—validate, measure impact, improve
2. Use `waza check` to ensure compliance before shipping
3. Publish to [microsoft/skills](https://github.com/microsoft/skills) with Waza evals
4. Optionally submit to SkillsBench for broader research validation

---

## The SkillsBench Insight We Should Amplify

SkillsBench's key finding:

> **Software engineering skills improve agent pass rate by +4.5 percentage points — the lowest of any domain.**

**This is Waza's strongest positioning lever.** Frame it this way:

---

### **In Software Engineering, Skill Quality is Everything**

SkillsBench measured skill effectiveness across 11 domains. Healthcare skills improved agents by +51.9pp. Cloud infrastructure by +18.2pp. Finance by +15.1pp.

**Software engineering: +4.5pp.**

That razor-thin margin means:
- ✅ A high-quality skill helps
- ❌ A medium skill doesn't move the needle
- ❌ A bad skill hurts more than it helps

This is *why* Waza matters. In domains with thicker margins, roughly-written skills might still work. In software engineering, they don't.

**Language to use:**

> "Software engineering is where skills are most fragile. A vague trigger, conflicting anti-triggers, or poor routing can erase the skill's value entirely. Sensei compliance scoring is built for that constraint. When margins are thin, quality is non-negotiable."

---

## How to Position in Different Contexts

### In Presentations / Talks

**Opening:**
> "Waza is the developer tool for building Copilot skills faster. SkillsBench is the research benchmark. We complement each other."

**When showing compliance scoring:**
> "SkillsBench found that software engineering skills have the thinnest margin for error (+4.5pp improvement). This is why compliance checking matters—we catch routing bugs, trigger conflicts, and scope creep before they ship."

**When discussing speed:**
> "Iteration speed wins for developers. Waza runs evals in seconds because you're running dozens per day. SkillsBench runs evals in minutes because they're building reference benchmarks."

### In Documentation

Always mention SkillsBench in the README's comparison table and in docs/GETTING-STARTED.md when discussing publishing skills.

Link to the SkillsBench paper when discussing "Why does skill quality matter?"

Reference the +4.5pp finding whenever explaining Sensei compliance.

### In Community Conversations

If someone asks: "Why not just use SkillsBench?"

> "SkillsBench is fantastic for research. For development, the Docker overhead kills iteration speed—we optimized for tight feedback loops instead. Most skills are built with Waza, some get published to SkillsBench."

If someone dismisses Waza as "just another eval tool":

> "We're not a benchmark, we're a dev environment. The value is in compliance scoring, token tracking, and fast iteration. We help skill authors build *better* skills before they ship."

---

## Communication Checklist

When positioning Waza:

- [ ] **Acknowledge SkillsBench.** Never pretend it doesn't exist.
- [ ] **Frame as complementary.** "Developer tool" ≠ "research benchmark"
- [ ] **Lead with speed.** Seconds vs. minutes. That's the punchline for developers.
- [ ] **Amplify compliance.** Sensei is our unique differentiator—lean on it.
- [ ] **Use the +4.5pp finding.** It's research-backed, it explains *why* we care about quality.
- [ ] **Don't attack Docker.** It's right for benchmarks, wrong for iteration loops.
- [ ] **Own Copilot.** We're built for microsoft/skills—that's a strength, not a limitation.

---

## Future Positioning Opportunities

As Waza grows, these positioning opportunities emerge:

1. **A/B Impact Measurement (#194)** — Once we can measure "does my skill help?", that becomes a major differentiator vs. SkillsBench's passive measurement model.

2. **Post-Failure Skill Generation** — HN discussion validated that skills generated *after* failures help. Waza could automate that—turn failures into improvements. Neither SkillsBench nor any other tool does this.

3. **Copilot Skill Composition** — Tasks in SkillsBench intentionally require 2+ skills composed. Waza could offer composition-aware graders and guided composition patterns. That's unique value for the microsoft/skills ecosystem.

4. **Token Normalization Across Agents** — As we add more agents, token budgets become portable—a 300-token Copilot skill could be validated on Claude Code. This is unique to multi-agent Waza; SkillsBench doesn't have token awareness.

---

## Summary

**Waza's positioning is clear and defensible:**

- **Speed:** Developer iteration (seconds vs. minutes)
- **Quality:** Sensei compliance scoring (SkillsBench finds skills matter in thin margins)
- **Practicality:** Token management, local validation, Copilot-native
- **Respect:** SkillsBench is research; we're tools. Different problems, different solutions.

**The key insight:** SkillsBench's +4.5pp finding in software engineering is *why* Waza's quality tooling matters. Use it to anchor every positioning conversation.

Team members should feel empowered to use the language in this doc directly in presentations, READMEs, and conversations. Consistency builds credibility.
