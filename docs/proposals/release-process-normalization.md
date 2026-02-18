# Release Process Normalization

**Issue:** #223 — Normalize release process — standalone CLI + azd extension + consolidated workflows

**Proposed By:** Rusty (Architecture Lead)

**Date:** 2026-02-20

---

## Executive Summary

The waza project currently has **fragmented release automation** across three separate workflows (`go-release.yml`, `azd-ext-release.yml`, and `squad-release.yml`) with overlapping concerns, manual version management, and inconsistent triggers. This proposal consolidates the release process into a **single unified workflow** that:

- Produces both **standalone CLI binaries** (cross-platform) and **azd extension packages**
- Keeps `version.txt`, `extension.yaml`, and `registry.json` **automatically in sync**
- Provides **clear separation** between PR validation (CI) and release automation (CD)
- Enables both **tag-push triggers** and **manual dispatch** with clear semantics
- Eliminates duplicate build logic and version file management

---

## Current State Analysis

### 1. Existing Workflows

| Workflow | Trigger | Purpose | Version Source | Output | Issues |
|----------|---------|---------|-----------------|--------|--------|
| `go-ci.yml` | PR + push to main/develop | Go unit tests, lint, build, integration test | N/A (CI only) | Binary + coverage | ✅ Working |
| `go-release.yml` | Tag push (v*.*.*) or manual dispatch | Build standalone CLI binaries (6 platforms) | Git tag or manual input | GitHub Release + 6 binaries | ❌ No extension |
| `azd-ext-release.yml` | Push to main on `version.txt` or `extension.yaml` change or manual dispatch | Build + pack + publish azd extension | `version.txt` | Registry.json PR (auto-merged) | ❌ Separate from CLI release |
| `squad-release.yml` | Push to main | Release squad frontend | `package.json` | GitHub Release | ✅ Independent (squad) |

### 2. Version File Status

| File | Current Version | Owner | Purpose | Updated By |
|------|-----------------|-------|---------|------------|
| `version.txt` | `0.4.0-alpha.1` | azd extension release workflow | Extension versioning | Automated (azd x publish) |
| `extension.yaml` | `0.3.0` | Manual | Extension manifest | ⚠️ Manual edit required |
| `registry.json` | Multiple versions (latest: `0.3.0` in v3 entry) | azd extension workflow | Extension registry (azd discovery) | Auto-updated by `azd x publish` |
| (Implicit in Go) | From git tag or manual input | go-release workflow | CLI binary versioning | Tag-based |

**Problems:**
1. **Version mismatch:** `extension.yaml` (0.3.0) != `version.txt` (0.4.0-alpha.1) != latest registry entry (0.3.0)
2. **Dual workflows:** CLI and extension released separately, no coordination
3. **Manual version bumps:** `extension.yaml` requires human update
4. **No single source of truth:** Three different version indicators
5. **Path-based triggers for extension:** Fragile (any edit to version.txt triggers full extension release)

### 3. Build Scripts

| Script | Purpose | Platforms | Notes |
|--------|---------|-----------|-------|
| `build.sh` | Cross-platform binary builds for extension | 6 platforms | Uses `BINARY_NAME=microsoft-azd-waza`, supports `PLATFORM` env var |
| `build.ps1` | Windows equivalent of build.sh | 6 platforms | Identical logic, PowerShell syntax |
| `Makefile` | Local dev build (single platform) | Current OS | Used by developers, `make build` |

**Analysis:**
- `build.sh` and `build.ps1` are **duplicates** (same logic, different shells)
- Only used by the `azd-ext-release.yml` workflow
- CLI release uses inline Go build commands

### 4. Workflow Dependencies & Coordination

```
PR/push to main/develop
    ↓
go-ci.yml (build + test)
    ↓
[Two parallel paths on main]:
    ├─ go-release.yml (triggered by tag v*.*.*)
    │   └─ Produces: waza-linux-amd64, waza-darwin-*, waza-windows-*
    │   └─ Creates: GitHub Release
    │
    └─ azd-ext-release.yml (triggered by version.txt/extension.yaml change)
        └─ Builds using build.sh/build.ps1
        └─ Uses `azd x publish` to update registry.json
        └─ Creates: registry.json PR (auto-merged)
        └─ Creates: GitHub Release (separate from CLI)
```

