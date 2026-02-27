# Decision: Azure Blob Storage Backend Implementation

**By:** Virgil (Azure/Cloud Integration Specialist)
**Date:** 2026-02-26
**Status:** IMPLEMENTED
**Branch:** squad/azure-storage-results

## What

Implemented Azure Blob Storage backend for waza result storage (`internal/storage/azure_blob.go`).

### Key Design Decisions

1. **DefaultAzureCredential with auto-login fallback**: If `NewDefaultAzureCredential()` fails, automatically run `az login` and retry once. This provides a seamless developer experience — no manual login needed.

2. **Blob path structure**: `{skill-name}/{run-id}.json` — organizes results by skill, making prefix-based filtering efficient.

3. **Metadata-based listing**: Store skill, model, passrate, timestamp, and runid as blob metadata. List operations read metadata only, avoiding full blob downloads for summary views.

4. **Download by metadata scan**: Since we don't know the skill name from just a run ID, Download lists all blobs and finds the match by runid metadata. This is slower than a direct path lookup but handles the general case.

5. **Context-based API**: All operations accept `context.Context` for cancellation and timeouts, following Azure SDK and Go best practices.

6. **Error wrapping**: All Azure errors are wrapped with `fmt.Errorf("azure blob <operation>: %w", err)` for clear error context.

## Why

**Problem:** Linus created the `ResultStore` interface and local filesystem implementation. We needed a cloud-backed implementation for remote storage.

**Requirements from charter:**
- Always use `DefaultAzureCredential` — never connection strings ✓
- Wrap Azure errors with context ✓
- All cloud operations accept `context.Context` ✓
- Graceful degradation: cloud failures warn but never fail local operations (handled by factory fallback to LocalStore)

**Auto-login rationale:** Developers often forget to run `az login`. Auto-triggering it once on credential failure reduces friction.

**Metadata rationale:** Blob metadata is returned by list operations without downloading blob bodies. This makes List() fast even with thousands of results.

## Impact

- **Files changed:**
  - `internal/storage/azure_blob.go` — new AzureBlobStore implementation (349 lines)
  - `internal/storage/store.go` — updated NewStore factory to create AzureBlobStore when provider is "azure-blob"

- **Dependencies:** Uses existing Azure SDK imports (already transitive deps via azd in go.mod)

- **Testing:** Build and `go vet` pass. Manual testing needed for actual Azure operations (requires storage account).

## Follow-up

- Add unit tests with mock azblob.Client
- Add integration test against Azurite (local Azure Storage emulator)
- Document .waza.yaml storage configuration in site docs
- Consider caching List results (similar to LocalStore's load pattern) if performance becomes an issue

