### `program` - External Program Grader

Runs an external program or script to grade agent output. The agent output is passed via stdin, and the workspace directory is available as the `WAZA_WORKSPACE_DIR` environment variable.

```yaml
- type: program
  name: custom_validator
  config:
    command: python3
    args: ["graders/validate.py"]
    timeout: 60
```

**Options:**
| Option | Type | Description |
|--------|------|-------------|
| `command` | string | **Required.** The program to execute |
| `args` | list[str] | Arguments to pass to the program |
| `timeout` | int | Maximum execution time in seconds (default: 30) |

**Scoring:** Binary — exit code `0` means pass (`1.0`), non-zero means fail (`0.0`).

**Environment:**
- **stdin**: The agent's output text
- **`WAZA_WORKSPACE_DIR`**: Path to the post-execution workspace directory

**stdout** from the program is captured and used as the grader feedback message on success.

**Example: Shell script grader**

```yaml
- type: program
  name: lint_check
  config:
    command: bash
    args: ["-c", "cd $WAZA_WORKSPACE_DIR && npm run lint"]
    timeout: 120
```

**Example: Python grader reading workspace files**

```yaml
- type: program
  name: file_validator
  config:
    command: python3
    args: ["graders/check_output.py"]
```

```python
#!/usr/bin/env python3
import os, sys

output = sys.stdin.read()
workspace = os.environ.get("WAZA_WORKSPACE_DIR", "")

# Check agent output
if "error" in output.lower():
    print("Output contains errors", file=sys.stderr)
    sys.exit(1)

# Check workspace files
readme = os.path.join(workspace, "README.md")
if not os.path.exists(readme):
    print("README.md not found", file=sys.stderr)
    sys.exit(1)

print("All checks passed")
```
