# JoinQuest agent integration playbook

> **Note:** The backend embeds a copy at `backend/internal/developer/agent_playbook.md` for the `developerAgentPlaybook` GraphQL query and MCP. Edit this file first, then sync the embed copy before release. The bundled agent skill at `.agents/skills/joinquest-integration/playbook.md` should match.

End-to-end workflow for AI agents helping a developer integrate a multiplayer game with JoinQuest — from first conversation through public catalog release.

**Companion docs (load on demand, do not duplicate inline):**

- Technical reference: `developer-integration-guide.md` (MCP: `joinquest_integration_get_integration_guide`)
- Discovery interview (Phase 1 detail): MCP `joinquest_integration_get_discovery_prompt`
- Product spec: `developer-self-service.md`

**Execution:** JoinQuest integration MCP (`joinquest-integration`) — see `mcp-setup.md` in the agent skill folder or the developer dashboard.

---

## Agent rules (always)

1. **Interview before implementing.** Understand the game before writing API code or catalog copy.
2. **Human gates.** Never call `updateMyGameMetadata` or `requestPublicRelease` without explicit developer approval of the draft.
3. **Public HTTPS only.** JoinQuest cannot reach `localhost` — use staging, Render, Fly, ngrok, etc.
4. **One phase at a time.** Finish the current phase (or get explicit skip) before moving on.
5. **MCP first when connected.** Prefer MCP tools over guessing dashboard state. If MCP is not configured, tell the developer how to set it up (Phase 2).
6. **Fail closed, fix forward.** When checks fail, read the error message, pull the integration guide section, fix the game API, re-run checks.

---

## Phase 1 — Discover the game

**Goal:** Understand what the developer is building so you can draft catalog copy and seatTemplate guidance.

**Ask:**

1. **One-liner** — What do players do together in one sentence?
2. **Player count** — Typical group size? Min/max? Fixed teams or variable?
3. **Structure** — Head-to-head duel, free-for-all, teams, or roles (e.g. clue-giver + guessers)?
4. **Social mode** — Competitive, cooperative, or party/social?
5. **Session length** — Quick rounds (~5 min) or longer sessions?
6. **Vibe** — How should the catalog card feel (casual, brainy, chaotic, tactical)?
7. **API URL** — Where will the game API be hosted publicly? (Must be HTTPS.)

**Draft (show to developer, do not save yet):**

| Output | Constraints |
|--------|-------------|
| `shortDescription` | ~120 chars, warm JoinQuest tone — what players *do*, not tech stack |
| `longDescription` | 2–4 short paragraphs for the detail page |
| `howToPlay` | 3–6 bullet steps for a first-time player |
| `tags` | 1–3 IDs from `joinquest_integration_get_catalog_tag_taxonomy` |
| `seatTemplate` plan | Point to integration guide §4 + seat-templates cookbook (duel, teams, roles) |

**Voice:** Plain, player-first. “Find your group. Play together.” — not “enter matchmaking” or JWT jargon.

**Done when:** Developer confirms the drafts are directionally correct (edits OK).

---

## Phase 2 — Connect JoinQuest MCP

**Goal:** Agent can call JoinQuest on the developer's behalf.

**Developer actions (guide them):**

