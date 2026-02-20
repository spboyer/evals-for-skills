### `llm` - LLM-as-Judge Grader

Uses an AI model to evaluate quality.

> **Status:** `llm` is not currently implemented in this repository. This document is retained for reference.

```yaml
- type: llm
  name: quality_judge
  config:
    model: gpt-4o-mini
    rubric: |
      Score the skill execution from 1-5:

      1. Correctness: Did it accomplish the task?
      2. Completeness: Were all requirements addressed?
      3. Quality: Was the approach appropriate?

      Return JSON: {"score": N, "reasoning": "...", "passed": true/false}
```

**Options:**
| Option | Type | Description |
|--------|------|-------------|
| `model` | str | Model to use (default: gpt-4o-mini) |
| `rubric` | str | Evaluation rubric (inline or file path) |
| `threshold` | float | Pass threshold (default: 0.75) |

**Score Normalization:** Raw scores 1-5 are normalized to 0-1:
- Score 1 → 0.0
- Score 3 → 0.5
- Score 5 → 1.0
