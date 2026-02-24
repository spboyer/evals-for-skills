# Decision: LinkScorer architecture

**By:** Linus (Backend Developer)
**Date:** 2026-02-21
**Issue:** #406

## What

The `LinkScorer` lives in `cmd/waza/dev/links.go` alongside the other scorers (`HeuristicScorer`, `SpecScorer`, `McpScorer`). It follows the same pattern: takes `*skill.Skill`, returns a `*LinkResult` with categorized issues.

For the MCP server (`internal/mcp/server.go`), a lightweight `quickLinkCheck` helper does regex-based link extraction instead of importing `cmd/waza/dev`. This avoids an `internal/ → cmd/` import cycle while still surfacing broken links in the MCP tool response.

## Why

`internal/` packages cannot import `cmd/` packages — that's a Go anti-pattern. Rather than extracting a new `internal/linkcheck/` package (which would split the scorer pattern), the MCP server does a simpler scan (regex on SKILL.md only, local links only, no BFS orphan detection). The full analysis lives in `waza check`; the MCP tool is a lightweight readiness signal.

## Impact

- `cmd/waza/dev/links.go` — full LinkScorer with goldmark AST, concurrent URL checking, BFS orphan detection
- `cmd/waza/cmd_check.go` — new linkResult field, 📎 Links display section, readiness gate, next steps
- `internal/mcp/server.go` — new fields on `skillCheckResult`, `quickLinkCheck` helper
- `internal/mcp/tools.go` — updated description
- goldmark promoted from indirect to direct dependency in go.mod
