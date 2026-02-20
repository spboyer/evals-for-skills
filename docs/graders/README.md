# Grader Reference

Complete reference for all available grader types in waza.

## Overview

Graders evaluate skill execution and produce scores. Each grader returns:
- `score`: 0.0 to 1.0
- `passed`: boolean
- `message`: human-readable result
- `details`: additional metadata

## Grader Types

- [`code` - Assertion-Based Grader](code.md)
- [`regex` - Pattern Matching Grader](regex.md)
- [`tool_calls` - Tool Usage Grader (not implemented)](tool_calls.md)
- [`script` - External Script Grader (not implemented)](script.md)
- [`diff` - Workspace File Comparison](diff.md)
- [`llm` - LLM-as-Judge Grader (not implemented)](llm.md)
- [`llm_comparison` - Reference Comparison Grader (not implemented)](llm_comparison.md)
- [`human` - Manual Review Grader (not implemented)](human.md)
- [`human_calibration` - Calibration Grader (not implemented)](human_calibration.md)
- [`behavior` - Agent Behavior Validation](behavior.md)
- [`action_sequence` - Tool Call Sequence Validation](action_sequence.md)
- [`skill_invocation` - Skill Invocation Sequence Validation](skill_invocation.md)
- [`prompt` - LLM-Based Evaluation](prompt.md)

## Inline vs Script Graders

Graders can be defined in two ways:

### Inline Graders (in eval.yaml or task files)

Best for simple validation logic that fits in YAML:

```yaml
graders:
  - type: code
    name: basic_check
    config:
      assertions:
        - "len(output) > 10"
        - "'success' in output.lower()"

  - type: regex
    name: format_check
    config:
      must_match:
        - "deployed to .+"
```

### Script Graders (in graders/ directory)

Best for complex, multi-criteria evaluation logic:

```
my-skill/
├── eval.yaml
├── tasks/
└── graders/
    └── quality_checker.py    # Complex custom logic
```

Reference in eval.yaml:
```yaml
graders:
  - type: script
    name: quality_checker
    config:
      script: graders/quality_checker.py
```

**When to use script graders:**
- Multi-criteria scoring (5+ checks)
- Domain-specific business logic
- Reusable across multiple evals
- Complex pattern matching or analysis
- Integration with external services

See the [code-explainer example](../../examples/code-explainer/graders/explanation_quality.py) for a complete script grader implementation.

## Azure ML Evaluator Integration

Waza's grader system is designed to be compatible with [Azure AI Evaluation SDK](https://learn.microsoft.com/azure/ai-studio/how-to/develop/evaluate-sdk) patterns. The `prompt` grader specifically implements the LLM-as-judge pattern used by Azure ML evaluators.

### Supported Azure ML Evaluator Patterns

The following Azure ML evaluator types map to Waza graders:

| Azure ML Evaluator | Waza Grader | Description || `RelevanceEvaluator` | `prompt` | Uses LLM to judge response relevance |
| `CoherenceEvaluator` | `prompt` | Evaluates logical flow and coherence |
| `FluencyEvaluator` | `prompt` | Assesses natural language quality |
| `GroundednessEvaluator` | `prompt` | Checks factual grounding |
| `ContentSafetyEvaluator` | `prompt` | Validates content safety |
| Custom evaluators | `prompt` | Any custom LLM-based evaluator |

### Adapting Azure ML Evaluators as Waza Rubrics

You can adapt Azure ML evaluators to Waza by converting their evaluation criteria into rubrics:

**Example: Azure ML Relevance Evaluator → Waza Rubric**

```yaml
# Azure ML pattern
evaluator = RelevanceEvaluator(model_config)
result = evaluator(query=query, response=response)

# Equivalent Waza configuration
- type: prompt
  name: relevance_check
  config:
    model: gpt-4o-mini
    rubric: |
      Evaluate how relevant the agent's response is to the user's query.

      Query: [Available in the test case prompt; the agent's response is provided as 'output']

      Rate relevance (1-5):
      1 - Completely irrelevant
      2 - Slightly relevant
      3 - Moderately relevant
      4 - Very relevant
      5 - Perfectly relevant and comprehensive

      Return JSON: {
        "relevance_score": <1-5>,
        "reasoning": "<explanation>",
        "passed": <true if score >= 4>
      }
    threshold: 0.75
```

**Example: Custom Azure ML Evaluator → Waza Rubric**

```python
# Azure ML custom evaluator
from promptflow.core import Prompty

evaluator = Prompty.load("security_evaluator.prompty")

# Convert to Waza rubric
```

```yaml
- type: prompt
  name: security_check
  config:
    model: gpt-4o
    rubric: |
      Evaluate the code for security vulnerabilities:

      Criteria:
      1. Input Validation: Are inputs properly validated?
      2. Authentication: Is auth implemented correctly?
      3. Data Protection: Is sensitive data protected?
      4. Error Handling: Are errors handled securely?

      For each criterion, rate 1-5 and provide specific findings.

      Return JSON: {
        "input_validation": <score>,
        "authentication": <score>,
        "data_protection": <score>,
        "error_handling": <score>,
        "overall_score": <average>,
        "findings": ["<issue 1>", "<issue 2>", ...],
        "passed": <true if overall >= 4>
      }
    threshold: 0.8
```

### Using Azure ML Prompt Flow Templates

Azure ML Prompt Flow `.prompty` files can be adapted to Waza rubrics:

1. **Extract the system prompt** from the `.prompty` file
2. **Convert to YAML rubric** in your grader config
3. **Map inputs** (context variables are available: `output`, `transcript`, `tool_calls`, etc.)
4. **Preserve scoring logic** from the original evaluator

