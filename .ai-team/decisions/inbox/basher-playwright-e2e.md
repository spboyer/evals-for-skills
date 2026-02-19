### 2025-07-18: Playwright E2E uses route interception, not Go server

**By:** Basher
**What:** Phase 1 E2E tests use Playwright route interception to mock API responses rather than spinning up the Go server with fixture data.
**Why:** Route interception is faster (no binary build), isolates frontend tests from backend changes, and runs reliably in CI without Go toolchain. Full-stack E2E against the real server can be added in Phase 2 when we need to test the complete request path.

### 2025-07-18: Use regex patterns for Playwright route matching

**By:** Basher
**What:** All API route interceptions use JavaScript regex (e.g., `/\/api\/runs(\?|$)/`) instead of glob patterns.
**Why:** Playwright glob patterns like `**/api/runs` don't match URLs with query strings (`/api/runs?sort=timestamp&order=desc`). This caused intermittent test failures. Regex gives precise control over URL matching.

### 2025-07-18: Tailwind v4 uses oklch color space

**By:** Basher
**What:** `getComputedStyle().backgroundColor` returns `oklch(0.21 0.006 285.885)` not `rgb(24, 24, 27)` in Tailwind v4.
**Why:** Theme tests must parse oklch values. Asserting on oklch lightness (< 0.3 for dark backgrounds) is more resilient than hardcoding exact color values.
