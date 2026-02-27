# Storage Package Test Strategy

**By:** Basher (Tester / QA)
**Date:** 2026-02-27
**Branch:** squad/azure-storage-results
**Status:** IMPLEMENTED

## What

Comprehensive test coverage for the storage package:

1. **Extended `local_test.go`** — 9 new edge case tests
2. **New `azure_blob_test.go`** — 10 mock-based tests (skipped until implementation)
3. **Extended `projectconfig/config_test.go`** — 7 new StorageConfig tests

Total: **26 new tests** (31 total storage tests, 7 storage-related config tests)

## Test Coverage by Category

### Edge Cases for LocalStore (9 tests)

- **Special characters:** RunID sanitization (`run/with:special\chars` → `run_with_special_chars.json`)
- **Concurrency:** 10 goroutines uploading simultaneously (goroutine safety)
- **Large files:** 1000 test outcomes in a single result (serialization/deserialization)
- **Invalid data:** Malformed JSON files in results directory (graceful skip)
- **Non-result JSON:** Other JSON structures in same directory (filter by schema)
- **Compare edge cases:** Compare with non-existent run ID (ErrNotFound propagation)
- **Complex filters:** All filter options combined (skill + model + since + limit)
- **Empty fields:** EvaluationOutcome with minimal/zero values
- **Permissions:** Read-only directory (error handling)

### Azure Blob Tests (10 tests, skipped)

All tests use a mock blob client interface to avoid Azure dependencies. Tests are structured to compile and pass once `azure_blob.go` is implemented:

- **Upload serialization:** Blob path construction (`{skill}/{timestamp}/{runID}.json`)
- **Metadata stamping:** skill/model/timestamp as blob metadata
- **Download deserialization:** JSON unmarshaling from blob storage
- **List filtering:** Metadata-based filtering (skill, model)
- **Compare delta calculation:** Download two blobs and compute deltas
- **Auth failure retry:** `az login` retry flow on auth errors
- **Network error handling:** Timeout/connection error handling
- **Graceful degradation:** Fallback to local storage on Azure unavailability
- **Blob path construction:** Date-based path generation
- **Special characters in skill name:** URL encoding/sanitization
- **Context cancellation:** Respect context.Context for cancellation

### StorageConfig Tests (7 tests)

- **Defaults:** ContainerName defaults to "waza-results"
- **Full Azure config:** Provider/AccountName/ContainerName/Enabled
- **Backward compat:** Config without storage section (defaults applied)
- **Partial config:** Provider + Enabled, other fields use defaults
- **Disabled state:** Enabled=false with provider configured
- **Custom container:** User-specified container name
- **Merge behavior:** File values override defaults, missing fields use defaults

## Why

The storage package is mission-critical for result persistence and comparison. Edge cases like concurrent uploads, invalid JSON, and permission errors can cause silent data loss or crashes. Azure integration will introduce auth, network, and cloud-specific failure modes.

**Key risks addressed:**
- **Data corruption:** Invalid JSON/non-result files could crash List()
- **Concurrency bugs:** Multiple uploads could corrupt cache or overwrite files
- **Auth failures:** Azure auth errors could silently fail without fallback
- **Path injection:** Special characters in RunID/skill names could break file paths or blob URLs

## Test Patterns Used

- **Table-driven tests:** `sanitizeFilename`, blob path construction
- **t.TempDir():** All filesystem tests use isolated temp directories
- **t.Parallel():** Concurrent upload test runs in parallel
- **Mock interfaces:** Azure blob client mock allows testing without cloud dependencies
- **t.Skip():** Azure tests skip until implementation exists (self-documenting)
- **Error checking:** `errors.Is()` for ErrNotFound propagation
- **Context cancellation:** Tests verify context.Context is respected

## Impact

- **Coverage:** Storage package now has 31 tests (12 original + 19 new)
- **Readiness:** Azure blob tests are pre-written and ready to activate
- **Safety:** Edge cases caught before production (concurrency, invalid data, permissions)
- **Documentation:** Tests serve as examples for LocalStore and future AzureBlobStore usage

## Next Steps

Once `azure_blob.go` is implemented:
1. Remove `t.Skip()` from Azure tests
2. Integrate mock blob client into AzureBlobStore (constructor injection)
3. Verify all 10 Azure tests pass
4. Add integration tests with real Azure storage (separate from unit tests)
