# Waza User Guide

Welcome to **Waza**, a CLI tool for evaluating AI agent skills. This guide covers everything you need to get started: installation, creating evaluations, running benchmarks, and reviewing results in the dashboard.

## What is Waza?

Waza helps you:

- **Define agent skills** with comprehensive documentation and behavioral requirements
- **Create test suites** with realistic test cases and validation rules  
- **Run evaluations** against different AI models to measure skill effectiveness
- **Compare results** across models and versions to track improvement
- **View metrics** in an interactive dashboard with live results, trends, and detailed analysis

Perfect for skill authors, platform teams, and developers building AI-powered applications.

---

## Installation

Choose one of three methods:

### 1. Binary Install (Recommended)

The fastest way to get started. The script auto-detects your OS and architecture:

```bash
curl -fsSL https://raw.githubusercontent.com/spboyer/waza/main/install.sh | bash
```

This downloads the latest release, verifies the checksum, and installs the `waza` binary to:
- `/usr/local/bin/waza` (if writable), or
- `~/bin/waza` (if `/usr/local/bin` is not writable)

After installation, verify it works:

```bash
waza --version
```

**Note:** If installed to `~/bin`, add it to your PATH:

```bash
export PATH="$HOME/bin:$PATH"
```

### 2. Install from Source

Requires **Go 1.25 or later**:

```bash
go install github.com/spboyer/waza/cmd/waza@latest
```

Verify installation:

```bash
waza --version
```

### 3. Azure Developer CLI Extension

If you use Azure Developer CLI (`azd`), install waza as an extension:

```bash
# Add the waza extension registry
azd ext source add -n waza -t url -l https://raw.githubusercontent.com/spboyer/waza/main/registry.json

# Install the extension
azd ext install microsoft.azd.waza

# Verify it works
azd waza --help
```

All waza commands are available under `azd waza`:

```bash
azd waza init my-project
azd waza run eval.yaml
azd waza serve
```

---

## Quick Start

Get a complete evaluation suite running in 5 minutes.

### Step 1: Initialize a Project

Create a new directory and initialize a waza project:

```bash
mkdir my-eval-suite
cd my-eval-suite
waza init
```

You'll be prompted to create your first skill. Enter a name like `code-explainer`, and the scaffolding is automatic. This creates:

```
my-eval-suite/
├── skills/                     # Skill definitions
│   └── code-explainer/
│       └── SKILL.md            # Skill metadata and description
├── evals/                      # Evaluation suites
│   └── code-explainer/
│       ├── eval.yaml           # Evaluation configuration
│       ├── tasks/              # Test case definitions
│       │   ├── basic-usage.yaml
│       │   └── edge-cases.yaml
│       └── fixtures/           # Test data and resources
│           └── sample.py
├── .github/workflows/
│   └── eval.yml                # CI/CD pipeline
├── .gitignore
└── README.md
```

**Skip the prompt?** Use `--no-skill`:

```bash
waza init --no-skill
```

### Step 2: Create a New Skill (if needed)

If you didn't create one during `init`, add a new skill:

```bash
waza new code-analyzer
```

This scaffolds a new skill with SKILL.md and eval suite. The interactive wizard will collect metadata (name, description, use cases).

### Step 3: Define Your Evaluation

Edit `evals/code-explainer/eval.yaml`:

```yaml
name: code-explainer
description: "Tests the agent's ability to explain code"
model: claude-sonnet-4.6
maxTokens: 4096
tasks:
  - id: basic-usage
    description: "Explain a simple Python function"
    fixture: sample.py
    expectedOutput:
      - contains: "function"
      - contains: "parameter"
      - regex: "returns.*value"
validators:
  - type: contains
    caseSensitive: false
  - type: regex
    caseSensitive: false
```

Add test case YAML files in `tasks/`:

```yaml
# tasks/basic-usage.yaml
id: basic-usage
description: "Explain a simple Python function"
prompt: "Explain what this function does:\n{{fixture:sample.py}}"
expectedOutput:
  - type: contains
    value: "function"
  - type: regex
    value: "returns.*value"
tags: ["basic", "core"]
```

Add fixtures in `fixtures/`:

```python
# fixtures/sample.py
def greet(name):
    """Return a greeting message."""
    return f"Hello, {name}!"
```

### Step 4: Run the Evaluation

Execute the benchmark:

```bash
waza run evals/code-explainer/eval.yaml --context-dir evals/code-explainer/fixtures -v
```

