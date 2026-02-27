## ResultStore Interface Design

**By:** Linus (Backend Developer)
**Date:** 2026-02-27
**Branch:** `squad/azure-storage-results`

### What

Created the `ResultStore` interface in `internal/storage/` with four operations: `Upload`, `List`, `Download`, `Compare`. All methods accept `context.Context` for cloud-backend compatibility. Two implementations planned:

1. **LocalStore** (implemented) — wraps the filesystem read/write pattern from `webapi.FileStore`, adds upload (write JSON), list with filtering (skill, model, since, limit), download, and compare with metric deltas.
2. **AzureBlobStore** (stub) — `NewStore` factory returns an error for `azure-blob` provider until Virgil implements it.

Added `StorageConfig` to `ProjectConfig` with `provider`, `accountName`, `containerName`, and `enabled` fields. Default container name is `"waza-results"`.

### Why

The storage layer needs a clean interface so `cmd_run.go` and `webapi` can persist results to either local disk or Azure Blob Storage through the same API. Building the interface and local adapter first lets us validate the contract before wiring up the cloud backend.

### Key Design Choices

1. **Separate from webapi.FileStore** — `ResultStore` is a different concern. `FileStore` is read-only and serves the dashboard API. `ResultStore` handles write + read + compare for the CLI pipeline. They can coexist and eventually `FileStore` could delegate to `LocalStore`.

2. **Context on all methods** — even `LocalStore` methods accept `context.Context` (and ignore it). This keeps the interface honest for the Azure implementation where context controls timeouts and cancellation.

3. **Factory pattern via `NewStore`** — single entry point decides local vs. cloud based on `StorageConfig`. Azure path errors cleanly until implemented.

4. **Lazy loading** — `LocalStore` uses the same `ensureLoaded()` pattern as `FileStore` to defer disk reads until first access.
