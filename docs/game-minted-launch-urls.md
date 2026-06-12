# Game-minted launch URLs

**Status:** Implemented (Lobby backend + reference games `demo-game-rps`, `word-hunt`).

## Model

**The game mints (or specifies) the launch URL; Lobby attaches auth.**

After a successful provision, the game response tells Lobby where to send each player (or a template Lobby fills per seat). Lobby **always** attaches the seat JWT via query parameters — regardless of whether `externalMatchId` lives in the path or query string.

```text
Provision (S2S)
  Lobby  →  POST /api/v1/matches  { assignment, lobby, … }
  Game   ←  201 { launchUrls?: { [lobbyUserId]: string }, launchUrlTemplate?: string, …match state… }

Link-out (browser)
  Lobby merges game URL base + fresh seat JWT (`token` query param)
  Player opens final URL
```

Games may return:

| Style | Example game returns | Lobby adds |
|-------|----------------------|------------|
| Query match id | `https://play.example.com/?match=abc&seat=1` | `&token=<jwt>` |
| Path match id | `https://play.example.com/m/abc?seat=clue` | `&token=<jwt>` |
| Per-player path | `https://play.example.com/join/abc/seat-red` | `?token=<jwt>` |
| Template | `https://play.example.com/m/{matchId}/s/{seatKey}` | `?token=<jwt>` after substitution |

**Required:** provision must return `launchUrls` (or `launchUrlTemplate`). There is no catalog `playUrl` fallback.

## Design principles

1. **Game owns link shape** — path vs query for match id, CDN hosts, A/B routes, etc.
2. **Lobby owns auth** — JWT is always appended by Lobby (query params), same signer as before.
3. **Partial game URLs fail the match** — if a game returns `launchUrls` but omits a seated player, finalize rolls back (no silent `MATCHED` without a link).
4. **Idempotent provision** — re-push same `externalMatchId` returns the same launch URL bases.
5. **No launch URLs, no match** — games that omit both `launchUrls` and `launchUrlTemplate` fail finalize.

## Resolved design decisions

| Question | Decision |
|----------|----------|
| Per-player vs per-match URL | **Per-player** — `launchUrls: Record<lobbyUserId, string>` |
| Response schema | `launchUrls` map; optional `launchUrlTemplate` with `{matchId}`, `{externalMatchId}`, `{seatKey}`, `{lobbyUserId}` |
| Catalog fallback | **No** — provision must supply launch URL bases |
| Persistence | Yes — `game_session_participants.launch_url_base` at finalize; refresh re-signs JWT only |
| URL validation | Public HTTPS (or localhost in dev) via `gameurl.ValidateOutboundURL`; no origin pinning |
| Capability discovery | `GET /api/v1/status` → `launchUrlsOnProvision: true` |

## JWT attachment

Lobby uses `gameurl.AttachSeatToken(base, token)` — merges `token` into the query string without breaking path- or query-style games. Game URL bases must **not** include a JWT.

## Logging (Lobby)

Structured `log.Printf` lines for ops/analytics (no PII):

| Event | Example prefix |
|-------|------------------|
| Provision start | `handoff: provision start session=… game=… seats=N` |
| Provision ok | `handoff: provision ok … latency_ms=… launch_urls=N source=game\|game-template\|none` |
| Provision fail / banned | `handoff: provision fail …` / `handoff: provision banned …` |
| Finalize | `handoff: finalize ok session=… notified=N` |
| Refresh mint | `handoff: launch url mint session=… user=… source=stored` |

Games log JSON: `{"event":"provision.launch_urls",…}` on each provision.

## Reference implementations

| Game | URL style | Config |
|------|-----------|--------|
| `demo-game-rps` | Query: `{GAME_PLAY_URL}/?match=&seat=` | `GAME_PLAY_URL` (defaults to `http://localhost:5174` locally) |
| `word-hunt` | Path: `{GAME_PLAY_URL}/m/{externalMatchId}?seat=` | same; client parses `/m/:id` |

## Related

- Wire contract: [`lobby-protocol-handoff.md`](./lobby-protocol-handoff.md)
- Partner verification: [`end-to-end-partner-checklist.md`](./end-to-end-partner-checklist.md)
- Migration: `backend/migrations/000024_session_launch_urls.up.sql`
- Drop catalog `play_url`: `backend/migrations/000035_drop_play_url.up.sql`
