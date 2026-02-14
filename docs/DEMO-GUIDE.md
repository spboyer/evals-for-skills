# Waza Demo Guide

Comprehensive guide for demonstrating waza evaluation framework features in various scenarios.

## Overview

This guide provides 7 ready-to-run demo scenarios showcasing different aspects of waza:
1. [Basic Evaluation Run](#demo-1-basic-evaluation-run)
2. [Multiple Grader Types](#demo-2-multiple-grader-types)
3. [Model Comparison](#demo-3-model-comparison)
4. [CI/CD Integration](#demo-4-cicd-integration)
5. [Token Management](#demo-5-token-management)
6. [Skill Generation](#demo-6-skill-generation)
7. [Multi-Skill Orchestration](#demo-7-multi-skill-orchestration) ⭐ **NEW**

Each demo includes setup instructions, commands, and expected output.

---

## Demo 1: Basic Evaluation Run

**Duration**: 2-3 minutes  
**Goal**: Show the simplest eval workflow

### Setup

```bash
cd /home/runner/work/waza/waza
make build
```

### Run the Demo

```bash
# Run the code-explainer example with mock executor
./waza run examples/code-explainer/eval.yaml \
  --executor mock \
  --context-dir examples/code-explainer/fixtures \
  -v
```

### Expected Output

```
Running evaluation: code-explainer v1.0
Executor: mock
Tasks: 4
Trials per task: 3

Task 1/4: explain-python-recursion
  Trial 1/3... ✓ (0.5s)
  Trial 2/3... ✓ (0.4s)
  Trial 3/3... ✓ (0.5s)
  Pass rate: 3/3 (100%)

...

Results:
  Pass rate: 12/12 (100%)
  Average score: 0.95
```

### Key Points to Highlight

- Simple CLI interface
- Mock executor for fast testing
- Multiple trials for reliability
- Clear pass/fail results

---

## Demo 2: Multiple Grader Types

**Duration**: 3-4 minutes  
**Goal**: Show different validation methods

### Setup

```bash
cd /home/runner/work/waza/waza
./waza run examples/grader-showcase/eval.yaml -v
```

### Available Graders

The grader showcase demonstrates 6 grader types:

| Grader | Purpose | Example Use Case |
|--------|---------|------------------|
| `code` | Python assertions | Complex validation logic |
| `regex` | Pattern matching | URL format, keywords |
| `file` | File validation | Check created files |
| `behavior` | Tool usage limits | Efficiency, safety |
| `action_sequence` | Tool call order | Workflow enforcement |
| `skill_invocation` | Skill orchestration | Multi-skill workflows |

### Demo Commands

```bash
# Run all grader demos
./waza run examples/grader-showcase/eval.yaml -v

# Run specific grader demo
./waza run examples/grader-showcase/eval.yaml --task="regex-demo" -v

# View a specific grader config
cat examples/grader-showcase/tasks/regex-task.yaml
```

### Example: Regex Grader

```yaml
graders:
  - type: regex
    name: pattern_check
    config:
      must_match:
        - "(?i)deployed to https?://.+"
        - "Resource group: .+"
      must_not_match:
        - "error|failed|exception"
```

### Key Points to Highlight

- Multiple grader types for different validation needs
- Combine graders for comprehensive testing
- Task-specific vs global graders
- See `docs/GRADERS.md` for full reference

---

## Demo 3: Model Comparison

**Duration**: 4-5 minutes  
**Goal**: Compare performance across different models

### Setup

```bash
cd /home/runner/work/waza/waza

# Run eval with different models (requires Copilot SDK)
./waza run examples/code-explainer/eval.yaml \
  --executor copilot-sdk \
  --model gpt-4o \
  --context-dir examples/code-explainer/fixtures \
  -o results-gpt4.json

./waza run examples/code-explainer/eval.yaml \
  --executor copilot-sdk \
  --model claude-sonnet-4-20250514 \
  --context-dir examples/code-explainer/fixtures \
  -o results-claude.json
```

### Compare Results

```bash
./waza compare results-gpt4.json results-claude.json
```

### Expected Output

```
Model Comparison
================

                    gpt-4o    claude-sonnet-4-20250514
Pass Rate           92%       95%
Average Score       0.88      0.92
Avg Duration        12.3s     9.7s
Avg Tool Calls      5.2       4.8

Task Breakdown:
  explain-recursion    0.90 → 0.95  (+0.05) ✓
  explain-async        0.85 → 0.88  (+0.03) ✓
  explain-comprehension 0.90 → 0.93  (+0.03) ✓
  explain-sql-join     0.87 → 0.92  (+0.05) ✓

Winner: claude-sonnet-4-20250514
```

### Key Points to Highlight

- Compare models objectively
- Task-by-task score deltas
- Performance metrics (duration, tool calls)
- Data-driven model selection

---

## Demo 4: CI/CD Integration

**Duration**: 3-4 minutes  
**Goal**: Show how waza integrates into CI/CD pipelines

### Setup

Create a GitHub Actions workflow:

```yaml
# .github/workflows/eval-skill.yml
name: Evaluate Skill

on:
  pull_request:
    paths:
      - 'skills/**'
      - 'eval/**'

jobs:
  evaluate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Install waza
        run: |
          cd waza-go
          make build
          sudo mv waza /usr/local/bin/
      
      - name: Run evaluation
        id: eval
        run: |
          waza run eval/eval.yaml \
            --executor mock \
            --format github-comment > comment.md
          echo "exit_code=$?" >> $GITHUB_OUTPUT
      
      - name: Post PR comment
        if: always()
        uses: actions/github-script@v7
        with:
          script: |
            const fs = require('fs');
            const comment = fs.readFileSync('comment.md', 'utf8');
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: comment
            });
      
      - name: Fail on test failure
        if: steps.eval.outputs.exit_code == '1'
        run: exit 1
```

### Exit Codes

| Code | Meaning | CI Action |
|------|---------|-----------|
| 0 | All tests passed | ✓ Pass build |
| 1 | Test failure | ✗ Fail build |
| 2 | Config error | ✗ Fail build |

### Demo Command

```bash
# Test locally with GitHub comment format
./waza run examples/code-explainer/eval.yaml \
  --executor mock \
  --context-dir examples/code-explainer/fixtures \
  --format github-comment
```

### Expected Output (Markdown)

```markdown
## 🎯 Waza Evaluation Results

**Skill**: code-explainer v1.0  
**Status**: ✅ **PASSED** (12/12 tasks)  
**Score**: 0.95/1.0

### Summary
- Pass rate: 100% (12/12)
- Average score: 0.95
- Total duration: 14.2s

### Tasks
| Task | Status | Score | Duration |
|------|--------|-------|----------|
| explain-recursion | ✅ | 0.96 | 3.5s |
| explain-async | ✅ | 0.94 | 3.2s |
...
```

### Key Points to Highlight

- Exit codes for build integration
- GitHub comment formatting
- Automated PR feedback
- Mock executor for fast CI runs

---

## Demo 5: Token Management

**Duration**: 2-3 minutes  
**Goal**: Show token counting and optimization

### Count Tokens

```bash
cd /home/runner/work/waza/waza

# Count tokens in skill files
./waza tokens count skills/code-explainer/

# Count with breakdown
./waza tokens count skills/code-explainer/ --verbose
```

### Expected Output

```
Counting tokens in: skills/code-explainer/

File                    Tokens
SKILL.md               2,341
references/python.md   1,245
references/js.md       987
-----------------------------------
Total                  4,573 tokens
```

### Suggest Optimizations

```bash
# Get optimization suggestions
./waza tokens suggest skills/code-explainer/
```

### Expected Output

```
Token Optimization Suggestions
==============================

skills/code-explainer/SKILL.md (2,341 tokens)

  ⚠️ Large code block (lines 45-78)
     Current: 456 tokens
     Suggestion: Consider extracting to a reference file
     Savings: ~400 tokens

  ⚠️ Large table (lines 120-145)
     Current: 234 tokens
     Suggestion: Simplify table or move to reference
     Savings: ~180 tokens

Potential savings: ~580 tokens (25%)
```

### Key Points to Highlight

- Monitor skill context size
- Identify optimization opportunities
- Stay within model limits
- See `docs/TOKEN-LIMITS.md` for details

---

## Demo 6: Skill Generation

**Duration**: 3-4 minutes  
**Goal**: Show automatic eval generation from SKILL.md

### Generate from GitHub Repo

```bash
# Generate eval for a specific skill
./waza generate \
  --repo microsoft/GitHub-Copilot-for-Azure \
  --skill azure-functions \
  -o /tmp/azure-functions-eval
```

### Expected Output

```
Scanning GitHub repository: microsoft/GitHub-Copilot-for-Azure
✓ Found 15 skill(s)

Parsing: azure-functions
✓ Parsed skill: azure-functions
  Description: Build and deploy serverless Azure Functions
  Triggers: 15
  Anti-triggers: 4

✓ Created eval.yaml
✓ Created trigger_tests.yaml
✓ Created tasks/task-001.yaml
✓ Created tasks/task-002.yaml
✓ Created tasks/task-003.yaml
✓ Created fixtures/function_app.py

Generated eval suite at: /tmp/azure-functions-eval
```

### LLM-Assisted Generation

```bash
# Use LLM to generate better tasks
./waza generate \
  --repo microsoft/GitHub-Copilot-for-Azure \
  --skill azure-functions \
  -o /tmp/azure-functions-eval-ai \
  --assist
```

### Key Points to Highlight

- Zero-effort eval creation
- Generated from existing SKILL.md
- LLM assistance for realistic tasks
- Includes fixtures and graders

---

## Demo 7: Multi-Skill Orchestration

**Duration**: 5-6 minutes  
**Goal**: Show skill orchestration validation with real skill directories ⭐ **NEW**

### Overview

This demo shows how to validate that an orchestration skill correctly invokes dependent skills in the proper sequence. This is essential for complex workflows that coordinate multiple skills.

### Setup Real Skill Directories

First, let's set up a realistic scenario with three Azure deployment skills:

```bash
cd /tmp
mkdir -p orchestration-demo/skills

# Create azure-prepare skill
mkdir -p orchestration-demo/skills/azure-prepare
cat > orchestration-demo/skills/azure-prepare/SKILL.md << 'EOF'
---
name: azure-prepare
version: 1.0.0
description: Prepare Azure environment for deployment
---

# Azure Prepare Skill

Prepares the Azure environment by:
- Creating resource groups
- Setting up networking
- Configuring storage accounts
EOF

# Create azure-deploy skill
mkdir -p orchestration-demo/skills/azure-deploy
cat > orchestration-demo/skills/azure-deploy/SKILL.md << 'EOF'
---
name: azure-deploy
version: 1.0.0
description: Deploy application to Azure
---

# Azure Deploy Skill

Deploys the application to Azure infrastructure.
EOF

# Create azure-monitor skill
mkdir -p orchestration-demo/skills/azure-monitor
cat > orchestration-demo/skills/azure-monitor/SKILL.md << 'EOF'
---
name: azure-monitor
version: 1.0.0
description: Monitor Azure deployment status
---

# Azure Monitor Skill

Monitors the deployment status and health.
EOF
```

### Create Orchestration Eval

```bash
# Create eval spec with skill_directories and skill_invocation grader
cat > orchestration-demo/eval.yaml << 'EOF'
name: azure-orchestration-eval
description: Validate Azure deployment orchestration workflow
skill: azure-orchestrator
version: "1.0"

config:
  trials_per_task: 1
  timeout_seconds: 300
  executor: mock
  model: claude-sonnet-4-20250514
  
  # Point to our skill directories
  skill_directories:
    - ./skills/azure-prepare
    - ./skills/azure-deploy
    - ./skills/azure-monitor
  
  # Validate required skills are present
  required_skills:
    - azure-prepare
    - azure-deploy
    - azure-monitor

metrics:
  - name: orchestration_correctness
    weight: 0.8
    threshold: 0.9
  - name: efficiency
    weight: 0.2
    threshold: 0.7

# Global grader: no errors
graders:
  - type: regex
    name: no_errors
    config:
      must_not_match:
        - "(?i)error|exception|failed"

tasks:
  - "tasks/*.yaml"
EOF

# Create task directory
mkdir -p orchestration-demo/tasks

# Create orchestration task
cat > orchestration-demo/tasks/deployment-workflow.yaml << 'EOF'
name: deployment-workflow
description: Test full Azure deployment orchestration

prompt: |
  Deploy my application to Azure. First prepare the environment,
  then deploy the application, and finally monitor the deployment.

graders:
  # Exact sequence: prepare → deploy → monitor
  - type: skill_invocation
    name: exact_deployment_sequence
    config:
      mode: exact_match
      required_skills:
        - azure-prepare
        - azure-deploy
        - azure-monitor
      allow_extra: false
  
  # Verify key skills in order (more flexible)
  - type: skill_invocation
    name: flexible_deployment_flow
    config:
      mode: in_order
      required_skills:
        - azure-prepare
        - azure-deploy
      allow_extra: true
  
  # Check all required skills invoked (any order)
  - type: skill_invocation
    name: all_skills_invoked
    config:
      mode: any_order
      required_skills:
        - azure-prepare
        - azure-deploy
        - azure-monitor
      allow_extra: true
EOF
```

### Run the Orchestration Eval

```bash
cd orchestration-demo

# Run with waza
/home/runner/work/waza/waza/waza run eval.yaml -v
```

### Expected Output

```
Running evaluation: azure-orchestration-eval v1.0
Executor: mock

✓ Required skills validation passed (3/3 skills found)
  - azure-prepare
  - azure-deploy
  - azure-monitor

Task 1/1: deployment-workflow
  Trial 1/1...
    Grader: exact_deployment_sequence
      Mode: exact_match
      Required: [azure-prepare, azure-deploy, azure-monitor]
      Actual: [azure-prepare, azure-deploy, azure-monitor]
      ✓ Passed (score: 1.0)
    
    Grader: flexible_deployment_flow
      Mode: in_order
      Required: [azure-prepare, azure-deploy]
      Actual: [azure-prepare, azure-deploy, azure-monitor]
      ✓ Passed (score: 1.0)
    
    Grader: all_skills_invoked
      Mode: any_order
      Required: [azure-prepare, azure-deploy, azure-monitor]
      Actual: [azure-prepare, azure-deploy, azure-monitor]
      ✓ Passed (score: 1.0)
  
  ✓ (1.2s)

Results:
  Pass rate: 1/1 (100%)
  Average score: 1.0
  orchestration_correctness: 1.0
  efficiency: 1.0
```

### Test Missing Skills

```bash
# Test what happens when a required skill is missing
cat > orchestration-demo/eval-missing.yaml << 'EOF'
name: azure-orchestration-eval
skill: azure-orchestrator
version: "1.0"

config:
  executor: mock
  
  skill_directories:
    - ./skills/azure-deploy  # Only has azure-deploy
  
  required_skills:
    - azure-prepare  # Missing!
    - azure-deploy
    - azure-monitor  # Missing!

tasks:
  - "tasks/*.yaml"
EOF

# Run - should fail validation
/home/runner/work/waza/waza/waza run eval-missing.yaml
```

### Expected Error

```
Error: skill validation failed

Required skills not found:
  - azure-prepare
  - azure-monitor

Searched directories:
  - ./skills/azure-deploy

Found skills:
  - azure-deploy

Tip: Add missing skills to skill_directories or update required_skills list
```

### Grading Modes Comparison

| Mode | Requirement | Use Case |
|------|-------------|----------|
| `exact_match` | Exact sequence, same length | Strict workflows (demos, compliance) |
| `in_order` | Required skills in order, allows extras | Flexible orchestration |
| `any_order` | All required skills present, any order | Loose validation |

### Key Points to Highlight

- ✅ **Required skills validation** - Catches missing skills before eval starts
- ✅ **Real skill directories** - Works with actual SKILL.md files
- ✅ **Multiple matching modes** - Choose strictness level
- ✅ **Clear error messages** - Shows what's missing and where waza looked
- ✅ **skill_invocation grader** - Validates orchestration workflows
- ✅ **Backward compatible** - `required_skills` is optional

### Advanced: Combine with Other Graders

```yaml
graders:
  # Validate orchestration sequence
  - type: skill_invocation
    name: deployment_sequence
    config:
      mode: in_order
      required_skills:
        - azure-prepare
        - azure-deploy
  
  # Validate efficiency
  - type: behavior
    name: efficiency
    config:
      max_tool_calls: 20
      max_duration_ms: 30000
  
  # Validate output
  - type: regex
    name: success_message
    config:
      must_match:
        - "(?i)deployment.*success"
```

---

## Tips for Effective Demos

### General Guidelines

1. **Start Simple**: Use demo 1 or 2 for first-time audiences
2. **Know Your Audience**: Choose demos relevant to their needs
3. **Prepare Environment**: Pre-build waza and test commands
4. **Have Backups**: Keep expected outputs ready if live demos fail
5. **Use Mock Executor**: Much faster for demos, no API costs

### Time Management

| Scenario | Quick Version | Full Version |
|----------|---------------|--------------|
| Basic eval | 2 min | 5 min |
| Grader types | 3 min | 8 min |
| Model comparison | Skip | 10 min |
| CI/CD | 3 min | 6 min |
| Tokens | 2 min | 4 min |
| Generation | 3 min | 6 min |
| Orchestration | 5 min | 10 min |

### Common Questions

**Q: Can I use my own skills?**  
A: Yes! Use `waza generate` to create evals from your SKILL.md files.

**Q: How do I integrate with GitHub Actions?**  
A: See Demo 4 and `docs/SKILLS_CI_INTEGRATION.md` for complete workflows.

**Q: What if my tests are flaky?**  
A: Increase `trials_per_task` to run multiple trials and get pass rates.

**Q: Can I test orchestration skills?**  
A: Yes! See Demo 7 for multi-skill orchestration validation.

**Q: How do I validate skill invocation sequences?**  
A: Use the `skill_invocation` grader (Demo 7) with appropriate matching modes.

---

## Related Documentation

- **Full Command Reference**: `README.md`
- **Grader Reference**: `docs/GRADERS.md`
- **Tutorial**: `docs/TUTORIAL.md`
- **CI Integration**: `docs/SKILLS_CI_INTEGRATION.md`
- **Token Management**: `docs/TOKEN-LIMITS.md`

---

## Example Demo Flow (15 minutes)

For a comprehensive 15-minute demo, use this flow:

1. **Introduction** (2 min) - What is waza, why evaluate skills
2. **Demo 1: Basic Run** (2 min) - Quick win, show it works
3. **Demo 2: Graders** (3 min) - Show validation flexibility
4. **Demo 7: Orchestration** (5 min) - Advanced feature, multi-skill validation
5. **Demo 4: CI/CD** (2 min) - Real-world integration
6. **Q&A** (1 min) - Address questions

This flow shows progression from simple to advanced while covering key features.
