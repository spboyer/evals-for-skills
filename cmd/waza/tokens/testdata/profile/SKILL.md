---
name: profile-test-skill
description: |
  **PROFILE TEST** - A skill for testing the profile command.
  USE FOR: testing profile analysis, verifying structural metrics.
  DO NOT USE FOR: production use.
---

# Profile Test Skill

A test skill used for verifying the `waza tokens profile` command.

## Activation Triggers

**Use this skill when:**
- Running profile analysis tests
- Verifying section counting
- Checking code block detection

## Workflow

1. **Parse the file** - Read SKILL.md and extract frontmatter
2. **Analyze structure** - Count sections, code blocks, steps
3. **Generate summary** - Produce one-line profile output
4. **Check warnings** - Flag potential issues

## Code Examples

### Example 1: Basic Usage

```bash
waza tokens profile my-skill
```

### Example 2: JSON Output

```json
{
  "name": "my-skill",
  "tokens": 500,
  "sections": 5
}
```

### Example 3: Table Output

```
File          Tokens   Sections
-------------------------------
SKILL.md         500          5
```

## Configuration

Set up your skill with proper frontmatter:

```yaml
name: my-skill
description: A useful skill
```

## Guidelines

- Keep token count under 2500
- Include workflow steps
- Use clear section headings
