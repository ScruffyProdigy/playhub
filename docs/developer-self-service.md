# Developer self-service — game registration & integration

**Status:** Phase A shipped (registration, dashboard, private testing)  
**Next:** Phase B — provision/JWT checks, public release review, MCP  
**Related:** [game-minted launch URLs](./game-minted-launch-urls.md) ✅ · [player experience roadmap](./player-experience-roadmap.md) (optional catalog polish)  
**Related:** [game-catalog-architecture.md](./game-catalog-architecture.md) · [lobby-protocol-handoff.md](./lobby-protocol-handoff.md) · [end-to-end-partner-checklist.md](./end-to-end-partner-checklist.md)

Today, adding a game means backfilling the database or calling admin-only `registerGame`. That works for us, but it blocks the goal of **any web developer** plugging in their game. This spec is the v1 path from homepage curiosity → registered game → working integration → friends testing → public release.

---

## Product principles

1. **Register freely; list publicly when ready.** Integration is open. Showing up in the main player catalog is reviewed — so players see real titles, not half-finished experiments.
2. **Fail closed with friendly errors.** If we can't reach your API or provision fails, we say *why* and what to fix — not “something went wrong.”
3. **Private until you're proud.** Your game isn't mixed into the public arcade until you ask and we approve.
4. **Agents welcome.** The same checklist humans see in the dashboard should be runnable by Claude/Cursor via MCP — happy paths and error paths — without the developer clicking through docs.

---

## v1 scope decisions

| Topic | v1 | Later |
|-------|-----|-------|
| Ownership | One JoinQuest user owns each game | Studios / teams (v2) |
| Auth | Same magic-link account as players | Separate dev portal optional |
| API URL | Public HTTPS only — **no `localhost`** | Optional tunnel/CLI helper (v2+) |
| Play URL | From game provision response ([spec](./game-minted-launch-urls.md)) | — |
| Player counts | From manifest | — |
| Icons / tags | See [player experience roadmap](./player-experience-roadmap.md) — likely before or with Phase A | Required for public release? TBD |
| PDF guide | Hold or generate from same source as agent doc | Phase C |

**Why no localhost:** Lobby servers must reach the game API to sync manifests and provision matches. A URL only you can hit from your laptop can't work in production. Devs should use a staging deploy, ngrok, or similar for early integration — we'll say that plainly in copy and errors.

---

## Visibility model (recommendation)

**Do not show private games in the public player catalog.** That list is for games anyone can queue into. Mixing in “my unfinished thing” confuses players and makes the homepage feel broken.

Instead, use **two surfaces**:

```text
Homepage (signed in)
├── Public catalog          →  status = public, checklist green
└── “Your games” strip      →  owner’s games (draft / testing / pending review)
        └── link to Developer dashboard

Developer dashboard (/developers/games/:id)
├── Checklist + errors
├── Integration credentials (serviceToken, webhook secret)
└── Actions: test table, request public release

Room (invite link)
└── Tables for any game the room members are allowed to play
    (owner’s private games + public games)
```

| Surface | Who sees private games |
|---------|----------------------|
| Public `games` query / main catalog | Nobody |
| Homepage “Your games” (owner signed in) | Owner |
| Developer dashboard | Owner |
| Room tables | Everyone in that room |

**Table rule:** If the owner creates a table for their game, **everyone in the room** can see and join that table — same as today. The game stays off the public catalog; the room is the share link for friends.

**`createPrivateTable`** already creates a room when the user doesn't have one — onboarding can say “Create a test table” as one step, not three.

---

## User-facing copy (casual / friendly)

### Homepage CTA

Link text:

> **Making a game?** We'd love to host it — [get started →](/developers)

### Landing page (`/developers`)

**Headline:** Have an idea for a multiplayer game?

**Body:**

> JoinQuest handles the boring stuff — player accounts, rooms, tables, matchmaking, finding players — so you can focus on building the game.
>
> You'll need some basic web dev experience and a **public URL** where your game can run. That's it.
>
> Register below to get started. You're not committing to anything, and your game won't show up for other players until **you** say it's ready (and we've had a quick look).

### Registration form

**Required**

- Game name  
- Short description (1–2 sentences for your catalog card later)  
- API base URL (HTTPS origin, e.g. `https://mygame.example.com`)  
- Contact email (spec updates & review — can differ from sign-in email)  
- Slug (auto-suggested from name, editable)

**Optional**

- Website  
- Discord (or other community link)

**Not on form**

- `playUrl` — game returns launch URLs on provision  
- Min/max players — from manifest  
- Icons / tags — follow-up or pre-work item

**Submit behavior (two-phase)**

