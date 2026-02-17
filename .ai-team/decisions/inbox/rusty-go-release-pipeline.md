# Decision: Go Release Pipeline Architecture

**By:** Rusty (Lead)  
**Date:** 2026-02-16  
**Status:** Design Spec (ready for implementation — Linus will code)  
**Related:** E7 (AZD Extension P2), #62 (tracking)

## Problem

Waza now has two implementations:
- **Python** (legacy) — deprecated, no longer maintained
- **Go** (primary) — active development, intended for production use

The current release workflow (`release.yaml`) only handles Python. We need a dedicated Go release pipeline that:
1. Builds cross-platform binaries (6 platforms)
2. Generates checksums and release notes
3. Supports both tag-based and manual version triggers
4. Injects version info into the binary
5. Provides an easy install script for users

## Architecture Decisions

### Trigger Strategy
- **Primary:** Git push to `v*` tags (e.g., `v0.3.1`)
- **Secondary:** `workflow_dispatch` with manual version input
- Rationale: Mirrors standard Go release patterns; `workflow_dispatch` allows one-off releases

### Matrix: 6 Platform Targets
```
OS × Architecture:
- linux/amd64
- linux/arm64
- darwin/amd64 (macOS Intel)
- darwin/arm64 (macOS Apple Silicon)
- windows/amd64
- windows/arm64
```

Matches production diversity; all 6 combinations built in every release.

### Binary Naming Convention
```
waza-{os}-{arch}[.exe]
```
Examples:
- `waza-linux-amd64`
- `waza-darwin-arm64`
- `waza-windows-amd64.exe`
- `waza-windows-arm64.exe`

**Why not `microsoft-azd-waza-*`?** The `microsoft-azd-waza` prefix is for the *azd extension* (published to azd registry). The standalone CLI is `waza`.

### Version Injection
```bash
go build -ldflags "-X main.version=$VERSION" -o $BINARY ./cmd/waza
```
Writes version into the binary at build time. Users see version via `waza --version`.

### Checksums & Artifacts
**SHA256 checksums file:** `waza-$VERSION.sha256`  
Format (one per line):
```
<sha256-hash>  <filename>
```
Example:
```
abc123def456  waza-linux-amd64
789ghi012jkl  waza-darwin-arm64
...
```

Stored as a release artifact alongside binaries. Users can verify integrity:
```bash
sha256sum -c waza-0.3.1.sha256
```

### Release Notes Generation
**Source:** Merged PRs since the previous tag, extracted via GitHub API.

**Template:**
```markdown
# Waza v{VERSION}

## Changes
- [List of merged PRs with authors]

## Installation
See **Getting Started** section below.

## Getting Started
[Include install.sh usage example]

## Platform Support
- Linux x86-64 and ARM64
- macOS (Intel and Apple Silicon)
- Windows x86-64 and ARM64
```

**Implementation:** Query GitHub API for merged PRs in the commit range; format automatically.

### Install Script (`install.sh`)
Single user-friendly script that:
1. Detects OS (linux/darwin/windows)
2. Detects architecture (amd64/arm64)
3. Downloads the correct binary from latest release
4. Verifies SHA256 checksum (if `sha256sum` available)
5. Makes executable and moves to `~/.local/bin/` (or `PATH`)

Example usage:
```bash
curl -fsSL https://raw.githubusercontent.com/spboyer/waza/main/install.sh | bash
waza --version
```

**Error handling:** Clear messages if OS/arch not supported or download fails.

## Workflow Stages (for Linus to implement)

### Stage 1: Setup
- Checkout repo (full history for tag diffing)
- Setup Go environment (1.25+)
- Parse version from either `git ref tags/{version}` or `inputs.version`

### Stage 2: Build (matrix job)
For each platform/arch pair:
- Set `GOOS`, `GOARCH`, `VERSION` env vars
- Run: `go build -ldflags "-X main.version=$VERSION" -o waza-{os}-{arch}[.exe] ./cmd/waza`
- Verify binary execution: `./waza-{os}-{arch} --version`

### Stage 3: Checksum Generation
- Compute SHA256 for all 6 binaries
- Generate `waza-{VERSION}.sha256` file

### Stage 4: Release Notes
- Query GitHub API for merged PRs since previous tag
- Format into markdown
- Include install script instructions

### Stage 5: Create Release
- Upload 6 binaries + checksums + install.sh to GitHub Release
- Use generated release notes from Stage 4
- Set prerelease flag if version contains alpha/beta/rc

### Stage 6: Publish Registry (conditional)
- If this is an azd extension release, update `registry.json` and create PR
- **Note:** Waza CLI releases ≠ azd extension releases (separate workflows)

## File Changes Required

### New file: `.github/workflows/go-release.yml`
Main release pipeline (250-300 lines)

### New file: `install.sh` (in repo root)
Platform detection + download logic (50-70 lines)

### Updated: `.github/workflows/release-python-legacy.yaml`
- Renamed from `release.yaml`
- Added `if: false` to disable
- Added deprecation header comment

## Acceptance Criteria (for Linus to verify)

- [ ] Workflow triggers on `v*` tags
- [ ] Workflow_dispatch accepts manual version input
- [ ] 6 platform binaries built successfully
- [ ] Binary naming matches convention (with .exe for Windows)
- [ ] Version injected into binary (readable via `--version`)
- [ ] SHA256 checksums file created and correct
- [ ] Release notes generated from merged PRs
- [ ] install.sh detects platform correctly
- [ ] install.sh downloads and verifies binary
- [ ] All artifacts uploaded to GitHub Release
- [ ] Prerelease flag set for alpha/beta/rc versions
- [ ] No breaking changes to azd extension workflow (orthogonal)

## Integration with Existing Workflows

### `azd-ext-release.yml` (unchanged)
- Publishes azd *extension* to azd registry
- Uses azd build commands; produces .zip/.tar.gz archives
- Separate versioning and tagging pattern
- Waza Go CLI releases don't trigger this

### `go-ci.yml` (unchanged)
- Runs on every PR/push to main/develop
- Continues to run tests, linting, coverage
- Release workflow is a separate pipeline

### Version Authority
- **CLI version:** `version.txt` (read by release workflow)
- **Extension version:** Also `version.txt` (shared)
- **Rationale:** Single source of truth; both CLI and extension version together

## Future Considerations

### Publish to package registries
- Homebrew tap (`.homebrew/waza.rb`) for macOS users
- Scoop bucket for Windows users
- apt/rpm repositories for Linux distributions
- Deferred to future epic (E8)

### Code signing
- Sign binaries with organizational certificate
- Verify checksums against signed manifest
- Deferred to security hardening pass

### Version tagging strategy
- Decide on semver vs calendar versioning for extension vs CLI
- Currently: shared version in `version.txt`
- Consider split if extension versioning diverges

## Why This Design

1. **Simplicity over features** — Matrix build is standard Go pattern; no custom orchestration
2. **User-friendly** — `install.sh` is copy-paste; no need to visit releases page
3. **Orthogonal to extension** — CLI and extension releases are independent concepts
4. **Cross-platform tested** — Matrix ensures 6 combinations all build; no platform-specific surprises
5. **Verifiable** — SHA256 checksums provide integrity assurance
6. **Maintainable** — Clear separation from Python legacy workflow

## Next Steps

1. **Linus implements** `.github/workflows/go-release.yml` per acceptance criteria
2. **Test manually** with `workflow_dispatch` before cutting first tag release
3. **Document in README** — add installation section linking to `install.sh`
4. **Tag and release** v0.3.1 once design is merged