1. Sign in at [joinquest.cc/developers](https://joinquest.cc/developers).
2. Developer dashboard → **Connect an AI assistant** → **Show setup** → **Generate API key** (copy once).
3. Install MCP CLI is **not required** — configs use `npx -y @joinquest/mcp-integration` (Node.js 20+).

4. Add MCP config — **Cursor** (note the cursor-specific npx bin):

```json
{
  "mcpServers": {
    "joinquest-integration": {
      "type": "stdio",
      "command": "npx",
      "args": ["--yes", "--package", "@joinquest/mcp-integration", "joinquest-integration-mcp-cursor"],
      "env": {
        "JOINQUEST_API_KEY": "lq_dev_..."
      }
    }
  }
}
```

**Claude Code:** `claude mcp add --scope project --transport stdio --env JOINQUEST_API_KEY=lq_dev_... joinquest-integration -- npx -y @joinquest/mcp-integration`

Config paths: Cursor → `~/.cursor/mcp.json` or `.cursor/mcp.json`; Claude Code → `.mcp.json`.

5. **Fully quit** the agent client and reopen (Cursor: **Cmd+Q**, not just Reload Window). Start a **new Agent chat**.

**Optional — Agent Skill (no MCP required for workflow text):**

Copy `.agents/skills/joinquest-integration/` from the [JoinQuest repo](https://github.com/scruffyprodigy/playhub) into the developer's project at `.agents/skills/joinquest-integration/` (works across Claude Code, Codex, Cursor, Copilot, Gemini CLI, and others per [agentskills.io](https://agentskills.io/)).

**Verify MCP:**

In Agent mode, ask the agent to call `joinquest_integration_list_my_games`. If auth fails, re-check API key and URL. Some Cursor versions have no **Tools & MCP** settings tab — a successful tool call is the real test.

**Done when:** MCP tools respond successfully.

---

## Phase 3 — Implement the game API

**Goal:** Game server exposes the endpoints JoinQuest calls. Implement in the **game repo**, not the lobby repo.

**Required endpoints:**

| Endpoint | Purpose |
|----------|---------|
| `GET /healthz` | Returns `ok` — JoinQuest probes reachability |
| `GET /api/v1/status` | Game name, version, `launchUrlsOnProvision: true` |
| `GET /api/v1/game-modes` | Modes with `minPlayers`, `maxPlayers`, valid `seatTemplate` (no flat `seats[]`) |
| `POST /api/v1/matches` | Provision — auth via `Authorization: Bearer {serviceToken}` |
| `POST /api/v1/matches/{externalMatchId}/claim` | Validate Lobby seat JWT |

**Provision success:** `200/201` with `launchUrls` map (lobbyUserId → URL base) or `launchUrlTemplate`.

**Provision banlist:** `403` with `{ "error": "...", "bannedLobbyUserIds": ["..."] }`.

**JWT:** Verify `iss`, `aud` (your API base URL), `matchId`, `seatKey`, `sub`. Publish JWKS at `{lobbyIssuer}/.well-known/jwks.json`.

Use Phase 1 seatTemplate plan. Full details: integration guide §3–§7.

**Done when:** Endpoints exist on a **public HTTPS** URL and basic manual curl tests pass.

---

## Phase 4 — Register the game on JoinQuest

**Goal:** Game row exists; manifest sync attempted.

**Registration fields** (developer fills at [/developers](https://joinquest.cc/developers) or you guide them):

| Field | Required | Notes |
|-------|----------|-------|
| Game name | Yes | Display name |
| Slug | Yes | Lowercase URL slug (auto-suggested from name) |
| Short description | Yes | Initial card blurb (~1–2 sentences) |
| API base URL | Yes | Public HTTPS origin, e.g. `https://mygame.example.com` |
| Contact email | Yes | Review notifications |
| Website URL | No | |
| Community URL | No | Discord, etc. |

**After register:** JoinQuest calls `healthz` + `game-modes`. Success → `private_testing`. Failure → stays `draft` with `connectError`.

**MCP (after register):**

- `joinquest_integration_list_my_games` — find `gameId`
- `joinquest_integration_get_game_checks` — current state
- `joinquest_integration_get_game_credentials` — `serviceToken`, `webhookSecret` (sensitive)
- `joinquest_integration_get_example_provision_payload` — sample POST body for local debugging

**Done when:** Game visibility is `PRIVATE_TESTING` (or developer understands draft errors and is fixing API URL).

---

## Phase 5 — Run integration checks and fix failures

**Goal:** All required checks pass.

**MCP:** `joinquest_integration_run_game_checks` with `gameId`.

**Required checks (must be PASS for public release):**

- `manifest.reach_api`, `manifest.status`, `manifest.game_modes`, `manifest.sync_freshness`
- `provision.happy_path`, `provision.auth`, `provision.launch_urls`
- JWT checks where applicable

**On failure:**

1. Read `message` and `detail` from the check result.
2. Load integration guide §10 index → relevant section.
3. Fix the **game API**, deploy, re-run checks.

Common fixes:

| Check | Fix |
|-------|-----|
| `manifest.reach_api` | API must be public HTTPS; `/healthz` returns 200 |
| `manifest.game_modes` | Valid `seatTemplate` JSON per mode |
| `provision.auth` | Accept dashboard `serviceToken` on Authorization header |
| `provision.launch_urls` | Return `launchUrls` for every seated player |
| `provision.banlist` | Return 403 + `bannedLobbyUserIds` array |
| `jwt.*` | JWKS URL, claim endpoint, reject wrong aud/iss |

**Done when:** Required checks are green (re-run to confirm).

---

## Phase 6 — Save catalog metadata

**Goal:** Listing copy matches Phase 1 drafts (updated if needed).

**MCP:** `joinquest_integration_get_catalog_tag_taxonomy` for valid tag IDs.

**Show developer the final draft.** On explicit approval:

**MCP:** `joinquest_integration_update_game_metadata` with:

```json
{
  "gameId": "...",
  "shortDescription": "...",
  "longDescription": "...",
  "howToPlay": "...",
  "tags": ["party", "quick"]
}
```

**Done when:** Metadata saved and visible via `joinquest_integration_get_game_checks`.

---

## Phase 7 — Test with friends

**Goal:** Real players can join a table in a private room.

**Developer dashboard:** **Create test table** (or equivalent flow).

Friends join via the developer's room link. Game stays off the public catalog until release.

**Agent role:** Confirm checks are green; encourage a playtest before release.

**Done when:** Developer confirms a successful test session (or explicitly skips).

---

## Phase 8 — Request public release

**Goal:** Submit for JoinQuest catalog review.

**Gates (automated):**

- Required checks PASS
- `shortDescription`, `longDescription`, at least one tag
- Visibility `PRIVATE_TESTING`

**Confirm with developer:** They want to go public.

**MCP:** `joinquest_integration_request_public_release` with `gameId`.

**After submit:** Visibility → `PENDING_REVIEW`. JoinQuest reviews for spam/IP/offensive content.

**Done when:** Request submitted (or game already public).

---

## MCP tool quick reference

| Phase | Tools |
|-------|-------|
| Setup | `joinquest_integration_get_agent_playbook` (this doc), `joinquest_integration_get_integration_guide` |
| Discover | `joinquest_integration_get_discovery_prompt`, `joinquest_integration_get_catalog_tag_taxonomy` |
| Register / status | `joinquest_integration_list_my_games`, `joinquest_integration_get_game_checks`, `joinquest_integration_get_game_credentials`, `joinquest_integration_get_example_provision_payload` |
| Checks | `joinquest_integration_run_game_checks` |
| Metadata | `joinquest_integration_update_game_metadata` |
| Release | `joinquest_integration_request_public_release` |

---

## Phase checklist (copy for status updates)

```text
[ ] Phase 1 — Discover game (drafts approved)
[ ] Phase 2 — MCP connected
[ ] Phase 3 — Game API implemented (public HTTPS)
[ ] Phase 4 — Registered on JoinQuest (private_testing)
[ ] Phase 5 — Integration checks passing
[ ] Phase 6 — Catalog metadata saved
[ ] Phase 7 — Test table played
[ ] Phase 8 — Public release requested
```
