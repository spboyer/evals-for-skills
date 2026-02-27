# Decision: Results CLI & Auto-Upload Wiring

**By:** Linus (Backend Developer)
**Date:** 2026-02-27
**Branch:** `squad/azure-storage-results`

## What

Wired the Phase 1 `ResultStore` interface into the CLI:

1. **Auto-upload in `cmd_run.go`:** After all outcomes are saved locally, `autoUploadOutcomes()` checks `cfg.Storage` and uploads via `ResultStore.Upload()`. Errors are non-fatal warnings — local results are always preserved. Works for both single-skill and multi-skill runs.

2. **`waza results list` command:** Lists stored evaluation results with filtering by `--skill`, `--model`, `--since`, and `--limit`. Displays a formatted table with Run ID, Skill, Model, Pass Rate, and Timestamp.

3. **`waza results compare <id1> <id2>` command:** Compares two stored runs showing pass rate delta, score delta, and per-metric deltas with color-coded indicators (green for improvements, red for regressions).

## Key Design Choices

- **Auto-upload is fire-and-forget:** Upload failures are warnings, not errors. The run command's exit code is never affected by storage failures.
- **`NewStore` factory routes on config:** Both `results` commands and auto-upload use the same `storage.NewStore()` factory — no separate store construction paths.
- **Results commands require `storage.enabled: true`:** If no storage section is in `.waza.yaml`, the commands return a clear error message pointing users to configure it.
- **Table output uses fixed-width columns with truncation** for readability in terminal output.

## Why

These are the user-facing integration points that make the Phase 1 storage layer useful. Without CLI wiring, the `ResultStore` interface sits unused.
