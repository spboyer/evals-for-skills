# Token Limits Configuration

Token limits are resolved in priority order:

1. **`.waza.yaml`** `tokens.limits` section — primary, workspace-level config
2. **`.token-limits.json`** in the skill directory — fallback for standalone use
3. **Built-in defaults** when neither config exists

You can also configure `tokens.warningThreshold` and `tokens.fallbackLimit` in `.waza.yaml`.
Patterns support workspace-root-relative paths (e.g., `plugin/skills/**/SKILL.md`).

## .waza.yaml tokens section (recommended)

```yaml
tokens:
  warningThreshold: 2500
  fallbackLimit: 2000
  limits:
    defaults:
      "SKILL.md": 500
      "references/**/*.md": 1000
      "*.md": 2000
    overrides:
      "README.md": 3000
```

## .token-limits.json (fallback)

The `.token-limits.json` file can still define token budgets per skill directory. It is checked when `.waza.yaml` does not provide limits.

## File Structure

```json
{
  "description": "Optional human-readable description",
  "defaults": {
    "<pattern>": <limit>,
    ...
  },
  "overrides": {
    "<exact-path>": <limit>,
    ...
  }
}
```

### Fields

| Field | Required | Type | Description |
|-------|----------|------|-------------|
| `description` | No | string | Documentation for humans; ignored by the CLI |
| `defaults` | Yes | object | Glob patterns mapped to token limits |
| `overrides` | No | object | Exact file paths mapped to token limits |

## Pattern Matching

### Overrides (Exact Matches)

Entries in `overrides` match files by suffix. The path is matched against the end of the normalized file path:

```json
"overrides": {
  "README.md": 4000,
  "docs/API.md": 3000
}
```

- `README.md` matches `./README.md` and `subdir/README.md`
- `docs/API.md` matches `./docs/API.md` and `subdir/docs/API.md`

Overrides are checked before defaults and take precedence.

### Defaults (Glob Patterns)

Entries in `defaults` use glob patterns:

| Pattern | Matches |
|---------|---------|
| `*.md` | Any `.md` file at any depth (e.g., `README.md`, `foo/bar.md`) |
| `SKILL.md` | Files named exactly `SKILL.md` in any directory |
| `references/*.md` | `.md` files directly in `references/` |
| `references/**/*.md` | `.md` files in subdirectories of `references/` (not directly in it) |
| `docs/**/*.md` | `.md` files in subdirectories of `docs/` (not directly in it) |

### Glob Syntax

| Syntax | Meaning |
|--------|---------|
| `*` | Matches any characters except `/` |
| `**` | Matches any characters including `/` (recursive); in `**/*` the `/*` requires at least one subdirectory level |
| `/` | Directory separator; patterns containing `/` are anchored to the project root |
| `.` | Literal dot (automatically escaped) |

### Pattern Specificity

When multiple patterns match a file, the most specific pattern wins. Specificity is calculated as:

1. **Exact match** (no wildcards): +10000 points
2. **Path depth**: +100 points per `/` in the pattern
3. **Single wildcards** (`*`): +10 points each
4. **Globstars** (`**`): -50 points each
5. **Pattern length**: +1 point per character

Example resolution for `references/test-templates/jest.md`:

| Pattern | Specificity | Result |
|---------|-------------|--------|
| `*.md` | Low | Fallback |
| `references/*.md` | Medium | Doesn't match (file is nested) |
| `references/**/*.md` | Medium-High | Matches |
| `references/test-templates/*.md` | Higher | Matches, **wins** |

## Complete Example

```json
{
  "description": "Token limits for skill repository",
  "defaults": {
    "SKILL.md": 5000,
    "references/*.md": 2000,
    "references/**/*.md": 2000,
    "references/test-templates/*.md": 1500,
    "*.md": 4000
  },
  "overrides": {
    "README.md": 4000
  }
}
```

### Resolution Examples

| File | Matching Pattern | Limit |
|------|------------------|-------|
| `SKILL.md` | `SKILL.md` (defaults) | 5000 |
| `README.md` | `README.md` (overrides) | 4000 |
| `references/scoring.md` | `references/*.md` | 2000 |
| `references/test-templates/jest.md` | `references/test-templates/*.md` | 1500 |
| `assets/guide.md` | `*.md` | 4000 |

## Default Configuration

When neither `.waza.yaml` nor `.token-limits.json` provides limits, these built-in defaults apply:

```json
{
  "defaults": {
    "SKILL.md": 500,
    "references/**/*.md": 1000,
    "docs/**/*.md": 1500,
    "*.md": 2000
  },
  "overrides": {
    "README.md": 3000,
    "CONTRIBUTING.md": 2500
  }
}
```

## Validation Behavior

The `waza tokens check` command:

1. Discovers all markdown files (`.md`, `.mdx`)
2. Excludes `node_modules`, `.git`, `dist`, `coverage` directories
3. For each file, finds the applicable limit using override → defaults
4. Reports files exceeding their limits
5. With `--strict`, exits with code 1 if any file exceeds its limit
