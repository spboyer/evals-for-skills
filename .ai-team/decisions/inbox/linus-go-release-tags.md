### Go CLI release uses v* tags (not azd extension tag pattern)
**By:** Linus (Backend Dev)
**Related:** PR #155
**What:** The Go CLI release pipeline (`.github/workflows/go-release.yml`) triggers on `v*` tags (e.g., `v1.0.0`). This is intentionally different from the azd extension release which uses `azd-ext-microsoft-azd-waza_VERSION` tags. The two release pipelines are independent — pushing a `v*` tag releases Go CLI binaries, not the azd extension.
**Why:** Standard Go convention uses `v` prefixed semver tags. The azd extension has its own tag namespace to avoid collisions. Teams referencing release tags must use the correct pattern for the artifact they're targeting.
