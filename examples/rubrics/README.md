# Tool Call Evaluation Rubrics

Pre-built rubric configurations for evaluating how well an AI agent uses tools.
These are adapted from [Azure ML's built-in evaluators](https://github.com/Azure/azureml-assets/tree/main/assets/evaluators/builtin)
for use with waza's `prompt` grader.

> **Note:** These rubrics require the `prompt` grader ([#104](https://github.com/spboyer/waza/issues/104)),
> which is not yet merged. The YAML files are ready to use once that grader lands.

## Rubrics

| File | What it evaluates | Scale |
|------|-------------------|-------|
| `tool_call_accuracy.yaml` | Overall tool call effectiveness (relevance, parameters, completeness, efficiency, execution) | 1–5 ordinal → 0.0–1.0 |
| `tool_selection.yaml` | Whether the right tools were chosen and no essential tools were missed | Binary pass/fail → 0.0/1.0 |
| `tool_input_accuracy.yaml` | Whether all parameters passed to tools are correct (grounded, typed, formatted) | Binary pass/fail → 0.0/1.0 |
| `tool_output_utilization.yaml` | Whether the agent correctly uses data returned by tools in its response | Binary pass/fail → 0.0/1.0 |

## How they work together

These four rubrics decompose tool-use quality into orthogonal dimensions:

```
tool_call_accuracy (umbrella)
├── tool_selection        — Did the agent pick the right tools?
├── tool_input_accuracy   — Did it pass the right parameters?
└── tool_output_utilization — Did it use the results correctly?
```

`tool_call_accuracy` is the composite evaluator — it covers all three concerns
in a single 1–5 score. The other three are focused evaluators that give you
granular pass/fail signals on each dimension.

## YAML structure

Each rubric file includes:

| Field | Purpose |
|-------|---------|
| `name` | Identifier for the rubric |
| `description` | What it evaluates |
| `version` | Semantic version of the rubric definition |
| `evaluation_criteria` | Detailed instructions for the LLM judge |
| `rating_levels` | Scoring scale with descriptions and examples |
| `score_normalization` | How raw scores map to the 0.0–1.0 range waza expects |
| `input_mapping` | How waza's `graders.Context` fields feed into the rubric |
| `chain_of_thought` | Step-by-step reasoning instructions for the LLM judge |
| `output_schema` | Expected JSON structure from the LLM judge |

## Usage with the prompt grader

Once the `prompt` grader is available, reference a rubric in your eval spec:

```yaml
# eval.yaml
tasks:
  - name: "tool-use-quality"
    prompt: "Book a flight from SEA to JFK for next Friday"
    graders:
      - type: prompt
        rubric: rubrics/tool_call_accuracy.yaml

      # Or use the focused evaluators for granular scoring:
      - type: prompt
        rubric: rubrics/tool_selection.yaml
      - type: prompt
        rubric: rubrics/tool_input_accuracy.yaml
      - type: prompt
        rubric: rubrics/tool_output_utilization.yaml
```

## Input mapping

The rubrics expect these fields from waza's grading context:

| Rubric field | Waza context field | Description |
|-------------|-------------------|-------------|
| `query` / `session_transcript` | `session_transcript` | The full conversation history |
| `tool_calls` | `tool_calls` | Tool calls made by the agent |
| `tool_definitions` | `tool_definitions` | Available tool schemas |
| `response` / `agent_response` | `agent_response` | The agent's latest message (tool_output_utilization only) |

## Score normalization

All rubrics normalize to waza's standard 0.0–1.0 range:

- **Binary rubrics** (tool_selection, tool_input_accuracy, tool_output_utilization): `0 → 0.0`, `1 → 1.0`
- **Ordinal rubric** (tool_call_accuracy): `(raw - 1) / 4` maps `1 → 0.0`, `5 → 1.0`

## Attribution

These rubrics are adapted from the Azure ML evaluator prompty templates, licensed
under the MIT License. See the [source repository](https://github.com/Azure/azureml-assets)
for the original implementations.
