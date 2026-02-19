# Release Process

This document describes how the Waza release process works. All releases are handled by the unified workflow at `.github/workflows/release.yml`.

## Cutting a Release

### Tag Push (recommended)

```bash
git tag v1.2.3
git push origin v1.2.3
```

This triggers both flows in parallel: CLI build + release, and extension build + publish.

### Manual Dispatch

Go to **Actions → Release → Run workflow** and fill in:

| Input | Description | Default |
|-------|-------------|---------|
| `version` | Semver without `v` prefix (e.g. `1.2.3`) | *required* |
| `build_cli` | Build and release standalone CLI binaries | `true` |
| `build_extension` | Build, release, and publish azd extension | `true` |

Manual dispatch creates the git tag automatically if it doesn't exist.

## Two Independent Flows

### Flow 1: Standalone CLI Release

1. **build-cli** — Matrix build for 6 platforms (linux, darwin, windows × amd64, arm64). Builds the web UI then produces `waza-{os}-{arch}` binaries with version injected via `-ldflags`.
2. **release-cli** — Downloads CLI artifacts, generates SHA256 checksums, creates a **GitHub Release** (`Waza vX.Y.Z`) with standalone binaries attached.

### Flow 2: azd Extension Release

A single `release-extension` job runs all steps sequentially in one workspace:

1. **Sync version files** — Updates `version.txt` and `extension.yaml` to the target version **before** any azd commands run, ensuring all tools see the correct version.
2. **Build & pack** — Runs `azd x build` and `azd x pack` to produce platform-specific archives.
3. **Release** — Runs `azd x release` to create an Extension GitHub Release.
4. **Publish** — Runs `azd x publish` to update `registry.json` with artifact URLs and checksums.
5. **Commit back** — Creates a PR with the updated `registry.json`, `version.txt`, and `extension.yaml`, reports required CI checks as successful, and auto-merges it.

Running everything in a single job ensures all azd commands operate on the same workspace with consistent, synced version files.

## Version File Locations

| File | Purpose |
|------|---------|
| `version.txt` | Canonical version string used by build scripts |
| `extension.yaml` | `version:` field for the azd extension manifest |
| `registry.json` | Extension registry with download URLs and checksums (updated by publish step) |

## How Automated PR Merging Works

GitHub intentionally does not trigger workflows on PRs created by `GITHUB_TOKEN` (to prevent recursive loops). Since the registry/version PR is created by `github-actions[bot]`, the required CI checks (`Build and Test Go Implementation`, `Lint Go Code`) never run. To satisfy branch protection, the release workflow **reports these status checks as successful** via the commit status API before merging. This is safe because:

- The PR only contains machine-generated changes (checksums, URLs, version strings)
- The release artifacts have already been built and validated earlier in the pipeline
- The content is deterministic and derived from the release that just completed
- The PR provides full traceability of what changed and when

## Deprecated Workflows

The following workflows are superseded by `release.yml` and kept for reference only:

- `go-release.yml` — Previously handled standalone CLI releases
- `azd-ext-release.yml` — Previously handled azd extension releases
