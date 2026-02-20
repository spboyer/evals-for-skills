### `keyword` - Keyword Matching Grader

Checks for keyword presence or absence in the agent output using case-insensitive matching.

```yaml
- type: keyword
  name: topic_check
  config:
    must_contain:
      - "authentication"
      - "authorization"
    must_not_contain:
      - "hardcoded password"
      - "plaintext secret"
```

**Options:**
| Option | Type | Description |
|--------|------|-------------|
| `must_contain` | list[str] | Keywords that must appear in the output (case-insensitive) |
| `must_not_contain` | list[str] | Keywords that must NOT appear in the output (case-insensitive) |

**Scoring:** `passed_checks / total_checks`

**Difference from `regex` grader:** The keyword grader uses simple case-insensitive substring matching, while the regex grader uses regular expression patterns. Use `keyword` for simple word/phrase checks and `regex` when you need pattern matching.