**Output:**
- `✓ Passed` — Task passed all validators
- `✗ Failed` — Task failed one or more validators
- Summary: `X passed, Y failed`

Save results to a JSON file:

```bash
waza run evals/code-explainer/eval.yaml \
  --context-dir evals/code-explainer/fixtures \
  -o results.json
```

### Step 5: View Results in the Dashboard

Launch the web dashboard:

```bash
waza serve
```

The browser opens automatically to `http://localhost:3000`. You'll see:

- **Dashboard Overview** — Run history, pass rate, model comparison
- **Run Details** — Individual task results, validation output
- **Compare** — Side-by-side model performance
- **Trends** — Pass rate over time

---

## Command Reference

### `waza init [directory]`

Initialize a waza project with directory structure.

**Flags:**
- `--no-skill` — Skip the first-skill creation prompt

**What it creates:**
- `skills/` — Skill definitions directory
- `evals/` — Evaluation suites directory
- `.github/workflows/eval.yml` — CI/CD pipeline
- `.gitignore` — With waza-specific entries
- `README.md` — Getting started guide

**Example:**
```bash
waza init my-skills --no-skill
```

---

### `waza new [skill-name]`

Create a new skill with its evaluation suite.

**Flags:**
- `--template, -t` — Template pack (coming soon)

**Modes:**

**Project mode** (inside a `skills/` directory):
```bash
cd my-skills-repo
waza new code-explainer
```
Creates `skills/code-explainer/SKILL.md` and `evals/code-explainer/`.

**Standalone mode** (no `skills/` directory):
```bash
cd my-project
waza new my-skill
```
Creates `my-skill/` with `SKILL.md`, `evals/`, CI/CD pipeline, and README.

**Examples:**
```bash
# Interactive wizard in a terminal
waza new code-analyzer

# Non-interactive (CI/CD)
waza new code-analyzer << EOF
Code Analyzer
Analyzes code for patterns and issues
code, analysis
EOF
```

---

### `waza run [eval.yaml | skill-name]`

Run an evaluation benchmark.

**Arguments:**
- `[eval.yaml]` — Path to evaluation spec file
- `[skill-name]` — Skill name (auto-detects eval.yaml)
- *(none)* — Auto-detect using workspace detection

**Flags:**
- `--context-dir` — Fixtures directory (default: `./fixtures`)
- `--output, -o` — Save results JSON
- `--verbose, -v` — Detailed progress output
- `--parallel` — Run tasks concurrently
- `--workers <n>` — Number of concurrent workers (default: 4)
- `--task <pattern>` — Filter tasks by name (glob pattern, can repeat)
- `--tags <pattern>` — Filter tasks by tags (glob pattern, can repeat)
- `--model <model>` — Override model (can repeat for multi-model runs)
- `--cache` — Enable result caching
- `--no-cache` — Disable caching (default)
- `--cache-dir <path>` — Cache directory (default: `.waza-cache`)
- `--format <format>` — Output format: `default`, `github-comment`
- `--interpret` — Print plain-language interpretation of results
- `--baseline` — Run A/B comparison: with skills vs without

**Examples:**

Run an evaluation from a spec file:
```bash
waza run evals/code-explainer/eval.yaml --context-dir evals/code-explainer/fixtures -v
```

Run a specific skill:
```bash
waza run code-explainer
```

Auto-detect in workspace:
```bash
# Single-skill workspace → runs that skill's eval
# Multi-skill workspace → runs all evals with summary
waza run
```

Run specific tasks:
```bash
waza run evals/code-explainer/eval.yaml --task "basic*" --task "edge*"
```

Run with specific models:
```bash
waza run evals/code-explainer/eval.yaml --model "gpt-4o" --model "claude-sonnet-4.6"
```

Save results:
```bash
waza run evals/code-explainer/eval.yaml -o results.json
```

---

### `waza check [skill-name | skill-path]`

Validate that a skill is ready for submission.

**Arguments:**
- `[skill-name]` — Skill name (e.g., `code-explainer`)
- `[skill-path]` — Path to skill directory (e.g., `skills/my-skill`)
- *(none)* — Auto-detect using workspace detection

**What it checks:**
1. **Compliance scoring** — Validates SKILL.md frontmatter (Low/Medium/Medium-High/High)
2. **Token budget** — Ensures SKILL.md is within limits
3. **Evaluation presence** — Confirms eval.yaml exists

