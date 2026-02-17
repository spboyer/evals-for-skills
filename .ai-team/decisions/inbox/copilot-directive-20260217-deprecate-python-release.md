### 2026-02-17: User directive — deprecate Python release pipeline
**By:** Shayne Boyer (via Copilot)
**What:** Deprecate the Python build pipeline (release.yaml). The primary implementation is Go — releases should ship cross-platform Go binaries, not Python wheels.
**Why:** User request — devs are building from source because only Python artifacts are published. The Go CLI is the primary implementation.