**Problems:**
- No coordination between the two release workflows
- Manual version bumps required in extension.yaml
- `azd x publish` has its own artifact uploading logic (separate from manual binary builds in go-release.yml)
- No way to atomically release both CLI and extension together

---

## Root Causes

1. **Two independent release pipelines** — CLI (go-release.yml) and extension (azd-ext-release.yml) don't talk to each other
2. **No single version source** — version.txt, extension.yaml, and implicit tag versions are separate
3. **Build script duplication** — build.sh and build.ps1 have identical logic
4. **Unclear semantics** — when should CLI be released vs extension? Both? In sync?
5. **Manual version management** — extension.yaml requires hand-editing

---

## Proposed Solution

### Design Principles

1. **Single Source of Truth:** Version lives in `version.txt` (read by extension.yaml, go-release.yml, and azd extension workflow)
2. **Atomic Releases:** A single tag (v*.*.* format) triggers **both** CLI binary builds and extension packaging/publishing
3. **Clear Separation:** PR checks (go-ci.yml) remain separate from CD (release workflows)
4. **No Duplicate Logic:** Consolidate build scripts; extend existing patterns
5. **Backwards Compatible:** Maintain existing CLI release semantics; enhance extension flow

### New Unified Release Workflow

**Name:** `.github/workflows/release.yml` (consolidates go-release.yml and azd-ext-release.yml)

**Trigger:**
```yaml
on:
  push:
    tags:
      - 'v*.*.*'           # Semantic version tags (e.g., v1.2.3)
  workflow_dispatch:       # Manual trigger (for edge cases)
    inputs:
      version:
        description: 'Version to release (e.g., 1.2.3 — without v prefix)'
        required: true
        type: string
      build_cli: true      # Option to build just CLI or extension separately
      build_extension: true
      publish_extension: false  # Default: don't auto-publish until manual confirmation
```

### Workflow Structure

```
release.yml
  │
  ├─ Job: setup-version
  │   ├─ Extract version from tag or manual input
  │   ├─ Validate semver format
  │   └─ Output: ${{ steps.version.outputs.version }}
  │
  ├─ Job: build-cli (needs: setup-version, matrix: [linux/amd64, linux/arm64, darwin/*, windows/*])
  │   ├─ Checkout
  │   ├─ Setup Go
  │   ├─ Build binary for platform
  │   │   Build flags: -ldflags "-X main.version=${VERSION}"
  │   └─ Upload artifact
  │
  ├─ Job: build-extension (needs: setup-version)
  │   ├─ Checkout
  │   ├─ Setup Go + azd
  │   ├─ Build extension: azd x build --all --skip-install
  │   └─ Upload artifacts
  │
  ├─ Job: create-cli-release (needs: build-cli)
  │   ├─ Download CLI artifacts
  │   ├─ Generate checksums
  │   └─ Create GitHub Release (v*.*.*)
  │
  ├─ Job: publish-extension (needs: [setup-version, build-extension], if: inputs.publish_extension)
  │   ├─ Download extension artifacts
  │   ├─ Run: azd x publish (updates registry.json, creates release)
  │   ├─ Create PR with registry.json update (auto-merge)
  │   └─ Update extension.yaml version field
  │
  └─ Job: sync-versions (needs: [setup-version, publish-extension])
      ├─ Update version.txt
      ├─ Update extension.yaml (via update-version script)
      └─ Commit & push if changes detected
```

### Version File Updates

**Before Release:**
- `version.txt` contains next version (e.g., 0.4.0)
- `extension.yaml` manually synced (e.g., version: 0.4.0)
- `registry.json` reflects current published versions

**Release Trigger (tag v0.4.0):**
1. Extract version: 0.4.0
2. Build CLI binaries, azd extension
3. Create GitHub Release with CLI binaries
4. Publish extension:
   - Run `azd x publish` (updates registry.json automatically)
   - Create PR with registry.json (auto-merge)
5. Sync-versions job:
   - Ensure extension.yaml version: matches 0.4.0 (via find-and-replace)
   - No need to update version.txt (already set)

**Result:**
- `version.txt`, `extension.yaml`, and `registry.json` all at v0.4.0
- Both CLI and extension published atomically
- Single GitHub Release with all artifacts (CLI + extension)

### Fallback Behavior