**Output:**
- ✓ Compliance level (e.g., "Medium-High")
- ✓ Token count / limit
- ✓ Evaluation status
- Suggestions for improvement

**Examples:**

Check a specific skill:
```bash
waza check code-explainer
```

Check by path:
```bash
waza check skills/my-skill
```

Auto-detect in workspace:
```bash
# Single skill → checks that skill
# Multiple skills → summary table for all
waza check
```

---

### `waza serve`

Launch the web dashboard to view evaluation results.

**Flags:**
- `--port <n>` — HTTP server port (default: `3000`)
- `--no-browser` — Don't auto-open the browser
- `--results-dir <path>` — Directory to read results from (default: `.`)
- `--tcp <address>` — JSON-RPC TCP server (e.g., `:9000`)
- `--tcp-allow-remote` — Bind to all interfaces (default: loopback only)
- `--http` — Start HTTP dashboard (default)

**Dashboard Pages:**
- **Overview** — Run history, pass rates, model comparison
- **Run Details** — Individual task results, validation logs, timestamps
- **Compare** — Side-by-side model performance
- **Trends** — Pass rate and timing trends over time
- **Live View** — Real-time results during active runs

**Examples:**

Start the dashboard:
```bash
waza serve
```

Custom port:
```bash
waza serve --port 8080
```

Load results from a specific directory:
```bash
waza serve --results-dir ./archived-runs
```

Don't auto-open browser:
```bash
waza serve --no-browser
```

JSON-RPC TCP server (for IDE integration):
```bash
waza serve --tcp :9000
```

---

## Advanced Usage

### Caching and Reproducibility

Cache evaluation results to avoid redundant runs:

```bash
# Enable caching (results stored in .waza-cache/)
waza run evals/code-explainer/eval.yaml --cache -o results.json
```

The cache stores results by task ID and model, keyed on the prompt and expected output. Same inputs = same cached results.

**Use cases:**
- Development — Iterate on validators without re-running tasks
- CI/CD — Avoid redundant API calls for unchanged evals
- Offline — Use cached results when connection is unavailable

Disable caching:
```bash
waza run evals/code-explainer/eval.yaml --no-cache
```

Specify cache location:
```bash
waza run evals/code-explainer/eval.yaml --cache --cache-dir ./my-cache
```

---

### Filtering and Selective Runs

Run only specific tasks:

```bash
# By task name (glob patterns)
waza run evals/code-explainer/eval.yaml --task "basic*" --task "edge*"

# By tag
waza run evals/code-explainer/eval.yaml --tags "critical"
```

**Task names and tags** are defined in your task YAML files:

```yaml
# tasks/basic-usage.yaml
id: basic-usage
description: "Explain a simple function"
tags: ["basic", "core"]
```

---

### Parallel Execution

Run tasks concurrently to speed up evaluation:

```bash
# Use default 4 workers
waza run evals/code-explainer/eval.yaml --parallel

# Use custom worker count
waza run evals/code-explainer/eval.yaml --parallel --workers 8
```

**When to use parallel execution:**
- Many tasks (50+)
- Long-running validators
- Development machines with spare CPU

---

### Multi-Model Comparison

Run the same evaluation against multiple models:

```bash
waza run evals/code-explainer/eval.yaml \
  --model "gpt-4o" \
  --model "claude-sonnet-4.6" \
  --model "gpt-4-turbo"
```

Results are grouped by model. Compare them in the dashboard or with:

```bash
waza compare results-gpt4.json results-sonnet.json
```

---

### CI/CD Integration

Run waza as part of your GitHub Actions workflow:

```yaml
# .github/workflows/eval.yml
name: Evaluate Skills

on: [push, pull_request]

jobs:
  evaluate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
        with:
          go-version: '1.25'
      - name: Install waza
        run: |
          curl -fsSL https://raw.githubusercontent.com/spboyer/waza/main/install.sh | bash
      - name: Run evals
        run: |
          waza run --output results.json -v
      - name: Comment PR
        if: github.event_name == 'pull_request'
        uses: actions/github-script@v7
        with:
          script: |
            const fs = require('fs');
            const results = JSON.parse(fs.readFileSync('results.json'));
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: `Evaluation Results: ${results.passCount}/${results.totalCount} passed`
            });
```

---

### Session Logging

Capture detailed session logs for debugging:

