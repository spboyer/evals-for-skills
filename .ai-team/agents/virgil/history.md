# History — Virgil (Azure/Cloud Integration Specialist)

## Project Learnings (from import)

- **Project:** waza — CLI tool for evaluating Agent Skills
- **Owner:** Shayne Boyer (spboyer@live.com)
- **Stack:** Go (primary), TypeScript/React 19 dashboard
- **Repository:** spboyer/waza
- **Key files:**
  - `internal/models/outcome.go` — `EvaluationOutcome` struct (results JSON)
  - `cmd/waza/cmd_run.go` — `saveOutcome()` writes results to disk
  - `internal/webapi/store.go` — `FileStore` reads results for dashboard
  - `internal/projectconfig/config.go` — `ProjectConfig` struct (.waza.yaml)
  - `go.mod` — already has transitive Azure SDK deps via azd

## Learnings

## 2026-02-27: Dashboard Azure Storage Integration

**Branch:** squad/azure-storage-results
**Request:** Extend waza serve dashboard to read from Azure Storage when configured

### Implementation approach

1. **Adapter pattern**: Created `StorageAdapter` to bridge `storage.ResultStore` → `webapi.RunStore` interface. This preserves the existing webapi handlers while allowing storage backend swapping.

2. **Configuration flow**: `cmd_serve.go` loads `.waza.yaml` → passes `StorageConfig` to webserver → webserver creates appropriate store (Azure or local).

3. **Graceful degradation**: If Azure storage is configured but fails to initialize, server logs a warning and falls back to local FileStore. This ensures `waza serve` never fails due to Azure being unavailable.

4. **Storage status endpoint**: Added `GET /api/storage/status` returning `{ configured, provider, account }` so the dashboard can display which backend is active.

5. **Source metadata**: Added `source` field to `RunSummary` ("local" or "azure-blob") so API consumers know where data came from.

### Key learnings

- **Interface bridging**: The adapter pattern let us reuse existing webapi handlers without changing their logic. `StorageAdapter` translates `storage.ResultSummary` → `webapi.RunSummary`.

- **Context timeouts**: All storage adapter calls use 30-second timeouts to prevent hanging requests when Azure is slow or unavailable.

- **Summary performance trade-off**: Computing aggregate metrics requires downloading all outcomes (to get token counts). This is slower than FileStore's in-memory cache, but correct for remote storage.

- **Backward compatibility**: When storage is not configured, behavior is identical to before — local FileStore only. No breaking changes.

### Charter adherence

✓ DefaultAzureCredential (all Azure auth via storage layer)
✓ Graceful degradation (Azure failure → local fallback with warning)
✓ Context propagation (30s timeouts on all Azure calls)
✓ Error wrapping (storage errors map to HTTP error responses)

### Files touched

- `internal/webapi/storage_adapter.go` — new adapter bridging storage → webapi (124 lines)
- `internal/webapi/types.go` — added `Source` field to `RunSummary`, `StorageStatusResponse` type
- `internal/webapi/store.go` — set `Source: "local"` in `outcomeToSummary()`
- `internal/webapi/handlers.go` — added `HandleStorageStatus`, `StorageConfig` struct, `RegisterRoutesWithStorage`
- `internal/webserver/server.go` — added `StorageConfig` to webserver `Config`
- `internal/webserver/routes.go` — storage initialization logic with fallback, calls `RegisterRoutesWithStorage`
- `cmd/waza/cmd_serve.go` — load `.waza.yaml`, pass `StorageConfig` to webserver
- `.waza.yaml.example-storage` — example config showing how to enable Azure Storage

### Example .waza.yaml

```yaml
storage:
  enabled: true
  provider: azure-blob
  accountName: myaccount
  containerName: waza-results
```

When enabled, `waza serve` reads results from Azure Blob Storage. When not configured, falls back to local FileStore reading from `--results-dir`.

## 2026-02-26: Azure Blob Storage Backend Implementation

**Branch:** squad/azure-storage-results
**Issue context:** Linus created ResultStore interface + LocalStore implementation

### Key implementation choices

1. **Auto-login flow**: `DefaultAzureCredential` → if fails → `az login` → retry once. This avoids users hitting "credentials not found" errors.

2. **Blob organization**: `{skill}/{runid}.json` paths with metadata (skill, model, passrate, timestamp, runid). Metadata enables fast List() without downloading blobs.

3. **Download scan pattern**: Since we don't know skill from runID alone, Download() lists all blobs and matches by runid metadata. Trade-off: slower than direct path lookup, but correct.

4. **Error wrapping**: All Azure SDK errors wrapped with operation context: `fmt.Errorf("azure blob upload: %w", err)`.

5. **Context propagation**: All operations accept `context.Context` per Azure SDK patterns and Go best practices.

### Azure SDK learnings

- `azidentity.NewDefaultAzureCredential()` tries multiple credential sources: env vars → managed identity → az CLI → etc.
- No `CredentialUnavailableError` type exported — just check if error is non-nil and attempt login
- `azblob.UploadBuffer()` is simpler than stream-based upload for JSON payloads
- Blob metadata is `map[string]*string` (pointers) — use helper `stringPtr()` for values
- `ListBlobsFlatPager` is the iterator pattern for blob listing

### Charter adherence

✓ DefaultAzureCredential (no connection strings)
✓ Context-based APIs
✓ Error wrapping with context
✓ Graceful degradation (factory falls back to LocalStore on config errors)

### Files touched

- `internal/storage/azure_blob.go` — new AzureBlobStore implementation (349 lines)
- `internal/storage/store.go` — updated NewStore factory
- `.ai-team/decisions/inbox/virgil-azure-blob.md` — design decision doc


