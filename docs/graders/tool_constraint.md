### `tool_constraint` - Tool Usage & Resource Constraint Grader

Validates which tools an agent used (or avoided), plus turn and token limits.

```yaml
- type: tool_constraint
  name: guardrails
  config:
    expect_tools:
      - "bash"
      - "view"
    reject_tools:
      - "web_fetch"
    max_turns: 15
    max_tokens: 8000
```

**Options:**
| Option | Type | Description |
|--------|------|-------------|
| `expect_tools` | list[str] | Tool names that MUST have been called |
| `reject_tools` | list[str] | Tool names that MUST NOT have been called |
| `max_turns` | int | Maximum number of conversation turns allowed |
| `max_tokens` | int | Maximum total token usage allowed |

At least one option must be configured.

**Scoring:** `passed_checks / total_checks`

Each `expect_tools` and `reject_tools` entry counts as one check. `max_turns` and `max_tokens` each count as one check when configured.

**Difference from `behavior` grader:** The `behavior` grader covers `max_tool_calls`, `max_duration_ms`, `required_tools`, and `forbidden_tools`. The `tool_constraint` grader replaces the tool fields with `expect_tools`/`reject_tools` and adds `max_turns` instead of `max_tool_calls`/`max_duration_ms`. Use whichever matches the constraints you care about.