1. **Save draft** — name, slug, contact, URLs stored; game `visibility = draft`.  
2. **Connect API** — Lobby calls `healthz` + `game-modes`; on success → `visibility = private_testing` and dashboard unlocks. On failure → stay draft, show actionable error (“We couldn't reach `https://…` — is it up? We can't use localhost.”).

### Post-registration (`/developers/games/:id/welcome`)

> **Nice — your game is registered.**
>
> Next up: teach your game how to talk to JoinQuest.
>
> **Using Claude, Cursor, or another AI assistant?** Turn on the JoinQuest MCP and point your agent at our integration guide — it can run the same checks as the dashboard and fix issues for you.
>
> **Prefer doing it yourself?** Same guide works for humans — step by step, no agent required.
>
> When things look good, spin up a **test table** and invite friends into your room. They can play before your game goes public.
>
> Ready for the world? Hit **Request public release** on your dashboard when the checklist is green.

### Public listing expectation (footer on dashboard + vision)

> Registering is free and instant. **Showing up in the main JoinQuest catalog** means a quick review — mostly so names aren't spam, offensive, or obvious IP issues. Your game still runs on your site either way.

---

## Game lifecycle states

```text
draft
  →  (API connect succeeds)  →  private_testing
  →  (developer requests)    →  pending_review
  →  (admin approves)        →  public
  →  (admin rejects)         →  private_testing  (with reason)
```

| State | Public catalog | Owner dashboard | Room tables |
|-------|----------------|-----------------|-------------|
| `draft` | Hidden | Yes | No (manifest not synced) |
| `private_testing` | Hidden | Yes | Yes (owner + room members) |
| `pending_review` | Hidden | Yes (read-only release btn) | Yes |
| `public` | Listed | Yes | Yes |

**v1 ownership:** `games.owner_user_id` → registering user. Admin (`LOBBY_ADMIN_EMAILS`) retains override + `registerGame` for internal titles.

---

## Registration & catalog (technical)

Extend today’s admin `registerGame` into owner-scoped **`registerMyGame`** (or same mutation with auth branch):

- Input aligns with form above (no `playUrl` once game-minted URLs ship).  
- Output: `game`, `webhookSecret`, `serviceToken` (same as today).  
- Manifest sync: existing `RegisterGame` / `ApplyGameManifest` pipeline.  
- `ListCatalogGames`: only `visibility = public` and active mode queues.  
- New: `myGames` query for owner dashboard + homepage strip.

**Contact email:** `games.contact_email` — notifications for spec changes and review outcomes.

**Optional links:** `website_url`, `community_url` (Discord etc.).

---

## Developer dashboard

Route: `/developers/games/:id` (signed-in owner only).

### Checklist sections

Each row: **status** (pass / fail / not run), **last checked**, **plain-language explanation**, **technical detail** (expandable), **Re-run** button.

#### 1. Manifest

| Check | Pass | Fail message (example) |
|-------|------|-------------------------|
| Reach API | `GET {apiBaseUrl}/healthz` → `ok` | “We couldn't reach your server. Check the URL and that it's HTTPS and publicly reachable (not localhost).” |
| Status | `GET /api/v1/status` | “Your game didn't return a status payload. See integration guide § Status.” |
| Game modes | `GET /api/v1/game-modes` valid `seatTemplate` | “Your seat manifest has an error: …” |
| Sync freshness | Manifest hash / ETag | “Your cached manifest is stale — …” |

#### 2. Provisioning

| Check | Pass | Fail message (example) |
|-------|------|-------------------------|
| Happy path | `POST /api/v1/matches` 200, idempotent re-push | “Provision failed: {status} — …” |
| Auth | Bearer `serviceToken` accepted | “Your server rejected our service token. Compare with dashboard credential.” |
| Banlist | `403` + `bannedLobbyUserIds` shape | “When banning a test user, return 403 with `bannedLobbyUserIds` array — we got …” |
| Launch URLs | Response includes per-player `launchUrls` (required when opted in; catalog fallback if omitted) | “Provision succeeded but didn't return launch URLs for all seats.” |

#### 3. JWT verification

| Check | Pass | Fail message (example) |
|-------|------|-------------------------|
| JWKS | Lobby JWKS reachable from game’s perspective (or we document game-side fetch) | “We couldn't verify JWKS setup — …” |
| Claim happy path | Test token → `POST …/claim` 200 | “Claim failed: …” |
| Wrong audience | Token with bad `aud` rejected | “Your game accepted a token with the wrong audience — `aud` must be your API base URL.” |
| Unknown match | Claim on bogus match → 404 | “Unknown match should return 404, got …” |

**Implementation note:** Run checks server-side on demand and on a schedule (e.g. every few hours for `private_testing` / `public` games). Store last result JSON on `game_integration_checks` or similar.

### Dashboard actions

