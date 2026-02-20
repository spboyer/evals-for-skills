### `prompt` - LLM-Based Evaluation

Uses a language model to evaluate skill execution quality based on a rubric. This grader follows the Azure ML evaluator pattern for LLM-as-judge evaluation.

> **Note:** This grader requires implementation and is currently planned for a future release.

```yaml
- type: prompt
  name: quality_judge
  config:
    model: gpt-4o-mini
    rubric: |
      Evaluate the code explanation on these criteria:

      1. **Correctness** (1-5): Is the explanation technically accurate?
      2. **Completeness** (1-5): Are all key concepts covered?
      3. **Clarity** (1-5): Is it easy to understand?

      Provide:
      - A score for each criterion (1-5)
      - Overall assessment (1-5)
      - Brief reasoning

      Return JSON:
      {
        "correctness": <1-5>,
        "completeness": <1-5>,
        "clarity": <1-5>,
        "overall_score": <1-5>,
        "reasoning": "...",
        "passed": <true/false>
      }
    threshold: 0.75
    score_type: normalized
    response_format: json
```

**Options:**
| Option | Type | Description |
|--------|------|-------------|
| `model` | string | Model to use for evaluation (default: gpt-4o-mini) |
| `rubric` | string | Evaluation rubric (inline text or file path) |
| `threshold` | float | Pass threshold (default: 0.75) |
| `score_type` | string | How to interpret scores: `normalized` (1-5 → 0-1) or `raw` |
| `response_format` | string | Expected response format: `json` or `text` |

**How Rubrics Work:**

A rubric is a structured evaluation criteria that guides the LLM's assessment:

1. **Define Criteria**: List specific aspects to evaluate (correctness, completeness, clarity, etc.)
2. **Rating Scale**: Specify the scale (typically 1-5 or 1-10)
3. **Guidelines**: Describe what each rating level means
4. **Output Format**: Request structured output (JSON preferred) with scores and reasoning

**Example: Multi-Criteria Rubric**

```yaml
- type: prompt
  name: comprehensive_quality
  config:
    model: gpt-4o
    rubric: |
      Evaluate the agent's performance:

      **Criteria:**
      1. Task Completion (1-5): Did the agent accomplish the stated goal?
         - 1: Failed completely
         - 3: Partially completed
         - 5: Fully completed with excellence

      2. Approach Quality (1-5): Was the approach appropriate and efficient?
         - 1: Poor approach with significant issues
         - 3: Adequate but could be improved
         - 5: Excellent, optimal approach

      3. Code Quality (1-5): Is the code well-structured and maintainable?
         - 1: Poor structure, hard to maintain
         - 3: Acceptable quality
         - 5: Excellent quality, follows best practices

      **Output Format (JSON):**
      {
        "task_completion": <score>,
        "approach_quality": <score>,
        "code_quality": <score>,
        "overall_score": <average of above>,
        "reasoning": "<2-3 sentence explanation>",
        "passed": <true if overall_score >= 3.5>
      }

      Think step-by-step and provide honest, critical evaluation.
    threshold: 0.7
    score_type: normalized
```

**Writing Custom Rubrics:**

Follow these best practices:

1. **Be Specific**: Define exactly what you're evaluating
2. **Use Clear Scales**: Explain what each rating means
3. **Request Reasoning**: Ask for chain-of-thought explanation
4. **Structured Output**: Use JSON for reliable parsing
5. **Include Pass Threshold**: Make the LLM decide pass/fail based on criteria
6. **Avoid Ambiguity**: Be explicit about edge cases

**Common Rubric Patterns:**

```yaml
# Binary pass/fail evaluation
rubric: |
  Does the output meet these requirements?
  1. Contains user authentication
  2. Follows security best practices
  3. Includes error handling

  Return JSON: {"score": 1 or 0, "passed": true/false, "reasoning": "..."}

# Comparative evaluation
rubric: |
  Compare the agent's solution to this reference approach:
  [reference description]

  Rate similarity (1-5) and quality improvement (1-5).
  Return JSON with scores and reasoning.

# Style compliance check
rubric: |
  Evaluate code style compliance:
  - Naming conventions
  - Documentation completeness
  - Code organization

  Return JSON with per-criterion scores and overall assessment.
```

**Shipped Rubric Templates:**

Waza ships with 8 pre-built rubric templates in [`examples/rubrics/`](../examples/rubrics/), adapted from [Azure ML evaluators](https://github.com/Azure/azureml-assets):

| Rubric | Category | Scale | Description |
|--------|----------|-------|-------------|
| `tool_call_accuracy` | Tool call | 1–5 ordinal | Overall tool call effectiveness |
| `tool_selection` | Tool call | Binary | Whether the right tools were chosen |
| `tool_input_accuracy` | Tool call | Binary | Whether parameters are correct |
| `tool_output_utilization` | Tool call | Binary | Whether tool results are used correctly |
| `task_completion` | Task eval | Binary | Whether the task was fully completed |
| `task_adherence` | Task eval | Binary flag | Goal, rule, and procedural adherence |
| `intent_resolution` | Task eval | 1–5 ordinal | How well user intent was resolved |
| `response_completeness` | Task eval | 1–5 ordinal | How thoroughly the response covers ground truth |

See [`examples/rubrics/README.md`](../examples/rubrics/README.md) for usage details and input mapping.

---
