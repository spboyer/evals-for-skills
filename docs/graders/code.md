### `code` - Assertion-Based Grader

Evaluates Python expressions against the execution context.

```yaml
- type: code
  name: my_grader
  config:
    assertions:
      - "len(output) > 0"
      - "'success' in output.lower()"
      - "len(errors) == 0"
```

**Available Context Variables:**
| Variable | Type | Description |
|----------|------|-------------|
| `output` | str | Final skill output |
| `outcome` | dict | Outcome state |
| `transcript` | list | Full execution transcript |
| `tool_calls` | list | Tool calls from transcript |
| `errors` | list | Errors from transcript |
| `duration_ms` | int | Execution duration |

**Available Functions:**
`len`, `any`, `all`, `str`, `int`, `float`, `bool`, `list`, `dict`, `re` (regex module)

**Scoring:** `passed_assertions / total_assertions`

**⚠️ Important:** Do NOT use generator expressions in assertions. They don't work with Python's `eval()` in restricted scope.

```yaml
# ❌ WRONG - generator expressions fail
assertions:
  - "any(kw in output for kw in ['azure', 'deploy'])"

# ✅ CORRECT - use explicit or chains
assertions:
  - "'azure' in output.lower() or 'deploy' in output.lower()"
```
