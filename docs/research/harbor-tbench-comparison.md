# Harbor & TBench vs Waza: Comparative Analysis

**Date:** January 2025  
**Author:** Saul, Documentation Lead  
**Last Updated:** 2025-01-15  
**Status:** Research & Strategic Planning

---

## Executive Summary

Waza and Harbor/TBench operate in the AI agent evaluation space but serve different purposes and audiences. TBench is an industry-standard leaderboard for terminal-based coding tasks (103 entries, top score 75.1%), while Harbor is the open-source infrastructure powering it. Waza is a specialized evaluation framework purpose-built for GitHub Copilot Skill validation—focusing on trigger accuracy, token budgets, behavior constraints, and multi-skill orchestration.

**Strategic insight:** Rather than compete with Harbor, Waza should complement it. The recommendation is to develop a Harbor adapter that allows Waza evals to run on Harbor infrastructure and contribute to the TBench leaderboard, positioning Waza as the skill-specific grading engine for general-purpose agent benchmarks.

**Key gaps today:**
- No containerized execution (temp dirs instead of Docker isolation)
- Multi-model execution blocked (#39)
- No benchmark registry for cross-repo eval sharing
- Limited agent plugin system (Copilot + Mock only)

---

## TBench 2.0: Industry Standard

### Overview

**TBench 2.0** (https://www.tbench.ai/leaderboard/terminal-bench/2.0) is the industry-standard leaderboard for evaluating AI coding agents on real-world terminal tasks.

- **103 benchmark entries** across all models and agents
- **Top score:** 75.1% (OpenAI Simple Codex + GPT-5.3-Codex)
- **Developers:** Laude Institute + Stanford University
- **Task format:** Docker-containerized terminal environments with automated verification scripts

### Task Design

TBench focuses on real-world terminal scenarios:

- **Project compilation** — C/C++, Python, Java projects
- **Server configuration** — Linux environment setup, package management
- **Security tasks** — Vulnerability discovery, remediation
- **Data science** — Data preprocessing, analysis pipelines
- **ML training** — Model fine-tuning, hyperparameter optimization

Each task:
1. Runs in a Docker container with a fresh sandbox
2. Provides context (readme, scripts, configuration files)
3. Expects terminal commands as output from the agent
4. Verifies success via automated test scripts (bash/Python)
5. Scores as **binary pass/fail** (✓ or ✗)

### Scoring Model

- **Per-task:** Binary pass/fail evaluation
- **Aggregate:** Percentage of tasks passed across the entire benchmark
- **Leaderboard:** Ranked by aggregate score, with model/agent breakdown

---

## Harbor: Open-Source Infrastructure

### Overview

**Harbor** (https://github.com/laude-institute/harbor) is the open-source evaluation framework maintained by Laude Institute that powers TBench 2.0.

**Philosophy:** Harbor is an extensible harness for benchmarking agents at scale with full trajectory capture for analysis and RL optimization.

### Architecture

#### Task Format

Tasks are directory-based:

```
tasks/task-name/
├── instruction.md         # Human-readable task description
├── task.toml              # Task metadata and configuration
├── tests/                 # Verification scripts (bash/Python)
│   ├── test_compilation.sh
│   └── verify_output.py
└── solution/              # Reference implementation (optional)
    └── solution.sh
```

#### Agent Abstraction

Harbor uses an **adapter pattern** for agents:

- **Claude Code** adapter
- **Codex CLI** adapter
- **OpenHands** adapter
- **Custom agents** via pluggable adapter interface

Agents communicate via a standardized interface (environment variables, stdin/stdout, action logs).

#### Execution Model

- **Containerized sandboxes** — Docker per task, isolated filesystem
- **Parallel execution** — Daytona or Modal cloud providers for scale
- **Local fallback** — Can run locally with Docker Compose

#### Trajectory Logging

Harbor captures full action trajectories:

```json
{
  "task_id": "compile_python_project",
  "model": "claude-3.5-sonnet",
  "agent": "Claude Code",
  "actions": [
    {"type": "read_file", "path": "setup.py", "timestamp": "..."},
    {"type": "run_command", "cmd": "pip install -r requirements.txt", "output": "...", "timestamp": "..."},
    {"type": "run_command", "cmd": "python -m pytest", "output": "...", "timestamp": "..."}
  ],
  "result": {"status": "pass", "verified_at": "..."}
}
```

This trajectory data enables:
- Post-hoc analysis and debugging
- RL training data collection
- Failure pattern identification
- Agent comparison beyond pass/fail

#### CLI Usage

```bash
# Run TBench 2.0 with Claude Opus on 5 tasks
harbor run -d terminal-bench@2.0 -a "Claude Code" -m "claude-opus" -k 5

# Run custom benchmark
harbor run -d ./my-tasks -a "OpenHands" -m "gpt-4o" --parallel 4
```

### Ecosystem

- **Task repositories** — Git-based, versioned
- **Registry** — Central task registry with git refs for reproducibility
- **Providers** — Daytona and Modal integrations for cloud-scale execution
- **Output formats** — JSON trajectories, leaderboard reports, CI annotations

---

## Waza: Copilot Skill Evaluation Framework

### Overview

Waza is a specialized evaluation framework designed for GitHub Copilot Skills. It focuses on:

- **Trigger accuracy** — Does the skill activate when appropriate?
- **Token budgets** — Does the skill comply with token constraints?
- **Behavior validation** — Does the skill behave as documented?
- **Multi-skill scenarios** — Can skills work together in orchestrated sequences?

### Unique Design

#### Grader Types (7 Types)

Waza provides rich, specialized validators beyond generic test scripts:

1. **Regex grader** — Pattern matching on output
2. **File grader** — File creation/modification verification
3. **Code grader** — Syntax/AST validation
4. **Behavior grader** — Multi-step interaction sequences
5. **Action sequence grader** — API call order validation
6. **Skill invocation grader** — Verify child skill execution
7. **Prompt grader** — Prompt compliance and quality

#### Task Orchestration

Unlike Harbor's flat task structure, Waza supports multi-skill scenarios:

```yaml
benchmark:
  skill_directories:
    - path: ./skills/auth
      required_skills:
        - auth-login
        - auth-token-refresh
  tasks:
    - name: multi_skill_flow
      steps:
        - invoke: auth-login
          input: ...
        - invoke: auth-token-refresh
          expect: skill_invocation(auth-token-refresh, token_valid=true)
```

#### Execution Model

- **Temp workspace per task** — Fresh copy of fixtures, isolated filesystem (no Docker)
- **Parallel with goroutines** — Local parallelism (not cloud-scale)
- **Session events** — Copilot SDK session recording for trajectory capture

#### CLI Usage

```bash
# Run with verbose output
waza run ./eval.yaml -v

# Run against fixtures directory
waza run ./eval.yaml --context-dir ./fixtures

# Save results
waza run ./eval.yaml -o results.json
```

### CI/CD Integration

- **PR reporter** — Comments results on pull requests
- **Exit codes** — Non-zero exit on eval failure for CI automation
- **GitHub Actions** — Integrated workflow for skill testing

---

## Gap Analysis: Waza vs Harbor/TBench

| Capability | Harbor/TBench | Waza Today | Status | Gap Description |
|-----------|--------------|-----------|--------|-----------------|
| **Containerized Execution** | Docker sandboxes per task | Temp workspace dirs (no isolation) | 🔴 Major Gap | Harbor uses full container isolation; Waza uses filesystem-based isolation only |
| **Task Format** | `instruction.md` + `task.toml` + `tests/` | `eval.yaml` + YAML tasks + fixtures | 🟡 Different Format | Both store metadata + verification logic, different syntax |
| **Verification Logic** | Arbitrary scripts (bash/Python) | 7 grader types (regex, file, code, etc.) | 🟡 Different Approach | Waza richer semantic graders; Harbor more flexible script-based |
| **Multi-Model Comparison** | Run same agent with different models | Single model per run (#39) | 🔴 Blocked | Waza can't easily compare Claude vs GPT for same task |
| **Agent Abstraction** | Pluggable adapters (Claude, Codex, OpenHands, custom) | 2 engines (CopilotEngine, MockEngine) | 🟡 Limited | Harbor is generic; Waza is Copilot-specific |
| **Cloud-Scale Parallelism** | Daytona, Modal, Kubernetes | Local goroutines | 🟡 Local Only | Harbor scales to 1000s of tasks; Waza scales to tens |
| **Trajectory Logging** | Full JSON trajectories with action logs | Session events + SkillInvocation tracking | 🟢 Comparable | Waza captures SkillInvocations; Harbor captures all actions |
| **Leaderboard/Reporting** | Public leaderboard with rankings | PR comments + exit codes | 🟡 No Leaderboard | Waza reports locally; Harbor feeds public leaderboards |
| **Skill-Specific Grading** | Generic pass/fail | 7 specialized graders for Copilot Skills | 🟢 Waza Advantage | Waza validates skill-specific constraints (tokens, triggers, behavior) |
| **Multi-Skill Orchestration** | Single agents only | skill_directories + required_skills + skill_invocation grader | 🟢 Waza Advantage | Waza can validate multi-skill workflows; Harbor is single-agent |
| **Benchmark Registry** | Git-based task registry with versioning | Eval YAML files in individual repos | 🔴 Missing | Harbor enables cross-repo task sharing; Waza is per-repo |
| **RL/Fine-Tuning Data Export** | Trajectories optimized for RL | Not designed for RL pipelines | 🟡 Not Priority | Harbor trajectory format enables RL; Waza doesn't focus on this |

### Legend

- 🟢 **Waza Advantage** — Waza is superior or unique in this dimension
- 🟡 **Different Approach** — Both work, different design choices
- 🔴 **Major Gap** — Waza lags significantly and blocks use cases

---

## Waza Advantages Over Harbor

### 1. Purpose-Built for Copilot Skills

Waza is designed specifically for validating GitHub Copilot Skills:

- **Trigger validation** — Verify skill activates on correct prompts
- **Token budgets** — Enforce token constraints per skill
- **Behavior constraints** — Validate skill follows documented behavior
- **Skill composition** — Multi-skill orchestration and dependency tracking

Harbor is agent-agnostic; it treats all agents the same. Waza's graders are specialized for the skill lifecycle.

### 2. Rich Semantic Graders

While Harbor uses generic test scripts, Waza provides 7 typed graders:

| Grader | Use Case |
|--------|----------|
| **Regex** | Output pattern matching (e.g., "error not in output") |
| **File** | File creation/modification (e.g., "file.txt created") |
| **Code** | Syntax/AST validation (e.g., "Python code is valid") |
| **Behavior** | Multi-step interaction sequences |
| **Action Sequence** | API call order validation (e.g., must call auth before query) |
| **Skill Invocation** | Child skill execution verification |
| **Prompt** | Prompt compliance and quality checks |

This enables precise, declarative validation without custom scripts.

### 3. Multi-Skill Orchestration

Unique to Waza:

```yaml
benchmark:
  skill_directories:
    - path: ./skills
      required_skills:
        - skill-a
        - skill-b
  tasks:
    - name: workflow
      steps:
        - invoke: skill-a
        - invoke: skill-b with output from skill-a
        - expect: skill_invocation(skill-b, success=true)
```

Harbor doesn't support skill composition; Waza's multi-skill grader enables testing skill interaction.

### 4. CI/CD Integration

Waza is built for developer workflows:

- **PR reporter** — Comment results directly on pull requests
- **Exit codes** — Fail CI pipelines on eval failure
- **GitHub Actions** — Native workflow integration for skill teams
- **Sensei dev loop** — Iterative skill improvement with automated compliance checking

Harbor is designed for research benchmarking; Waza is designed for developer iteration.

---

## Strategic Recommendation

### Don't Compete — Complement

**Thesis:** Waza and Harbor solve different problems for different audiences. Rather than try to match Harbor's breadth, Waza should specialize and integrate.

- **Waza = Skill-specific validator** — Focuses on trigger accuracy, token budgets, multi-skill scenarios
- **Harbor = General agent benchmarker** — Focuses on scaling, trajectory capture, leaderboard reporting

### Proposed Integration Path

#### Phase 1: Containerized Execution (P1)

Add optional Docker-per-task execution to Waza:

```bash
waza run ./eval.yaml --container docker    # Use Docker sandboxes
waza run ./eval.yaml --container none      # Fallback to temp dirs
```

This enables Waza evals to run in isolated environments like Harbor.

#### Phase 2: Multi-Model Execution (P0)

Implement #39 to support model comparison:

```yaml
benchmark:
  agents:
    - name: claude
      model: claude-opus
    - name: gpt
      model: gpt-4o
  tasks:
    - name: test_task
      inputs: ...
      expect: ...
```

Run the same task against multiple models and aggregate results.

#### Phase 3: Harbor Adapter (P2)

Build a Harbor adapter that translates Waza evals to Harbor task format:

```bash
# Waza eval → Harbor task format
waza to-harbor ./skill-eval.yaml --output ./harbor-task/

# Then run on Harbor infrastructure
harbor run -d ./harbor-task -a "Copilot" -m "claude-opus"
```

This allows Waza skill evals to contribute to cross-model benchmarks and appear on public leaderboards.

#### Phase 4: Benchmark Registry (P2)

Develop a registry for sharing evals across teams:

```yaml
# registry.json
{
  "evals": [
    {
      "id": "skill-auth-system",
      "repo": "github.com/org/skill-auth",
      "ref": "v1.0.0",
      "eval_path": "evals/auth-system.yaml"
    }
  ]
}
```

Enable teams to discover and import evals from other skill repos.

---

## Proposed Work Items

### P0 — Critical Path

| ID | Title | Description | Impact |
|----|-------|-------------|--------|
| #39 | Multi-model execution | Support running same eval against multiple models (Claude, GPT, etc.) | Unblocks cross-model comparison; needed for leaderboard contribution |

### P1 — High Priority

| ID | Title | Description | Related |
|----|-------|-------------|---------|
| TBD | Container isolation for task execution | Docker-per-task with fallback to temp dirs | Matches Harbor's execution model; improves security/isolation |

### P2 — Medium Priority

| ID | Title | Description | Related |
|----|-------|-------------|---------|
| TBD | Benchmark registry | Central registry for sharing evals across repos | Enables eval reuse; foundation for cross-repo analysis |
| TBD | Harbor adapter | Translate Waza evals to Harbor task format | Allows Waza evals to run on Harbor infrastructure + TBench leaderboard |
| TBD | Agent plugin system | Pluggable engine abstractions beyond Copilot/Mock | Opens Waza to other agents (Claude CLI, OpenHands, etc.) |
| TBD | Trajectory export in Harbor format | Save full trajectories as Harbor JSON | Enables cross-tool analysis; supports RL data collection |

### Tracking

See **GitHub Issue #66** (Waza Platform Roadmap) for the full epic breakdown and prioritization.

---

## Related GitHub Issues

- **#39** — Multi-model execution (blocked, P0)
- **#66** — Waza Platform Roadmap (tracking, contains proposed work items)
- **#138** — Container isolation research (pending spec)

---

## Conclusion

Waza and Harbor represent different design philosophies:

- **Harbor:** General-purpose agent benchmarking at scale with public leaderboards
- **Waza:** Specialized skill validation with rich semantic graders and multi-skill orchestration

Rather than compete, Waza should **complement Harbor** by:

1. Adopting Harbor-compatible execution models (containers, multi-model, trajectory format)
2. Building adapters to contribute Waza evals to Harbor infrastructure
3. Specializing in skill-specific validation that Harbor doesn't handle (triggers, tokens, skill composition)

This positions Waza as the **skill-specific grading engine** for the broader agent evaluation ecosystem, enabling GitHub Copilot Skills to participate in industry benchmarks while maintaining focus on skill developer needs.

---

## References

- **TBench 2.0 Leaderboard:** https://www.tbench.ai/leaderboard/terminal-bench/2.0
- **Harbor GitHub:** https://github.com/laude-institute/harbor
- **TBench Paper:** [Laude Institute + Stanford research]
- **Waza Roadmap:** #66
- **Waza Docs:** `docs/` directory in this repo
