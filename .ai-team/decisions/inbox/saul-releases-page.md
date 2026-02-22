# Decision: Releases page pattern

**By:** Saul (Documentation Lead)
**Date:** 2026-02-21
**Issue:** #383
**PR:** #384

**What:** Created a releases reference page at `site/src/content/docs/reference/releases.mdx` that shows the current release (v0.8.0) with changelog highlights, download table, install commands, and azd extension info. Older releases link out to GitHub Releases rather than duplicating content.

**Why:** The docs site should be a self-contained starting point for users downloading waza. Having binaries, install commands, and changelog highlights in one place reduces friction. Linking to GitHub Releases for history avoids maintaining two changelog surfaces.

**Pattern for future releases:** When cutting a new version, update the releases.mdx page — change the version number, update the changelog highlights, and update download URLs. The CHANGELOG.md remains the source of truth; the releases page is a curated summary of the latest.
