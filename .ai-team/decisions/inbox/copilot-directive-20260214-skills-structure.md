### 2026-02-14: User directive — Skills repo plugin bundle structure
**By:** Shayne Boyer (via Copilot)
**What:** The microsoft/skills repo is being reorganized from a flat `.github/skills/` layout (133 items) into plugin bundles (`.github/plugins/<bundle>/skills/<name>/`). Waza CI compatibility (#60) and any future skills integration must support both the current flat layout and the new nested plugin bundle layout. Key bundles: azure-skills (18 orchestration), azure-sdk-python (41), azure-sdk-dotnet (29), azure-sdk-typescript (24), azure-sdk-java (26), azure-sdk-rust (7), azure-core (6).
**Why:** User shared the distribution strategy gist (https://gist.github.com/spboyer/011190893f33d82d967180cdc5a2432d) — this is the planned future state and all CI work should be forward-compatible.
