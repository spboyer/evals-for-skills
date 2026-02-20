# Linus — Backend Developer History

## Learnings

### 2025-01-26: FileStore Recursive Scanning and Validation (PR #280)

**Context:** The FileStore was using `os.ReadDir()` which only scans the top-level directory. When users ran `waza run --output-dir ./results` which creates subdirectories like `results/code-explainer/gpt-4o.json`, the serve command couldn't find any results. Additionally, `summary.json` files were being parsed as phantom evaluation runs because `json.Unmarshal` silently produces zero-value structs for unknown fields.

**Implementation:**
- Replaced flat `os.ReadDir()` with recursive `filepath.WalkDir()` to scan all subdirectories
- Added post-unmarshal validation: valid outcomes must have either `BenchName != ""` or `Digest.TotalTests > 0`
- Changed RunID fallback to use relative path from results-dir (e.g., `code-explainer/gpt-4o`) to avoid collisions

**Tests added:**
- Recursive scanning with nested directory structure
- summary.json filtering validation
- Invalid/empty/malformed JSON handling
- Deep nesting support (a/b/c/result.json)
- RunID collision prevention with same filenames in different subdirs

**Key insight:** Go's `json.Unmarshal` is permissive by design—it ignores unknown fields and produces zero-value structs. For file scanning code, always validate that unmarshaled data actually represents the expected schema before accepting it. A simple check for non-zero required fields prevents phantom entries.

**Pattern established:** Use `filepath.WalkDir()` for directory traversal (not `os.ReadDir()` with manual recursion), and always validate unmarshaled JSON beyond just checking for unmarshal errors.
