### `human` - Manual Review Grader

Marks tasks for human review.

> **Status:** `human` is not currently implemented in this repository. This document is retained for reference.

```yaml
- type: human
  name: expert_review
  config:
    instructions: "Review for security best practices"
    criteria:
      - "Uses managed identity"
      - "No hardcoded secrets"
      - "Follows least privilege"
```

**Output:** Returns `pending` status until human submits review.
