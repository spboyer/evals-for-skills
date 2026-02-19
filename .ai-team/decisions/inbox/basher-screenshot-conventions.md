# Decision: Screenshot spec conventions

**By:** Basher (Tester / QA)
**Date:** 2026-02-19
**Issue:** #251

## What

Screenshot tests live in `web/e2e/screenshots.spec.ts` and output to `docs/images/`. Conventions:

1. **Viewport:** 1280×720, chromium only (no firefox — screenshots must be pixel-consistent)
2. **Paths:** Use `../docs/images/` (relative to Playwright config root `web/`), NOT relative to the test file
3. **Mock data:** Reuse `mockAllAPIs` and existing fixtures — no screenshot-specific mock data
4. **Views requiring interaction:** Set up state (select options, expand rows) before capturing
5. **Naming:** kebab-case matching the view name: `dashboard-overview.png`, `run-detail.png`, `compare.png`, `trends.png`

## Why

Reproducible screenshots from mock data mean docs images stay consistent regardless of when/where they're generated. Running `npx playwright test e2e/screenshots.spec.ts --project=chromium` regenerates all four images deterministically.