If `publish_extension` is `false` (default for manual dispatch):
- CLI binaries are built and released
- Extension is built but **not** published (artifacts available for manual review)
- No registry.json update
- Allows safe testing before publishing

---

## Implementation Plan

### Phase 1: Create Unified Workflow

**Files to create:**
- `.github/workflows/release.yml` — New consolidated release workflow (consolidates go-release.yml + azd-ext-release.yml)

**Key sections:**
1. Version setup (extract from tag or input)
2. Parallel build jobs (CLI matrix + extension)
3. Release jobs (GitHub Release + azd extension publish)
4. Version sync (update extension.yaml to match version.txt)

**Compatibility:**
- Existing `go-release.yml` → Keep for one release cycle (deprecate in favor of unified)
- Existing `azd-ext-release.yml` → Keep for one release cycle (deprecate in favor of unified)

### Phase 2: Consolidate Version Management

**Create helper scripts:**
- `scripts/sync-versions.sh` — Update extension.yaml version field to match version.txt

**Update documentation:**
- Add release process docs to `docs/RELEASE.md`
- Document version.txt as source of truth
- Update `AGENTS.md` with release instructions

### Phase 3: Deprecate Old Workflows

After one successful release using `release.yml`:
1. Rename old workflows (e.g., `go-release.yml.deprecated`)
2. Remove triggers from old workflows
3. Update CI configuration (if any) to use new workflow
4. Document migration path for developers

---

## Key Benefits

| Benefit | Impact |
|---------|--------|
| **Atomic Releases** | Both CLI and extension released from single tag; no version mismatch |
| **Single Version Source** | version.txt is truth; extension.yaml and registry.json derived |
| **Reduced Complexity** | One workflow instead of two; shared build logic |
| **Clearer Semantics** | Tag push = full release (CLI + extension); manual dispatch = testing |
| **Better Coordination** | Both artifacts in single GitHub Release |
| **Faster Iteration** | Less boilerplate; easier to add new platforms or targets |

---

## Migration Path

1. **Week 1:** Implement `release.yml` (parallel with existing workflows)
2. **Week 2:** Test new workflow on develop branch (manual dispatch)
3. **Week 3:** First production release using `release.yml` (tag v0.4.0-beta.1)
4. **Week 4:** Deprecate old workflows; update documentation
5. **Week 5+** Remove old workflows in next major version

---

## Risks & Mitigation

| Risk | Mitigation |
|------|------------|
| **Breaking existing CLI releases** | New workflow uses same build logic + flags; test on develop first |
| **Extension publish fails** | Add `publish_extension: false` default; manual approval for first release |
| **Version mismatch between files** | Sync-versions job verifies all three files post-release |
| **azd x publish side effects** | Document exact mutations; test in isolated branch first |
| **Artifact naming inconsistencies** | Use consistent naming: `{binary-name}-{os}-{arch}{.exe}` |

---

## Deliverables

1. ✅ `docs/proposals/release-process-normalization.md` (this document)
2. 📋 `.github/workflows/release.yml` (implementation)
3. 📋 `docs/RELEASE.md` (release process documentation)
4. 📋 `scripts/sync-versions.sh` (version management helper)
5. 📋 Update `AGENTS.md` with release instructions
6. 📋 Deprecation notice in old workflow files

---

## Approval Checklist

- [ ] Architecture: Reviewable design (clear separation of concerns)
- [ ] Implementation: Workflow handles both CLI + extension from single tag
- [ ] Testing: Manual test on develop branch (no side effects)
- [ ] Documentation: Release process clearly documented
- [ ] Compatibility: No breaking changes to existing CLI release process
- [ ] Team consensus: Reviewed by Wallace Breza + squad leads

---

## Questions for Review

1. Should we add a dry-run mode (build artifacts but don't create releases)?
2. Should CLI and extension versions always be in sync, or can they diverge?
3. Should we support releasing just CLI or just extension separately? (Currently: only atomic)
4. How should pre-release versions (e.g., 0.4.0-alpha.1) be handled?

---

## References

- **Issue:** #223 — Normalize release process
- **Existing workflows:** `.github/workflows/{go-release,azd-ext-release,squad-release}.yml`
- **Build scripts:** `build.sh`, `build.ps1`, `Makefile`
- **Version files:** `version.txt`, `extension.yaml`, `registry.json`
- **Custom instructions:** `.ai-team/agents/rusty/` (Rusty architecture directives)
