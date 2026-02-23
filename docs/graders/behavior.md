### `behavior` - Agent Behavior Validation

Validates agent behavior patterns like tool call counts, token usage, and execution duration.

```yaml
- type: behavior
  name: efficiency_check
  config:
    max_tool_calls: 20
    max_tokens: 10000
    max_duration_ms: 300000
    required_tools:
      - "bash"
      - "view"
    forbidden_tools:
      - "create"
      - "web_fetch"
```

**Options:**
| Option | Type | Description |
|--------|------|-------------|
| `max_tool_calls` | int | Maximum allowed tool calls |
| `max_tokens` | int | Maximum token usage allowed |
| `max_duration_ms` | int64 | Maximum execution time in milliseconds |
| `required_tools` | list[str] | Tool names (exact matches) that MUST be called |
| `forbidden_tools` | list[str] | Tool names (exact matches) that MUST NOT be called |

**Note:** `required_tools` and `forbidden_tools` use exact string matching on tool names; patterns, wildcards, or regular expressions are not supported.

At least one option must be configured.

**Scoring:** `passed_checks / total_checks`

Each configured option counts as one check, except `required_tools` and `forbidden_tools` which contribute one check per entry.
