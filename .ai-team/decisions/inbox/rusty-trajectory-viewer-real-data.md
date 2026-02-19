# Decision: Trajectory viewer uses real TranscriptEvent data

**By:** Rusty (Lead / Architect)
**Date:** 2025-07-25
**PR:** #243

**What:** The trajectory replay viewer now consumes real `TranscriptEvent` data from the API rather than regex-heuristic guessing from grader messages. When transcript data is unavailable, it falls back gracefully to the old grader-based summary.

**Why:** The old approach fabricated fake timestamps (`new Date().toISOString()`) and used regex to guess event types from grader messages — unreliable and inaccurate. With #237 adding transcript and session digest to the API, we can render the actual agent execution timeline.

**Pattern established:** New transcript-aware components (`SessionDigestCard`, `ToolCallDetail`) follow the existing dashboard conventions: zinc-900 dark theme, Tailwind CSS v4, lucide-react icons, collapsible detail panels. Tool call correlation by `toolCallId` is the canonical way to match start/complete events.
