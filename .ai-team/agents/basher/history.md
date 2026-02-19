# History — Basher

## Project Context
- **Project:** waza — CLI tool for evaluating Agent Skills
- **Stack:** Go (primary), React 19 + Tailwind CSS v4 (web UI)
- **User:** Shayne Boyer (spboyer)
- **Repo:** spboyer/waza
- **Universe:** The Usual Suspects

## Key Learnings

### Testing Strategy
- **Model directive:** Coding in Claude Opus 4.6 (same as production code)
- **Test types:** Go unit tests (*_test.go), integration tests, Playwright E2E
- **Fixture isolation:** Original fixtures never modified — tests work in temp workspace
- **Coverage goal:** Non-negotiable (Rusty's requirement)

### Waza-specific Tests
- TestCase execution scenarios
- BenchmarkSpec validation
- Validator registry functionality
- CLI flag handling
- Agent execution mocking

### CI/CD
- Branch protection requires tests to pass
- Go CI workflow in .github/workflows/go-ci.yml
- Test results tracked for quality assurance

### Playwright E2E (PR #241, Issue #208)
- **Config:** `web/playwright.config.ts` — Chromium, vite preview on port 4173, screenshots/video on failure
- **Test dir:** `web/e2e/` with `fixtures/mock-data.ts`, `helpers/api-mock.ts`, and spec files
- **Script:** `npm run test:e2e` (pre-builds with `npm run build`)
- **Route interception:** Must use regex patterns (not globs) to handle query strings — e.g. `/\/api\/runs(\?|$)/`
- **Tailwind v4 colors:** `getComputedStyle` returns `oklch()` not `rgb()` — assert lightness < 0.3 for dark theme
- **react-query retries:** Default 3 retries with backoff; error state tests need ~15s timeout
- **page.request vs page.evaluate:** `page.request.get()` bypasses route interception; use `page.evaluate(() => fetch(...))` instead
- **Hash routing:** App uses hash-based routing (`#/runs/:id`), not react-router. URLs in tests are like `/#/runs/run-001`
- **Previous branch was broken:** Old `squad/208-playwright-e2e` committed node_modules and diverged from main. Force-pushed clean version.
