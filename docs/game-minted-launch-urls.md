# Game-minted launch URLs

**Status:** HIGH PRIORITY TODO — not implemented. Protocol + Lobby change so games tell us where to send players.

## Problem

Today Lobby builds launch URLs locally after provision:

1. `POST /api/v1/matches` — push roster (game accepts/rejects)
2. Lobby signs a seat JWT
3. Lobby concatenates catalog `playUrl` + `?match=<id>&token=<jwt>`

That assumes every game uses the same link shape. Games that want match id in the path, custom slugs, or environment-specific play URLs cannot express that during provision. Refresh paths also re-provision and re-sign just to rebuild a link we already had at match time.

## Target model

**The game mints (or specifies) the launch URL; Lobby attaches auth.**

After a successful provision, the game response tells Lobby where to send each player (or a template Lobby fills per seat). Lobby **always** attaches the seat JWT via query parameters — regardless of whether `externalMatchId` lives in the path or query string.

```text
Provision (S2S)
  Lobby  →  POST /api/v1/matches  { assignment, lobby, … }
  Game   ←  200 { launchUrls?: { [lobbyUserId]: string }, launchUrlTemplate?: string, … }

Link-out (browser)
  Lobby merges game URL + JWT query param(s)
  Player opens final URL
```

Games may return:

| Style | Example game returns | Lobby adds |
|-------|----------------------|------------|
| Query match id | `https://play.example.com/?match=abc` | `&token=<jwt>` (and `&seat=` if needed) |
| Path match id | `https://play.example.com/m/abc` | `?token=<jwt>` |
| Per-player path | `https://play.example.com/join/abc/seat-red` | `?token=<jwt>` |
| Template | `https://play.example.com/m/{matchId}/s/{seatKey}` | `?token=<jwt>` after substitution |

Lobby must **not** assume `playUrl` from catalog is the only base — catalog `playUrl` becomes fallback for games that do not yet return URLs.

## Design principles

1. **Game owns link shape** — path vs query for match id, CDN hosts, A/B routes, etc.
2. **Lobby owns auth** — JWT is always appended by Lobby (query params), same signer as today.
3. **Provision failure = no match** — if the game does not return a usable URL for a notified player, treat as provision failure (no silent `MATCHED` without a link).
4. **Idempotent provision** — re-push same `externalMatchId` returns the same launch URLs.
5. **Backward compatible** — games that omit launch URLs keep today's Lobby-built links until they opt in (capability flag or non-empty response field).

## Open design questions

- [ ] **Per-player vs per-match URL** — one URL + seat claim only, or distinct URL per `lobbyUserId`?
- [ ] **Response schema** — `launchUrls: Record<lobbyUserId, string>` vs `launchUrlTemplate` with `{matchId}`, `{seatKey}`, `{lobbyUserId}` placeholders?
- [ ] **Catalog fallback** — when game returns nothing, keep building from `playUrl` (current behavior)?
- [ ] **Persistence** — store minted URLs on session at finalize time so refresh never re-provisions just to rebuild links?
- [ ] **URL validation** — allowlist hosts against catalog `playUrl` origin to prevent open redirects?

## Implementation checklist

### Protocol & docs
- [ ] Extend provision response schema in [`lobby-protocol-handoff.md`](./lobby-protocol-handoff.md) and partner checklist
- [ ] Document JWT attachment rules (always query param; merge with existing query string)
- [ ] Add capability discovery (`/api/v1/status` or game-modes) for `launchUrlsOnProvision`

### Backend (Lobby)
- [ ] `gameclient.ProvisionMatch` — parse launch URL fields from response body
- [ ] `handoff.go` — `finalizeMatchedSession` prefers game URLs; fallback to local build
- [ ] URL merge helper — attach `token` (and optional `seat`, `lobby_user`) without breaking path/query games
- [ ] Stop swallowing launch URL errors on `myActiveIntent` / persist URLs at finalize
- [ ] Reconcile commit order: provision + URLs before session commit, or rollback session on failure

### Reference game(s)
- [ ] Update `demo-game-rps` (or first partner game) to return launch URLs on provision
- [ ] Integration tests: query-style URL, path-style URL, template URL

### Frontend
- [ ] No link-shape assumptions beyond opening the URL Lobby returns
- [ ] Banner/subscriptions use stored or event URLs only (already moving this direction)

## Related

- Current handoff: [`lobby-protocol-handoff.md`](./lobby-protocol-handoff.md) — link-out row assumes Lobby-built `{playUrl}?match=&token=`
- [`end-to-end-partner-checklist.md`](./end-to-end-partner-checklist.md) — step 5–6 (match → launch)
- Intent banner join URL delivery fixes (frontend + subscription) — interim; this spec is the durable fix
