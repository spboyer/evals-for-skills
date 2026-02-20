### `script` - External Script Grader

Runs a custom Python script for complex validation.

> **Status:** `script` is not currently implemented in this repository. This document is retained for reference.

```yaml
- type: script
  name: custom_logic
  config:
    script: graders/my_grader.py
```

**Script Format:**
```python
#!/usr/bin/env python3
import json
import sys

def grade(context: dict) -> dict:
    output = context.get("output", "")

    # Your custom logic here
    score = 1.0 if "success" in output else 0.0

    return {
        "score": score,
        "passed": score >= 0.5,
        "message": "Custom grading complete",
        "details": {"custom_field": "value"}
    }

if __name__ == "__main__":
    context = json.load(sys.stdin)
    print(json.dumps(grade(context)))
```