- Copy `serviceToken`, webhook secret (masked + reveal)  
- **Run all checks**  
- **Create test table** → `createPrivateTable` + open room panel  
- **Request public release** (enabled when required checks pass + name/description present)

---

## Public release

**Developer prerequisites (automated gate):**

- Required checklist items green  
- Game name + description set  
- HTTPS API URL  
- (Optional later) icon present  

**Human review (admin queue):**

- Name not junk / offensive / obvious infringement  
- Description reasonable  
- Approve → `visibility = public`  
- Reject → back to `private_testing` with reason email to `contact_email`

**Admin UI v1:** Could be GraphQL + script or minimal `/admin/games/pending` — detail TBD.

---

## JoinQuest MCP (Phase B)

Goal: agent runs the **same probes** as the dashboard without the developer in the loop.

**Auth:** Developer session token or API key tied to `owner_user_id`.

**Tools (sketch):**

| Tool | Purpose |
|------|---------|
| `joinquest_get_integration_guide` | Returns markdown guide (single source of truth) |
| `joinquest_list_my_games` | Owner's games + states |
| `joinquest_get_game_checks` | Latest checklist results + errors |
| `joinquest_run_game_checks` | Run manifest / provision / JWT suite |
| `joinquest_get_game_credentials` | serviceToken, webhook URL (masked) |
| `joinquest_get_example_provision_payload` | Sample assignment for copy-paste |

MCP calls backend endpoints shared with dashboard — **one implementation**, two clients (UI + MCP).

**Cursor/Claude setup:** Published config snippet on `/developers` (“Add to Cursor”).

---

## Integration guide (single source of truth)

**File:** `docs/developer-integration-guide.md` (to be written alongside Phase B MCP)

**Consumers:**

- Human-readable web page (`/developers/guide`)  
- MCP `joinquest_get_integration_guide`  
- PDF export (Phase C — generate from same markdown)

**Structure:**

1. Quick start (health, status, game-modes)  
2. Seat manifest / `seatTemplate`  
3. Provision contract + banlist  
4. JWT + claim + JWKS  
5. Launch URLs (game-minted)  
6. Testing with a private table + room invite  
7. Requesting public release  
8. Troubleshooting index mapped 1:1 to dashboard check IDs  

Dashboard error strings link to anchor IDs in this doc (`#provision-403-banned`).

---

## Phasing

### Phase A — Critical path ✅ Shipped

Registration, private testing, and the developer dashboard:

- Homepage CTA + `/developers` landing  
- Sign-in-gated registration form + draft / connect API  
- `owner_user_id`, visibility states, `myGames` query  
- Developer dashboard with checklist (manifest live; provision/JWT stubbed until Phase B)  
- Homepage “Your games” strip for owners  
- Private testing: room tables, not public catalog  
- Post-register welcome + integration guide stub  
- **Test table** CTA (`createPrivateTable`)  

**Dependency:** [game-minted launch URLs](./game-minted-launch-urls.md) — shipped; handoff requires provision `launchUrls`.

**Optional polish:** game icons + catalog tags (see [player experience roadmap](./player-experience-roadmap.md)).

### Phase B — Go public + agents

- Request public release + admin review queue  
- JoinQuest MCP + full `developer-integration-guide.md`  
- Scheduled re-checks + email on spec-breaking manifest changes  

### Phase C — PDF (optional / defer)

- Generate printable guide from same markdown as MCP  
- Hold until Phase B feedback settles  

---

## Open questions (minor)

- [ ] Icons required before public release, or optional v1?  
- [ ] Admin review UI: in-app vs email + GraphQL script?  
- [ ] Rate limits on `run checks` / synthetic provision (per game)?  
- [ ] Webhook URL shown on dashboard for manifest-changed setup?  

---

## Migration sketch

```sql
-- games
ALTER TABLE games ADD COLUMN owner_user_id UUID REFERENCES users(id);
ALTER TABLE games ADD COLUMN visibility TEXT NOT NULL DEFAULT 'draft'
  CHECK (visibility IN ('draft', 'private_testing', 'pending_review', 'public'));
ALTER TABLE games ADD COLUMN contact_email TEXT;
ALTER TABLE games ADD COLUMN website_url TEXT;
ALTER TABLE games ADD COLUMN community_url TEXT;

-- optional: game_integration_check_runs (game_id, check_id, status, message, detail_json, ran_at)
```

Existing catalog games: backfill `owner_user_id = NULL`, `visibility = public` (admin-owned).

---

## Success criteria

- [ ] New developer registers without DB access or admin email  
- [ ] Dashboard shows a failing manifest check with an actionable message  
- [ ] Owner creates test table; friend joins room and plays  
- [ ] Game stays off public catalog until approved  
- [ ] MCP runs full checklist; agent can fix a deliberate 403 banlist mistake from error output  
- [ ] Integration guide section IDs match dashboard deep links  
