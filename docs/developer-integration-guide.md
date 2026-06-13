# JoinQuest developer integration guide

> **Note:** The backend embeds a copy at `backend/internal/developer/integration_guide.md` for the `developerIntegrationGuide` GraphQL query. Edit this file first, then sync the embed copy before release.

Single source of truth for humans, dashboard copy, and MCP agents (`joinquest_integration_get_integration_guide`).

**Agent workflow:** [developer-agent-playbook.md](./developer-agent-playbook.md) — end-to-end phases for vibe coding (`joinquest_integration_get_agent_playbook`). Bundled as an [Agent Skill](../.agents/skills/joinquest-integration/SKILL.md) at `.agents/skills/joinquest-integration/`.

**Product context:** [developer-self-service.md](./developer-self-service.md) · Wire contract: [lobby-protocol-handoff.md](./lobby-protocol-handoff.md) · Seat layouts: [seat-templates-and-matchmaking.md](./seat-templates-and-matchmaking.md)

---

## 0. Discover your game (agents + developers)

Before writing API code or catalog copy, make sure you (or your agent) understand what you're building. Use this script in conversation — agents should **ask**, listen, then **draft** metadata and manifest suggestions for developer approval.

### Discovery questions

1. **One-liner** — What do players do together in one sentence?
2. **Player count** — Typical group size? Min/max? Fixed teams or variable?
3. **Structure** — Head-to-head duel, free-for-all, teams, or roles (e.g. clue-giver + guessers)?
4. **Social mode** — Competitive, cooperative, or party/social (low stakes, lots of laughs)?
5. **Session length** — Quick rounds (~5 min) or longer sessions?
6. **Vibe** — How should the catalog card *feel* (casual, brainy, chaotic, tactical)?

### From answers → manifest + tags

| If the game is… | Start with | Example tags |
|-----------------|------------|--------------|
| 1v1 competitive | `{ "count": 2 }` duel mode | `competitive`, `1v1`, `quick` |
| Team vs team | `Team: { count: 2, Seat: { count: N } }` | `competitive`, `party` |
| Co-op in one group | `{ "count": N }` or single team template | `cooperative`, `party` |
| Roles / composition | Role buckets in `seatTemplate` | `competitive` or `cooperative` + role-specific copy |

