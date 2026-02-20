### `tool_calls` - Tool Usage Grader

Validates which tools were called and how.

> **Status:** `tool_calls` is not currently implemented in this repository. This document is retained for reference.

```yaml
- type: tool_calls
  name: tool_validator
  config:
    required:
      - pattern: "azd up"
      - pattern: "git commit"
    forbidden:
      - pattern: "rm -rf"
      - pattern: "sudo"
    max_calls: 20
```

**Options:**
| Option | Type | Description |
|--------|------|-------------|
| `required` | list | Patterns that MUST appear in tool calls |
| `forbidden` | list | Patterns that MUST NOT appear |
| `max_calls` | int | Maximum allowed tool calls |
