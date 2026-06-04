# Composition queues and “join as” options

This clarifies **composition matchmaking** (“Join as DPS / Tank”) vs a **character-select screen** inside the game.

**Seat / queue spec:** [seat-templates-and-matchmaking.md](./seat-templates-and-matchmaking.md)

---

## What “Join as” means in JoinQuest (today / near term)

For **composition** modes, the template defines **queue paths** (e.g. `Offense.Wizard`, `Defense`). The player picks a **role bucket** before entering matchmaking — not necessarily a specific character skin.

That is **lobby-side queue selection**, similar to picking a role in a MOBA queue. It is **not** a full character-select UI with animations and loadouts.

For **fifo** modes (simple 1v1 / FFA), the JoinQuest UI shows **“Look for group”** with no role choice (API field remains `joinQueue`).

---

## What belongs on the game site

Many games need **eligibility** that Lobby cannot guess from a static manifest:

- Which classes or decks this **account** may play
- **DLC** or season pass gates
- **Progression unlocks** (“after 10 games, unlock Advanced class”)
- Balance rules (“no hard support in solo queue”)

Those should come from the **game server**, not from a fixed list in `seatTemplate`.

### Recommended pattern (target)

1. Player chooses mode on JoinQuest (and role bucket if composition).
2. Before or during join, Lobby calls the game (authenticated as that player), e.g.  
   `GET /api/v1/players/{lobbyUserId}/queue-options?modeKey=…&queuePath=…`
3. Game returns allowed options: `{ "choices": [{ "id": "wizard", "label": "Wizard", "locked": false }, …] }`
4. JoinQuest UI shows only allowed choices; `joinQueue` sends the selected `queuePath` / preference in `PartyInput` (future).
5. After matchmaking, provision still sends **`seatKey`** — the game maps role → concrete character in-client if needed.

Lobby caches responses **briefly** (seconds), not as source of truth. DLC and unlock changes stay on the game.

---

## Why not only a manifest?

`seatTemplate` describes **layout and affinity** (teams, slots, capacities). It does not know player inventory, MMR brackets, or unlock progress. Duplicating unlock logic in the manifest would drift from the game and break DLC/progression models.

---

## Character select screen?

| Screen | Where | When |
|--------|-------|------|
| Role / queue bucket | JoinQuest | Before queue, drives matchmaking |
| Eligibility list | JoinQuest (data from game API) | Before queue, optional step |
| Cosmetic / loadout / map character | **Game client** | After JWT claim, before or during match |

So: **not** replacing the game’s character select — **optional** pre-queue eligibility on JoinQuest, fed by the game.

---

## v1 scope

- **Shipped:** fifo default queue per mode; no composition UI.
- **Spec’d:** composition paths from template; `queuePath` on party input.
- **Deferred:** game `queue-options` API, DLC-aware polling, unlock progression.

When implementing composition UI, start with **static paths from manifest** for demos; add **game-backed options** before production titles with progression.