**Example Conversion:**

```promptyname: CodeQualityEvaluator
description: Evaluates code quality
inputs:
  code: string
outputs:
  score: integer
  reasoning: stringsystem:
Evaluate the code quality on a scale of 1-5...
[evaluation criteria]
```

Becomes:

```yaml
- type: prompt
  name: code_quality
  config:
    rubric: |
      Evaluate the code quality on a scale of 1-5...
      [evaluation criteria - copied from .prompty]

      Return JSON: {"score": <1-5>, "reasoning": "..."}
```

### Creating Custom LLM-as-Judge Graders

To create graders that match Azure ML evaluator patterns:

1. **Define clear criteria**: What aspects are you evaluating?
2. **Use consistent scales**: 1-5 or 1-10 (Waza normalizes to 0-1)
3. **Request chain-of-thought**: Ask the LLM to explain its reasoning
4. **Structure output**: Use JSON for reliable parsing
5. **Set appropriate thresholds**: Define what "passing" means for your criteria

**Template:**

```yaml
- type: prompt
  name: my_custom_evaluator
  config:
    model: gpt-4o-mini
    rubric: |
      Evaluate [what you're assessing] based on:

      1. [Criterion 1] (1-5): [Description]
      2. [Criterion 2] (1-5): [Description]
      3. [Criterion 3] (1-5): [Description]

      For each criterion:
      - Consider [specific aspects]
      - Rate honestly and critically
      - Provide specific examples

      Return JSON: {
        "criterion1_score": <1-5>,
        "criterion2_score": <1-5>,
        "criterion3_score": <1-5>,
        "overall_score": <average>,
        "reasoning": "<detailed explanation>",
        "passed": <true/false based on threshold>
      }

      Context available:
      - output: The agent's final response
      - tool_calls: List of tools the agent used
      - duration_ms: Execution time
    threshold: 0.7
    score_type: normalized
```

## Task-Level Graders

You can also define graders per-task:

```yaml
# In task YAML
graders:
  - name: task_specific_check
    type: code
    assertions:
      - "specific_condition"
    weight: 0.5  # Weight within this task
```

## Grader Weights

When multiple graders are used, results are combined:

```yaml
graders:
  - type: code
    name: basic_check
    # Default weight: 1.0

  - type: llm
    name: quality_check
    # Default weight: 1.0
```

**Final Score:** Average of all grader scores (weighted if specified)

## Trigger Tests

Trigger tests measure whether a skill activates for the right prompts and stays
silent for the wrong ones. They run automatically when a `trigger_tests.yaml`
file exists alongside `eval.yaml`.

### File Format

```yaml
skill: code-explainer

should_trigger_prompts:
  - prompt: "Explain this code to me"
    reason: "Direct explanation request"        # optional, for documentation
    confidence: high                            # high (default) or medium

  - prompt: "I don't understand what this code is doing"
    confidence: medium

should_not_trigger_prompts:
  - prompt: "Write me a function to sort a list"
    reason: "Code writing request, not explanation"
    confidence: high

  - prompt: "Fix the bug in my code"
    confidence: medium
```

**Fields:**

| Field | Required | Description || `skill` | yes | Skill name to check for invocation |
| `should_trigger_prompts` | at least one of the two prompt lists | Prompts where the skill should activate |
| `should_not_trigger_prompts` | at least one of the two prompt lists | Prompts where the skill should stay silent |
| `prompt` | yes | The test prompt text |
| `reason` | no | Human-readable explanation (not used in scoring) |
| `confidence` | no | `high` (default) or `medium` — controls scoring weight |

### Confidence Weighting

Each prompt's result is weighted by its confidence level:

- **`high`** (or omitted): weight **1.0** — a clear-cut case where the expected
  behavior is unambiguous.
- **`medium`**: weight **0.5** — an edge case or borderline prompt where the
  expected behavior is less certain.

This lets you include borderline prompts without letting them dominate the score.
For example a "medium" false positive penalizes accuracy half as much as a "high"
one.

### Metrics

Trigger tests produce standard classification metrics:

| Metric | Description || **Accuracy** | (TP + TN) / total |
| **Precision** | TP / (TP + FP) — how often activation was correct |
| **Recall** | TP / (TP + FN) — how often it activated when it should have |
| **F1** | Harmonic mean of precision and recall |
| **Errors** | Prompts that failed to execute (counted as incorrect) |

### Using `trigger_accuracy` as a Metric

Add `trigger_accuracy` to the `metrics` section of your `eval.yaml` to set a
pass/fail threshold:

```yaml
metrics:
  - name: trigger_accuracy
    threshold: 0.9
    weight: 30
```

When configured, trigger accuracy is included in the benchmark outcome and the
run fails if accuracy falls below the threshold.

### Error Handling

When a prompt fails to execute (engine error), it counts as an incorrect
classification — a false negative for should-trigger prompts or a false positive
for should-not-trigger prompts. The error count is reported separately so you can
distinguish engine failures from genuine misclassifications.

## Creating Custom Graders

Extend the `Grader` base class:

```python
from waza.graders.base import Grader, GraderContext, GraderType, GraderRegistry
from waza.schemas.results import GraderResult

@GraderRegistry.register("my_custom")
class MyCustomGrader(Grader):
    @property
    def grader_type(self) -> GraderType:
        return GraderType.CODE

    def grade(self, context: GraderContext) -> GraderResult:
        # Your logic here
        return GraderResult(
            name=self.name,
            type=self.grader_type.value,
            score=1.0,
            passed=True,
            message="Custom grading complete",
        )
```

Then use in eval.yaml:
```yaml
graders:
  - type: my_custom
    name: special_check
```
