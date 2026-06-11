# End-to-end partner checklist

Use this to verify the **“sign in → queue → play → return”** story on local or staging before inviting external game authors.

**Product context:** [vision.md](./vision.md) · **Wire contract:** [lobby-protocol-handoff.md](./lobby-protocol-handoff.md)

---

## Prerequisites

- [ ] PostgreSQL migrated through `000024` (session launch URL persistence)
- [ ] JoinQuest API + frontend running (`./scripts/dev.sh` or deployed stack)
- [ ] Reference game registered with valid `play_url` and `api_base_url` (e.g. `demo-game-rps`)
- [ ] Game exposes `GET /healthz`, `GET /api/v1/game-modes`, accepts `POST /api/v1/matches`
- [ ] Game returns `launchUrls` for every seated player on provision (or documents catalog fallback); `/api/v1/status` reports `launchUrlsOnProvision: true` when opted in
- [ ] Game API has `GAME_PLAY_URL` matching catalog `play_url` (browser origin for minted links)
- [ ] Game verifies Lobby JWT via `/.well-known/jwks.json`
- [ ] `LOBBY_GAME_TOKEN_PEPPER` set in prod/staging (per-game `serviceToken` on provision)

---

## Player flow (manual)

1. **Sign in** on JoinQuest (magic link or login code).
2. **Catalog** lists the game with an active default queue.
3. **Look for group** — only one waiting search globally; joining another game **moves** the player out of the first and shows `joinQueue.message`.
4. **Second browser / incognito** — second player joins the same queue.
5. **Match** — both receive `joinUrl` (subscription or `myQueueStatus` / `myActiveIntent`).
6. **Launch** — each player opens `joinUrl`, game claims `seatKey` from JWT.
7. **Play** — match runs on the game origin (not JoinQuest).
8. **Return** — game sends player to `{lobby.returnUrl}?match={externalMatchId}`; Lobby `/return` resolves per-player destination. See [player-return-routing.md](./player-return-routing.md).
9. **Callbacks** — game calls `reportMatchResult` (and optionally `reportPlayerFinished`); see [match-lifecycle-callbacks.md](./match-lifecycle-callbacks.md).

---

## Failure paths to exercise

- [ ] **403 banned** on provision — Lobby re-matchmakes without other players stuck mid-flow.
- [ ] **Switch queue** — join game B while waiting in game A; A’s queue count drops; player sees removal message on B.
- [ ] **Leave queue** — `leaveQueue` clears waiting.
- [ ] **Matched + try another queue** — rejected until match is left/expired or completed.
- [ ] **Stale matched row** — `expireStaleMatchedModeQueue` allows re-queue after TTL (integration tests cover this).

---

## Automated coverage in this repo

- `backend/graph/handoff_integration_test.go` — provision payload, seat keys, `returnUrl`, `serviceToken`
- `backend/graph/handoff_game_urls_test.go` — game-minted URL merge, JWT attach, persistence
- `backend/graph/queue_integration_test.go` — join, websocket `queueUpdated`, two-player match
- `backend/internal/store/mode_queue_one_global_test.go` — one waiting queue across games

---

## Staging sign-off

- [ ] Public JoinQuest URL serves frontend + `/graphql`
- [ ] Game `api_base_url` is reachable from Lobby pods (not localhost)
- [ ] JWKS URL reachable from game pods
- [ ] `play_url` uses HTTPS in production
- [ ] Smoke test with two real accounts on mobile + desktop
