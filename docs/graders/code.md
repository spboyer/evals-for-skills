### `code` - Assertion-Based Grader

Evaluates expressions against the execution context using an inline script runner.

```yaml
- type: code
  name: my_grader
  config:
    assertions:
      - "len(output) > 0"
      - "'success' in output.lower()"
      - "len(errors) == 0"
```

**Options:**
| Option | Type | Description |
|--------|------|-------------|
| `assertions` | list[str] | Expressions to evaluate |
| `language` | string | Script language: `python` (default) or `javascript` |

**Available Context Variables (Python):**
| Variable | Type | Description |
|----------|------|-------------|
| `output` | str | Final skill output |
| `outcome` | dict | Outcome state |
| `transcript` | list | Full execution transcript |
| `tool_calls` | list | Tool calls from transcript |
| `errors` | list | Transcript events containing errors |
| `duration_ms` | int | Execution duration |

**Available Functions (Python):**
`len`, `any`, `all`, `str`, `int`, `float`, `bool`, `list`, `dict`, `re` (regex module)

**Available Context Variables (JavaScript):**
The same variables (`output`, `outcome`, `transcript`, `tool_calls`, `errors`, `duration_ms`) are available, plus built-in JS globals: `Array`, `Object`, `String`, `Number`, `Boolean`, `Math`, `JSON`, `RegExp`, `parseInt`, `parseFloat`.

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
