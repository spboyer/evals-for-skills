# Required Skills Validation Demo

This example demonstrates the `required_skills` preflight validation feature in waza.

## Overview

The `required_skills` field in the eval config allows you to specify which skills must be present before the evaluation starts. This is particularly useful for orchestration skills that call other skills during execution.

## Example Scenario

Consider an orchestration skill like `azure-deploy` that calls other skills:
- `azure-prepare` - prepares the Azure environment
- `azure-validate` - validates the Azure configuration

If these dependent skills aren't available in the `skill_directories`, the eval will run but the agent won't be able to invoke them, leading to confusing failures.

## Usage

### With All Required Skills Present

```yaml
config:
  skill_directories:
    - ./skills/azure-deploy
    - ./skills/azure-prepare
    - ./skills/azure-validate
  required_skills:
    - azure-deploy
    - azure-prepare
    - azure-validate
```

**Result:** ✓ Validation passes, eval runs normally

```
✓ Required skills validation passed (3/3 skills found)
```

### With Missing Required Skills

```yaml
config:
  skill_directories:
    - ./skills/azure-deploy  # Only has azure-deploy
  required_skills:
    - azure-deploy
    - azure-prepare
    - azure-validate
```

**Result:** ✗ Clear error before eval starts

```
Error: skill validation failed:
required skills not found:
  - azure-prepare
  - azure-validate

Searched directories:
  - ./skills/azure-deploy

Found skills:
  - azure-deploy
```

### Backward Compatible (No required_skills)

```yaml
config:
  skill_directories:
    - ./skills/azure-deploy
  # No required_skills specified
```

**Result:** ✓ Validation skipped, eval runs normally (backward compatible)

## Benefits

1. **Early Error Detection** - Catches missing skills before eval starts, not during execution
2. **Clear Error Messages** - Shows exactly what's missing and where waza looked
3. **Better Developer Experience** - Orchestration skill authors can document dependencies clearly
4. **Backward Compatible** - If `required_skills` is omitted, behavior is unchanged

## Implementation Details

- **Field:** `required_skills` in config section (array of strings)
- **Validation:** Runs after engine initialization, before loading test cases
- **Discovery:** Scans each `skill_directories` for SKILL.md files
- **Parsing:** Extracts skill name from SKILL.md frontmatter `name` field
- **Error Format:** Lists missing skills, searched directories, and found skills
- **Backward Compatible:** Empty or omitted `required_skills` skips validation

## Testing

Run the tests to see the feature in action:

```bash
go test ./internal/orchestration/... -run TestValidateRequiredSkills
```
