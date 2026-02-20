### `action_sequence` - Tool Call Sequence Validation

Validates that the agent's tool calls match an expected action sequence. Supports three matching modes and calculates precision, recall, and F1 scores.

```yaml
- type: action_sequence
  name: deployment_workflow
  config:
    matching_mode: in_order_match
    expected_actions:
      - "bash"
      - "edit"
      - "bash"
      - "report_progress"
```

**Options:**
| Option | Type | Description |
|--------|------|-------------|
| `matching_mode` | string | How to match sequences (see modes below) |
| `expected_actions` | list[str] | List of expected tool names in sequence |

**Matching Modes:**

1. **`exact_match`** - Perfect match required
   - Actual sequence must match expected sequence exactly
   - Same length, same order, same tools
   - Example: Expected `["bash", "edit"]` only matches actual `["bash", "edit"]`

2. **`in_order_match`** - Actions must appear in order
   - All expected actions must appear in actual sequence
   - Can have extra actions between expected ones
   - Order must be preserved
   - Example: Expected `["bash", "edit"]` matches actual `["bash", "view", "edit", "report_progress"]`

3. **`any_order_match`** - All actions present regardless of order
   - All expected actions must appear in actual sequence
   - Order doesn't matter
   - Frequency must match (if expected has 2x "bash", actual must have at least 2x "bash")
   - Example: Expected `["edit", "bash"]` matches actual `["bash", "view", "edit"]`

**Scoring:**

The grader calculates three metrics:
- **Precision**: `true_positives / len(actual_actions)` - What fraction of actual actions were expected?
- **Recall**: `true_positives / len(expected_actions)` - What fraction of expected actions were performed?
- **F1 Score**: `2 * precision * recall / (precision + recall)` - Harmonic mean (used as the final score)

The `passed` field is based on the matching mode constraint, while the `score` field always uses F1.

**Example Use Cases:**

```yaml
# Ensure exact workflow for reproducible demos
- type: action_sequence
  name: demo_script
  config:
    matching_mode: exact_match
    expected_actions: ["bash", "view", "edit", "bash", "report_progress"]

# Verify key steps happen in order (allows flexibility)
- type: action_sequence
  name: deployment_flow
  config:
    matching_mode: in_order_match
    expected_actions: ["bash", "edit", "report_progress"]

# Check that required tools were used (any order)
- type: action_sequence
  name: required_tools
  config:
    matching_mode: any_order_match
    expected_actions: ["bash", "view", "edit"]
```