See the [seat template cookbook](./seat-templates-and-matchmaking.md#cookbook) for JSON examples.

Agents: after discovery, call `updateMyGameMetadata` with drafted copy (developer confirms) and implement the suggested `seatTemplate` on the game API.

---

## 1. Catalog voice (JoinQuest tone)

Lobby copy is **warm, plain, and player-first** — not enterprise, not hype.

**Do**

- Use short sentences. Card blurbs ~120 characters.
- Say “look for a group” / “play together”, not “enter matchmaking queue”.
- Describe what players *do*, not tech stack.
- Match the friendly tone of JoinQuest: *“Find your group. Play together.”*

**Don't**

- Stack buzzwords (“next-gen social competitive ecosystem”).
- Mention JWT, GraphQL, or provision in player-facing copy.
- Use ALL CAPS or excessive exclamation marks.

**Examples**

| Weak | Better (short description) |
|------|------------------------------|
| “Real-time multiplayer word game leveraging WebSocket architecture” | “Race friends to find words on a shared grid.” |
| “The ultimate 1v1 RPS experience!!!” | “Classic rock paper scissors — best of three, no accounts on your side.” |

**Field guide**

| Field | Audience | Length / notes |
|-------|----------|----------------|
| `shortDescription` | Catalog card | 1–2 sentences, ~120 chars |
| `longDescription` | Detail page | 2–4 short paragraphs; what it is, why it's fun |
| `howToPlay` | Detail page | 3–6 bullet steps; first-time player |
| `tags` | Catalog chips | 1–3 from [taxonomy](#catalog-tag-taxonomy); max 3 shown in UI |

---

## 2. Catalog tag taxonomy

Use **only** these tag IDs (labels shown in UI):

| ID | Label | Use when |
|----|-------|----------|
| `competitive` | Competitive | Winners/losers, rankings, direct opposition |
| `cooperative` | Co-op | Players win or lose together |
| `party` | Party | Social, groups, casual fun |
| `1v1` | 1v1 | Exactly two players head-to-head |
| `quick` | Quick | Sessions under ~10 minutes |
| `words` | Words | Word/language puzzles |
| `strategy` | Strategy | Planning, hidden info, tactics |
| `casual` | Casual | Low pressure, easy to pick up |

Query `catalogTagTaxonomy` for the machine-readable list.

---

## 3. Quick start (API)

1. **Health** — `GET {apiBaseUrl}/healthz` → `ok`
2. **Status** — `GET {apiBaseUrl}/api/v1/status` → `{ game, version, launchUrlsOnProvision: true }`
3. **Game modes** — `GET {apiBaseUrl}/api/v1/game-modes` → modes with valid `seatTemplate`

Register at [/developers](/developers). Lobby syncs your manifest on connect.

---

## 4. Seat manifest (`seatTemplate`)

Each mode needs `minPlayers`, `maxPlayers`, and a `seatTemplate` Lobby expands into join buckets. Games receive a **final seat map** — you don't run matchmaking.

Flat `seats[]` arrays are **rejected**. Use `count` or nested `Team` / role nodes.

Full reference: [seat-templates-and-matchmaking.md](./seat-templates-and-matchmaking.md).

---

## 5. Provision contract

`POST {apiBaseUrl}/api/v1/matches`

- **Auth:** `Authorization: Bearer {serviceToken}` from dashboard (or provision body `lobby.serviceToken`)
- **Idempotent** on `assignment.externalMatchId`
- **Success:** `200/201` with `launchUrls` map (lobbyUserId → URL base, no JWT) or `launchUrlTemplate`
- **Banned roster:** `403` with `{ "error": "...", "bannedLobbyUserIds": ["..."] }`

### Troubleshooting

<a id="provision-403-banned"></a>

**`provision.banlist` fail** — Return 403 with a `bannedLobbyUserIds` string array. To verify, temporarily ban test user `a0000000-0000-4000-8000-000000000099` during checklist runs.

<a id="provision-auth"></a>

**`provision.auth` fail** — Accept the dashboard `serviceToken` on the `Authorization` header.

<a id="provision-launch-urls"></a>

**`provision.launch_urls` fail** — Every seated player in the assignment needs a `launchUrls` entry (or a usable `launchUrlTemplate`).

---

## 6. JWT + claim

Lobby mints seat JWTs (`iss`, `aud` = your API base URL, `matchId`, `seatKey`, `sub` = lobby user id).

- Publish JWKS at `{lobbyIssuer}/.well-known/jwks.json`
- **Claim:** `POST {apiBaseUrl}/api/v1/matches/{externalMatchId}/claim` with `Authorization: Bearer {jwt}`
- Reject wrong `aud`, wrong `iss`, unknown match (`404`), bad token (`401`)

---

## 7. Launch URLs (game-minted)

Required. Lobby attaches `token=<jwt>` to each URL base you return. See [game-minted-launch-urls.md](./game-minted-launch-urls.md).

---

## 8. Testing with friends

When manifest checks pass:

1. Open your developer dashboard → **Create test table**
2. Invite friends via your room link
3. They can play while the game stays off the public catalog

---

## 9. Public release

When integration checks are green and catalog metadata is complete (`shortDescription`, `longDescription`, at least one tag):

1. Dashboard → **Request public release**
2. JoinQuest reviews (spam/IP/offensive names)
3. Approved games appear in the main catalog

---

## 10. Checklist ID index

| Check ID | Section |
|----------|---------|
| `manifest.reach_api` | §3 Quick start |
| `manifest.status` | §3 Quick start |
| `manifest.game_modes` | §4 Seat manifest |
| `manifest.sync_freshness` | §3 Quick start |
| `provision.happy_path` | §5 Provision |
| `provision.auth` | [#provision-auth](#provision-auth) |
| `provision.banlist` | [#provision-403-banned](#provision-403-banned) |
| `provision.launch_urls` | [#provision-launch-urls](#provision-launch-urls) |
| `jwt.jwks` | §6 JWT |
| `jwt.claim_happy_path` | §6 JWT |
| `jwt.wrong_audience` | §6 JWT |
| `jwt.unknown_match` | §6 JWT |
