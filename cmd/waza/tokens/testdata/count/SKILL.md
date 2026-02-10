```skill
---
name: test-skill
description: |
  **TEST SKILL** - A sample skill for testing waza token counting.
  USE FOR: testing token count command, verifying markdown scanning.
  DO NOT USE FOR: actual agent operations.
---

# Test Skill

A demonstration skill used for testing the `waza tokens` command suite.

## Skill Activation Triggers

**Use this skill when:**
- Running unit tests for token counting
- Verifying file scanning behavior
- Testing output formats (table, JSON)

## Instructions

1. **Scan files** - Find all markdown files in the skill directory
2. **Count tokens** - Estimate token usage for each file
3. **Report results** - Output in requested format

## Output Format

Structure results as:

| File | Tokens | Characters | Lines |
|------|--------|------------|-------|
| SKILL.md | 150 | 600 | 40 |
| README.md | 50 | 200 | 15 |

## Examples

### Example 1: Count All Files

```bash
waza tokens count
```

### Example 1b: Count a Specific Directory

```bash
waza tokens count ~/my-skill
waza tokens count references
```

**Output:**
```
File          Tokens   Chars  Lines
------------------------------------
README.md         50     200     15
SKILL.md         150     600     40
```

### Example 2: JSON Format

```bash
waza tokens count --format json
```

**Output:**
```json
{
  "totalTokens": 200,
  "totalFiles": 2,
  "files": {
    "README.md": {"tokens": 50},
    "SKILL.md": {"tokens": 150}
  }
}
```

## Behavior Guidelines

- Only count markdown files (.md, .mdx)
- Skip binary and non-text files
- Respect .gitignore patterns
- Use streaming for large file sets
```
