### `llm_comparison` - Reference Comparison Grader

Compares output against a reference using LLM.

> **Status:** `llm_comparison` is not currently implemented in this repository. This document is retained for reference.

```yaml
- type: llm_comparison
  name: reference_check
  config:
    model: gpt-4o-mini
    reference: |
      Expected output should include:
      - Confirmation of deployment
      - URL of deployed resource
      - Next steps for the user
```
