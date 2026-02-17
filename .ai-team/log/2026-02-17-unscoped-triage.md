# Session: 2026-02-17 Unscoped Triage

**Requested by:** Shayne Boyer  
**Lead:** Rusty  
**Date:** 2026-02-17

## Summary

Rusty triaged 5 unscoped issues (#2, #10, #14, #16, #21):
- Recommended **closing #2** (superseded by #14)
- **Labeled rest P2:** #10, #14, #16, #21 assigned P2, epic labels, release:backlog
- **Commented on all** with detailed analysis

## Key Finding

#16 (JSON-RPC server) is **independent from #14** (Web UI). JSON-RPC is a standalone integration layer that enables IDE plugins (VS Code, JetBrains) independently — #14 depends on it, not the reverse.

## Actions

1. Triage analysis commented on all 5 issues
2. P2 + epic labels applied to #10, #14, #16, #21
3. #2 marked for closure (awaiting user approval)
4. #14 identified for rename (remove JSON-RPC scope to avoid overlap)
