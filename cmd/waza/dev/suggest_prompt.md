## Instructions

You are analyzing an Agent Skill to identify improvements that would help an agent decide whether to use this Skill given an appropriate user prompt.

Read the entire Skill so you understand it well, but focus your suggestions on the Skill's frontmatter. Frontmatter is critical because it's the only content agents read when deciding whether to use the Skill.

Frontmatter is a YAML block the beginning of the Skill content and looks like this:

```yaml
---
name: code-explainer
description: |
  **UTILITY SKILL** - Explain code snippets, functions, and algorithms in plain language.
  USE FOR: explain code, what does this code do, break down this function,
  help me understand this, walk through this algorithm, clarify this logic,
  explain this snippet, describe what happens here.
  DO NOT USE FOR: writing new code (use code generation), fixing bugs (use debugging),
  refactoring (use refactoring skills), code review with action items.
  INVOKES: file reading tools to access code, language detection for tailored explanations.
  FOR SINGLE OPERATIONS: If the user just needs to see file contents, use file reading tools directly.
---
```

Consider the following potential improvements:
1. Making trigger phrases more specific and varied; does the Skill body suggest new trigger phrases to incorporate into the frontmatter?
2. Strengthening the description to clearly communicate the skill's purpose
3. Adding or improving USE FOR / DO NOT USE FOR sections
4. Improving routing clarity by mentioning tools the Skill may invoke (INVOKES, FOR SINGLE OPERATIONS)
5. Removing ambiguity that could cause false positives or missed triggers

The description must not be too long (~200 words is too long).

Return a numbered list of specific, actionable suggestions. Don't comment on or editorialize about the task; simply return suggestions. Each suggestion should describe exactly what to change and why it will improve trigger accuracy.

## Skill to analyze