```bash
waza run evals/code-explainer/eval.yaml --session-log --session-dir ./logs
```

Logs are stored in NDJSON format (one event per line):
```json
{"event":"task_started","id":"basic-usage","timestamp":"2024-01-15T10:30:00Z"}
{"event":"task_completed","id":"basic-usage","passed":true,"timestamp":"2024-01-15T10:30:05Z"}
```

---

### Output Formats

**Default format** (human-readable):
```bash
waza run evals/code-explainer/eval.yaml
```

**GitHub Comment format** (for PR comments):
```bash
waza run evals/code-explainer/eval.yaml --format github-comment
```

**JSON output** (for automation):
```bash
waza run evals/code-explainer/eval.yaml -o results.json
```

---

## Dashboard

The `waza serve` command launches an interactive web dashboard for viewing and analyzing evaluation results.

### Starting the Dashboard

```bash
waza serve
```

The dashboard opens automatically at `http://localhost:3000`.

### Navigation

- **Overview** — Evaluation summary, pass rates, recent runs
- **Run Details** — Click a run to see task-by-task results
- **Compare** — Select multiple runs to compare models
- **Trends** — Historical pass rates and performance over time
- **Live View** — Real-time results during active evaluations

### Dashboard Features

- **Live updates** — See results as tasks complete
- **Search** — Find runs, tasks, and models
- **Filtering** — Filter by status (passed/failed), tags, date range
- **Export** — Download results as CSV or JSON
- **Dark mode** — Switch to dark theme for comfortable viewing

### Advanced Dashboard Flags

**Custom port:**
```bash
waza serve --port 8080
```

**Skip auto-open:**
```bash
waza serve --no-browser
```

**Load results from archive:**
```bash
waza serve --results-dir ./previous-runs
```

**JSON-RPC for IDE integration:**
```bash
waza serve --tcp :9000
```

---

## Troubleshooting

### Port Already in Use

If you see "address already in use," the default port 3000 is taken.

**Solution:** Use a different port:
```bash
waza serve --port 3001
```

Or find and stop the process using port 3000:
```bash
# macOS/Linux
lsof -i :3000
kill -9 <PID>

# Windows
netstat -ano | findstr :3000
taskkill /PID <PID> /F
```

---

### No Results Found in Dashboard

If the dashboard shows no results:

1. **Check results directory:** Make sure you've saved results with `-o`:
   ```bash
   waza run evals/code-explainer/eval.yaml -o results.json
   ```

2. **Specify results directory:** If results are in a subdirectory:
   ```bash
   waza serve --results-dir ./eval-results
   ```

3. **Verify file format:** Results must be JSON files (`*.json`).

---

### Fixture Directory Not Found

If you see "fixture directory not found":

1. **Check the path:** Is your `--context-dir` correct?
   ```bash
   waza run eval.yaml --context-dir ./evals/fixtures
   ```

2. **Use relative paths:** Path is relative to the spec file:
   ```
   eval.yaml
   fixtures/
   ```

3. **From any directory:** Specify absolute paths:
   ```bash
   waza run /path/to/eval.yaml --context-dir /path/to/fixtures
   ```

---

### Validation Failures

If tasks fail validation unexpectedly:

1. **Enable verbose output:** See what the model returned:
   ```bash
   waza run eval.yaml -v
   ```

2. **Check your validators:** Are the rules too strict?
   ```yaml
   validators:
     - type: contains
       value: "function"
       caseSensitive: false  # Don't be too strict
   ```

3. **Adjust expected output:** Make rules realistic:
   ```yaml
   expectedOutput:
     - type: contains
       value: "important keyword"  # Not every word needs to match
   ```

---

## Next Steps

- **[Create your first skill](./GETTING-STARTED.md)** — Step-by-step walkthrough
- **[Skill best practices](./SKILL-BEST-PRACTICES.md)** — Guidelines for effective skills
- **[Validators reference](./GRADERS.md)** — All validator types and options
- **[Token limits](./TOKEN-LIMITS.md)** — Optimize SKILL.md for size
- **[Examples](../examples/)** — Real eval suites to explore

---

## Getting Help

- **[GitHub Issues](https://github.com/spboyer/waza/issues)** — Report bugs or request features
- **[Discussions](https://github.com/spboyer/waza/discussions)** — Ask questions, share ideas
- **[Waza Examples](../examples/)** — Runnable evaluation suites

---

**Happy evaluating! 🚀**
