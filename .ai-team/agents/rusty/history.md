# History — Rusty

## Project Context
- **Project:** waza — CLI tool for evaluating Agent Skills
- **Stack:** Go (primary), React 19 + Tailwind CSS v4 (web UI)
- **User:** Shayne Boyer (spboyer)
- **Repo:** spboyer/waza
- **Universe:** The Usual Suspects

## Key Learnings

### Architecture
- **Model selection directive (2026-02-18):** Coding in Claude Opus 4.6, reviews in GPT-5.3-Codex, design in Gemini Pro 3
- **Web UI styling:** Keep clean and functional — colors close to DevEx dashboard, no fancy gradients
- **Agent execution:** Go engine drives CLI, web UI for visualization

### Code Quality
- Test coverage is non-negotiable
- Interface-based design for flexibility (AgentEngine, Validator patterns in Go)
- Functional options for configuration (Go convention)

### Team Structure
- Linus owns Go backend implementation
- Basher owns all testing strategy
- Livingston/Saul own documentation
- Richard Park available for Copilot SDK questions

## Work Log

### 2026-02-19: #80 — BPE Tokenizer (PR #260)
- **Reviewed PR** by Charles Lowell (chlowell) — ported BPE tokenizer from Microsoft/Tokenizer
- **Architecture:** `Counter` interface preserved. `NewCounter(tokenizer)` factory replaces `NewEstimatingCounter()`. BPE is new default via `TokenizerDefault`.
- **New package:** `internal/tokens/bpe/` — BinaryMap, LRU cache, byte-pair encoder, tokenizer, builder
- **Embedded model:** `o200k_base.tiktoken` (~3.6MB) via `go:embed` — adds to binary size
- **Flag design:** `--tokenizer` flag only on `count` command; `check`/`compare`/`suggest` hardcode `TokenizerDefault`
- **Findings:** `regex` field on Tokenizer struct is dead code (set but never read); `NewTokenizerFromFile` is dead code (defined but never called); `Cache` field is exported unnecessarily; LRU cache is not thread-safe (fine for CLI but should be documented)
- **Verdict:** APPROVE with comments — architecture is sound, implementation correct, concerns are improvements not blockers

### 2025-07-25: #238 — True trajectory replay viewer (PR #243)
- **Branch:** `squad/238-trajectory-viewer`
- Full rewrite of `TrajectoryViewer.tsx` to consume real `TranscriptEvent` data
- Created `SessionDigestCard.tsx` (digest stats + tools used badges + errors)
- Created `ToolCallDetail.tsx` (expandable JSON viewers for args/result)
- Timeline: color-coded dots (blue=tool start, green/red=complete, emerald=turn, red=error)
- `toolCallId` correlation links Start ↔ Complete events
- Graceful fallback to grader-based heuristic when transcript is empty
- Depends on #237 (transcript + session digest in API)

## Learnings

### Release Infrastructure Audit (Issue #223, Feb 20)

**Version Management State:**
- `version.txt` (source-of-truth candidate): 0.4.0-alpha.1
- `extension.yaml`: 0.3.0 (STALE, 2 patch releases behind)
- `registry.json`: max version 0.3.0 (missing 0.4.0-alpha.1, blocking extension users)
- Go binary version: injected via ldflags `-X main.version=${VERSION}` during build

**Key Files & Their Purpose:**
- `.github/workflows/go-release.yml` — Standalone CLI release (builds 6 platform matrix, creates GitHub Release on `v*` tag) — **ACTIVE, currently used**
- `.github/workflows/azd-ext-release.yml` — Extension release (triggers on version.txt/extension.yaml changes, publishes to azd registry via `azd x publish`) — **ACTIVE, currently used**
- `.github/workflows/release.yml` — NEW unified workflow (both CLI + extension from single trigger, includes version sync job) — **PREPARED BUT NOT YET ACTIVATED**
- `.github/workflows/squad-release.yml` — Package.json-based release (NOT relevant to waza-go)
- `Makefile` — Local build with `VERSION?=0.1.0` (default, overridable)
- `build.sh` — Cross-platform binary builder for extension (VERSION env var driven)
- `cmd/waza/root.go` — Version variable: `var version = "dev"` (overwritten at build-time)

**Critical Issues Identified:**
1. **No unified release trigger** — CLI and extension release independently, easy to desync
2. **Version sync failure** — extension.yaml not bumped when CLI version.txt is updated
3. **registry.json desynchronization** — stale, depends on manual azd-ext-release workflow execution
4. **release.yml has logical flaw** — sync-versions job runs AFTER build jobs, too late to affect artifact versions
5. **Dual tag schemes** — CLI uses `v*`, extension uses `azd-ext-microsoft-azd-waza_*`, confusing
6. **No pre-flight validation** — easy for files to drift; no check before build

**Recommended Architecture:**
- Single canonical trigger: git tag `v*.*.*` 
- release.yml should be the sole release coordinator (retire go-release.yml + azd-ext-release.yml once stable)
- Pre-flight job in release.yml to validate version.txt == tag before proceeding
- Sync-versions job should run BEFORE builds, not after
- Document the flow in docs/RELEASE.md

**Immediate Action:** registry.json is blocking extension users on 0.4.0-alpha.1. Should manually trigger release.yml or azd-ext-release.yml to sync.


## 📌 Team update (2026-02-20): Model policy overhaul

All code roles now use `claude-opus-4.6`. Docs/Scribe/diversity use `gemini-3-pro-preview`. Heavy code gen uses `gpt-5.2-codex`. Decided by Scott Boyer. See decisions.md for full details.
