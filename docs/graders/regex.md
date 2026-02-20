### `regex` - Pattern Matching Grader

Matches output against regex patterns.

```yaml
- type: regex
  name: format_checker
  config:
    must_match:
      - "deployed to https?://.+"
      - "Resource group: .+"
    must_not_match:
      - "error|failed|exception"
      - "permission denied"
```

**Options:**
| Option | Type | Description |
|--------|------|-------------|
| `must_match` | list[str] | Patterns that MUST appear |
| `must_not_match` | list[str] | Patterns that MUST NOT appear |

**Scoring:** `passed_checks / total_checks`
