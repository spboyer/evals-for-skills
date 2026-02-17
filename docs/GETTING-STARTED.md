# Getting Started with Waza

Waza is a CLI tool for evaluating AI agent skills. This guide walks you through creating and testing your first skill evaluation suite.

## Quick Start

### 1. Initialize a new project

```bash
waza init my-skills-project
cd my-skills-project
```

This creates the project structure:

```
my-skills-project/
├── skills/                    # Your skill definitions
├── evals/                     # Your evaluation suites  
├── .github/workflows/eval.yml # CI pipeline (runs on PR)
├── .gitignore
└── README.md
```

`waza init` is idempotent — run it again to verify or repair your project structure.

### 2. Create your first skill

```bash
waza new my-skill
```

This creates:

```
skills/my-skill/SKILL.md              # Skill definition
evals/my-skill/eval.yaml              # Evaluation spec
evals/my-skill/tasks/basic-usage.yaml  # Happy path test
evals/my-skill/tasks/edge-case.yaml    # Edge case test
evals/my-skill/tasks/should-not-trigger.yaml  # Anti-trigger test
evals/my-skill/fixtures/sample.py      # Example fixture
```

### 3. Edit your skill

Open `skills/my-skill/SKILL.md` and define:

- **What the skill does** (description)
- **When it should activate** (USE FOR section)
- **When it should NOT activate** (DO NOT USE FOR section)

Example:

```markdown
# my-skill

A skill that analyzes Python code and suggests optimizations.

## USE FOR
- Code performance analysis
- Identifying potential bottlenecks
- Suggesting refactoring opportunities

## DO NOT USE FOR
- Syntax error checking (use a linter)
- Code formatting (use Black or similar)
- Architectural design decisions
```

### 4. Run evaluations

```bash
waza run                    # run all evals
waza run my-skill           # run one skill's evals
```

Watch the output to see which tests pass and fail.

### 5. Check compliance

```bash
waza check                  # check all skills
waza dev my-skill           # improve with real-time scoring
waza tokens check           # verify token budgets
```

### 6. Push to trigger CI

```bash
git add . && git commit -m "feat: add my-skill"
git push
```

CI automatically runs evaluations on pull requests.

## Project Layouts

### Single-Skill Project

For standalone skills (one skill per repo):

```
my-skill/
├── SKILL.md
├── evals/
│   ├── eval.yaml
│   ├── tasks/
│   │   ├── basic-usage.yaml
│   │   ├── edge-case.yaml
│   │   └── should-not-trigger.yaml
│   └── fixtures/
│       └── sample.py
├── .github/workflows/eval.yml
└── README.md
```

### Multi-Skill Project

For skill collections (like microsoft/skills):

```
my-skills-repo/
├── skills/
│   ├── azure-prepare/SKILL.md
│   ├── azure-deploy/SKILL.md
│   └── azure-monitor/SKILL.md
├── evals/
│   ├── azure-prepare/
│   │   ├── eval.yaml
│   │   ├── tasks/
│   │   └── fixtures/
│   ├── azure-deploy/
│   │   ├── eval.yaml
│   │   ├── tasks/
│   │   └── fixtures/
│   └── azure-monitor/
│       ├── eval.yaml
│       ├── tasks/
│       └── fixtures/
├── .github/workflows/eval.yml
└── README.md
```

## Workspace Detection

Waza automatically detects your project layout. All commands adapt:

| Context | No args | Skill name arg |
|---------|---------|----------------|
| Single-skill dir | Runs against this skill | N/A |
| Multi-skill repo | Runs against ALL skills | Runs against named skill |

This means you don't need to think about structure — just run `waza run` and it works!

## Understanding Skills

A skill is a reusable agent behavior that can be:

- **Invoked** — Called by other skills or the user
- **Scoped** — Clear trigger conditions (USE FOR / DO NOT USE FOR)
- **Evaluated** — Tested with task scenarios

Skills are described in `SKILL.md` files with frontmatter metadata:

```markdown
# skill-name

Brief description (< 150 characters recommended).

## USE FOR
- When to use this skill
- Specific scenarios where it activates

## DO NOT USE FOR
- When NOT to use this skill
- Common misconceptions
- Related tools that do similar things

[Optional: detailed explanation, examples, etc.]
```

## Understanding Evaluations

An evaluation suite tests a skill with multiple scenarios:

- **eval.yaml** — Configuration (model, graders, timeout)
- **tasks/** — Individual test cases (YAML files)
- **fixtures/** — Context files (code samples, test data)

### Task Format

Each task is a YAML file that defines:

```yaml
name: basic-usage
description: Test the skill with a simple, happy path scenario
input: |
  Write a Python function that adds two numbers.
expected_outputs:
  - pattern: "def add"
  - pattern: "return"
validators:
  - type: regex
    config:
      must_match: ["def add"]
```

### Evaluation Spec Format

The `eval.yaml` file configures how tests run:

```yaml
name: my-skill-eval
skill: my-skill
version: "1.0"

config:
  trials_per_task: 3
  timeout_seconds: 300
  executor: copilot-sdk          # or mock for testing
  model: claude-sonnet-4-20250514

graders:
  - type: code
    name: assertions
    config:
      language: python
      code: |
        assert "def " in result
        assert "return" in result
  
  - type: regex
    name: pattern_check
    config:
      must_match: ["function.*add"]

tasks:
  - "tasks/*.yaml"
```

## Commands Reference

### Project Initialization

| Command | Purpose |
|---------|---------|
| `waza init [dir]` | Initialize project structure (idempotent) |
| `waza new <name>` | Create a new skill with eval suite |

### Running Evaluations

| Command | Purpose |
|---------|---------|
| `waza run [name]` | Run evaluations (all skills or named skill) |
| `waza run <eval.yaml>` | Run a specific eval spec file |
| `waza cache clear` | Clear cached evaluation results |

### Skill Compliance

| Command | Purpose |
|---------|---------|
| `waza check [name]` | Check skill readiness for submission |
| `waza dev [name]` | Iterative improvement with real-time scoring |
| `waza generate <SKILL.md>` | Generate eval suite from skill definition |

### Token Management

| Command | Purpose |
|---------|---------|
| `waza tokens count [paths...]` | Count tokens in markdown files |
| `waza tokens suggest [paths...]` | Suggest token optimizations |
| `waza tokens compare` | Compare token usage across git refs |

### Comparison & Analysis

| Command | Purpose |
|---------|---------|
| `waza compare <file1> <file2> ...` | Compare results across models |

## Evaluation Patterns

### Per-Skill Evals (default)

Each skill has its own eval suite testing it in isolation:

```
skills/
├── azure-prepare/SKILL.md
└── evals/
    └── azure-prepare/eval.yaml
skills/
├── azure-deploy/SKILL.md
└── evals/
    └── azure-deploy/eval.yaml
```

Use `waza run` to test each skill independently.

### Orchestration Evals

Test how multiple skills work together using `skill_directories` and `required_skills`:

```yaml
config:
  skill_directories:
    - ./skills/azure-prepare
    - ./skills/azure-deploy
  required_skills:
    - azure-prepare
    - azure-deploy

graders:
  - type: skill_invocation
    name: workflow_check
    config:
      mode: in_order
      required_skills:
        - azure-prepare
        - azure-deploy

tasks:
  - "tasks/*.yaml"
```

This tests that:
1. The first task invokes `azure-prepare`
2. The second task invokes `azure-deploy` (after azure-prepare completes)
3. Both skills are invoked in order

## Common Workflows

### Quick Evaluation of a Single Skill

```bash
cd skills/my-skill
waza run
```

### Check Skill Readiness

```bash
waza check skills/my-skill
```

This performs:
1. **Compliance scoring** — Validates SKILL.md frontmatter (Low/Medium/Medium-High/High)
2. **Token budget** — Checks if SKILL.md is within limits
3. **Evaluation suite** — Looks for eval.yaml

### Improve Skill Compliance

```bash
waza dev skills/my-skill --target medium-high
```

Iteratively scores and improves your SKILL.md based on:
- Description length and clarity
- Presence of USE FOR / DO NOT USE FOR sections
- Trigger clarity and routing indicators

### Compare Results Across Models

```bash
waza run eval.yaml --model claude-opus-4-20250514 --output results-opus.json
waza run eval.yaml --model claude-sonnet-4-20250514 --output results-sonnet.json
waza compare results-opus.json results-sonnet.json
```

### Run Evals in CI

Add to `.github/workflows/eval.yml`:

```yaml
name: Evaluate Skills

on: [push, pull_request]

jobs:
  evaluate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Install waza
        run: |
          curl -fsSL https://raw.githubusercontent.com/spboyer/waza/main/install.sh | bash
      - name: Run evaluations
        run: waza run --verbose --output results.json
      - name: Upload results
        uses: actions/upload-artifact@v4
        with:
          name: evaluation-results
          path: results.json
```

## Grader Types

Waza supports multiple grader types for comprehensive evaluation:

| Grader | Purpose |
|--------|---------|
| `code` | Python/JavaScript assertion-based validation |
| `regex` | Pattern matching in output |
| `file` | File existence and content validation |
| `diff` | Workspace file comparison with snapshots |
| `behavior` | Agent behavior constraints (tool calls, tokens) |
| `action_sequence` | Tool call sequence validation |
| `skill_invocation` | Skill orchestration sequence validation |
| `prompt` | LLM-as-judge evaluation with rubrics |

See [GRADERS.md](GRADERS.md) for detailed configuration options.

## Troubleshooting

### Evaluation fails with "context not found"

Make sure the `fixtures/` directory exists and is referenced correctly:

```yaml
tasks:
  - "tasks/*.yaml"
```

Fixtures are automatically discovered relative to the eval.yaml file.

### Skill check shows "Low" compliance

Add more detail to your SKILL.md:

- Write a description longer than 150 characters
- Add specific USE FOR scenarios
- Add DO NOT USE FOR (anti-triggers) to clarify scope

```bash
waza dev skills/my-skill --target medium-high --auto
```

### Token count exceeds budget

Skill.md has a 500-token soft limit (5000 hard limit). Suggestions:

- Move detailed examples to separate files
- Use concise language for trigger/anti-trigger lists
- Link to external documentation instead of duplicating

```bash
waza tokens suggest skills/my-skill/SKILL.md
```

### Results are cached when I don't want them

Clear the cache before running:

```bash
waza cache clear
waza run eval.yaml
```

Or disable caching:

```bash
waza run eval.yaml --no-cache
```

## Next Steps

1. **[Review the Tutorial](TUTORIAL.md)** — Deep dive into eval spec format and grader configuration
2. **[Explore Graders](GRADERS.md)** — Complete reference for all 7 grader types
3. **[Check Demo Guide](DEMO-GUIDE.md)** — Live demo scenarios and presentations
4. **[CI Integration](SKILLS_CI_INTEGRATION.md)** — Set up automated skill evaluation in GitHub Actions

## Installation

### Binary Install (recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/spboyer/waza/main/install.sh | bash
```

### Install from Source

Requires Go 1.25+:

```bash
go install github.com/spboyer/waza/cmd/waza@latest
```

### Azure Developer CLI (azd) Extension

```bash
azd ext source add -n waza -t url -l https://raw.githubusercontent.com/spboyer/waza/main/registry.json
azd ext install microsoft.azd.waza
azd waza --help
```

## Resources

- **Repository:** [github.com/spboyer/waza](https://github.com/spboyer/waza)
- **Tracking Issue:** [#66 — Waza Platform Roadmap](https://github.com/spboyer/waza/issues/66)
- **Example Projects:**
  - [code-explainer](../examples/code-explainer/) — Simple grader showcase
  - [grader-showcase](../examples/grader-showcase/) — All grader types demo
